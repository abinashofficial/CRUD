package auth

import (
	"crud/model"
	"crud/tapcontext"
)

type Repository interface {
	FindAppByKey(ctx tapcontext.TContext, appName, key, secret string) (model.OpenApplication, error)
}
