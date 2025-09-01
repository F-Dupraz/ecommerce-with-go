package model

import (
    "time"
    "database/sql"

	"github.com/google/uuid"
)

type Session struct {
    ID               uuid.UUID
    UserID           uuid.UUID
    RefreshTokenHash string
    ParentSessionID  uuid.NullUUID
    FamilyID         uuid.UUID
    IPAddress        string
    UserAgent        string
    DeviceInfo       sql.NullString
    ExpiresAt        time.Time
    RevokedAt        sql.NullTime
    RotatedAt        sql.NullTime
    WasRotated       bool
    RefreshCount     int
    CreatedAt        time.Time
    LastUsedAt       sql.NullTime
    LastIP           sql.NullString
    LastUserAgent    sql.NullString
    User             *User
}

type SessionActivity struct {
    ID        uuid.UUID
    SessionID uuid.UUID
    Action    string
    IPAddress string
    UserAgent string
    Metadata  map[string]interface{}
    CreatedAt time.Time
}
