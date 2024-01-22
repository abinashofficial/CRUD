package app

import (
	"context"
	"crud/handlers"
	"crud/handlers/fields"
	"crud/log"
	"crud/model"
	"crud/redismgr"
	mongoStore "crud/store"
	authRepo "crud/store/auth"
	"crud/store/mongo"
	"crud/tapcontext"
	"os"
)

var h handlers.Store
var repos mongoStore.Store

func setupRepos(client model.MongoClient, redisClient redismgr.CacheManager, IsRedisRequired bool) {
	repos = mongoStore.Store{
		AuthRepo: authRepo.New(client),
	}
}

func setupHandlers() {
	h = handlers.Store{
		FieldsHandler: fields.New(),
	}
}

func Start() {
	envPort := os.Getenv("PORT")
	redisHost := os.Getenv("REDIS_URL")
	mongoURL := os.Getenv("MONGO_URL")

	//envPort := "8080"
	//redisHost := "redis://default:5ft4zGUqto38vDS6R0WOcjuudJk6LKlf@redis-11532.c264.ap-south-1-1.ec2.cloud.redislabs.com:11532"
	//mongoURL := "mongodb+srv://abinash1411999:Abinash1234@cluster0.fpdirql.mongodb.net/?retryWrites=true&w=majority"

	ctx := tapcontext.TContext{
		Context:    context.Background(),
		TapContext: tapcontext.TapContext{},
	}
	var IsRedisRequired bool
	allMongoURL := map[string]string{
		"mongoURL":          mongoURL,
		"secondaryMongoURL": "",
	}
	client, mongoErr := mongo.Init(ctx, allMongoURL)
	defer mongo.Disconnect(ctx)

	if mongoErr.PrimaryClientError != nil {
		log.GenericError(ctx, mongoErr.PrimaryClientError, log.FieldsMap{"error": "MongoURL Connection Failed"})
		log.FatalLog(ctx, mongoErr.PrimaryClientError, nil)
	}

	redisClient, err := redismgr.NewRedisMgr(ctx, redisHost, true)
	if err != nil {
		log.FatalLog(ctx, err, nil)
	}
	setupRepos(client, redisClient, IsRedisRequired)

	setupHandlers()
	runServer(envPort, h, ctx)
}
