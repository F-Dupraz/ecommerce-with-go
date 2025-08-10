package application

import (
  "context"
  "fmt"
  "net/http"
  "time"
)

type App struct {
  router http.Handler
}

func New() *App {
  app := &App{
    router: loadRoutes(),
  }

  return app
}

func (a *App) Start(ctx context.Context) error {
  server := &http.Server{
    Addr:         ":3000",
    Handler:      a.router,
    ReadTimeout:  15 * time.Second,
    WriteTimeout: 15 * time.Second,
    IdleTimeout:  60 * time.Second,
  }

  fmt.Println("Server starting on port 3000...")

  err := server.ListenAndServe()
  if err != nil {
    return fmt.Errorf("failed to start server: %w", err)
  }

  return nil
}
