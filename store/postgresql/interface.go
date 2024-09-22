package postgresql

import (
	"context"
	"crud/model"
	"github.com/go-redis/redis/v8"
	"database/sql"

)

type SqlManager interface {
	SetUserData(ctx context.Context, db *sql.DB , user model.User) error
	// GetUserData(ctx context.Context, client *redis.Client, userID string) (model.User, error)
	// UpdateUserData(ctx context.Context, client *redis.Client, user model.User) error
	// DeleteUserData(ctx context.Context, client *redis.Client, userID string) error
	GetAll(ctx context.Context, client *redis.Client) (model.StudentInfo, error)
}
