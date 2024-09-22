package app

import (
	"crud/handlers"
	"crud/handlers/fields"
	mongoStore "crud/store"
	"crud/store/postgresql"
	"crud/store/redismanager"
	"database/sql"
	"fmt"
	"github.com/go-redis/redis/v8"
	_ "github.com/lib/pq"
)

var h handlers.Store
// var db *sql.DB

var repos mongoStore.Store

func setupRepos() {
	repos = mongoStore.Store{
		CacheStore: redismanager.New(),
	}
}

func setupHandlers(cacheRepo redismanager.CacheManager, client *redis.Client, sqlRepo postgresql.SqlManager, sqlDB *sql.DB) {
	h = handlers.Store{
		FieldsHandler: fields.New(cacheRepo, client,sqlRepo, sqlDB),
	}
}

func Start() {
	// envPort := os.Getenv("PORT")
	// rediUrl := os.Getenv("REDIS_URL")
	// sqlUrl := os.Getenv("SQL_URL")
	envPort := "8080"
	rediUrl := "redis://default:Zbk8sTu9N6zakmvletQG2mT4LfAZ034b@redis-11301.c212.ap-south-1-1.ec2.redns.redis-cloud.com:11301"
	sqlUrl := "postgresql://develop_owner:fkdK1b9vzohQ@ep-green-feather-a1lerkc8.ap-southeast-1.aws.neon.tech/develop?sslmode=require"
	sqlDB, err := sql.Open("postgres", sqlUrl)
	if err != nil {
		fmt.Println(err, "sql")
	}
	options, err := redis.ParseURL(rediUrl)
	if err != nil {
		fmt.Println(err, "redis")
	}
	redisClient := redis.NewClient(options)
	setupRepos()
	setupHandlers(repos.CacheStore, redisClient,repos.SqlStore, sqlDB)

	runServer(envPort, h)
}
