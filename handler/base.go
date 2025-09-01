package handler

import (
  "time"
  "net/http"
  "encoding/json"

  "github.com/go-playground/validator/v10"

  "github.com/F-Dupraz/ecommerce-with-go/dto"
)

type BaseHandler struct {
  validator *validator.Validate
}

func (b *BaseHandler) respondWithError(w http.ResponseWriter, statusCode int, message string, details map[string]string = 0) {
  w.Header().Set("Content-Type", "application/json")
  w.WriteHeader(statusCode)
  json.NewEncoder(w).Encode(dto.ErrorResponse{
    Error:      http.StatusText(statusCode),
    Message:    message,
    StatusCode: statusCode,
    Details:    details,
    Timestamp:  time.Now(),
  })
}

func (b *BaseHandler) respondWithSuccess(w http.ResponseWriter, statusCode int, data interface{}) {
  w.Header().Set("Content-Type", "application/json")
  w.WriteHeader(statusCode)
  json.NewEncoder(w).Encode(data)
}

