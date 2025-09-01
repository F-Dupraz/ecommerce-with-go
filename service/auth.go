package service

import (
  "context"
  "errors"
  "time"
  "fmt"
  "log"

  "github.com/F-Dupraz/ecommerce-with-go/dto"
  "github.com/F-Dupraz/ecommerce-with-go/model"
  "github.com/F-Dupraz/ecommerce-with-go/auth"
  "github.com/F-Dupraz/ecommerce-with-go/security"
  "github.com/F-Dupraz/ecommerce-with-go/repository"
)

type AuthService struct {
    userRepo      repository.UserRepository
    sessionRepo   repository.SessionRepository
    jwtManager    *auth.JWTManager
    hasher        *security.PasswordHasher
    rateLimiter   *security.RateLimiter
}

type RefreshTokenResponse struct {
    AccessToken     string
    NewRefreshToken string
	ExpiresIn       int
    User           *dto.UserInfo
}

func (s *AuthService) Login(ctx context.Context, email, password string) (*AuthResponse, error) {
    if err := s.rateLimiter.CheckLoginAttempt(ctx, email); err != nil {
        return nil, ErrTooManyAttempts
    }
    
    user, err := s.userRepo.FindByEmail(ctx, email)
    if err != nil {
        s.rateLimiter.RecordFailedAttempt(ctx, email)
        return nil, ErrInvalidCredentials
    }
    
    if !s.hasher.Verify(password, user.PasswordHash) {
        s.rateLimiter.RecordFailedAttempt(ctx, email)
        return nil, ErrInvalidCredentials
    }
    
    if !user.IsActive {
        return nil, ErrAccountInactive
    }
    if user.RequiresPasswordChange {
        return nil, ErrPasswordChangeRequired
    }
    
    s.rateLimiter.ResetAttempts(ctx, email)
    
    accessToken, refreshToken, err := s.jwtManager.GenerateTokenPair(
        user.ID,
        user.Email,
        user.Role,
    )
    
    session := &model.Session{
        UserID:           user.ID,
        RefreshTokenHash: hashToken(refreshToken),
        UserAgent:        getUserAgent(ctx),
        IPAddress:        getClientIP(ctx),
        ExpiresAt:        time.Now().Add(7 * 24 * time.Hour),
    }
    
    if err := s.sessionRepo.Create(ctx, session); err != nil {
        return nil, err
    }
    
    return &AuthResponse{
        AccessToken:  accessToken,
        RefreshToken: refreshToken,
        ExpiresIn:    900,
        User: UserInfo{
            ID:    user.ID,
            Email: user.Email,
            Name:  user.Name,
            Role:  user.Role,
        },
    }, nil
}

func (s *AuthService) RevokeRefreshToken(ctx context.Context, refreshToken string) error {
    tokenHash := hashToken(refreshToken)
    
    session, err := s.sessionRepo.FindByTokenHash(ctx, tokenHash)
    if err != nil {
        return nil    }
    
    return s.sessionRepo.RevokeSession(ctx, session.ID)
}

func (s *AuthService) RefreshTokens(ctx context.Context, refreshToken string, clientInfo dto.ClientInfo) (*RefreshTokenResponse, error) {
    claims, err := s.jwtManager.ParseRefreshToken(refreshToken)
    if err != nil {
        log.Warn("Invalid refresh token structure", 
            "error", err,
            "ip", clientInfo.IPAddress)
        return nil, ErrInvalidRefreshToken
    }
    
    if time.Now().After(claims.ExpiresAt.Time) {
        return nil, ErrRefreshTokenExpired
    }
    
    tokenHash := s.hashToken(refreshToken)
    session, err := s.sessionRepo.FindByTokenHash(ctx, tokenHash)
    if err != nil {
        if errors.Is(err, pgx.ErrNoRows) {
            log.Warn("Refresh token not found in DB",
                "user_id", claims.Subject,
                "ip", clientInfo.IPAddress)
            return nil, ErrInvalidRefreshToken
        }
        return nil, fmt.Errorf("database error: %w", err)
    }
    
    if err := s.validateSession(session, clientInfo); err != nil {
        return nil, err
    }
    
    user, err := s.userRepo.GetByID(ctx, session.UserID)
    if err != nil {
        return nil, fmt.Errorf("user not found: %w", err)
    }
    
    if !user.IsActive {
        s.sessionRepo.RevokeAllUserSessions(ctx, user.ID)
        return nil, ErrUserDeactivated
    }
    
    accessToken, err := s.jwtManager.GenerateAccessToken(
        user.ID,
        user.Email,
        user.Role,    )
    if err != nil {
        return nil, fmt.Errorf("failed to generate access token: %w", err)
    }
    
    updateData := map[string]interface{}{
        "last_used_at": time.Now(),
        "last_ip":      clientInfo.IPAddress,
        "last_user_agent": clientInfo.UserAgent,
        "refresh_count": session.RefreshCount + 1,
    }
    
    if err := s.sessionRepo.UpdateSession(ctx, session.ID, updateData); err != nil {
        log.Warn("Could not update session metadata", "session_id", session.ID)
    }
    
    var newRefreshToken string
    if s.shouldRotateRefreshToken(session) {
        newRefreshToken, err = s.rotateRefreshToken(ctx, session, user.ID)
        if err != nil {
            log.Error("Could not rotate refresh token", "error", err)
        }
    }
    
    return &RefreshTokenResponse{
        AccessToken:     accessToken,
        NewRefreshToken: newRefreshToken,
        ExpiresIn:       900, // 15 minutos
        User: &dto.UserInfo{
            ID:    user.ID,
            Email: user.Email,
            Name:  user.Name,
            Role:  user.Role,
        },
    }, nil
}

func (s *AuthService) validateSession(session *model.Session, clientInfo dto.ClientInfo) error {
    if session.RevokedAt != nil {
        log.Warn("Attempt to use revoked session",
            "session_id", session.ID,
            "revoked_at", session.RevokedAt)
        return ErrSessionRevoked
    }
    
    if time.Now().After(session.ExpiresAt) {
        return ErrRefreshTokenExpired
    }
    
    if s.isSessionAnomaly(session, clientInfo) {
        log.Warn("Session anomaly detected",
            "session_id", session.ID,
            "stored_ip", session.IPAddress,
            "current_ip", clientInfo.IPAddress)
    }
    
    if session.RefreshCount > 100 {
        log.Warn("Excessive refresh count",
            "session_id", session.ID,
            "count", session.RefreshCount)
        return ErrTooManyRefreshes
    }
    
    if session.LastUsedAt != nil {
        timeSinceLastRefresh := time.Since(*session.LastUsedAt)
        if timeSinceLastRefresh < 10*time.Second {
            return ErrRefreshTooSoon
        }
    }
    
    return nil
}

func (s *AuthService) shouldRotateRefreshToken(session *model.Session) bool {
    if time.Since(session.CreatedAt) > 24*time.Hour {
        return true
    }
    
    if session.RefreshCount > 10 {
        return true
    }
    
    if s.config.AlwaysRotateRefreshTokens {
        return true
    }
    
    return false
}


func (s *AuthService) rotateRefreshToken(ctx context.Context, oldSession *model.Session, userID string) (string, error) {
    newRefreshToken, err := s.jwtManager.GenerateRefreshToken(userID)
    if err != nil {
        return "", fmt.Errorf("failed to generate new refresh token: %w", err)
    }
    
    rotationData := &repository.TokenRotation{
        OldSessionID:     oldSession.ID,
        NewTokenHash:     s.hashToken(newRefreshToken),
        UserID:           userID,
        FamilyID:         oldSession.FamilyID,
        IPAddress:        oldSession.IPAddress,
        UserAgent:        oldSession.UserAgent,
        ExpiresAt:        time.Now().Add(s.config.RefreshTokenTTL),
    }
    
    newSessionID, err := s.sessionRepo.RotateToken(ctx, rotationData)
    if err != nil {
        return "", fmt.Errorf("failed to rotate token: %w", err)
    }
    
    log.Info("Token rotated successfully",
        "old_session", oldSession.ID,
        "new_session", newSessionID,
        "user_id", userID)
    
    return newRefreshToken, nil
}


func (s *AuthService) isSessionAnomaly(session *model.Session, clientInfo dto.ClientInfo) bool {
    if session.IPAddress != clientInfo.IPAddress {
        if !s.isSameNetwork(session.IPAddress, clientInfo.IPAddress) {
            return true
        }
    }
    
    if session.UserAgent != clientInfo.UserAgent {
        if !s.isSimilarUserAgent(session.UserAgent, clientInfo.UserAgent) {
            return true
        }
    }
    
    return false
}
