package database

import (
	"os"

	"github.com/go-redis/redis/v8"
	"github.com/joho/godotenv"
	"gorm.io/gorm"
)

var (
	DBConn *gorm.DB
	_      = godotenv.Load()
	Rdb    = redis.NewClient(&redis.Options{
		Addr: os.Getenv("REDIS_ADDR"),
	})
)
