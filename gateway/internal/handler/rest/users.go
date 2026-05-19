package rest

import (
	"net/http"

	userv1 "github.com/eventhub/proto/gen/user/v1"
	"github.com/gorilla/mux"
)

// ListUsers godoc
// @Summary      List users
// @Description  Returns all registered users
// @Tags         users
// @Produce      json
// @Success      200  {array}   User
// @Failure      500  {object}  ErrorResponse
// @Router       /api/v1/users [get]
func (h *Handler) ListUsers(w http.ResponseWriter, r *http.Request) {
	resp, err := h.Clients.User.ListUsers(r.Context(), &userv1.ListUsersRequest{})
	if err != nil {
		writeGRPCError(w, err)
		return
	}
	users := make([]User, 0, len(resp.GetUsers()))
	for _, u := range resp.GetUsers() {
		users = append(users, mapProtoUser(u))
	}
	writeJSON(w, http.StatusOK, users)
}

// GetUser godoc
// @Summary      Get user by ID
// @Description  Returns a single user profile
// @Tags         users
// @Produce      json
// @Param        id   path      string  true  "User ID (UUID)"
// @Success      200  {object}  User
// @Failure      404  {object}  ErrorResponse
// @Failure      500  {object}  ErrorResponse
// @Router       /api/v1/users/{id} [get]
func (h *Handler) GetUser(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	resp, err := h.Clients.User.GetUser(r.Context(), &userv1.GetUserRequest{Id: id})
	if err != nil {
		writeGRPCError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, mapProtoUser(resp.GetUser()))
}

func mapProtoUser(u *userv1.User) User {
	return User{
		ID:        u.GetId(),
		Email:     u.GetEmail(),
		Name:      u.GetName(),
		Role:      u.GetRole(),
		CreatedAt: u.GetCreatedAt(),
	}
}
