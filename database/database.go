package database

import (
	"fmt"
	"os"

	"github.com/go-redis/redis/v8"
	"github.com/joho/godotenv"
	"gorm.io/gorm"
)

var (
	DBConn *gorm.DB

	_ = godotenv.Load()

	//Redis connection
	Rdb = redis.NewClient(&redis.Options{
		Addr: os.Getenv("REDIS_ADDR"),
	})

	//RabbitMq url for amqp dial
	RabbitMqUrl = fmt.Sprintf("amqp://%s:%s@localhost:%s/", os.Getenv("RABBITMQ_USER"), os.Getenv("RABBITMQ_PASSWORD"), os.Getenv("RABBITMQ_PORT"))
)
