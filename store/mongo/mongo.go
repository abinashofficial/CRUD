package mongo

import (
	"crud/consts"
	"crud/log"
	"crud/model"
	"crud/tapcontext"
	"crud/utils"
	"go.elastic.co/apm/module/apmmongo"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"sync"

	"context"
	"errors"
)

var (
	initMongo sync.Once
	client    model.MongoClient
	err       model.MongoError
)

// Init MongoConn instance
func Init(ctx tapcontext.TContext, server map[string]string) (model.MongoClient, model.MongoError) {
	initMongo.Do(func() {
		for index, val := range server {
			clientOpts := options.Client().ApplyURI(val).SetMinPoolSize(10).SetMaxPoolSize(100).SetMonitor(apmmongo.CommandMonitor()).SetAppName("go_support_portal")
			if index == "mongoURL" {
				client.PrimaryClient, err.PrimaryClientError = mongo.Connect(context.TODO(), clientOpts)
			} else {
				client.SecondaryClient, err.SecondaryClientError = mongo.Connect(context.TODO(), clientOpts)
			}
		}
	})
	if err.PrimaryClientError != nil {
		errMsg := utils.GetError(consts.ErrMongoConn, ctx.Locale) + err.PrimaryClientError.Error()
		log.GenericError(ctx, errors.New(errMsg), log.FieldsMap{"error": "MongoURL Connection Failed"})
	} else if err.SecondaryClientError != nil {
		client.SecondaryClient = client.PrimaryClient
		log.GenericInfo(ctx, "MongoInit", log.FieldsMap{"message": "Secondary DB Connection failed, " +
			"Marking secondary db client also to primary"})
	}

	return client, err
}

// Disconnect mongo instance
func Disconnect(ctx tapcontext.TContext) model.MongoError {
	if tempErr := client.PrimaryClient.Disconnect(context.TODO()); tempErr != nil {
		errMsg := utils.GetError(consts.ErrMongoDisconnect, ctx.Locale) + tempErr.Error()
		log.GenericError(ctx, tempErr, log.FieldsMap{"Primary Mongo Disconnect": errMsg})
		err.PrimaryClientError = tempErr
	}
	if tempErr := client.SecondaryClient.Disconnect(context.TODO()); tempErr != nil {
		errMsg := utils.GetError(consts.ErrMongoDisconnect, ctx.Locale) + tempErr.Error()
		log.GenericError(ctx, tempErr, log.FieldsMap{"Secondary Mongo Disconnect": errMsg})
		err.PrimaryClientError = tempErr
	}
	return err
}
