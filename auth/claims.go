package auth

import (
	"time"

	"github.com/F-Dupraz/ecommerce-with-go/model"

	"github.com/golang-jwt/jwt/v5"
)

type Claims struct {
    UserID   string         `json:"user_id"`
    Email    string         `json:"email"`
    Role     model.UserRole `json:"role"`
	IsAdmin  bool           `json:"is_admin"`
    jwt.RegisteredClaims
}

func NewClaims(userID, email, role string) *Claims {
    now := time.Now()
    return &Claims{
        UserID:  userID,
        Email:   email,
        Role:    role,
        IsAdmin: role == "admin",
        RegisteredClaims: jwt.RegisteredClaims{
            ExpiresAt: jwt.NewNumericDate(now.Add(15 * time.Minute)),
			IssuedAt:  jwt.NewNumericDate(now),
            NotBefore: jwt.NewNumericDate(now),
            Issuer:    "your-ecommerce",
            Subject:   userID,
        },
    }
}
