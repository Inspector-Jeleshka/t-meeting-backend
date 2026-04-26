package route

import (
	"t-meeting-backend/internal/api/controller"
	"t-meeting-backend/internal/api/middleware"
	"t-meeting-backend/internal/jwt"

	"github.com/go-chi/chi/v5"
)

func SetupAuthRoutes(r chi.Router, authController *controller.AuthController, jwtManager *jwt.JWTManager) {
	r.Route("/auth", func(r chi.Router) {
		r.Post("/register", authController.Register)
		r.Post("/login", authController.Login)
	})

	r.Group(func(r chi.Router) {
		r.Use(middleware.JWTVerifier(jwtManager))
		r.Get("/user", authController.Me)
		r.Post("/auth/refresh", authController.Refresh)
	})
}
