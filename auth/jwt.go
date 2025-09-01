package auth

import (
    "crypto/rand"
    "encoding/base64"
    "errors"
    "fmt"
    "time"

    "github.com/golang-jwt/jwt/v5"
)

type JWTManager struct {
    secretKey     []byte
    refreshSecret []byte
    accessTTL     time.Duration
    refreshTTL    time.Duration
}

func NewJWTManager(secretKey string) *JWTManager {
    return &JWTManager{
		secretKey:     []byte(secretKey),
		refreshSecret: []byte(secretKey + "_refresh"),
		accessTTL:     15 * time.Minute,
		refreshTTL:    7 * 24 * time.Hour,
    }
}

func (j *JWTManager) GenerateTokenPair(userID, email, role string) (access, refresh string, err error) {
    claims := NewClaims(userID, email, role)

    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    access, err = token.SignedString(j.secretKey)
    if err != nil {
        return "", "", fmt.Errorf("failed to sign access token: %w", err)
    }

    refreshClaims := jwt.RegisteredClaims{
        Subject:   userID,
        ExpiresAt: jwt.NewNumericDate(time.Now().Add(j.refreshTTL)),
        IssuedAt:  jwt.NewNumericDate(time.Now()),
    }

    refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims)
    refresh, err = refreshToken.SignedString(j.refreshSecret)
    if err != nil {
        return "", "", fmt.Errorf("failed to sign refresh token: %w", err)
    }

    return access, refresh, nil
}

func (j *JWTManager) ValidateAccessToken(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
	  if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
            return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
        }
        return j.secretKey, nil
    })
    
    if err != nil {
        return nil, fmt.Errorf("invalid token: %w", err)
    }
    
    claims, ok := token.Claims.(*Claims)
    if !ok || !token.Valid {
        return nil, errors.New("invalid token claims")
    }
    
    return claims, nil
}

func (j *JWTManager) RefreshAccessToken(refreshToken string) (string, error) {
    token, err := jwt.Parse(refreshToken, func(token *jwt.Token) (interface{}, error) {
        if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
            return nil, fmt.Errorf("unexpected signing method")
        }
        return j.refreshSecret, nil
    })

    if err != nil || !token.Valid {
        return "", errors.New("invalid refresh token")
    }

    subject, err := token.Claims.GetSubject()
    if err != nil {
        return "", err
    }

    // TODO: Acá deberías ir a la DB a buscar el user actual
    // user := getUserFromDB(subject)
    // Por ahora hardcodeamos
	fmt.Println(subject)


    claims := NewClaims(subject, "user@email.com", "customer")
    newToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

    return newToken.SignedString(j.secretKey)
}

