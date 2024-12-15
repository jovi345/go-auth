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
	"github.com/jovi345/login-register/input"
	"github.com/jovi345/login-register/models"
	"github.com/jovi345/login-register/response"
	"golang.org/x/crypto/bcrypt"
)

func CheckEmailAvailability(userInput input.UserRegisterInput) bool {
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
	var userInput input.UserRegisterInput
	err := json.NewDecoder(r.Body).Decode(&userInput)
	if err != nil {
		log.Printf("Failed to parse input: %v", err)
		response.SendResponse(w, http.StatusBadRequest, "Invalid input format")
		return
	}

	validate := validator.New(validator.WithRequiredStructEnabled())
	err = validate.Struct(userInput)
	if err != nil {
		log.Printf("Error: %v", err.Error())
		if validationErrors, ok := err.(validator.ValidationErrors); ok {
			errorMessages := make(map[string]string)
			for _, fieldErr := range validationErrors {
				fieldName := response.GetJSONFieldName(fieldErr.Field(), input.UserRegisterInput{})
				errorMessages[fieldName] = "Validation failed on tag: " + fieldName
			}
			response.SendResponse(w, http.StatusBadRequest, errorMessages)
			return
		}
		response.SendResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	status := CheckEmailAvailability(userInput)
	if status {
		response.SendResponse(w, http.StatusBadRequest, "Email is not available")
		return
	}

	if userInput.Password != userInput.ConfirmPassword {
		response.SendResponse(w, http.StatusBadRequest, "Password do not match")
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(userInput.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Printf("Error: %v", err)
		response.SendResponse(w, http.StatusInternalServerError, "Internal server error")
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
		response.SendResponse(w, http.StatusInternalServerError, "Internal server error")
		return
	}

	response.SendResponse(w, http.StatusCreated, "User registered successfully")
}

func Login(w http.ResponseWriter, r *http.Request) {
	var userInput input.UserLoginInput
	err := json.NewDecoder(r.Body).Decode(&userInput)
	if err != nil {
		log.Printf("Error: %v", err.Error())
		response.SendResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	validate := validator.New(validator.WithRequiredStructEnabled())
	err = validate.Struct(userInput)
	if err != nil {
		log.Printf("Error: %v", err.Error())
		if validationErrors, ok := err.(validator.ValidationErrors); ok {
			errorMessages := make(map[string]string)
			for _, fieldErr := range validationErrors {
				fieldName := response.GetJSONFieldName(fieldErr.Field(), input.UserLoginInput{})
				errorMessages[fieldName] = "Validation failed on tag: " + fieldName
			}
			response.SendResponse(w, http.StatusBadRequest, errorMessages)
			return
		}
		response.SendResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	query := "SELECT id, email, password FROM users WHERE email = ?"
	row := config.DB.QueryRow(query, userInput.Email)

	var userID, email, hashedPassword string
	err = row.Scan(&userID, &email, &hashedPassword)
	if err == sql.ErrNoRows {
		response.SendResponse(w, http.StatusUnauthorized, "Email not found")
		return
	}
	if err != nil {
		log.Printf("Error: %v", err.Error())
		response.SendResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(userInput.Password))
	if err != nil {
		response.SendResponse(w, http.StatusUnauthorized, "Wrong password")
		return
	}

	response.SendResponse(w, http.StatusBadRequest, []string{userID, email, hashedPassword})
}
