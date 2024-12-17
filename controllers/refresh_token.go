package controllers

import (
	"database/sql"
	"net/http"
	"os"

	"github.com/jovi345/login-register/config"
	"github.com/jovi345/login-register/helper"
	"github.com/jovi345/login-register/utils"
)

func RefreshToken(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("refresh_token")
	if err != nil {
		helper.SendResponse(w, http.StatusUnauthorized, "Refresh token not found")
		return
	}

	secretKey := os.Getenv("JWT_SECRET_REFRESH")
	refreshToken := cookie.Value

	query := "SELECT id FROM users WHERE refresh_token = ?"
	row := config.DB.QueryRow(query, refreshToken)

	var userID string
	err = row.Scan(&userID)
	if err == sql.ErrNoRows {
		helper.SendResponse(w, http.StatusForbidden, "User not found")
		return
	}

	_, err = utils.ValidateToken(refreshToken, secretKey)
	if err != nil {
		helper.SendResponse(w, http.StatusForbidden, "Invalid or expired token")
		return
	}

	accessToken, err := utils.GenerateAccessToken(userID)
	if err != nil {
		helper.SendResponse(w, http.StatusInternalServerError, "Failed to generate access token")
		return
	}

	helper.SendResponse(w, http.StatusOK, map[string]string{
		"access_token": accessToken,
	})
}
