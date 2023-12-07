package main

import (
	"context"
	"fmt"

	"github.com/jinzhu/gorm"

	_ "github.com/lib/pq"
	redis "github.com/redis/go-redis/v9"
)

func InitDatabase() (*gorm.DB, error) {

	var db *gorm.DB

	var err error

	dbURI := fmt.Sprintf("host=localhost user=root dbname=postgres sslmode=disable")

	db, err = gorm.Open("postgres", dbURI)

	if err != nil {
		fmt.Println("Error in establishing database connection", err)
		return nil, err
	}

	err = db.DB().Ping()
	if err != nil {
		fmt.Println("Error in establishing database connection", err)
		return nil, err
	}

	Db = db

	fmt.Println("Database connection established")

	return db, nil
}

func InitRedis() (*redis.Client, error) {
	redisClient := redis.NewClient(
		&redis.Options{
			Addr: "localhost:6379",
			DB: 0,
			Password: "",
		},
	);
	if err := redisClient.Ping(context.Background()).Err(); err != nil {
		return nil, err
	}
	fmt.Println("Redis connectin established")
	return redisClient, nil
}