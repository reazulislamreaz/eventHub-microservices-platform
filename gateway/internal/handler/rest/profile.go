package rest

import (
	"encoding/json"
	"net/http"

	userv1 "github.com/eventhub/proto/gen/user/v1"
	"github.com/eventhub/gateway/pkg/auth"
)

// UpdateProfile godoc
// @Summary      Update profile
// @Description  Updates the authenticated user's display name
// @Tags         users
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        body  body      UpdateProfileRequest  true  "Profile"
// @Success      200   {object}  User
// @Router       /api/v1/users/me [patch]
func (h *Handler) UpdateProfile(w http.ResponseWriter, r *http.Request) {
	claims, ok := auth.ClaimsFromContext(r.Context())
	if !ok {
		writeError(w, http.StatusUnauthorized, "authentication required")
		return
	}
	var req UpdateProfileRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid JSON body")
		return
	}
	resp, err := h.Clients.User.UpdateProfile(r.Context(), &userv1.UpdateProfileRequest{
		Id: claims.UserID, Name: req.Name,
	})
	if err != nil {
		writeGRPCError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, mapProtoUser(resp.GetUser()))
}
