package handler

import (
  "fmt"
  "net/http"
)

type User struct {}

func (u *User) Create(w http.ResponseWriter, r *http.Request) {
  fmt.Println("Create a user")
}

func (u *User) ListAllUsers(w http.ResponseWriter, r *http.Request) {
  fmt.Println("List all users")
}

func (u *User) GetByID(w http.ResponseWriter, r *http.Request) {
  fmt.Println("Get a user by ID")
}

func (u *User) UpdateByID(w http.ResponseWriter, r *http.Request) {
  fmt.Println("Update a user by ID")
}

func (u *User) DeleteByID(w http.ResponseWriter, r *http.Request) {
  fmt.Println("Delete a user by ID")
}
