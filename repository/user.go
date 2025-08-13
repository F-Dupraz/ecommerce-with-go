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
  tx, err := r.db.Begin(ctx)
  if err != nil {
	return err
  }

  defer tx.Rollback(ctx)

  var existsUsername bool
  var existsEmail bool

  err := tx.QueryRow(ctx, "SELECT EXISTS(SELECT 1 FROM users WHERE username = $1)", user.Username).Scan(&existsUsername)
  if err != nil {
	return err
  }

  if existsUsername {
	return ErrDuplicateUsername
  }

  err := tx.QueryRow(ctx, "SELECT EXISTS(SELECT 1 FROM users WHERE email = $1)", user.Email).Scan(&existsEmail)
  if err != nil {
	return nil
  }

  if existsEmail {
	return ErrDuplicateEmail
  }

  err := tx.QueryRow(ctx, "INSERT INTO users (id, email, username, password, address, city, country) VALUES ($1, $2, $3, $4, $5, $6) RETURNING created_at, updated_at", 
	user.ID, user.Email, user.Username, user.Password, user.Address, user.Country).Scan(&user.CreatedAt, &user.UpdatedAt)

  if err != nil {
	return r.translateError(err)
  }

  return tx.Commit(ctx)
}

func (r *UserRepository) GetByID(ctx context.Context, id string) (*model.User, error) {
    // TODO: SELECT * FROM users WHERE id = $1 AND deleted_at IS NULL
    // TODO: Manejar pgx.ErrNoRows → ErrUserNotFound
    
  return nil, nil
}

func (r *UserRepository) GetByEmail(ctx context.Context, email string) (*model.User, error) {
    // TODO: SELECT * FROM users WHERE email = $1 AND deleted_at IS NULL
    // TODO: Manejar pgx.ErrNoRows → ErrUserNotFound
    
  return nil, nil
}

func (r *UserRepository) Update(ctx context.Context, id string, user *model.User) error {
    // TODO: Build dynamic UPDATE query based on non-nil fields
    // TODO: RETURNING updated_at para actualizar el timestamp en el struct
    // TODO: Check rows affected
    
  return nil
}

func (r *UserRepository) Delete(ctx context.Context, id string) (time.Time, error) {
    // TODO: UPDATE users SET deleted_at = NOW() WHERE id = $1 AND deleted_at IS NULL
    // TODO: RETURNING deleted_at
    // TODO: Check rows affected
    
  return time.Time{}, nil
}

func (r *UserRepository) ExistsByEmail(ctx context.Context, email string) (bool, error) {
    // TODO: SELECT EXISTS(SELECT 1 FROM users WHERE email = $1 AND deleted_at IS NULL)
    
  return false, nil
}

func (r *UserRepository) ExistsByUsername(ctx context.Context, username string) (bool, error) {
    // TODO: SELECT EXISTS(SELECT 1 FROM users WHERE username = $1 AND deleted_at IS NULL)
    
  return false, nil
}

func (r *UserRepository) List(ctx context.Context, limit, offset int) ([]*model.User, int64, error) {
    // TODO: Query para obtener users paginados
    // TODO: Query separada para COUNT total
    // TODO: Retornar slice de users y total count
    
  return nil, 0, nil
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
