package services

import (
	"context"
	"fmt"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/sergicanet9/go-hexagonal-api/config"
	"github.com/sergicanet9/go-hexagonal-api/core/entities"
	"github.com/sergicanet9/go-hexagonal-api/core/models"
	"github.com/sergicanet9/go-hexagonal-api/core/ports"
	"golang.org/x/crypto/bcrypt"
)

// userService adapter of an user service
type userService struct {
	config     config.Config
	repository ports.UserRepository
}

// NewUserService creates a new user service
func NewUserService(cfg config.Config, repo ports.UserRepository) ports.UserService {
	return &userService{
		config:     cfg,
		repository: repo,
	}
}

// Login user
func (s *userService) Login(ctx context.Context, credentials models.LoginUserReq) (models.LoginUserResp, error) {
	filter := map[string]interface{}{"email": credentials.Email}
	result, err := s.repository.Get(ctx, filter, nil, nil)
	if err != nil {
		return models.LoginUserResp{}, err
	}
	if len(result) < 1 {
		return models.LoginUserResp{}, fmt.Errorf("email not found")
	}

	user := models.UserResp(*result[0].(*entities.User))

	if checkPasswordHash(credentials.Password, user.PasswordHash) {
		token, err := createToken(user.ID, s.config.JWTSecret, user.Claims)
		if err != nil {
			return models.LoginUserResp{}, err
		}

		result := models.LoginUserResp{
			User:  user,
			Token: token,
		}
		return result, nil
	}
	return models.LoginUserResp{}, fmt.Errorf("incorrect password")
}

func checkPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
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
		addClaims[entities.Claim(claim).String()] = true
	}

	add := jwt.NewWithClaims(jwt.SigningMethodHS256, addClaims)
	token, err := add.SignedString([]byte(jwtSecret))
	if err != nil {
		return "", err
	}
	return token, nil
}

// Create user
func (s *userService) Create(ctx context.Context, u models.UserReq) (models.CreationResp, error) {
	err := hashPassword(&u.PasswordHash)
	if err != nil {
		return models.CreationResp{}, err
	}

	err = validateClaims(u.Claims)
	if err != nil {
		return models.CreationResp{}, err
	}

	now := time.Now().UTC()
	u.CreatedAt = now
	u.UpdatedAt = now
	insertedID, err := s.repository.Create(ctx, entities.User(u))
	if err != nil {
		return models.CreationResp{}, err
	}
	return models.CreationResp{InsertedID: insertedID}, nil
}

func hashPassword(password *string) error {
	bytes, err := bcrypt.GenerateFromPassword([]byte(*password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	*password = string(bytes)
	return nil
}

func validateClaims(claims []int64) error {
	for _, claim := range claims {
		if ok := entities.Claim(claim).IsValid(); !ok {
			return fmt.Errorf("not valid claim detected: %d", claim)
		}
	}
	return nil
}

// GetAll users
func (s *userService) GetAll(ctx context.Context) ([]models.UserResp, error) {
	result, err := s.repository.Get(ctx, map[string]interface{}{}, nil, nil)
	if err != nil {
		return []models.UserResp{}, err
	}

	users := make([]models.UserResp, len(result))
	for i, v := range result {
		users[i] = models.UserResp(*(v.(*entities.User)))
	}

	return users, nil
}

// GetByEmail user
func (s *userService) GetByEmail(ctx context.Context, email string) (models.UserResp, error) {
	filter := map[string]interface{}{"email": email}
	result, err := s.repository.Get(ctx, filter, nil, nil)
	if err != nil {
		return models.UserResp{}, err
	}
	if len(result) < 1 {
		return models.UserResp{}, fmt.Errorf("email not found")
	}
	return models.UserResp(*(result[0].(*entities.User))), nil
}

// GetByID user
func (s *userService) GetByID(ctx context.Context, ID string) (models.UserResp, error) {
	user, err := s.repository.GetByID(ctx, ID)
	if err != nil {
		return models.UserResp{}, err
	}
	return models.UserResp(*user.(*entities.User)), nil
}

// Update user
func (s *userService) Update(ctx context.Context, ID string, u models.UpdateUserReq) error {
	result, err := s.repository.GetByID(ctx, ID)
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
func (s *userService) Delete(ctx context.Context, ID string) error {
	err := s.repository.Delete(ctx, ID)
	return err
}

// Get user claims
func (s *userService) GetClaims(ctx context.Context) (map[int]string, error) {
	return entities.GetClaims(), nil
}

// AtomicTransationProof creates two users atomically
func (s *userService) AtomicTransationProof(ctx context.Context) error {
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
		entities.User{
			Name:         "Entity1",
			Surnames:     "Entity1",
			Email:        "Entity1",
			PasswordHash: user1Hash,
			Claims:       nil,
			CreatedAt:    now,
			UpdatedAt:    now,
		},
		entities.User{
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
