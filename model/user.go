package model

import (
  "time"
)

type User struct {
    ID        string     `db:"id"`
    Username  string     `db:"username"`
    Email     string     `db:"email"`
    Password  string     `db:"password" json:"-"`
    Address   string     `db:"address"`
    City      string     `db:"city"`
    Country   string     `db:"country"`
    CreatedAt time.Time  `db:"created_at"`
    UpdatedAt time.Time  `db:"updated_at"`
    DeletedAt *time.Time `db:"deleted_at"`
}

