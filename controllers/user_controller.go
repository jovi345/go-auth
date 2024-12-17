package controllers

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/jovi345/login-register/config"
	"github.com/jovi345/login-register/helper"
	"github.com/jovi345/login-register/models"
	"github.com/jovi345/login-register/utils"
	"golang.org/x/crypto/bcrypt"
)

func CheckEmailAvailability(userInput models.UserRegisterInput) bool {
	query := "SELECT email FROM users WHERE email = ?"
	row := config.DB.QueryRow(query, userInput.Email)

	var userID string
	err := row.Scan(&userID)
	if err != nil {
		if err.Error() == "sql: no rows in result set" {
			return false
		}

		log.Println(err.Error())
		return true
	}

	return true
}

func Register(w http.ResponseWriter, r *http.Request) {
	var userInput models.UserRegisterInput
	err := json.NewDecoder(r.Body).Decode(&userInput)
	if err != nil {
		log.Printf("Failed to parse input: %v", err)
		helper.SendResponse(w, http.StatusBadRequest, "Invalid input format")
		return
	}

	validate := validator.New(validator.WithRequiredStructEnabled())
	err = validate.Struct(userInput)
	if err != nil {
		log.Printf("Error: %v", err.Error())
		if validationErrors, ok := err.(validator.ValidationErrors); ok {
			errorMessages := make(map[string]string)
			for _, fieldErr := range validationErrors {
				fieldName := helper.GetJSONFieldName(fieldErr.Field(), models.UserRegisterInput{})
				errorMessages[fieldName] = "Validation failed on tag: " + fieldName
			}
			helper.SendResponse(w, http.StatusBadRequest, errorMessages)
			return
		}
		helper.SendResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	status := CheckEmailAvailability(userInput)
	if status {
		helper.SendResponse(w, http.StatusBadRequest, "Email is not available")
		return
	}

	if userInput.Password != userInput.ConfirmPassword {
		helper.SendResponse(w, http.StatusBadRequest, "Password do not match")
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(userInput.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Printf("Error: %v", err)
		helper.SendResponse(w, http.StatusInternalServerError, "Internal server error")
		return
	}

	user := models.User{
		ID:         "user-" + uuid.NewString(),
		FirstName:  userInput.FirstName,
		MiddleName: userInput.MiddleName,
		LastName:   userInput.LastName,
		Email:      userInput.Email,
		Password:   string(hashedPassword),
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
		Role:       "user",
	}

	query := "INSERT INTO users (id, first_name, middle_name, last_name, email, password, created_at, updated_at, role) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)"

	_, err = config.DB.Exec(query, user.ID, user.FirstName, user.MiddleName, user.LastName, user.Email, user.Password, user.CreatedAt, user.UpdatedAt, user.Role)

	if err != nil {
		log.Printf("Error: %v", err)
		helper.SendResponse(w, http.StatusInternalServerError, "Internal server error")
		return
	}

	helper.SendResponse(w, http.StatusCreated, "User registered successfully")
}

func Login(w http.ResponseWriter, r *http.Request) {
	var userInput models.UserLoginInput
	err := json.NewDecoder(r.Body).Decode(&userInput)
	if err != nil {
		log.Printf("Error: %v", err.Error())
		helper.SendResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	validate := validator.New(validator.WithRequiredStructEnabled())
	err = validate.Struct(userInput)
	if err != nil {
		log.Printf("Error: %v", err.Error())
		if validationErrors, ok := err.(validator.ValidationErrors); ok {
			errorMessages := make(map[string]string)
			for _, fieldErr := range validationErrors {
				fieldName := helper.GetJSONFieldName(fieldErr.Field(), models.UserLoginInput{})
				errorMessages[fieldName] = "Validation failed on tag: " + fieldName
			}
			helper.SendResponse(w, http.StatusBadRequest, errorMessages)
			return
		}
		helper.SendResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	query := "SELECT id, email, password FROM users WHERE email = ?"
	row := config.DB.QueryRow(query, userInput.Email)

	var userID, email, hashedPassword string
	err = row.Scan(&userID, &email, &hashedPassword)
	if err == sql.ErrNoRows {
		helper.SendResponse(w, http.StatusUnauthorized, "Email not found")
		return
	}
	if err != nil {
		log.Printf("Error: %v", err.Error())
		helper.SendResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(userInput.Password))
	if err != nil {
		helper.SendResponse(w, http.StatusUnauthorized, "Wrong password")
		return
	}

	accessToken, err := utils.GenerateAccessToken(userID)
	if err != nil {
		helper.SendResponse(w, http.StatusInternalServerError, "Failed to generate access token")
		return
	}

	refreshToken, err := utils.GenerateRefreshToken(userID)
	if err != nil {
		helper.SendResponse(w, http.StatusInternalServerError, "Failed to generate refresh token")
		return
	}

	lastLogin := time.Now()

	query = "UPDATE users SET refresh_token = ?, last_login = ? WHERE id = ?"
	_, err = config.DB.Exec(query, refreshToken, lastLogin, userID)
	if err != nil {
		log.Printf("Error: %v", err.Error())
		helper.SendResponse(w, http.StatusInternalServerError, "Failed to save refresh token")
		return
	}

	utils.SetRefreshTokenCookie(w, refreshToken)

	helper.SendResponse(w, http.StatusOK, map[string]string{
		"access_token": accessToken,
	})
}

func Logout(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("refresh_token")
	if err != nil {
		helper.SendResponse(w, http.StatusNoContent, "No content")
		return
	}

	refreshToken := cookie.Value
	query := "SELECT id FROM users WHERE refresh_token = ?"
	row := config.DB.QueryRow(query, refreshToken)

	var userID string
	err = row.Scan(&userID)
	if err == sql.ErrNoRows {
		helper.SendResponse(w, http.StatusNoContent, "No content")
		return
	}

	query = "UPDATE users SET refresh_token = ? WHERE id = ?"
	_, err = config.DB.Exec(query, sql.NullString{}, userID)
	if err != nil {
		helper.SendResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	utils.ClearCookie(w)

	helper.SendResponse(w, http.StatusOK, "Successfully logged out")
}
