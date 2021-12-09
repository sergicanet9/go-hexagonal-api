package user

import (
	"context"
	"fmt"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/scanet9/go-mongo-restapi/config"
	"github.com/scanet9/go-mongo-restapi/models/entities"
	"github.com/scanet9/go-mongo-restapi/models/requests"
	"github.com/scanet9/go-mongo-restapi/models/responses"
	infrastructure "github.com/scanet9/scv-go-framework/v2/infrastructure/mongo"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readconcern"
	"go.mongodb.org/mongo-driver/mongo/writeconcern"
	"golang.org/x/crypto/bcrypt"
)

//Service struct
type Service struct {
	config config.Config
	db     *mongo.Database
	repo   infrastructure.MongoRepository
}

// UserService interface represents a UserService
type UserService interface {
	Login(credentials requests.Login) (responses.Login, error)
	Create(u requests.User) (responses.Creation, error)
	GetAll() ([]responses.User, error)
	GetByEmail(email string) (responses.User, error)
	GetByID(ID string) (responses.User, error)
	Update(ID string, u requests.Update) error
	Delete(ID string) error
	AtomicTransationProof() error
}

// NewUserService creates a new user service
func NewUserService(cfg config.Config, db *mongo.Database) *Service {
	return &Service{
		config: cfg,
		db:     db,
		repo:   *infrastructure.NewMongoRepository(db.Collection(entities.CollectionNameUser), &entities.User{}),
	}
}

// Login user
func (s *Service) Login(credentials requests.Login) (responses.Login, error) {
	filter := bson.M{"email": credentials.Email}
	result, err := s.repo.Get(context.Background(), filter)
	if err != nil {
		return responses.Login{}, err
	}
	if len(result) < 1 {
		return responses.Login{}, fmt.Errorf("email not found")
	}
	user := responses.User(**result[0].(**entities.User))

	if checkPasswordHash(credentials.Password, user.PasswordHash) {
		token, err := createToken(user.ID.Hex(), s.config.JWTSecret)
		if err != nil {
			return responses.Login{}, err
		}

		result := responses.Login{
			User:  user,
			Token: token,
		}
		return result, nil
	}
	return responses.Login{}, fmt.Errorf("incorrect password")
}

//Create user
func (s *Service) Create(u requests.User) (responses.Creation, error) {
	err := hashPassword(&u.PasswordHash)
	if err != nil {
		return responses.Creation{}, err
	}

	u.ID = primitive.NewObjectID()
	insertedID, err := s.repo.Create(context.Background(), entities.User(u))
	if err != nil {
		return responses.Creation{}, err
	}
	return responses.Creation{InsertedID: insertedID}, nil
}

// GetAll users
func (s *Service) GetAll() ([]responses.User, error) {
	result, err := s.repo.Get(context.Background(), bson.M{})
	if err != nil {
		return []responses.User{}, err
	}

	users := make([]responses.User, len(result))
	for i, v := range result {
		users[i] = responses.User(**(v.(**entities.User)))
	}

	return users, nil
}

//GetByEmail users
func (s *Service) GetByEmail(email string) (responses.User, error) {
	filter := bson.M{"email": email}
	result, err := s.repo.Get(context.Background(), filter)
	if err != nil {
		return responses.User{}, err
	}
	user := responses.User(**(result[0].(**entities.User)))
	return user, nil
}

// GetByID user
func (s *Service) GetByID(ID string) (responses.User, error) {
	user, err := s.repo.GetByID(context.Background(), ID)
	if err != nil {
		return responses.User{}, err
	}
	return responses.User(*user.(*entities.User)), nil
}

// Update user
func (s *Service) Update(ID string, u requests.Update) error {
	result, err := s.repo.GetByID(context.Background(), ID)
	if err != nil {
		return err
	}
	user := *result.(*entities.User)
	if u.Name != nil {
		user.Name = *u.Name
	}
	if u.Surnames != nil {
		user.Surnames = *u.Surnames
	}
	if u.Email != nil {
		user.Email = *u.Email
	}
	if u.NewPassword != nil {
		if checkPasswordHash(*u.OldPassword, user.PasswordHash) {
			err = hashPassword(u.NewPassword)
			if err != nil {
				return err
			}

			user.PasswordHash = *u.NewPassword
		} else {
			return fmt.Errorf("old password incorrect")
		}
	}

	if err := s.repo.Update(context.Background(), ID, user, false); err != nil {
		return err
	}
	return nil
}

// Delete user
func (s *Service) Delete(ID string) error {
	if err := s.repo.Delete(context.Background(), ID); err != nil {
		return err
	}
	return nil
}

// AtomicTransationProof creates two entities atomically, creating a sessionContext
func (s *Service) AtomicTransationProof() error {
	wc := writeconcern.New(writeconcern.WMajority())
	rc := readconcern.Snapshot()
	txnOpts := options.Transaction().SetWriteConcern(wc).SetReadConcern(rc)

	session, err := s.db.Client().StartSession()
	if err != nil {
		return err
	}
	defer session.EndSession(context.Background())

	user1Hash := "Entity1"
	err = hashPassword(&user1Hash)
	if err != nil {
		return err
	}
	user2Hash := "Entity2"
	err = hashPassword(&user2Hash)
	if err != nil {
		return err
	}

	callback := func(sessionContext mongo.SessionContext) (interface{}, error) {
		s.repo.Create(sessionContext,
			entities.User{
				Name:         "Entity1",
				Surnames:     "Entity1",
				Email:        "Entity1",
				PasswordHash: user1Hash,
			})

		s.repo.Create(sessionContext,
			entities.User{
				Name:         "Entity2",
				Surnames:     "Entity2",
				Email:        "Entity2",
				PasswordHash: user2Hash,
			})
		return nil, nil
	}

	_, err = session.WithTransaction(context.Background(), callback, txnOpts)
	return err
}

func createToken(userid string, jwtSecret string) (string, error) {
	var err error
	atClaims := jwt.MapClaims{}
	atClaims["authorized"] = true
	atClaims["user_id"] = userid
	atClaims["exp"] = time.Now().Add(time.Hour * 168).Unix()
	at := jwt.NewWithClaims(jwt.SigningMethodHS256, atClaims)
	token, err := at.SignedString([]byte(jwtSecret))
	if err != nil {
		return "", err
	}
	return token, nil
}

func hashPassword(password *string) error {
	bytes, err := bcrypt.GenerateFromPassword([]byte(*password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	*password = string(bytes)
	return nil
}

func checkPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}