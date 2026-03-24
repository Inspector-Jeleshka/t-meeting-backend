package route

import (
	"github.com/go-chi/chi/v5"

	"t-meeting-backend/internal/controller"
)

func SetupAuthRoutes(r chi.Router, AuthController *controller.AuthController) {
	r.Route("/auth", func(r chi.Router) {
		r.Post("/register", AuthController.Register)
		r.Post("/login", AuthController.Login)
		r.Post("/refresh", AuthController.Refresh)
		r.Post("/logout", AuthController.Logout)
		r.Get("/me", AuthController.Me)
	})
}
