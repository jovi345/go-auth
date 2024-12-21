package main

import (
	"log"
	"net/http"

	"github.com/gorilla/handlers"
	"github.com/joho/godotenv"
	"github.com/jovi345/login-register/config"
	"github.com/jovi345/login-register/route"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	config.Connect()

	corsOptions := handlers.CORS(
		handlers.AllowedOrigins([]string{"http://localhost:5173"}),
		handlers.AllowedHeaders([]string{"Content-Type", "Authorization"}),
		handlers.AllowedMethods([]string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}),
		handlers.AllowCredentials(),
	)

	r := route.RegisterRoutes()
	log.Println("Server is running on port 8080...")
	http.ListenAndServe(":8080", corsOptions(r))
}
