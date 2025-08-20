package repository

import (
  "context"
  "errors"
  "fmt"
  "strings"
  "time"

  "github.com/F-Dupraz/ecommerce-with-go/model"

  "github.com/jackc/pgx/v5"
  "github.com/jackc/pgx/v5/pgconn"
  "github.com/jackc/pgx/v5/pgxpool"
)

var (
  ErrUserNotFound = errors.New("user not found")
  ErrDuplicateEmail = errors.New("email already exists")
  ErrDuplicateUsername = errors.New("username already exists")
  ErrNoRowsAffected = errors.New("no rows affected")
)

type UserRepository struct {
  db *pgxpool.Pool
}

func NewUserRepository(db *pgxpool.Pool) *UserRepository {
  return &UserRepository{
	db: db,
  }
}

func (r *UserRepository) CreateUserAtomic(ctx context.Context, user *model.User) error {
  err := r.db.QueryRow(ctx,
	`INSERT INTO users (id, email, username, password, address, city, country) 
	VALUES ($1, $2, $3, $4, $5, $6, $7) 
	RETURNING created_at, updated_at`,
	user.ID, user.Email, user.Username, user.Password, 
	user.Address, user.City, user.Country,
  ).Scan(&user.CreatedAt, &user.UpdatedAt)

  if err != nil {
	return r.translateError(err)  // translateError ya maneja los duplicados
  }

  return nil
}

func (r *UserRepository) GetByID(ctx context.Context, id string) (*model.User, error) {
  var user model.User
  err := r.db.QueryRow(ctx, "SELECT id, email, password, username, address, city, country, created_at, updated_at FROM users WHERE id = $1 AND deleted_at IS NULL", id).Scan(
  &user.ID, &user.Email, &user.Password, &user.Username, &user.Address, &user.City, &user.Country, &user.CreatedAt, &user.UpdatedAt)

  if err != nil {
	if errors.Is(err, pgx.ErrNoRows) {
	  return nil, ErrUserNotFound
	}

	return nil, fmt.Errorf("failed to get user by id: %w", err)
  }

  return &user, nil
}

func (r *UserRepository) GetByEmail(ctx context.Context, email string) (*model.User, error) {
  var user model.User
  err := r.db.QueryRow(ctx, "SELECT id, email, password, username, address, city, country, created_at, updated_at FROM users WHERE email = $1 AND deleted_at IS NULL", email).Scan(
  &user.ID, &user.Email, &user.Username, &user.Password, &user.Address, &user.City, &user.Country, &user.CreatedAt, &user.UpdatedAt)

  if err != nil {
	if errors.Is(err, pgx.ErrNoRows) {
	  return nil, ErrUserNotFound
	}

	return nil, fmt.Errorf("failed to get user by id: %w", err)
  }

  return &user, nil
}

func (r *UserRepository) Update(ctx context.Context, id string, updates map[string]interface{}) (*model.User, error) {
  var user model.User

  setClauses := []string{}
  args := []interface{}{}
  argID := 1

  for field, value := range updates {
	setClauses = append(setClauses, fmt.Sprintf("%s = $%d", field, argID))
	args = append(args, value)
	argID++
  }

  setClauses = append(setClauses, fmt.Sprintf("updated_at = $%d", argID))
  args = append(args, time.Now())
  argID++

  args = append(args, id)

  query := fmt.Sprintf(
	"UPDATE users SET %s WHERE id = $%d AND deleted_at IS NULL RETURNING id, email, username, address, city, country, created_at, updated_at",
	strings.Join(setClauses, ", "),
	argID,
  )

  err := r.db.QueryRow(ctx, query).Scan(&user.ID, &user.Email, &user.Username, &user.Address, &user.City, &user.Country, &user.CreatedAt, &user.UpdatedAt)
  if err != nil {
	if errors.Is(err, pgx.ErrNoRows) {
	  return nil, ErrUserNotFound
	}

	return nil, fmt.Errorf("failed to update user: %w", err)
  }

  return &user, nil
}

func (r *UserRepository) Delete(ctx context.Context, id string) (time.Time, error) {
  var userDeletedAt time.Time
  err := r.db.QueryRow(ctx, "UPDATE users SET deleted_at = NOW() WHERE id = $1 AND deleted_at IS NULL RETURNING deleted_at", id).Scan(&userDeletedAt)
  if err != nil {
	if errors.Is(err, pgx.ErrNoRows) {
	  return time.Time{}, ErrUserNotFound
	}

	return time.Time{}, fmt.Errorf("failed to delete user: %w", err)
  }

  return userDeletedAt, nil
}

func (r *UserRepository) ExistsByEmail(ctx context.Context, email string) (bool, error) {
  var exists bool
  err := r.db.QueryRow(ctx, 
	"SELECT EXISTS(SELECT 1 FROM users WHERE email = $1 AND deleted_at IS NULL)", 
	email,
  ).Scan(&exists)

  if err != nil {
	return false, fmt.Errorf("failed to check email existence: %w", err)
  }

  return exists, nil
}

func (r *UserRepository) ExistsByUsername(ctx context.Context, username string) (bool, error) {
  var exists bool
  err := r.db.QueryRow(ctx, 
	"SELECT EXISTS(SELECT 1 FROM users WHERE email = $1 AND deleted_at IS NULL)", 
	username,
  ).Scan(&exists)

  if err != nil {
	return false, fmt.Errorf("failed to check email existence: %w", err)
  }

  return exists, nil
}

func (r *UserRepository) translateError(err error) error {
  if err == nil {
    return nil
  }

  if errors.Is(err, pgx.ErrNoRows) {
    return ErrUserNotFound
  }

  var pgErr *pgconn.PgError
  if errors.As(err, &pgErr) {
    if pgErr.Code == "23505" { // unique_violation
      if strings.Contains(pgErr.Detail, "email") {
        return ErrDuplicateEmail
      }
      if strings.Contains(pgErr.Detail, "username") {
        return ErrDuplicateUsername
      }
    }
  }

  return err
}
