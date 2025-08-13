package service

import (
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
  Update(ctx context.Context, id string, user *model.User) error
  Delete(ctx context.Context, id string) error
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

  return &responseUser{
	ID: newUser.ID,
	User: s.modelToResponse(&newUser),
	Message: "User created successfully!",
  }, nil
}

func (s *UserService) GetUserById(ctx context.Context, req dto.GetUserByIDRequest) (*dto.UserResponse, error) {
    // TODO: Validate UUID format if needed
    
    // TODO: Call repo.GetByID with req.ID
    
    // TODO: Handle ErrUserNotFound if user doesn't exist
    
    // TODO: Convert model.User to dto.UserResponse
    
    return nil, nil
}

func (s *UserService) GetUserByEmail(ctx context.Context, req dto.GetUserByEmailRequest) (*dto.UserResponse, error) {
    // TODO: Call repo.GetByEmail with req.Email
    
    // TODO: Handle ErrUserNotFound if user doesn't exist
    
    // TODO: Convert model.User to dto.UserResponse
    
    return nil, nil
}

func (s *UserService) UpdateUser(ctx context.Context, userID string, req dto.UpdateUserRequest) (*dto.UpdateUserResponse, error) {
    // TODO: Validate UUID format
    
    // TODO: Get existing user with repo.GetByID to ensure exists
    
    // TODO: If email is being updated, check it's not taken (excluding current user)
    
    // TODO: If username is being updated, check it's not taken (excluding current user)
    
    // TODO: If password is being updated, hash it
    
    // TODO: Build updated model.User with only changed fields
    
    // TODO: Call repo.Update
    
    // TODO: Get updated user and convert to response
    
    // TODO: Build and return dto.UpdateUserResponse
    
    return nil, nil
}

func (s *UserService) DeleteUser(ctx context.Context, req dto.DeleteUserRequest) (*dto.DeleteUserResponse, error) {
    // TODO: Validate UUID format
    
    // TODO: Check if user exists with repo.GetByID
    
    // TODO: Call repo.Delete (soft delete ideally)
    
    // TODO: Build and return dto.DeleteUserResponse with deletion timestamp
    
    return nil, nil
}

// Helpers

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

func generateUUID() string {
    // TODO: Implement UUID generation (use google/uuid library)
    return ""
}

func hashPassword(password string) (string, error) {
    // TODO: Wrapper for bcrypt.GenerateFromPassword with your cost setting
    return "", nil
}

func comparePasswords(hashedPassword, password string) error {
    // TODO: Wrapper for bcrypt.CompareHashAndPassword
    return nil
}
