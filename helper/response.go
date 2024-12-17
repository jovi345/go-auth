package helper

import (
	"encoding/json"
	"net/http"

	"github.com/jovi345/login-register/models"
)

func SendResponse(w http.ResponseWriter, statusCode int, result interface{}) {
	response := models.JSONResponse{
		Status: http.StatusText(statusCode),
		Result: result,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(response)
}
