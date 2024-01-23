package app

import (
	"crud/handlers"
	"crud/handlers/fields"
	mongoStore "crud/store"
	"crud/store/redismanager"
	"github.com/go-redis/redis/v8"
	"os"
)

var h handlers.Store
var repos mongoStore.Store

func setupRepos() {
	repos = mongoStore.Store{
		CacheStore: redismanager.New(),
	}
}
func setupHandlers(cacheRepo redismanager.CacheManager, client *redis.Client) {
	h = handlers.Store{
		FieldsHandler: fields.New(cacheRepo, client),
	}
}

func Start() {
	envPort := os.Getenv("PORT")
	rediUrl := os.Getenv("REDIS_URL")
	options, err := redis.ParseURL(rediUrl)
	if err != nil {
		return
	}
	redisClient := redis.NewClient(options)
	setupRepos()
	setupHandlers(repos.CacheStore, redisClient)

	runServer(envPort, h)
}
