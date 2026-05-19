package service

import (
	"context"
	"errors"
	"testing"

	"github.com/eventhub/user-service/internal/model"
	"github.com/eventhub/user-service/internal/repository"
	"github.com/google/uuid"
)

type mockUserRepo struct {
	users map[string]*model.User
}

func (m *mockUserRepo) Create(ctx context.Context, user *model.User) error {
	if _, ok := m.users[user.Email]; ok {
		return errors.New("duplicate")
	}
	m.users[user.Email] = user
	return nil
}

func (m *mockUserRepo) GetByID(ctx context.Context, id uuid.UUID) (*model.User, error) {
	for _, u := range m.users {
		if u.ID == id {
			return u, nil
		}
	}
	return nil, repository.ErrUserNotFound
}

func (m *mockUserRepo) GetByEmail(ctx context.Context, email string) (*model.User, error) {
	if u, ok := m.users[email]; ok {
		return u, nil
	}
	return nil, repository.ErrUserNotFound
}

func (m *mockUserRepo) List(ctx context.Context) ([]model.User, error) {
	out := make([]model.User, 0, len(m.users))
	for _, u := range m.users {
		out = append(out, *u)
	}
	return out, nil
}

func TestCreateUser_Validation(t *testing.T) {
	svc := NewUserService(&mockUserRepo{users: map[string]*model.User{}})

	_, err := svc.CreateUser(context.Background(), "", "Name", "password1", "user")
	if !errors.Is(err, ErrInvalidInput) {
		t.Fatalf("expected ErrInvalidInput, got %v", err)
	}

	_, err = svc.CreateUser(context.Background(), "a@b.com", "Name", "short", "user")
	if !errors.Is(err, ErrInvalidInput) {
		t.Fatalf("expected ErrInvalidInput for short password, got %v", err)
	}
}

func TestCreateUser_Success(t *testing.T) {
	svc := NewUserService(&mockUserRepo{users: map[string]*model.User{}})
	user, err := svc.CreateUser(context.Background(), "test@example.com", "Test", "password12", "user")
	if err != nil {
		t.Fatal(err)
	}
	if user.Email != "test@example.com" {
		t.Fatalf("unexpected email %s", user.Email)
	}
}

func TestValidateCredentials_Invalid(t *testing.T) {
	svc := NewUserService(&mockUserRepo{users: map[string]*model.User{}})
	_, err := svc.ValidateCredentials(context.Background(), "missing@example.com", "password12")
	if !errors.Is(err, ErrInvalidCredentials) {
		t.Fatalf("expected ErrInvalidCredentials, got %v", err)
	}
}
