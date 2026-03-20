package controller

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"

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
	var req dto.RefreshRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "decode body: "+err.Error(), http.StatusBadRequest)
		return
	}

	claims, err := ac.jwt.ParseToken(req.RefreshToken)
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
			http.Error(w, domain.ErrInvalidToken.Error(), http.StatusUnauthorized)
			return
		}
		http.Error(w, "get user by id: "+err.Error(), http.StatusInternalServerError)
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

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(dto.TokenPair{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	})
}

func (ac *AuthController) Me(w http.ResponseWriter, r *http.Request) {
	tokenString, err := exctractBearerToken(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	claims, err := ac.jwt.ParseToken(tokenString)
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

func (ac *AuthController) buildAuthResponse(user *domain.User) (*dto.AuthResponse, error) {
	accessToken, err := ac.jwt.GenerateAccessToken(user)
	if err != nil {
		return nil, err
	}

	refreshToken, err := ac.jwt.GenerateRefreshToken(user)
	if err != nil {
		return nil, err
	}

	return &dto.AuthResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		User: dto.User{
			Email: user.Email,
			Role:  user.Role,
		},
	}, nil
}

func exctractBearerToken(r *http.Request) (string, error) {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		return "", domain.ErrInvalidToken
	}

	const prefix = "Bearer "
	if !strings.HasPrefix(authHeader, prefix) {
		return "", domain.ErrInvalidToken
	}

	token := strings.TrimPrefix(authHeader, prefix)
	if token == "" {
		return "", domain.ErrInvalidToken
	}

	return token, nil
}
