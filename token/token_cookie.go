package token

import (
	"net/http"
	"time"
)

func SetRefreshTokenCookie(w http.ResponseWriter, refreshToken string) {
	cookie := &http.Cookie{
		Name:     "refresh_token",
		Value:    refreshToken,
		Expires:  time.Now().Add(3 * 24 * time.Hour),
		HttpOnly: true,
		// Secure:   true,
		SameSite: http.SameSiteDefaultMode,
		Path:     "/",
	}
	http.SetCookie(w, cookie)
}

func ClearCookie(w http.ResponseWriter) {
	cookie := &http.Cookie{
		Name:    "refresh_token",
		Value:   "",
		Expires: time.Unix(0, 0),
		Path:    "/",
	}

	http.SetCookie(w, cookie)
}
