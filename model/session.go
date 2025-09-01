package model

import (
    "time"
)

type Session struct {
    ID               string    `db:"id"`
    UserID           string    `db:"user_id"`
    RefreshToken     string    `db:"refresh_token"`
    AccessTokenJTI   string    `db:"access_token_jti"`
    IPAddress        string    `db:"ip_address"`
    UserAgent        string    `db:"user_agent"`
    ExpiresAt        time.Time `db:"expires_at"`
    LastActivityAt   time.Time `db:"last_activity_at"`
    IsValid          bool      `db:"is_valid"`
	CreatedAt        time.Time `db:"created_at"`
}
