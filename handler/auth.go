package handler

import (
  "log"
  "net/http"
  "context"
  "encoding/json"

  "github.com/go-chi/chi/v5"
  "github.com/go-playground/validator/v10"

  "github.com/F-Dupraz/ecommerce-with-go/dto"
  "github.com/F-Dupraz/ecommerce-with-go/service"
)

type AuthConfig struct {
    CookieDomain        string
    IsProduction        bool
    AccessTokenTTL      time.Duration
    RefreshTokenTTL     time.Duration
    SecureCookies       bool
    AllowedOrigins      []string
}

func NewAuthHandler(authService *service.AuthService, userSvc *service.UserService, validator *validator.Validate, config AuthConfig) *AuthHandler {
    return &AuthHandler{
        authService: authSvc,
        userService: userSvc,
		config:      config,
		BaseHandler: BaseHandler{validator: validator}
	}
}

func (h *AuthHandler) RegisterRoutes(router chi.Router) {
  router.Route("/auth", func (r chi.Router) {
	router.Post("/login", u.Login)
	router.Post("/signin", u.Signup)
	router.Post("/refresh", u.Refresh)
	router.Post("/logout", u.Logout)
  })
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
    var req dto.LoginRequest
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        h.respondWithError(w, "Invalid request", http.StatusBadRequest)
        return
    }

	if err := validator.Struct(req); err != nil {
		validationErrors := dto.FormatValidationErrors(err)
		h.respondWithError(w, http.StatusUnprocessableEntity, "Validation failed!", validationErrors)
		return
	}
    
    response, err := h.authService.Login(r.Context(), req.Email, req.Password)
    if err != nil {
        switch {
        case errors.Is(err, service.ErrInvalidCredentials):
            h.respondWithError(w, "Invalid email or password", http.StatusUnauthorized)
        case errors.Is(err, service.ErrAccountLocked):
            h.respondWithError(w, "Account temporarily locked", http.StatusTooManyRequests)
        default:
            h.respondWithError(w, "Internal error", http.StatusInternalServerError)
        }
        return
    }
	
	h.respondWithSuccess(w, http.StatusCreated, response)
}

func (h *AuthHandler) Refresh(w http.ResponseWriter, r *http.Request) {
    cookie, err := r.Cookie("refresh_token")
    if err != nil {
        h.respondWithError(w, "No refresh token provided", http.StatusUnauthorized)
        return
    }
    
    if cookie.Value == "" {
        h.respondWithError(w, "Empty refresh token", http.StatusUnauthorized)
        return
    }
    
    response, err := h.authService.RefreshTokens(r.Context(), cookie.Value, clientInfo)
    if err != nil {
        switch {
        case errors.Is(err, service.ErrInvalidRefreshToken):
            h.clearRefreshTokenCookie(w)
            h.respondWithError(w, http.StatusUnauthorized, "Invalid refresh token")
            
        case errors.Is(err, service.ErrRefreshTokenExpired):
            h.clearRefreshTokenCookie(w)
            h.respondWithError(w, http.StatusUnauthorized, "Refresh token expired. Please login again")
            
        case errors.Is(err, service.ErrSessionRevoked):
            h.clearRefreshTokenCookie(w)
            h.respondWithError(w, http.StatusUnauthorized, "Session has been revoked")
            
        case errors.Is(err, service.ErrRefreshTokenReused):
            h.clearRefreshTokenCookie(w)
            h.respondWithError(w, http.StatusUnauthorized, "Security violation: token reuse detected")
            
        case errors.Is(err, service.ErrUserDeactivated):
            h.clearRefreshTokenCookie(w)
            h.respondWithError(w, http.StatusForbidden, "Account has been deactivated")
            
        default:
            log.Error("Refresh token error", "error", err)
            h.respondWithError(w, http.StatusInternalServerError, "Could not refresh token")
        }
        return
    }
    
    if response.NewRefreshToken != "" {
        h.setRefreshTokenCookie(w, response.NewRefreshToken)
    }
    
	h.resopndWithSuccess(w, http.StatusOK, dto.RefresTokenResponse)
}

func (h *AuthHandler) Signup(w http.ResponseWriter, r *http.Request) {
    var req dto.SignupRequest
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        h.respondWithError(w, http.StatusBadRequest, "Invalid request", nil)
        return
    }
    
	if err := h.validator.Struct(req); err != nil {
		validationErrors := dto.FormatValidationErrors(err)
		h.respondWithError(w, http.StatusUnprocessableEntity, "Validation failed!", validationErrors)
		return
	}

    createUserReq := dto.CreateUserRequest{
        Email:    req.Email,
        Password: req.Password,
        Name:     req.Name,
		City:     req.City,
		Country:  req.Country,
		Role:     "customer",
	}
    
    user, err := h.userService.CreateUser(r.Context(), createUserReq)
    if err != nil {
        switch {
        case errors.Is(err, service.ErrEmailAlreadyExists):
            h.respondWithError(w, http.StatusConflict, "Email already registered")
        case errors.Is(err, service.ErrWeakPassword):
            h.respondWithError(w, http.StatusBadRequest, "Password too weak")
        default:
            h.respondWithError(w, http.StatusInternalServerError, "Could not create account")
        }
        return
    }
    
    authResponse, err := h.authService.CreateSessionForUser(r.Context(), user)
    if err != nil {
        log.Error("Could not auto-login new user", "user_id", user.ID, "error", err)
        h.respondWithError(w, http.StatusInternalServerError, "Account created but could not login. Please try logging in.")
        return
    }
    
    h.setRefreshTokenCookie(w, authResponse.RefreshToken)
    
	h.respondWithSuccess(w, http.StatusOK, authResponse)
}

func (h *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
    loggedOut := false
    
    cookie, err := r.Cookie("refresh_token")
    if err == nil && cookie.Value != "" {
        if err := h.authService.RevokeRefreshToken(r.Context(), cookie.Value); err != nil {
            log.Warn("Could not revoke refresh token", "error", err)
        } else {
            loggedOut = true
        }
    }
    
    if !loggedOut {
        if userID, ok := middleware.GetUserID(r.Context()); ok {
            if err := h.authService.RevokeUserSessions(r.Context(), userID); err != nil {
                log.Warn("Could not revoke user sessions", "user_id", userID, "error", err)
            } else {
                loggedOut = true
            }
        }
    }
    
    h.clearRefreshTokenCookie(w)
    
    if loggedOut {
        h.respondWithSuccess(w, http.StatusOK, "Logged out successfully")
    } else {
		h.respondWithSuccess(w, http.StatusOK, "No active session")
	}
}

func (h *AuthHandler) setRefreshTokenCookie(w http.ResponseWriter, refreshToken string) {
    cookie := &http.Cookie{
        Name:     "refresh_token",
        Value:    refreshToken,
        Path:     "/",
        Domain:   h.config.CookieDomain,
        MaxAge:   7 * 24 * 60 * 60,
        HttpOnly: true,
        Secure:   h.config.IsProduction,
        SameSite: http.SameSiteLaxMode,
    }
    
    http.SetCookie(w, cookie)
}

func (h *AuthHandler) clearRefreshTokenCookie(w http.ResponseWriter) {
    http.SetCookie(w, &http.Cookie{
        Name:     "refresh_token",
        Value:    "",
        HttpOnly: true,
        Secure:   true,
        SameSite: http.SameSiteLax,
        MaxAge:   -1,
		Path:     "/",
    })
}

