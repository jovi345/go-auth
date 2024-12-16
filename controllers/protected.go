package controllers

import (
	"net/http"

	"github.com/jovi345/login-register/response"
	"github.com/jovi345/login-register/token"
)

func ProtectedEndpoint(w http.ResponseWriter, r *http.Request) {
	claims, ok := token.GetClaimsFromContext(r.Context())
	if !ok {
		response.SendResponse(w, http.StatusUnauthorized, "Unable to extract user information")
		return
	}

	userID := claims["user_id"].(string)
	response.SendResponse(w, http.StatusOK, map[string]string{
		"message": "You have accessed a protected route!",
		"user_id": userID,
	})
}
