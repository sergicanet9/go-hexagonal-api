package services

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v4"
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

	token, err := createToken(user.ID, s.config.JWTSecret, user.ClaimIDs)
	if err != nil {
		return
	}

	resp = models.LoginUserResp{
		User:  user,
		Token: token,
	}

	return
}

func (s *userService) validateLogin(ctx context.Context, credentials models.LoginUserReq) (models.GetUserResp, error) {
	if err := credentials.Validate(); err != nil {
		return models.GetUserResp{}, err
	}

	user, err := s.GetByEmail(ctx, credentials.Email)
	if err != nil {
		return models.GetUserResp{}, err
	}

	err = validatePassword(credentials.Password, user.PasswordHash)
	if err != nil {
		return models.GetUserResp{}, err
	}

	return user, nil
}

func validatePassword(password, hash string) error {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
		err = fmt.Errorf("password incorrect")
	}
	return wrappers.NewValidationErr(err)
}

func createToken(userid string, jwtSecret string, claimsIDs []int64) (string, error) {
	var err error
	addClaims := jwt.MapClaims{}
	addClaims["authorized"] = true
	addClaims["user_id"] = userid
	addClaims["exp"] = time.Now().UTC().Add(time.Hour * 168).Unix()

	err = validateClaims(claimsIDs)
	if err != nil {
		return "", err
	}
	for _, claimID := range claimsIDs {
		addClaims[entities.UserClaim(claimID).String()] = true
	}

	add := jwt.NewWithClaims(jwt.SigningMethodHS256, addClaims)
	token, err := add.SignedString([]byte(jwtSecret))
	return token, err
}

func validateClaims(claimsIDs []int64) error {
	for _, claimID := range claimsIDs {
		if ok := entities.UserClaim(claimID).IsValid(); !ok {
			return wrappers.NewValidationErr(fmt.Errorf("claim %d is not valid", claimID))
		}
	}
	return nil
}

// Create user
func (s *userService) Create(ctx context.Context, user models.CreateUserReq) (resp models.CreateUserResp, err error) {
	entity, err := s.createUserEntity(user, time.Now().UTC())
	if err != nil {
		return
	}

	insertedID, err := s.repository.Create(ctx, entity)
	if err != nil {
		return
	}

	resp = models.CreateUserResp{
		InsertedID: insertedID,
	}

	return
}

func (s *userService) createUserEntity(user models.CreateUserReq, creationTime time.Time) (entity entities.User, err error) {
	if err = user.Validate(); err != nil {
		return
	}

	err = validateClaims(user.ClaimIDs)
	if err != nil {
		return
	}

	hash, err := hashPassword(user.Password)
	if err != nil {
		return
	}

	entity = entities.User{
		Name:         user.Name,
		Surnames:     user.Surnames,
		Email:        user.Email,
		PasswordHash: hash,
		ClaimIDs:     user.ClaimIDs,
		CreatedAt:    creationTime,
		UpdatedAt:    creationTime,
	}
	return
}

func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	hash := string(bytes)
	return hash, nil
}

// CreateMany users
func (s *userService) CreateMany(ctx context.Context, users []models.CreateUserReq) (resp models.CreateManyUserResp, err error) {
	var create []interface{}
	var entity entities.User
	creationTime := time.Now().UTC()

	for _, user := range users {
		entity, err = s.createUserEntity(user, creationTime)
		if err != nil {
			return
		}
		create = append(create, entity)
	}

	insertedIDs, err := s.repository.CreateMany(ctx, create)
	if err != nil {
		return
	}

	resp = models.CreateManyUserResp{
		InsertedIDs: insertedIDs,
	}
	return
}

// GetAll users
func (s *userService) GetAll(ctx context.Context) (resp []models.GetUserResp, err error) {
	result, err := s.repository.Get(ctx, map[string]interface{}{}, nil, nil)
	if err != nil {
		if errors.Is(err, wrappers.NonExistentErr) {
			err = nil
		}
		return
	}

	resp = make([]models.GetUserResp, len(result))
	for i, v := range result {
		resp[i] = models.GetUserResp(*(v.(*entities.User)))
	}

	return
}

// GetByEmail user
func (s *userService) GetByEmail(ctx context.Context, email string) (resp models.GetUserResp, err error) {
	filter := map[string]interface{}{"email": email}
	result, err := s.repository.Get(ctx, filter, nil, nil)
	if err != nil {
		if errors.Is(err, wrappers.NonExistentErr) {
			err = wrappers.NewNonExistentErr(fmt.Errorf("email %s not found", email))
		}
		return
	}

	resp = models.GetUserResp(*(result[0].(*entities.User)))

	return
}

// GetByID user
func (s *userService) GetByID(ctx context.Context, ID string) (resp models.GetUserResp, err error) {
	user, err := s.repository.GetByID(ctx, ID)
	if err != nil {
		if errors.Is(err, wrappers.NonExistentErr) {
			err = wrappers.NewNonExistentErr(fmt.Errorf("ID %s not found", ID))
		}
		return
	}

	resp = models.GetUserResp(*user.(*entities.User))

	return
}

// Update user
func (s *userService) Update(ctx context.Context, ID string, user models.UpdateUserReq) (err error) {
	dbUser, err := s.GetByID(ctx, ID)
	if err != nil {
		return
	}

	if user.Name != nil {
		dbUser.Name = *user.Name
	}
	if user.Surnames != nil {
		dbUser.Surnames = *user.Surnames
	}
	if user.Email != nil {
		dbUser.Email = *user.Email
	}
	if user.NewPassword != nil {
		err = validatePassword(*user.OldPassword, dbUser.PasswordHash)
		if err != nil {
			return
		}

		var hash string
		hash, err = hashPassword(*user.NewPassword)
		if err != nil {
			return
		}

		dbUser.PasswordHash = hash
	}
	if user.ClaimIDs != nil {
		err = validateClaims(*user.ClaimIDs)
		if err != nil {
			return err
		}
		dbUser.ClaimIDs = *user.ClaimIDs
	}
	dbUser.ID = ""
	dbUser.UpdatedAt = time.Now().UTC()

	err = s.repository.Update(ctx, ID, entities.User(dbUser))
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
func (s *userService) GetUserClaims(ctx context.Context) (claims map[int]string) {
	claims = entities.GetUserClaims()
	return
}
