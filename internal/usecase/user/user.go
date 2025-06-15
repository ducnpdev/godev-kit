package user

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/ducnpdev/godev-kit/internal/entity"
	"github.com/ducnpdev/godev-kit/internal/repo"
	"github.com/ducnpdev/godev-kit/internal/repo/persistent/models"
	"github.com/golang-jwt/jwt/v4"
	"github.com/jackc/pgx/v5"
	"golang.org/x/crypto/bcrypt"
)

// UseCase -.
type UseCase struct {
	repo      repo.UserRepo
	jwtSecret []byte
}

// New -.
func New(r repo.UserRepo, jwtSecret string) *UseCase {
	return &UseCase{
		repo:      r,
		jwtSecret: []byte(jwtSecret),
	}
}

// hashPassword hashes a password using bcrypt
func (uc *UseCase) hashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("UserUseCase - hashPassword - bcrypt.GenerateFromPassword: %w", err)
	}
	return string(hashedPassword), nil
}

// Create -.
func (uc *UseCase) Create(ctx context.Context, user entity.User) (entity.User, error) {
	// Hash password before creating user
	hashedPassword, err := uc.hashPassword(user.Password)
	if err != nil {
		return entity.User{}, fmt.Errorf("UserUseCase - Create - uc.hashPassword: %w", err)
	}

	userModel := models.UserModel{
		Email:    user.Email,
		Username: user.Username,
		Password: hashedPassword,
	}

	createdUser, err := uc.repo.Create(ctx, userModel)
	if err != nil {
		return entity.User{}, fmt.Errorf("UserUseCase - Create - uc.repo.Create: %w", err)
	}

	return createdUser, nil
}

// GetByID -.
func (uc *UseCase) GetByID(ctx context.Context, id int64) (entity.User, error) {
	user, err := uc.repo.GetByID(ctx, id)
	if err != nil {
		return entity.User{}, fmt.Errorf("UserUseCase - GetByID - uc.repo.GetByID: %w", err)
	}

	return user, nil
}

// Update -.
func (uc *UseCase) Update(ctx context.Context, user entity.User) error {
	userModel := models.UserModel{
		ID:       user.ID,
		Email:    user.Email,
		Username: user.Username,
	}

	// Only hash password if it's being updated
	if user.Password != "" {
		hashedPassword, err := uc.hashPassword(user.Password)
		if err != nil {
			return fmt.Errorf("UserUseCase - Update - uc.hashPassword: %w", err)
		}
		userModel.Password = hashedPassword
	}

	err := uc.repo.Update(ctx, userModel)
	if err != nil {
		return fmt.Errorf("UserUseCase - Update - uc.repo.Update: %w", err)
	}

	return nil
}

// Delete -.
func (uc *UseCase) Delete(ctx context.Context, id int64) error {
	err := uc.repo.Delete(ctx, id)
	if err != nil {
		return fmt.Errorf("UserUseCase - Delete - uc.repo.Delete: %w", err)
	}

	return nil
}

// List -.
func (uc *UseCase) List(ctx context.Context) (entity.UserHistory, error) {
	users, err := uc.repo.List(ctx)
	if err != nil {
		return entity.UserHistory{}, fmt.Errorf("UserUseCase - List - uc.repo.List: %w", err)
	}

	return entity.UserHistory{Users: users}, nil
}

// JWTClaims represents the claims in a JWT token
type JWTClaims struct {
	UserID int64  `json:"user_id"`
	Email  string `json:"email"`
	jwt.RegisteredClaims
}

// Login authenticates a user and returns a JWT token
func (uc *UseCase) Login(ctx context.Context, email, password string) (string, entity.User, error) {
	// Validate input
	if email == "" {
		return "", entity.User{}, fmt.Errorf("UserUseCase - Login - email is required")
	}
	if password == "" {
		return "", entity.User{}, fmt.Errorf("UserUseCase - Login - password is required")
	}

	// Get user by email
	user, err := uc.repo.GetByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return "", entity.User{}, fmt.Errorf("UserUseCase - Login - invalid email or password")
		}
		return "", entity.User{}, fmt.Errorf("UserUseCase - Login - failed to get user: %w", err)
	}

	// Compare passwords
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return "", entity.User{}, fmt.Errorf("UserUseCase - Login - invalid email or password")
		}
		return "", entity.User{}, fmt.Errorf("UserUseCase - Login - failed to compare passwords: %w", err)
	}

	// Set token expiration (24 hours from now)
	expirationTime := time.Now().Add(24 * time.Hour)

	// Create JWT claims
	claims := &JWTClaims{
		UserID: user.ID,
		Email:  user.Email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    "godev-kit",
			Subject:   fmt.Sprintf("%d", user.ID),
		},
	}

	// Create token with claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Generate encoded token
	tokenString, err := token.SignedString(uc.jwtSecret)
	if err != nil {
		return "", entity.User{}, fmt.Errorf("UserUseCase - Login - failed to sign token: %w", err)
	}

	// Clear password from user entity before returning
	user.Password = ""

	return tokenString, user, nil
}
