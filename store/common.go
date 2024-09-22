package store

import (
	"crud/store/postgresql"
	"crud/store/redismanager"
)

type Store struct {
	CacheStore redismanager.CacheManager
	SqlStore postgresql.SqlManager
}
