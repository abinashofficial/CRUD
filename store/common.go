package store

import (
	"crud/store/redismanager"
)

type Store struct {
	CacheStore redismanager.CacheManager
}
