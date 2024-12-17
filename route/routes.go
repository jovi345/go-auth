package route

import (
	"github.com/gorilla/mux"
	"github.com/jovi345/login-register/controllers"
	"github.com/jovi345/login-register/middleware"
)

func RegisterRoutes() *mux.Router {
	r := mux.NewRouter()

	r.HandleFunc("/api/v1/user/register", controllers.Register).Methods("POST")
	r.HandleFunc("/api/v1/user/login", controllers.Login).Methods("POST")
	r.HandleFunc("/api/v1/user/logout", controllers.Logout).Methods("DELETE")
	r.HandleFunc("/api/v1/protected", middleware.VerifyToken(controllers.ProtectedEndpoint)).Methods("GET")
	r.HandleFunc("/api/v1/token/refresh", controllers.RefreshToken)

	return r
}
