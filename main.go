package main

import (
  "fmt"
  "context"

  "github.com/F-Dupraz/ecommerce-with-go/application"
)

func main() {
  app := application.New()

  err := app.Start(context.TODO())
  if err != nil {
	fmt.Println("failed to start the server: %w", err)
  }
}

