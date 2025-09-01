package dto

import "time"

type UserInfo struct {
    ID    string `json:"id"`
    Email string `json:"email"`
    Name  string `json:"name,omitempty"`
    Role  string `json:"role"`
}

type LoginRequest struct {
    Email    string `json:"email" validate:"required,email"`
    Password string `json:"password" validate:"required,min=8"`
}

type LoginResponse struct {
    AccessToken string   `json:"access_token"`
    TokenType   string   `json:"token_type"`
    ExpiresIn   int      `json:"expires_in"`
    User        UserInfo `json:"user"`
}

type SignupRequest struct {
	Username string `json:"username" validate:"required,min=3,max=50,alphanum"`
	Email    string `json:"email" validate:"required,email"`
	Address  string `json:"address,omitempty" validate:"max=200"`
	City     string `json:"city,omitempty" validate:"max=100"`
	Country  string `json:"country,omitempty" validate:"omitempty,iso3166_1_alpha2"`
    Password        string `json:"password" validate:"required,min=8"`
    PasswordConfirm string `json:"password_confirm" validate:"required,eqfield=Password"`
    AcceptTerms     bool   `json:"accept_terms" validate:"required"`
}

type SignupResponse struct {
    Message     string   `json:"message"`
    AccessToken string   `json:"access_token"`
    TokenType   string   `json:"token_type"`
    ExpiresIn   int      `json:"expires_in"`
    User        UserInfo `json:"user"`
}

type RefreshRequest struct {
}

type RefreshResponse struct {
    AccessToken string    `json:"access_token"`
    TokenType   string    `json:"token_type"`
    ExpiresIn   int       `json:"expires_in"`
    User        *UserInfo `json:"user,omitempty"`
}

type ChangePasswordRequest struct {
    CurrentPassword string `json:"current_password" validate:"required"`
    NewPassword     string `json:"new_password" validate:"required,min=8"`
    ConfirmPassword string `json:"confirm_password" validate:"required,eqfield=NewPassword"`
}

type ForgotPasswordRequest struct {
    Email string `json:"email" validate:"required,email"`
}

type ResetPasswordRequest struct {
    Token           string `json:"token" validate:"required"`
    NewPassword     string `json:"new_password" validate:"required,min=8"`
    ConfirmPassword string `json:"confirm_password" validate:"required,eqfield=NewPassword"`
}

type VerifyEmailRequest struct {
    Token string `json:"token" validate:"required"`
}

type SessionInfo struct {
    ID           string    `json:"id"`
    IPAddress    string    `json:"ip_address"`
    UserAgent    string    `json:"user_agent"`
    LastUsedAt   time.Time `json:"last_used_at"`
    CreatedAt    time.Time `json:"created_at"`
    IsCurrent    bool      `json:"is_current"`
}

type ClientInfo struct {
    IPAddress string
    UserAgent string
}
