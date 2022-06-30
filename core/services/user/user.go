package user

import (
	"context"
	"fmt"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/sergicanet9/go-mongo-restapi/config"
	"github.com/sergicanet9/go-mongo-restapi/core/domain"
	"github.com/sergicanet9/go-mongo-restapi/core/dto/requests"
	"github.com/sergicanet9/go-mongo-restapi/core/dto/responses"
	"github.com/sergicanet9/go-mongo-restapi/core/ports"
	"golang.org/x/crypto/bcrypt"
)

//Service struct
type Service struct {
	config     config.Config
	repository ports.UserRepository
}

// NewUserService creates a new user service
func NewUserService(cfg config.Config, repo ports.UserRepository) *Service {
	return &Service{
		config:     cfg,
		repository: repo,
	}
}

// Login user
func (s *Service) Login(ctx context.Context, credentials requests.LoginUser) (responses.LoginUser, error) {
	filter := map[string]interface{}{"email": credentials.Email}
	result, err := s.repository.Get(ctx, filter, nil, nil)
	if err != nil {
		return responses.LoginUser{}, err
	}
	if len(result) < 1 {
		return responses.LoginUser{}, fmt.Errorf("email not found")
	}

	user := responses.User(**result[0].(**domain.User))

	if checkPasswordHash(credentials.Password, user.PasswordHash) {
		token, err := createToken(user.ID, s.config.JWTSecret, user.Claims)
		if err != nil {
			return responses.LoginUser{}, err
		}

		result := responses.LoginUser{
			User:  user,
			Token: token,
		}
		return result, nil
	}
	return responses.LoginUser{}, fmt.Errorf("incorrect password")
}

//Create user
func (s *Service) Create(ctx context.Context, u requests.User) (responses.Creation, error) {
	err := hashPassword(&u.PasswordHash)
	if err != nil {
		return responses.Creation{}, err
	}

	err = validateClaims(u.Claims)
	if err != nil {
		return responses.Creation{}, err
	}

	now := time.Now().UTC()
	u.CreatedAt = now
	u.UpdatedAt = now
	insertedID, err := s.repository.Create(ctx, domain.User(u))
	if err != nil {
		return responses.Creation{}, err
	}
	return responses.Creation{InsertedID: insertedID}, nil
}

// GetAll users
func (s *Service) GetAll(ctx context.Context) ([]responses.User, error) {
	result, err := s.repository.Get(ctx, map[string]interface{}{}, nil, nil)
	if err != nil {
		return []responses.User{}, err
	}

	users := make([]responses.User, len(result))
	for i, v := range result {
		users[i] = responses.User(**(v.(**domain.User)))
	}

	return users, nil
}

//GetByEmail user
func (s *Service) GetByEmail(ctx context.Context, email string) (responses.User, error) {
	filter := map[string]interface{}{"email": email}
	result, err := s.repository.Get(ctx, filter, nil, nil)
	if err != nil {
		return responses.User{}, err
	}
	if len(result) < 1 {
		return responses.User{}, fmt.Errorf("email not found")
	}
	return responses.User(**(result[0].(**domain.User))), nil
}

// GetByID user
func (s *Service) GetByID(ctx context.Context, ID string) (responses.User, error) {
	user, err := s.repository.GetByID(ctx, ID)
	if err != nil {
		return responses.User{}, err
	}
	return responses.User(*user.(*domain.User)), nil
}

// Update user
func (s *Service) Update(ctx context.Context, ID string, u requests.UpdateUser) error {
	result, err := s.repository.GetByID(ctx, ID)
	if err != nil {
		return err
	}
	user := *result.(*domain.User)
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
	if u.Claims != nil {
		err = validateClaims(*u.Claims)
		if err != nil {
			return err
		}
		user.Claims = *u.Claims
	}
	user.ID = ""
	user.UpdatedAt = time.Now().UTC()

	if err := s.repository.Update(ctx, ID, user, false); err != nil {
		return err
	}
	return nil
}

// Delete user
func (s *Service) Delete(ctx context.Context, ID string) error {
	err := s.repository.Delete(ctx, ID)
	return err
}

// Get user claims
func (s *Service) GetClaims(ctx context.Context) (map[int]string, error) {
	return domain.GetClaims(), nil
}

// // AtomicTransationProof creates two entities atomically, creating a sessionContext
// func (s *Service) AtomicTransationProof(ctx context.Context) error { //TODO!!!
// 	wc := writeconcern.New(writeconcern.WMajority())
// 	rc := readconcern.Snapshot()
// 	txnOpts := options.Transaction().SetWriteConcern(wc).SetReadConcern(rc)

// 	session, err := s.db.Client().StartSession()
// 	if err != nil {
// 		return err
// 	}
// 	defer session.EndSession(ctx)

// 	user1Hash := "Entity1"
// 	err = hashPassword(&user1Hash)
// 	if err != nil {
// 		return err
// 	}
// 	user2Hash := "Entity2"
// 	err = hashPassword(&user2Hash)
// 	if err != nil {
// 		return err
// 	}

// 	callback := func(sessionContext mongo.SessionContext) (interface{}, error) {
// 		_, err = s.repository.Create(sessionContext,
// 			domain.User{
// 				ID:           primitive.NewObjectID().String(),
// 				Name:         "Entity1",
// 				Surnames:     "Entity1",
// 				Email:        "Entity1",
// 				PasswordHash: user1Hash,
// 				Claims:       nil,
// 			})
// 		if err != nil {
// 			return nil, err
// 		}

// 		_, err = s.repository.Create(sessionContext,
// 			domain.User{
// 				ID:           primitive.NewObjectID(),
// 				Name:         "Entity2",
// 				Surnames:     "Entity2",
// 				Email:        "Entity2",
// 				PasswordHash: user2Hash,
// 				Claims:       nil,
// 			})
// 		if err != nil {
// 			return nil, err
// 		}

// 		return nil, nil
// 	}

// 	_, err = session.WithTransaction(context.Background(), callback, txnOpts)
// 	return err
// }

func createToken(userid string, jwtSecret string, claims []int) (string, error) {
	var err error
	addClaims := jwt.MapClaims{}
	addClaims["authorized"] = true
	addClaims["user_id"] = userid
	addClaims["exp"] = time.Now().UTC().Add(time.Hour * 168).Unix()

	err = validateClaims(claims)
	if err != nil {
		return "", err
	}
	for _, claim := range claims {
		addClaims[domain.Claim(claim).String()] = true
	}

	add := jwt.NewWithClaims(jwt.SigningMethodHS256, addClaims)
	token, err := add.SignedString([]byte(jwtSecret))
	if err != nil {
		return "", err
	}
	return token, nil
}

func validateClaims(claims []int) error {
	for _, claim := range claims {
		if ok := domain.Claim(claim).IsValid(); !ok {
			return fmt.Errorf("not valid claim detected: %d", claim)
		}
	}
	return nil
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
