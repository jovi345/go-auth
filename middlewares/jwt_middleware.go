package middlewares

import (
	"net/http"
	"os"
	"strings"

	"github.com/jovi345/login-register/response"
	"github.com/jovi345/login-register/token"
)

func JWTAuthMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			response.SendResponse(w, http.StatusUnauthorized, "Authorization header missing")
			return
		}

		tokenParts := strings.Split(authHeader, " ")
		if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
			response.SendResponse(w, http.StatusUnauthorized, "Invalid authorization header")
			return
		}

		accessToken := tokenParts[1]
		secretKey := os.Getenv("JWT_SECRET_ACCESS")
		claims, err := token.ParseToken(accessToken, secretKey)
		if err != nil {
			response.SendResponse(w, http.StatusForbidden, "Invalid or expired token")
			return
		}

		r = r.WithContext(token.ContextWithClaims(r.Context(), *claims))
		next(w, r)
	}
}
