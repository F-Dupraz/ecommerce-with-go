package handler

import (
    "encoding/json"
    "net/http"
    "errors"
    
    "github.com/go-chi/chi/v5"
    "github.com/go-playground/validator/v10"
    
    "github.com/F-Dupraz/ecommerce-with-go/dto"
    "github.com/F-Dupraz/ecommerce-with-go/service"
)

type AuthHandler struct {
    BaseHandler
    authService *service.AuthService
    userService *service.UserService
}

func NewAuthHandler(authService *service.AuthService, userService *service.UserService, validator *validator.Validate) *AuthHandler {
    return &AuthHandler{
        authService: authService,
        userService: userService,
        BaseHandler: BaseHandler{validator: validator},
    }
}

func (h *AuthHandler) RegisterRoutes(router chi.Router) {
    router.Route("/auth", func(r chi.Router) {
        r.Post("/login", h.Login)
        r.Post("/signup", h.Signup)
        r.Post("/refresh", h.Refresh)
        r.Post("/logout", h.Logout)
    })
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
    var req dto.LoginRequest
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        h.respondWithError(w, http.StatusBadRequest, "Invalid request format", nil)
        return
    }
    
    if err := h.validator.Struct(req); err != nil {
        validationErrors := dto.FormatValidationErrors(err)
        h.respondWithError(w, http.StatusUnprocessableEntity, "Validation failed", validationErrors)
        return
    }
    
    ctx := context.WithValue(r.Context(), "ip_address", r.RemoteAddr)
    ctx = context.WithValue(ctx, "user_agent", r.Header.Get("User-Agent"))
    
    response, err := h.authService.Login(ctx, req.Email, req.Password)
    if err != nil {
        if errors.Is(err, service.ErrInvalidCredentials) {
            h.respondWithError(w, http.StatusUnauthorized, "Invalid email or password", nil)
            return
        }
        h.respondWithError(w, http.StatusInternalServerError, "Login failed", nil)
        return
    }
    
    http.SetCookie(w, &http.Cookie{
        Name:     "refresh_token",
        Value:    response.RefreshToken,
        Path:     "/api/v1/auth/refresh",
        HttpOnly: true,
        Secure:   true, // false para desarrollo local
        SameSite: http.SameSiteStrictMode,
        MaxAge:   7 * 24 * 60 * 60, // 7 días
    })
    
    response.RefreshToken = ""
    h.respondWithSuccess(w, http.StatusOK, response)
}

func (h *AuthHandler) Signup(w http.ResponseWriter, r *http.Request) {
    var req dto.SignupRequest
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        h.respondWithError(w, http.StatusBadRequest, "Invalid request format", nil)
        return
    }
    
    if err := h.validator.Struct(req); err != nil {
        validationErrors := dto.FormatValidationErrors(err)
        h.respondWithError(w, http.StatusUnprocessableEntity, "Validation failed", validationErrors)
        return
    }
    
    createReq := dto.CreateUserRequest{
        Username: req.Username,
        Email:    req.Email,
        Password: req.Password,
        Address:  req.Address,
        City:     req.City,
        Country:  req.Country,
    }
    
    userResponse, err := h.userService.CreateUser(r.Context(), createReq)
    if err != nil {
        if errors.Is(err, service.ErrEmailAlreadyExists) {
            h.respondWithError(w, http.StatusConflict, "Email already registered", nil)
            return
        }
        h.respondWithError(w, http.StatusInternalServerError, "Could not create account", nil)
        return
    }
    
    ctx := context.WithValue(r.Context(), "ip_address", r.RemoteAddr)
    ctx = context.WithValue(ctx, "user_agent", r.Header.Get("User-Agent"))
    
    loginResponse, err := h.authService.Login(ctx, req.Email, req.Password)
    if err != nil {
        h.respondWithSuccess(w, http.StatusCreated, dto.SignupResponse{
            Message: "Account created successfully. Please login.",
            User:    userResponse.User,
        })
        return
    }
    
    http.SetCookie(w, &http.Cookie{
        Name:     "refresh_token",
        Value:    loginResponse.RefreshToken,
        Path:     "/api/v1/auth/refresh",
        HttpOnly: true,
        Secure:   true,
        SameSite: http.SameSiteStrictMode,
        MaxAge:   7 * 24 * 60 * 60,
    })
    
    h.respondWithSuccess(w, http.StatusCreated, dto.SignupResponse{
        Message:     "Account created successfully",
        AccessToken: loginResponse.AccessToken,
        TokenType:   "Bearer",
        ExpiresIn:   900,
        User:        loginResponse.User,
    })
}

func (h *AuthHandler) Refresh(w http.ResponseWriter, r *http.Request) {
    cookie, err := r.Cookie("refresh_token")
    if err != nil {
        h.respondWithError(w, http.StatusUnauthorized, "No refresh token provided", nil)
        return
    }
    
    response, err := h.authService.RefreshToken(r.Context(), cookie.Value)
    if err != nil {
        if errors.Is(err, service.ErrInvalidToken) || errors.Is(err, service.ErrSessionExpired) {
            http.SetCookie(w, &http.Cookie{
                Name:     "refresh_token",
                Value:    "",
                Path:     "/api/v1/auth/refresh",
                HttpOnly: true,
                MaxAge:   -1,
            })
            h.respondWithError(w, http.StatusUnauthorized, "Invalid or expired refresh token", nil)
            return
        }
        h.respondWithError(w, http.StatusInternalServerError, "Could not refresh token", nil)
        return
    }
    
    h.respondWithSuccess(w, http.StatusOK, response)
}

func (h *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
    userID, ok := r.Context().Value("user_id").(string)
    if ok && userID != "" {
        h.authService.Logout(r.Context(), userID)
    }
    
    http.SetCookie(w, &http.Cookie{
        Name:     "refresh_token",
        Value:    "",
        Path:     "/api/v1/auth/refresh",
        HttpOnly: true,
        MaxAge:   -1,
    })
    
    h.respondWithSuccess(w, http.StatusOK, map[string]string{
        "message": "Logged out successfully",
    })
}
