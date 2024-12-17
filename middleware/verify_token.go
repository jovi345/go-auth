package middleware

import (
	"net/http"
	"os"
	"strings"

	"github.com/jovi345/login-register/helper"
	"github.com/jovi345/login-register/utils"
)

func VerifyToken(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			helper.SendResponse(w, http.StatusUnauthorized, "Authorization header missing")
			return
		}

		tokenParts := strings.Split(authHeader, " ")
		if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
			helper.SendResponse(w, http.StatusUnauthorized, "Invalid authorization header")
			return
		}

		accessToken := tokenParts[1]
		secretKey := os.Getenv("JWT_SECRET_ACCESS")
		claims, err := utils.ValidateToken(accessToken, secretKey)
		if err != nil {
			helper.SendResponse(w, http.StatusForbidden, "Invalid or expired token")
			return
		}

		r = r.WithContext(utils.ContextWithClaims(r.Context(), *claims))
		next(w, r)
	}
}
