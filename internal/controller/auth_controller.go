package controller

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/google/uuid"

	"t-meeting-backend/internal/domain"
	"t-meeting-backend/internal/dto"
	"t-meeting-backend/internal/service"
)

type AuthController struct {
	us  *service.UserService
	jwt *service.JWTService
}

func NewAuthController(us *service.UserService, jwt *service.JWTService) *AuthController {
	return &AuthController{
		us:  us,
		jwt: jwt,
	}
}

func (ac *AuthController) Register(w http.ResponseWriter, r *http.Request) {
	var credentials dto.AuthCredentials
	if err := json.NewDecoder(r.Body).Decode(&credentials); err != nil {
		http.Error(w, "decode body: "+err.Error(), http.StatusBadRequest)
		return
	}

	user, err := ac.us.Register(r.Context(), &credentials)
	if err != nil {
		if errors.Is(err, service.ErrUserAlreadyExists) {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		http.Error(w, "register user: "+err.Error(), http.StatusInternalServerError)
		return
	}

	accessToken, err := ac.jwt.GenerateAccessToken(user)
	if err != nil {
		http.Error(w, "generate access token: "+err.Error(), http.StatusInternalServerError)
		return
	}
	refreshToken, err := ac.jwt.GenerateRefreshToken(user)
	if err != nil {
		http.Error(w, "generate refresh token: "+err.Error(), http.StatusInternalServerError)
		return
	}
	resp := dto.User{
		Email: user.Email,
		Role:  user.Role,
	}

	accessTokenCookie := &http.Cookie{
		Name:     "access_token",
		Value:    accessToken,
		Path:     "/",
		HttpOnly: true,
		SameSite: http.SameSiteStrictMode,
	}
	refreshTokeCookie := &http.Cookie{
		Name:     "refresh_token",
		Value:    refreshToken,
		Path:     "/",
		HttpOnly: true,
		SameSite: http.SameSiteStrictMode,
	}

	w.Header().Set("Content-Type", "application/json")
	http.SetCookie(w, accessTokenCookie)
	http.SetCookie(w, refreshTokeCookie)
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(resp)
}

func (ac *AuthController) Login(w http.ResponseWriter, r *http.Request) {
	var credentials dto.AuthCredentials
	if err := json.NewDecoder(r.Body).Decode(&credentials); err != nil {
		http.Error(w, "decode body: "+err.Error(), http.StatusBadRequest)
		return
	}

	user, err := ac.us.Login(r.Context(), &credentials)
	if err != nil {
		if errors.Is(err, domain.ErrInvalidCredentials) {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}
		http.Error(w, "login user: "+err.Error(), http.StatusInternalServerError)
		return
	}

	accessToken, err := ac.jwt.GenerateAccessToken(user)
	if err != nil {
		http.Error(w, "generate access token: "+err.Error(), http.StatusInternalServerError)
		return
	}
	refreshToken, err := ac.jwt.GenerateRefreshToken(user)
	if err != nil {
		http.Error(w, "generate refresh token: "+err.Error(), http.StatusInternalServerError)
		return
	}
	resp := dto.User{
		Email: user.Email,
		Role:  user.Role,
	}

	accessTokenCookie := &http.Cookie{
		Name:     "access_token",
		Value:    accessToken,
		Path:     "/",
		HttpOnly: true,
		SameSite: http.SameSiteStrictMode,
	}
	refreshTokeCookie := &http.Cookie{
		Name:     "refresh_token",
		Value:    refreshToken,
		Path:     "/",
		HttpOnly: true,
		SameSite: http.SameSiteStrictMode,
	}

	w.Header().Set("Content-Type", "application/json")
	http.SetCookie(w, accessTokenCookie)
	http.SetCookie(w, refreshTokeCookie)
	_ = json.NewEncoder(w).Encode(resp)
}

func (ac *AuthController) Refresh(w http.ResponseWriter, r *http.Request) {
	refreshToken, err := extractCookieToken(r, "refresh_token")
	if err != nil {
		http.Error(w, domain.ErrInvalidToken.Error(), http.StatusUnauthorized)
		return
	}

	claims, err := ac.jwt.ParseToken(refreshToken)
	if err != nil {
		http.Error(w, domain.ErrInvalidToken.Error(), http.StatusUnauthorized)
		return
	}

	if claims.Type != service.TokenTypeRefresh {
		http.Error(w, domain.ErrInvalidToken.Error(), http.StatusUnauthorized)
		return
	}

	userID, err := uuid.Parse(claims.Subject)
	if err != nil {
		http.Error(w, "invalid token subject: "+err.Error(), http.StatusUnauthorized)
		return
	}

	user, err := ac.us.GetByID(r.Context(), userID)
	if err != nil {
		if errors.Is(err, domain.ErrUserNotFound) {
			http.Error(w, domain.ErrUserNotFound.Error(), http.StatusNotFound)
			return
		}
		http.Error(w, "get user by id: "+err.Error(), http.StatusInternalServerError)
		return
	}

	newAccessToken, err := ac.jwt.GenerateAccessToken(user)
	if err != nil {
		http.Error(w, "generate access token: "+err.Error(), http.StatusInternalServerError)
		return
	}

	newRefreshToken, err := ac.jwt.GenerateRefreshToken(user)
	if err != nil {
		http.Error(w, "generate refresh token: "+err.Error(), http.StatusInternalServerError)
		return
	}

	accessTokenCookie := &http.Cookie{
		Name:     "access_token",
		Value:    newAccessToken,
		Path:     "/",
		HttpOnly: true,
		SameSite: http.SameSiteStrictMode,
	}

	refreshTokeCookie := &http.Cookie{
		Name:     "refresh_token",
		Value:    newRefreshToken,
		Path:     "/",
		HttpOnly: true,
		SameSite: http.SameSiteStrictMode,
	}

	http.SetCookie(w, accessTokenCookie)
	http.SetCookie(w, refreshTokeCookie)

	w.WriteHeader(http.StatusOK)
}

func (ac *AuthController) Me(w http.ResponseWriter, r *http.Request) {
	accessToken, err := extractCookieToken(r, "access_token")
	if err != nil {
		http.Error(w, domain.ErrInvalidToken.Error(), http.StatusUnauthorized)
		return
	}

	claims, err := ac.jwt.ParseToken(accessToken)
	if err != nil {
		http.Error(w, domain.ErrInvalidToken.Error(), http.StatusUnauthorized)
		return
	}

	if claims.Type != service.TokenTypeAccess {
		http.Error(w, domain.ErrInvalidToken.Error(), http.StatusUnauthorized)
		return
	}

	userID, err := uuid.Parse(claims.Subject)
	if err != nil {
		http.Error(w, "invalid token subject: "+err.Error(), http.StatusUnauthorized)
		return
	}

	user, err := ac.us.GetByID(r.Context(), userID)
	if err != nil {
		if errors.Is(err, domain.ErrUserNotFound) {
			http.Error(w, domain.ErrInvalidToken.Error(), http.StatusUnauthorized)
			return
		}
		http.Error(w, "get user by id: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(dto.User{
		Email: user.Email,
		Role:  user.Role,
	})
}


func extractCookieToken(r *http.Request, cookieName string) (string, error) {
	cookie, err := r.Cookie(cookieName)
	if err != nil {
		return "", domain.ErrInvalidToken
	}

	if cookie.Value == "" {
		return "", domain.ErrInvalidToken
	}

	return cookie.Value, nil
}
