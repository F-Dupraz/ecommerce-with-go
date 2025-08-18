package service

import (
  "fmt"
  "time"
  "context"
  "errors"

  "github.com/F-Dupraz/ecommerce-with-go/dto"
  "github.com/F-Dupraz/ecommerce-with-go/model"

  "golang.org/x/crypto/bcrypt"
  "github.com/google/uuid"
)

// Custom errors para el dominio
var (
  ErrUserNotFound = errors.New("user not found")
  ErrEmailAlreadyExists = errors.New("email already exists")
  ErrUsernameAlreadyExists = errors.New("username already exists")
  ErrInvalidCredentials = errors.New("invalid credentials")
  ErrInvalidUserID = errors.New("invalid user id")
)

// Repository interface - contrato con la capa de datos
type UserRepository interface {
  CreateUserAtomic(ctx context.Context, user *model.User) error
  GetByID(ctx context.Context, id string) (*model.User, error)
  GetByEmail(ctx context.Context, email string) (*model.User, error)
  Update(ctx context.Context, id string, updates map[string]interface{}) (*model.User, error)
  Delete(ctx context.Context, id string) error
  ExistsByEmail(ctx context.Context, email string) (bool, error)
  ExistsByUsername(ctx context.Context, username string) (bool, error)
}

type UserService struct {
    repo UserRepository
}

func NewUserService(repo UserRepository) *UserService {
    return &UserService{
        repo: repo,
    }
}

func (s *UserService) CreateUser(ctx context.Context, req dto.CreateUserRequest) (*dto.CreateUserResponse, error) {
  hashedPass, err := bcrypt.GenerateFromPassword(req.Password, 12)
  if err != nil {
	return nil, fmt.Errorf("Failed to hash password: %w", err)
  }

  newID := uuid.New().String()

  var newUser model.User{
	ID: newID,
	Username: req.Username,
	Email: req.Email,
	Password: string(hashedPass),
	Address: req.Address,
	City: req.City,
	Country: req.Country,
  }

  if err := s.repo.CreateUserAtomic(ctx, &newUser); err != nil {
	if errors.Is(err, repository.ErrDuplicateEmail) {
	  return nil, ErrEmailAlreadyExists
	}

	if errors.Is(err, repository.ErrDuplicateUsername) {
	  return nil, ErrUsernameAlreadyExists
	}

	return nil, fmt.Errorf("failed to create user: %w", err)
  }

  return &dto.CreateUserResponse{
	ID: newUser.ID,
	User: s.modelToResponse(&newUser),
	Message: "User created successfully!",
  }, nil
}

func (s *UserService) GetUserByID(ctx context.Context, req dto.GetUserByIDRequest) (*dto.UserResponse, error) {
  if err := uuid.Validate(userID); err != nil {
	return nil, ErrInvalidUserID
  }

  serId := req.ID

  var user model.User

  user, err := s.repo.GetUserByID(ctx, userID)
  if err != nil {
	return nil, ErrInvalidUserID
  }

  return &dto.UserResponse{
	ID: user.ID,
	Username: user.Username,
	Email: user.Email,
	Address: user.Email,
	City: user.City,
	Country: user.Country,
	CreatedAt: user.CreatedAt,
	UpdatedAt: user.UpdatedAt,
  }, nil
}

func (s *UserService) GetUserByEmail(ctx context.Context, req dto.GetUserByEmailRequest) (*dto.UserResponse, error) {
  userEmail := req.Email

  var user model.User

  user, err := s.repo.GetUserByEmail(ctx, userEmail)
  if err != nil {
	return nil, ErrUserNotFound
  }

  return &dto.UserResponse{
	ID: user.ID,
	Username: user.Username,
	Email: user.Email,
	Address: user.Email,
	City: user.City,
	Country: user.Country,
	CreatedAt: user.CreatedAt,
	UpdatedAt: user.UpdatedAt,
  }, nil
}

func (s *UserService) UpdateUser(ctx context.Context, userID string, req dto.UpdateUserRequest) (*dto.UpdateUserResponse, error) {
  if err := uuid.Validate(userID); err != nil {
	return nil, ErrInvalidUserID
  }

  var currentUser model.User

  currentUser, err := s.repo.GetUserByID(ctx, userID)
  if err != nil {
	return nil, ErrInvalidUserID
  }

  if req.Email != nil && *req.Email != currentUser.Email {
	  exists, err := s.repo.ExistsByEmail(ctx, *req.Email)
	  if err != nil {
		return nil, fmt.Errorf("failed to check email existence: %w", err)
	  }
	  if exists {
		return nil, ErrEmailAlreadyExists
	  }
  }

  if req.Username != nil && *req.Username != currentUser.Username {
	  exists, err := s.repo.ExistsByUsername(ctx, *req.Username)
	  if exists {
		return nil, ErrUsernameAlreadyExists
	  }
  }

  if req.Password != nil {
	hashedPass, err := bcrypt.GenerateFromPassword(*req.Password, 12)
	if err != nil {
	  return nil, fmt.Errorf("Failed to hash password: %w", err)
	}
  }

  updates := make(map[string]interface{})
  if req.Username != nil {
	updates["username"] = *req.Username
  }
  if req.Email != nil {
	updates["email"] = *req.Email
  }
  if req.Password != nil {
	hashedPass, err := bcrypt.GenerateFromPassword([]byte(*req.Password), 12)
	if err != nil {
	  return nil, fmt.Errorf("failed to hash password: %w", err)
	}
	updates["password"] = string(hashedPass)
  }
  if req.Address != nil {
	updates["address"] = *req.Address
  }
  if req.City != nil {
	updates["city"] = *req.City
  }
  if req.Country != nil {
	updates["country"] = *req.Country
  }

  updatedUser, err := s.repo.Update(ctx, userID, updates)
  if err != nil {
	return nil, fmt.Errorf("failed to update user: %w", err)
  }

  return &dto.UpdateUserResponse{
	User: s.modelToResponse(updatedUser),
	Message: "User updated successfully!",
  }, nil
}

func (s *UserService) DeleteUser(ctx context.Context, req dto.DeleteUserRequest) (*dto.DeleteUserResponse, error) {
  if err := uuid.Validate(req.ID); err != nil {
	return nil, ErrInvalidUserID
  }

  userID := req.ID

  userExists, err := s.repo.GetUserByID(ctx, userID)
  if err != nil {
	return nil, ErrInvalidUserID
  }

  s.repo.Delete(ctx, userID)

  return &dto.DeleteUserResponse{
	ID: userID,
	Message: "User deleted succesfully!",
	DeletedAt: time.Now()
  }, nil
}

// Helper

func (s *UserService) modelToResponse(user *model.User) *dto.UserResponse {
  return &dto.UserResponse{
	ID:        user.ID,
	Username:  user.Username,
	Email:     user.Email,
	Address:   user.Address,
	City:      user.City,
	Country:   user.Country,
	CreatedAt: user.CreatedAt,
	UpdatedAt: user.UpdatedAt,
  }

  return nil
}

