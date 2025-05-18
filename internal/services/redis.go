package services

import "time"

type IRedis interface {
	ExistsKey(key string) bool
	Save(key string, value interface{}, expirationTime int64) error
	Delete(key string) error
	Get(key string, result interface{}) error
	GetTTL(key string) (time.Duration, error)
}
