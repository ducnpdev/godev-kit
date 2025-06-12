package user

import (
	"context"
	"fmt"

	"github.com/ducnpdev/godev-kit/internal/entity"
	"github.com/ducnpdev/godev-kit/internal/repo"
	"github.com/ducnpdev/godev-kit/internal/repo/persistent/models"
	"golang.org/x/crypto/bcrypt"
)

// UseCase -.
type UseCase struct {
	repo repo.UserRepo
}

// New -.
func New(r repo.UserRepo) *UseCase {
	return &UseCase{
		repo: r,
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
