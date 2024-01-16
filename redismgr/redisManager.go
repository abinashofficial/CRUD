package redismgr

import (
	"crud/log"
	"crud/tapcontext"
	"crud/taputils"
	"errors"
	"fmt"
	"github.com/go-redis/redis"
	"strconv"
)

type CacheManager interface {
	Get(key string) (string, error)
	Set(key, value string) error
	Delete(key string) (int64, error)
	HashGet(key, field string) (string, error)
	HashSet(key, field string, value interface{}) error
	HashDel(key, field string) error
	HashGetAll(key string) (map[string]string, error)
}

type redisService struct {
	isStandAlone bool
	cluster      *redis.ClusterClient
	standalone   *redis.Client
}

// NewRedisMgr  returns mongo client to the host provided
func NewRedisMgr(ctx tapcontext.TContext, url string, isStandAlone bool, optionsMap ...map[string]interface{}) (CacheManager, error) {
	var mgr redisService
	if len(url) == 0 {
		log.GenericError(ctx, errors.New("env variable not set for redis address"), nil)
	}
	isLocal, _ := strconv.ParseBool(taputils.GetEnv(ctx, "local", true))
	if isLocal {
		isStandAlone = true // set to true  for local development
	}
	log.GenericInfo(ctx, "redis endpoint available: "+url, nil)
	if isStandAlone {
		options, err := redis.ParseURL(url)
		if err != nil {
			ctx := tapcontext.TContext{}
			log.GenericError(ctx, fmt.Errorf("URL parsing failed, err: %v", err), nil)
			return &mgr, err
		}
		if len(optionsMap) > 0 {
			optionMap := optionsMap[0]
			if poolSize, ok := optionMap["poolSize"]; ok {
				options.PoolSize, ok = poolSize.(int)
				if !ok {
					log.GenericError(ctx, errors.New("incorrect pool size data type"))
				}
			}
			if maxRetries, ok := optionMap["maxRetries"]; ok {
				options.MaxRetries, ok = maxRetries.(int)
				if !ok {
					log.GenericError(ctx, errors.New("incorrect max retries data type"))
				}
			}
		}
		client := redis.NewClient(options)
		_, err = client.Ping().Result()
		if err != nil {
			return &mgr, err
		}
		mgr.isStandAlone = true
		mgr.standalone = client
	} else {
		clusterOptions := &redis.ClusterOptions{
			Addrs: []string{url},
		}
		if len(optionsMap) > 0 {
			optionMap := optionsMap[0]
			if poolSize, ok := optionMap["poolSize"]; ok {
				clusterOptions.PoolSize, ok = poolSize.(int)
				if !ok {
					log.GenericError(ctx, errors.New("incorrect pool size data type"))
				}
			}
			if maxRetries, ok := optionMap["maxRetries"]; ok {
				clusterOptions.MaxRetries, ok = maxRetries.(int)
				if !ok {
					log.GenericError(ctx, errors.New("incorrect max retries data type"))
				}
			}
		}
		client := redis.NewClusterClient(clusterOptions)
		_, err := client.Ping().Result()
		if err != nil {
			return &mgr, err
		}
		mgr.cluster = client
	}

	return &mgr, nil
}
