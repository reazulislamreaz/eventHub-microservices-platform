package grpc

import (
	"context"
	"errors"

	userv1 "github.com/eventhub/proto/gen/user/v1"
	"github.com/eventhub/user-service/internal/model"
	"github.com/eventhub/user-service/internal/repository"
	"github.com/eventhub/user-service/internal/service"
	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type UserHandler struct {
	userv1.UnimplementedUserServiceServer
	svc service.UserService
}

func NewUserHandler(svc service.UserService) *UserHandler {
	return &UserHandler{svc: svc}
}

func (h *UserHandler) CreateUser(ctx context.Context, req *userv1.CreateUserRequest) (*userv1.CreateUserResponse, error) {
	user, err := h.svc.CreateUser(ctx, req.GetEmail(), req.GetName(), req.GetPassword(), req.GetRole())
	if err != nil {
		return nil, mapError(err)
	}
	return &userv1.CreateUserResponse{User: toProtoUser(user)}, nil
}

func (h *UserHandler) GetUser(ctx context.Context, req *userv1.GetUserRequest) (*userv1.GetUserResponse, error) {
	id, err := uuid.Parse(req.GetId())
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid user id")
	}
	user, err := h.svc.GetUser(ctx, id)
	if err != nil {
		return nil, mapError(err)
	}
	return &userv1.GetUserResponse{User: toProtoUser(user)}, nil
}

func (h *UserHandler) ListUsers(ctx context.Context, _ *userv1.ListUsersRequest) (*userv1.ListUsersResponse, error) {
	users, err := h.svc.ListUsers(ctx)
	if err != nil {
		return nil, mapError(err)
	}
	resp := &userv1.ListUsersResponse{Users: make([]*userv1.User, 0, len(users))}
	for i := range users {
		u := users[i]
		resp.Users = append(resp.Users, toProtoUser(&u))
	}
	return resp, nil
}

func (h *UserHandler) ValidateCredentials(ctx context.Context, req *userv1.ValidateCredentialsRequest) (*userv1.ValidateCredentialsResponse, error) {
	user, err := h.svc.ValidateCredentials(ctx, req.GetEmail(), req.GetPassword())
	if err != nil {
		if errors.Is(err, service.ErrInvalidCredentials) {
			return &userv1.ValidateCredentialsResponse{Valid: false}, nil
		}
		return nil, mapError(err)
	}
	return &userv1.ValidateCredentialsResponse{User: toProtoUser(user), Valid: true}, nil
}

func toProtoUser(u *model.User) *userv1.User {
	return &userv1.User{
		Id:        u.ID.String(),
		Email:     u.Email,
		Name:      u.Name,
		Role:      u.Role,
		CreatedAt: u.CreatedAt.UTC().Format("2006-01-02T15:04:05Z07:00"),
	}
}

func mapError(err error) error {
	switch {
	case errors.Is(err, service.ErrInvalidInput):
		return status.Error(codes.InvalidArgument, err.Error())
	case errors.Is(err, service.ErrEmailTaken):
		return status.Error(codes.AlreadyExists, err.Error())
	case errors.Is(err, repository.ErrUserNotFound):
		return status.Error(codes.NotFound, err.Error())
	default:
		return status.Error(codes.Internal, "internal error")
	}
}
