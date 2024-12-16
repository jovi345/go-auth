package main

import (
	"log"
	"net/http"

	"github.com/joho/godotenv"
	"github.com/jovi345/login-register/config"
	"github.com/jovi345/login-register/routes"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	config.Connect()

	r := routes.RegisterRoutes()
	log.Println("Server is running on port 8080...")
	http.ListenAndServe(":8080", r)
}
