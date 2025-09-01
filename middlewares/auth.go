package middleware

import (
    "context"
    "net/http"
    "strings"

	"github.com/F-Fupraz/ecommerce-with-go/auth"
)

type AuthMiddleware struct {
    jwtManager *auth.JWTManager
}

func NewAuthMiddleware(jwtManager *auth.JWTManager) *AuthMiddleware {
    return &AuthMiddleware{jwtManager: jwtManager}
}

func (am *AuthMiddleware) Authenticate(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        authHeader := r.Header.Get("Authorization")
        if authHeader == "" {
            next.ServeHTTP(w, r)
            return
        }
        
        const bearerPrefix = "Bearer "
        if !strings.HasPrefix(authHeader, bearerPrefix) {
            http.Error(w, "Invalid authorization header format", http.StatusUnauthorized)
            return
        }
        
        tokenString := authHeader[len(bearerPrefix):]
        
        claims, err := am.jwtManager.ValidateAccessToken(tokenString)
        if err != nil {
            http.Error(w, "Invalid or expired token", http.StatusUnauthorized)
            return
        }
        
        ctx := context.WithValue(r.Context(), UserIDKey, claims.UserID)
        ctx = context.WithValue(ctx, UserRoleKey, claims.Role)
        ctx = context.WithValue(ctx, IsAdminKey, claims.IsAdmin)
        ctx = context.WithValue(ctx, EmailKey, claims.Email)
        
        next.ServeHTTP(w, r.WithContext(ctx))
    })
}

func RequireAuth(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        userID, ok := GetUserID(r.Context())
        if !ok || userID == "" {
            http.Error(w, "Unauthorized", http.StatusUnauthorized)
            return
        }
        next.ServeHTTP(w, r)
    })
}

func RequireAdmin(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        if !IsAdmin(r.Context()) {
            http.Error(w, "Forbidden", http.StatusForbidden)
            return
        }
        next.ServeHTTP(w, r)
    })
}

