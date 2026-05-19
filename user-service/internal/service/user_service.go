package service

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/eventhub/user-service/internal/model"
	"github.com/eventhub/user-service/internal/repository"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrInvalidInput      = errors.New("invalid input")
	ErrEmailTaken        = errors.New("email already registered")
	ErrInvalidCredentials = errors.New("invalid credentials")
)

type UserService interface {
	CreateUser(ctx context.Context, email, name, password, role string) (*model.User, error)
	GetUser(ctx context.Context, id uuid.UUID) (*model.User, error)
	ListUsers(ctx context.Context) ([]model.User, error)
	ValidateCredentials(ctx context.Context, email, password string) (*model.User, error)
}

type userService struct {
	repo repository.UserRepository
}

func NewUserService(repo repository.UserRepository) UserService {
	return &userService{repo: repo}
}

func (s *userService) CreateUser(ctx context.Context, email, name, password, role string) (*model.User, error) {
	email = strings.TrimSpace(strings.ToLower(email))
	name = strings.TrimSpace(name)
	if email == "" || name == "" || len(password) < 8 {
		return nil, ErrInvalidInput
	}
	if role == "" {
		role = model.RoleUser
	}
	if role != model.RoleUser && role != model.RoleAdmin {
		return nil, ErrInvalidInput
	}

	if _, err := s.repo.GetByEmail(ctx, email); err == nil {
		return nil, ErrEmailTaken
	} else if !errors.Is(err, repository.ErrUserNotFound) {
		return nil, err
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("hash password: %w", err)
	}

	user := &model.User{
		Email:        email,
		Name:         name,
		PasswordHash: string(hash),
		Role:         role,
	}
	if err := s.repo.Create(ctx, user); err != nil {
		return nil, err
	}
	return user, nil
}

func (s *userService) GetUser(ctx context.Context, id uuid.UUID) (*model.User, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *userService) ListUsers(ctx context.Context) ([]model.User, error) {
	return s.repo.List(ctx)
}

func (s *userService) ValidateCredentials(ctx context.Context, email, password string) (*model.User, error) {
	user, err := s.repo.GetByEmail(ctx, strings.TrimSpace(strings.ToLower(email)))
	if err != nil {
		if errors.Is(err, repository.ErrUserNotFound) {
			return nil, ErrInvalidCredentials
		}
		return nil, err
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		return nil, ErrInvalidCredentials
	}
	return user, nil
}
