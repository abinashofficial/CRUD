package redismanager

import (
	"context"
	"crud/model"
	"github.com/go-redis/redis/v8"
)

type CacheManager interface {
	SetUserData(ctx context.Context, client *redis.Client, user model.StudentInfo) error
	GetUserData(ctx context.Context, client *redis.Client, userID string) (model.User, error)
	UpdateUserData(ctx context.Context, client *redis.Client, user model.User) error
	CreateUserData(ctx context.Context, client *redis.Client, user model.User) error
	DeleteUserData(ctx context.Context, client *redis.Client, userID string) error
	GetAll(ctx context.Context, client *redis.Client) (model.StudentInfo, error)
}
