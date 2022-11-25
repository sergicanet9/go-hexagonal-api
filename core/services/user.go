package services

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/sergicanet9/go-hexagonal-api/config"
	"github.com/sergicanet9/go-hexagonal-api/core/entities"
	"github.com/sergicanet9/go-hexagonal-api/core/models"
	"github.com/sergicanet9/go-hexagonal-api/core/ports"
	"github.com/sergicanet9/scv-go-tools/v3/wrappers"
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
func (s *userService) Login(ctx context.Context, credentials models.LoginUserReq) (resp models.LoginUserResp, err error) {
	user, err := s.validateLogin(ctx, credentials)
	if err != nil {
		return
	}

	token, err := createToken(user.ID, s.config.JWTSecret, user.Claims)
	if err != nil {
		return
	}

	resp = models.LoginUserResp{
		User:  user,
		Token: token,
	}

	return
}

func (s *userService) validateLogin(ctx context.Context, credentials models.LoginUserReq) (models.UserResp, error) {
	if err := credentials.Validate(); err != nil {
		return models.UserResp{}, err
	}

	user, err := s.GetByEmail(ctx, credentials.Email)
	if err != nil {
		return models.UserResp{}, err
	}

	err = validatePassword(credentials.Password, user.PasswordHash)
	if err != nil {
		return models.UserResp{}, err
	}

	return user, nil
}

func validatePassword(password, hash string) error {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return wrappers.NewValidationErr(err)
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

func validateClaims(claims []int64) error {
	for _, claim := range claims {
		if ok := entities.Claim(claim).IsValid(); !ok {
			return wrappers.NewValidationErr(fmt.Errorf("claim %d is not valid", claim))
		}
	}
	return nil
}

// Create user
func (s *userService) Create(ctx context.Context, u models.CreateUserReq) (resp models.CreationResp, err error) {
	if err := u.Validate(); err != nil {
		return resp, err
	}

	err = hashPassword(&u.PasswordHash)
	if err != nil {
		return
	}

	err = validateClaims(u.Claims)
	if err != nil {
		return
	}

	now := time.Now().UTC()
	u.CreatedAt = now
	u.UpdatedAt = now
	insertedID, err := s.repository.Create(ctx, entities.User(u))
	if err != nil {
		return
	}

	resp = models.CreationResp{
		InsertedID: insertedID,
	}

	return
}

func hashPassword(password *string) error {
	bytes, err := bcrypt.GenerateFromPassword([]byte(*password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	*password = string(bytes)
	return nil
}

// GetAll users
func (s *userService) GetAll(ctx context.Context) (resp []models.UserResp, err error) {
	result, err := s.repository.Get(ctx, map[string]interface{}{}, nil, nil)
	if err != nil {
		if errors.Is(err, wrappers.NonExistentErr) {
			err = nil
		}
		return
	}

	resp = make([]models.UserResp, len(result))
	for i, v := range result {
		resp[i] = models.UserResp(*(v.(*entities.User)))
	}

	return
}

// GetByEmail user
func (s *userService) GetByEmail(ctx context.Context, email string) (resp models.UserResp, err error) {
	filter := map[string]interface{}{"email": email}
	result, err := s.repository.Get(ctx, filter, nil, nil)
	if err != nil {
		if errors.Is(err, wrappers.NonExistentErr) {
			err = wrappers.NewNonExistentErr(fmt.Errorf("email %s not found", email))
		}
		return
	}

	resp = models.UserResp(*(result[0].(*entities.User)))

	return
}

// GetByID user
func (s *userService) GetByID(ctx context.Context, ID string) (resp models.UserResp, err error) {
	user, err := s.repository.GetByID(ctx, ID)
	if err != nil {
		if errors.Is(err, wrappers.NonExistentErr) {
			err = wrappers.NewNonExistentErr(fmt.Errorf("ID %s not found", ID))
		}
		return
	}

	resp = models.UserResp(*user.(*entities.User))

	return
}

// Update user
func (s *userService) Update(ctx context.Context, ID string, u models.UpdateUserReq) (err error) {
	user, err := s.GetByID(ctx, ID)
	if err != nil {
		return
	}

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
		err = validatePassword(*u.OldPassword, user.PasswordHash)
		if err != nil {
			return
		}

		err = hashPassword(u.NewPassword)
		if err != nil {
			return
		}

		user.PasswordHash = *u.NewPassword
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

	err = s.repository.Update(ctx, ID, entities.User(user))
	return err
}

// Delete user
func (s *userService) Delete(ctx context.Context, ID string) (err error) {
	err = s.repository.Delete(ctx, ID)
	if errors.Is(err, wrappers.NonExistentErr) {
		err = wrappers.NewNonExistentErr(fmt.Errorf("ID %s not found", ID))
	}

	return
}

// GetClaims user
func (s *userService) GetClaims(ctx context.Context) (claims map[int]string) {
	claims = entities.GetClaims()
	return
}

// AtomicTransationProof creates two users atomically
func (s *userService) AtomicTransationProof(ctx context.Context) (err error) {
	user1Hash := "Entity1"
	err = hashPassword(&user1Hash)
	if err != nil {
		return
	}
	user2Hash := "Entity2"
	err = hashPassword(&user2Hash)
	if err != nil {
		return
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

	return
}
