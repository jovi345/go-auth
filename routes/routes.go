package routes

import (
	"github.com/gorilla/mux"
	"github.com/jovi345/login-register/controllers"
)

func RegisterRoutes() *mux.Router {
	r := mux.NewRouter()

	r.HandleFunc("/api/v1/user/register", controllers.Register).Methods("POST")
	r.HandleFunc("/api/v1/user/login", controllers.Login).Methods("POST")

	return r
}
