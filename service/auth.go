// service/auth.go - Versión MVP
package service

import (
    "context"
    "crypto/sha256"
    "encoding/hex"
    "errors"
    "time"
    
    "github.com/google/uuid"
    "golang.org/x/crypto/bcrypt"
)

var (
    ErrInvalidCredentials = errors.New("invalid credentials")
    ErrInvalidToken = errors.New("invalid token")
    ErrSessionExpired = errors.New("session expired")
)

type AuthService struct {
    userRepo    *repository.UserRepository
    sessionRepo *repository.SessionRepository
    jwtManager  *auth.JWTManager
}

func NewAuthService(userRepo *repository.UserRepository, sessionRepo *repository.SessionRepository, jwtManager *auth.JWTManager) *AuthService {
    return &AuthService{
        userRepo:    userRepo,
        sessionRepo: sessionRepo,
        jwtManager:  jwtManager,
    }
}

func (s *AuthService) Login(ctx context.Context, email, password string) (*dto.LoginResponse, error) {
    user, err := s.userRepo.GetByEmail(ctx, email)
    if err != nil {
        return nil, ErrInvalidCredentials
    }
    
    if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
        return nil, ErrInvalidCredentials
    }
    
    s.sessionRepo.InvalidateByUserID(ctx, user.ID)
    
    accessToken, refreshToken, err := s.jwtManager.GenerateTokenPair(
        user.ID, 
        user.Email, 
        string(user.Role),
    )
    if err != nil {
        return nil, err
    }
    
    session := &model.Session{
        ID:           uuid.New().String(),
        UserID:       user.ID,
        RefreshToken: s.hashToken(refreshToken),
        IPAddress:    getIPFromContext(ctx),
        UserAgent:    getUserAgentFromContext(ctx),
        ExpiresAt:    time.Now().Add(7 * 24 * time.Hour),
    }
    
    if err := s.sessionRepo.Create(ctx, session); err != nil {
        return nil, err
    }
    
    return &dto.LoginResponse{
        AccessToken: accessToken,
        RefreshToken: refreshToken,
        TokenType:   "Bearer",
        ExpiresIn:   900, // 15 minutos
        User: dto.UserInfo{
            ID:    user.ID,
            Email: user.Email,
            Name:  user.Username,
            Role:  string(user.Role),
        },
    }, nil
}

func (s *AuthService) RefreshToken(ctx context.Context, refreshToken string) (*dto.RefreshResponse, error) {
    tokenHash := s.hashToken(refreshToken)
    session, err := s.sessionRepo.GetByRefreshToken(ctx, tokenHash)
    if err != nil {
        return nil, ErrInvalidToken
    }
    
    if time.Now().After(session.ExpiresAt) {
        s.sessionRepo.InvalidateByID(ctx, session.ID)
        return nil, ErrSessionExpired
    }
    
    user, err := s.userRepo.GetByID(ctx, session.UserID)
    if err != nil {
        return nil, err
    }
    
    accessToken, _, err := s.jwtManager.GenerateTokenPair(
        user.ID,
        user.Email,
        string(user.Role),
    )
    if err != nil {
        return nil, err
    }
    
    s.sessionRepo.UpdateLastActivity(ctx, session.ID)
    
    return &dto.RefreshResponse{
        AccessToken: accessToken,
        TokenType:   "Bearer",
        ExpiresIn:   900,
    }, nil
}

func (s *AuthService) Logout(ctx context.Context, userID string) error {
    return s.sessionRepo.InvalidateByUserID(ctx, userID)
}

func (s *AuthService) hashToken(token string) string {
    hash := sha256.Sum256([]byte(token))
    return hex.EncodeToString(hash[:])
}
