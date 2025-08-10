package handler

import (
  "context"
  "fmt"
  "net/http"
  "encoding/json"

  "github.com/go-chi/chi/v5"
  "github.com/go-playground/validator/v10"

  "github.com/F-Dupraz/ecommerce-with-go/dto"
)

type UserService interface {
  CreateUser(ctx context.Context, usr dto.CreateUserRequest) (*dto.CreateUserResponse, error)
  GetUserById(ctx context.Context, user_id dto.GetUserByIDRequest) (*dto.UserResponse, error)
  GetUserByEmail(ctx context.Context, email dto.GetUserByEmailRequest) (*dto.UserResponse, error)
  UpdateUser(ctx context.Context, user_id string, usr dto.UpdateUserRequest) (*dto.UpdateUserResponse, error)
  DeleteUser(ctx context.Context, user_id dto.DeleteUserRequest) (*dto.DeleteUserResponse, error)
}

type UserHandler struct {
  BaseHandler
  userService UserService
}

func NewUserHandler(userService UserService, validator *validator.Validate) *UserHandler {
  return &UserHandler{
	userService: userService,
	BaseHandler: BaseHandler{validator: validator},
  }
}

func (u *UserHandler) RegisterRoutes(router chi.Router) {
  router.Route("/users", func (r chi.Router) {
	router.Get("/", u.GetUserByEmail)
	router.Get("/{id}", u.GetUserByID)
	router.Post("/", u.CreateUser)
	router.Put("/{id}", u.UpdateUser)
	router.Delete("/{id}", u.DeleteUser)
  })
}

func (u *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
  var req dto.CreateUserRequest

  if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
	u.respondWithError(w, http.StatusBadRequest, "Cannot parse JSON: " + err.Error(), nil )
	return
  }

  if err := u.validator.Struct(req); err != nil {
	validationErrors := dto.FormatValidationErrors(err)
	u.respondWithError(w, http.StatusUnprocessableEntity, "Validation failed!", validationErrors)
	return
  }

  response, err := u.userService.CreateUser(r.Context(), req)
  if err != nil {
	u.respondWithError(w, http.StatusInternalServerError, "Failed to create user", nil)
	return
  }

  u.respondWithSuccess(w, http.StatusCreated, response)
}

func (u *UserHandler) GetUserByEmail(w http.ResponseWriter, r *http.Request) {
  fmt.Println("List all users")
}

func (u *UserHandler) GetUserByID(w http.ResponseWriter, r *http.Request) {
  fmt.Println("Get a user by ID")
}

func (u *UserHandler) UpdateUser(w http.ResponseWriter, r *http.Request) {
  fmt.Println("Update a user by ID")
}

func (u *UserHandler) DeleteUser(w http.ResponseWriter, r *http.Request) {
  fmt.Println("Delete a user by ID")
}
