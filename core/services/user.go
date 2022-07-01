package services

import (
	"context"
	"fmt"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/sergicanet9/go-hexagonal-api/config"
	"github.com/sergicanet9/go-hexagonal-api/core/domain"
	"github.com/sergicanet9/go-hexagonal-api/core/dto/requests"
	"github.com/sergicanet9/go-hexagonal-api/core/dto/responses"
	"github.com/sergicanet9/go-hexagonal-api/core/ports"
	"golang.org/x/crypto/bcrypt"
)

//UserService struct
type UserService struct {
	config     config.Config
	repository ports.UserRepository
}

// NewUserService creates a new user service
func NewUserService(cfg config.Config, repo ports.UserRepository) *UserService {
	return &UserService{
		config:     cfg,
		repository: repo,
	}
}

// Login user
func (s *UserService) Login(ctx context.Context, credentials requests.LoginUser) (responses.LoginUser, error) {
	filter := map[string]interface{}{"email": credentials.Email}
	result, err := s.repository.Get(ctx, filter, nil, nil)
	if err != nil {
		return responses.LoginUser{}, err
	}
	if len(result) < 1 {
		return responses.LoginUser{}, fmt.Errorf("email not found")
	}

	user := responses.User(*result[0].(*domain.User))

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
func (s *UserService) Create(ctx context.Context, u requests.User) (responses.Creation, error) {
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
func (s *UserService) GetAll(ctx context.Context) ([]responses.User, error) {
	result, err := s.repository.Get(ctx, map[string]interface{}{}, nil, nil)
	if err != nil {
		return []responses.User{}, err
	}

	users := make([]responses.User, len(result))
	for i, v := range result {
		users[i] = responses.User(*(v.(*domain.User)))
	}

	return users, nil
}

//GetByEmail user
func (s *UserService) GetByEmail(ctx context.Context, email string) (responses.User, error) {
	filter := map[string]interface{}{"email": email}
	result, err := s.repository.Get(ctx, filter, nil, nil)
	if err != nil {
		return responses.User{}, err
	}
	if len(result) < 1 {
		return responses.User{}, fmt.Errorf("email not found")
	}
	return responses.User(*(result[0].(*domain.User))), nil
}

// GetByID user
func (s *UserService) GetByID(ctx context.Context, ID string) (responses.User, error) {
	user, err := s.repository.GetByID(ctx, ID)
	if err != nil {
		return responses.User{}, err
	}
	return responses.User(*user.(*domain.User)), nil
}

// Update user
func (s *UserService) Update(ctx context.Context, ID string, u requests.UpdateUser) error {
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

	err = s.repository.Update(ctx, ID, user)
	return err
}

// Delete user
func (s *UserService) Delete(ctx context.Context, ID string) error {
	err := s.repository.Delete(ctx, ID)
	return err
}

// Get user claims
func (s *UserService) GetClaims(ctx context.Context) (map[int]string, error) {
	return domain.GetClaims(), nil
}

// AtomicTransationProof creates two users atomically
func (s *UserService) AtomicTransationProof(ctx context.Context) error {
	user1Hash := "Entity1"
	err := hashPassword(&user1Hash)
	if err != nil {
		return err
	}
	user2Hash := "Entity2"
	err = hashPassword(&user2Hash)
	if err != nil {
		return err
	}
	now := time.Now().UTC()

	var users = []interface{}{
		domain.User{
			Name:         "Entity1",
			Surnames:     "Entity1",
			Email:        "Entity1",
			PasswordHash: user1Hash,
			Claims:       nil,
			CreatedAt:    now,
			UpdatedAt:    now,
		},
		domain.User{
			Name:         "Entity2",
			Surnames:     "Entity2",
			Email:        "Entity2",
			PasswordHash: user2Hash,
			Claims:       nil,
			CreatedAt:    now,
			UpdatedAt:    now,
		},
	}

	err = s.repository.InsertMany(ctx, users)
	return err

}

func createToken(userid string, jwtSecret string, claims []int64) (string, error) {
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

func validateClaims(claims []int64) error {
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
