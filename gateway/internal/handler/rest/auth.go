package rest

import (
	"encoding/json"
	"net/http"

	userv1 "github.com/eventhub/proto/gen/user/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// Register godoc
// @Summary      Register a new user
// @Description  Creates a user account and returns a JWT token
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        body  body      RegisterRequest  true  "Registration details"
// @Success      201   {object}  AuthResponse
// @Failure      400   {object}  ErrorResponse
// @Failure      409   {object}  ErrorResponse
// @Failure      500   {object}  ErrorResponse
// @Router       /api/v1/auth/register [post]
func (h *Handler) Register(w http.ResponseWriter, r *http.Request) {
	var req RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid JSON body")
		return
	}
	if req.Email == "" || req.Name == "" || len(req.Password) < 8 {
		writeError(w, http.StatusBadRequest, "email, name required; password min 8 characters")
		return
	}

	resp, err := h.Clients.User.CreateUser(r.Context(), &userv1.CreateUserRequest{
		Email:    req.Email,
		Name:     req.Name,
		Password: req.Password,
		Role:     "user",
	})
	if err != nil {
		writeGRPCError(w, err)
		return
	}

	user := mapProtoUser(resp.GetUser())
	token, err := h.JWT.Generate(user.ID, user.Email, user.Role)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to generate token")
		return
	}
	writeJSON(w, http.StatusCreated, AuthResponse{Token: token, User: user})
}

// Login godoc
// @Summary      Login
// @Description  Authenticates with email/password and returns a JWT token
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        body  body      LoginRequest  true  "Login credentials"
// @Success      200   {object}  AuthResponse
// @Failure      401   {object}  ErrorResponse
// @Failure      400   {object}  ErrorResponse
// @Failure      500   {object}  ErrorResponse
// @Router       /api/v1/auth/login [post]
func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid JSON body")
		return
	}

	resp, err := h.Clients.User.ValidateCredentials(r.Context(), &userv1.ValidateCredentialsRequest{
		Email:    req.Email,
		Password: req.Password,
	})
	if err != nil {
		writeGRPCError(w, err)
		return
	}
	if !resp.GetValid() {
		writeError(w, http.StatusUnauthorized, "invalid email or password")
		return
	}

	user := mapProtoUser(resp.GetUser())
	token, err := h.JWT.Generate(user.ID, user.Email, user.Role)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to generate token")
		return
	}
	writeJSON(w, http.StatusOK, AuthResponse{Token: token, User: user})
}

func writeGRPCError(w http.ResponseWriter, err error) {
	if st, ok := status.FromError(err); ok {
		switch st.Code() {
		case codes.InvalidArgument:
			writeError(w, http.StatusBadRequest, st.Message())
		case codes.AlreadyExists:
			writeError(w, http.StatusConflict, st.Message())
		case codes.NotFound:
			writeError(w, http.StatusNotFound, st.Message())
		case codes.FailedPrecondition:
			writeError(w, http.StatusConflict, st.Message())
		case codes.PermissionDenied:
			writeError(w, http.StatusForbidden, st.Message())
		default:
			writeError(w, http.StatusInternalServerError, st.Message())
		}
		return
	}
	writeError(w, http.StatusInternalServerError, "internal server error")
}
