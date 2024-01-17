package handler

import (
	"fmt"
	"github.com/redis/go-redis/v9"
	"sync"
	"time"
)

//go:generate mockgen -source=C:/Users/prsok/GolandProjects/CRUD_API/pkg/handler/cache.go -destination=mocks/mocks_cache.go

type Caches interface {
	Set(cacheKey string, data interface{}, duration time.Duration) error
	Get(key string) interface{}
}

// Time to live cacheTTl
const (
	cacheTTl = 60 * time.Hour
)

type CacheData struct { //"{\"id\":5,\"name\":\"5\",\"description\":\"5\"}"
	data      interface{}
	expiresAt time.Time
}

type CacheKey struct { //product:5
	cache map[string]CacheData //Cache-значение нашего ключа
}

func (c *CacheKey) Clear() { //очистка кэша
	c.cache = make(map[string]CacheData)
}

func NewCacheKey() *CacheKey {
	return &CacheKey{
		cache: make(map[string]CacheData),
	}
}

var once sync.Once //пакет предост. использование ф-ции только 1 раз
var singleton *CacheKey

// GetSingleton нужен для использования уже созданной переменной,чтобы неск.раз не создавалась новая!
func GetSingleton() *CacheKey {
	once.Do(func() {
		singleton = NewCacheKey()
	})
	return singleton
}

func (c *CacheKey) Set(cacheKey string, data interface{}, duration time.Duration) error {
	c.cache[cacheKey] = CacheData{
		data:      data,
		expiresAt: time.Now().Add(duration),
	}
	return nil
}

func (c *CacheKey) Get(key string) interface{} {
	cacheEntry, ok := c.cache[key]
	if !ok {
		// знач не найдены по ключу
		return nil
	}

	if time.Now().After(cacheEntry.expiresAt) {
		// Если время кэша истекло, удаляем его
		delete(c.cache, key)
		fmt.Println("cache is expired!")
		return nil
	}

	return cacheEntry.data
}

var RedisClient *redis.Client

func InitRedis() {
	RedisClient = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})
}
