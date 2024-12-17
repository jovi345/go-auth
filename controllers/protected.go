package controllers

import (
	"net/http"

	"github.com/jovi345/login-register/helper"
)

func ProtectedEndpoint(w http.ResponseWriter, r *http.Request) {
	helper.SendResponse(w, http.StatusOK, map[string]string{
		"message": "You have accessed a protected route!",
	})
}
