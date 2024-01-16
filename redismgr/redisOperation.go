package redismgr

// Get gets the value for given key
func (r *redisService) Get(key string) (string, error) {
	if r.isStandAlone {
		return r.standalone.Get(key).Result()
	}
	return r.cluster.Get(key).Result()
}

// Set sets the value for given key
func (r *redisService) Set(key, value string) error {
	if r.isStandAlone {
		return r.standalone.Set(key, value, 0).Err()
	}
	return r.cluster.Set(key, value, 0).Err()
}

// HashSet to set up the hash with key and field value
func (r *redisService) HashSet(key, field string, value interface{}) error {
	if r.isStandAlone {
		return r.standalone.HSet(key, field, value).Err()
	}
	return r.cluster.HSet(key, field, value).Err()

}

// HashDel to delete the key from hash
func (r *redisService) HashDel(key, field string) error {
	if r.isStandAlone {
		return r.standalone.HDel(key, field).Err()
	}
	return r.cluster.HDel(key, field).Err()
}

// HashGet  to get the specific value of hash key
func (r *redisService) HashGet(key, field string) (string, error) {
	if r.isStandAlone {
		return r.standalone.HGet(key, field).Result()
	}
	return r.cluster.HGet(key, field).Result()
}

// HashGetAll returns all the values of hast set
func (r *redisService) HashGetAll(key string) (map[string]string, error) {
	if r.isStandAlone {
		return r.standalone.HGetAll(key).Result()
	}
	return r.cluster.HGetAll(key).Result()
}

// Delete deletes a key
func (r *redisService) Delete(key string) (int64, error) {
	if r.isStandAlone {
		cmd := r.standalone.Del(key)
		return cmd.Val(), cmd.Err()
	}
	cmd := r.cluster.Del(key)
	return cmd.Val(), cmd.Err()
}
