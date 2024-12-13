package main

import (
	"log"
	"net/http"

	"github.com/jovi345/login-register/config"
	"github.com/jovi345/login-register/routes"
)

func main() {
	config.Connect()

	r := routes.RegisterRoutes()
	log.Println("Server is running on port 8080...")
	http.ListenAndServe(":8080", r)
}
