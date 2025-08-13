package dto

import (
  "time"
)

// Requests

type CreateUserRequest struct {
  Username string `json:"username" validate:"required,min=3,max=50,alphanum"`
  Email    string `json:"email" validate:"required,email"`
  Password string `json:"password" validate:"required,min=8,max=72,password"`
  Address  string `json:"address,omitempty" validate:"max=200"`
  City     string `json:"city,omitempty" validate:"max=100"`
  Country  string `json:"country,omitempty" validate:"omitempty,iso3166_1_alpha2"`
}

type UpdateUserRequest struct {
  Username *string `json:"username,omitempty" validate:"omitempty,min=3,max=50,alphanum"`
  Email    *string `json:"email,omitempty" validate:"omitempty,email"`
  Password *string `json:"password,omitempty" validate:"omitempty,min=8,max=72,password"`
  Address  *string `json:"address,omitempty" validate:"omitempty,max=200"`
  City     *string `json:"city,omitempty" validate:"omitempty,max=100"`
  Country  *string `json:"country,omitempty" validate:"omitempty,iso3166_1_alpha2"`
}

type GetUserByIDRequest struct {
  ID string `param:"id" validate:"required,uuid"`
}

type GetUserByEmailRequest struct {
  Email string `query:"email" validate:"required,email"`
}

type DeleteUserRequest struct {
  ID string `param:"id" validate:"required,uuid"`
}

// Responses

type UserResponse struct {
  ID        string    `json:"id"`
  Username  string    `json:"username"`
  Email     string    `json:"email"`
  Address   string    `json:"address,omitempty"`
  City      string    `json:"city,omitempty"`
  Country   string    `json:"country,omitempty"`
  CreatedAt time.Time `json:"created_at"`
  UpdatedAt time.Time `json:"updated_at"`
}

type CreateUserResponse struct {
  ID      string        `json:"id"`
  User    *UserResponse `json:"user"`
  Message string        `json:"message"`
}

type UpdateUserResponse struct {
  User    *UserResponse `json:"user"`
  Message string        `json:"message"`
}

type DeleteUserResponse struct {
  ID        string    `json:"id"`
  Message   string    `json:"message"`
  DeletedAt time.Time `json:"deleted_at"`
}

type ErrorResponse struct {
  Error      string            `json:"error"`
  Message    string            `json:"message"`
  StatusCode int               `json:"status_code"`
  RequestID  string            `json:"request_id,omitempty"`
  Details    map[string]string `json:"details,omitempty"`
  Timestamp  time.Time         `json:"timestamp"`
}

