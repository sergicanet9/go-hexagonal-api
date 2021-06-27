package business

import (
	"context"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/scanet9/go-mongo-restapi/config"
	"github.com/scanet9/go-mongo-restapi/models/entities"
	"github.com/scanet9/go-mongo-restapi/models/requests"
	"github.com/scanet9/go-mongo-restapi/models/responses"
	"github.com/scanet9/scv-go-framework/infrastructure"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readconcern"
	"go.mongodb.org/mongo-driver/mongo/writeconcern"
	"golang.org/x/crypto/bcrypt"
)

// UserService interface represents a UserService
type UserService interface {
	Login(credentials requests.Login) responses.Login
	Create(user entities.User) responses.Creation
	GetAll() []entities.User
	GetByEmail(email string) entities.User
	GetByID(ID string) entities.User
	Update(ID string, user entities.User)
	Delete(ID string)
	AtomicTransationProof()
}

// NewUserService creates a new user service
func NewUserService(db *mongo.Database) *Service {
	return &Service{
		db:   db,
		repo: *infrastructure.NewMongoRepository(db.Collection(entities.CollectionNameUser), func() interface{} { return &entities.User{} }),
	}
}

// Login user
func (s *Service) Login(credentials requests.Login) responses.Login {
	filter := bson.M{"email": credentials.Email}
	result := s.repo.Get(context.Background(), filter)
	if len(result) < 1 {
		panic("Email not found")
	}
	user := *result[0].(*entities.User)

	if checkPasswordHash(credentials.Password, user.PasswordHash) {
		token := createToken(user.ID.Hex())
		result := responses.Login{
			User:  user,
			Token: token,
		}
		return result
	}
	panic("Incorrect password")
}

//Create user
func (s *Service) Create(user entities.User) responses.Creation {
	user.PasswordHash = hashPassword(user.PasswordHash)
	insertedID := s.repo.Create(context.Background(), user)
	return responses.Creation{InsertedID: insertedID}
}

// GetAll users
func (s *Service) GetAll() []entities.User {
	result := s.repo.Get(context.Background(), bson.M{})

	users := make([]entities.User, len(result))
	for i, v := range result {
		users[i] = *(v.(*entities.User))
	}

	return users
}

//GetByEmail users
func (s *Service) GetByEmail(email string) entities.User {
	filter := bson.M{"email": email}
	result := s.repo.Get(context.Background(), filter)
	user := *(result[0].(*entities.User))
	return user
}

// GetByID user
func (s *Service) GetByID(ID string) entities.User {
	user := s.repo.GetByID(context.Background(), ID)
	return *user.(*entities.User)
}

// Update user
func (s *Service) Update(ID string, user entities.User) {
	s.repo.Update(context.Background(), ID, user)
}

// Delete user
func (s *Service) Delete(ID string) {
	s.repo.Delete(context.Background(), ID)
}

// AtomicTransationProof creates two entities atomically, creating a sessionContext
func (s *Service) AtomicTransationProof() {
	wc := writeconcern.New(writeconcern.WMajority())
	rc := readconcern.Snapshot()
	txnOpts := options.Transaction().SetWriteConcern(wc).SetReadConcern(rc)

	session, err := s.db.Client().StartSession()
	if err != nil {
		panic(err)
	}
	defer session.EndSession(context.Background())

	callback := func(sessionContext mongo.SessionContext) (interface{}, error) {
		s.repo.Create(sessionContext,
			entities.User{
				Name:         "Entity1",
				Surnames:     "Entity1",
				Email:        "Entity1",
				PasswordHash: "Entity1",
			})

		s.repo.Create(sessionContext,
			entities.User{
				Name:         "Entity1",
				Surnames:     "Entity1",
				Email:        "Entity1",
				PasswordHash: "Entity1",
			})
		return nil, nil
	}

	_, err = session.WithTransaction(context.Background(), callback, txnOpts)
	if err != nil {
		panic(err)
	}
}

func createToken(userid string) string {
	var err error
	atClaims := jwt.MapClaims{}
	atClaims["authorized"] = true
	atClaims["user_id"] = userid
	atClaims["exp"] = time.Now().Add(time.Hour * 168).Unix()
	at := jwt.NewWithClaims(jwt.SigningMethodHS256, atClaims)
	token, err := at.SignedString([]byte(config.JWTSecret))
	if err != nil {
		panic(err)
	}
	return token
}

func hashPassword(password string) string {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	if err != nil {
		panic(err)
	}
	return string(bytes)
}

func checkPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
