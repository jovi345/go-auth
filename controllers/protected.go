package controllers

import (
	"net/http"

	"github.com/jovi345/login-register/helper"
	"github.com/jovi345/login-register/utils"
)

func ProtectedEndpoint(w http.ResponseWriter, r *http.Request) {
	claims, ok := utils.GetClaimsFromContext(r.Context())
	if !ok {
		helper.SendResponse(w, http.StatusUnauthorized, "Unable to extract user information")
		return
	}

	userID := claims["user_id"].(string)
	email := claims["email"].(string)
	firstName := claims["first_name"].(string)
	lastName := claims["last_name"].(string)
	role := claims["role"].(string)

	helper.SendResponse(w, http.StatusOK, map[string]string{
		"message":    "You have accessed a protected route!",
		"user_id":    userID,
		"email":      email,
		"first_name": firstName,
		"last_name":  lastName,
		"role":       role,
	})
}
