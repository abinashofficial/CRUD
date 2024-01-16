package auth

import (
	"crud/consts"
	"crud/model"
	"crud/tapcontext"
	"crud/utils"
	"errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func New(db model.MongoClient) Repository {
	repo := db.SecondaryClient.Database(consts.SupportPortalDB).Collection(consts.AuthApplications)
	return &authRepo{
		repo: repo,
	}
}

type authRepo struct {
	repo *mongo.Collection
}

func (r *authRepo) FindAppByKey(ctx tapcontext.TContext, appName, key, secret string) (app model.OpenApplication, err error) {
	filter := bson.M{"app_name": appName, "app_key": key, "app_secret": secret}
	result := r.repo.FindOne(ctx.Context, filter)

	if result.Err() != nil {
		return app, errors.New(utils.GetError(consts.ErrAppNotFound, ctx.Locale))
	}

	err = result.Decode(&app)
	if err != nil {
		return app, errors.New(utils.GetError(consts.ErrDecodeApp, ctx.Locale))
	}

	return app, nil
}
