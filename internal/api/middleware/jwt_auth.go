package middleware

import (
	"net/http"
	"t-meeting-backend/internal/domain"
	"t-meeting-backend/internal/jwt"
	"time"
)

func JWTVerifier(jwtManager *jwt.JWTManager) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		verifier := func(w http.ResponseWriter, r *http.Request) {
			cookie, err := r.Cookie("access_token")
			if err != nil {
				http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
				return
			}
			claims, err := jwtManager.ParseToken(cookie.Value)
			if err != nil {
				http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
				return
			}
			if claims.Role != domain.AdminRole || !claims.ExpiresAt.Time.After(time.Now()) {
				http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
				return
			}
			next.ServeHTTP(w, r)
		}
		return http.HandlerFunc(verifier)
	}
}
