 package middleware

 import (
  "context"
 )

 type contextKey string

 const (
  UserIDKey contextKey = "user_id"
  UserRoleKey contextKey = "user_role"
  IsAdminKey contextKey = "is_admin"
  EmailKey contextKey = "user_email"
 )

 func GetUserID (ctx context.Context) (string, bool) {
  userID, ok := ctx.Value(UserIDKey).(string)
  return userID, ok
 }

func IsAdmin(ctx context.Context) bool {
  isUserAdmin := ctx.Value(IsAdminKey).(bool)
  return isUserAdmin
}

