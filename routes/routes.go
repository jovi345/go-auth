package routes

import (
	"github.com/gorilla/mux"
	"github.com/jovi345/login-register/controllers"
)

func RegisterRoutes() *mux.Router {
	r := mux.NewRouter()

	r.HandleFunc("/api/v1/user/register", controllers.Register).Methods("POST")

	return r
}
