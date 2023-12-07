package controller

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"example.com/goFirst/models"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	"github.com/redis/go-redis/v9"
)

type Controller struct {
	Db    *gorm.DB
	Redis *redis.Client
}

func (ctr *Controller) AddBook(w http.ResponseWriter, r *http.Request) {
	book := &models.Books{
		Id:     uuid.New().String(),
		Name:   "DDIA",
		Author: "O'Reilly",
		Price:  2000,
	}
	ctr.Db.Create(&book)
	fmt.Println("Successfully inserted a new book in db")
	bookByte, err := json.Marshal(book)
	if err != nil {
		fmt.Println("Error in marshaling object", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Write(bookByte)
}

func (ctr *Controller) GetBooks(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	var result []models.Books

	redisVal := ctr.Redis.Get(ctx, "REDIS:BOOKS")
	if redisVal == nil || redisVal.Val() == "" {
		fmt.Println("Key is Missing")
	} else {
		fmt.Println("Key found")
		err := json.Unmarshal([]byte(redisVal.Val()), &result)
		if err != nil {
			fmt.Println("Error in unmarshaling value")
		}
		w.Write([]byte(redisVal.Val()))
	}

	ctr.Db.Raw("Select * from books").Scan(&result)
	fmt.Println(result)

	resultB, err := json.Marshal(result)
	if err != nil {
		fmt.Println("Error in getting books", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	redisSet := ctr.Redis.Set(ctx, "REDIS:BOOKS", resultB, 10*time.Minute)
	if redisSet == nil || redisSet.Val() == "" {
		fmt.Println("Unable to set key to redis")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Write(resultB)
}

func (ctr *Controller) GetBookById(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id := params["id"]
	if id == "" {
		fmt.Println("Id is missing in requent")
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	var result models.Books
	ctx := context.Background()
	redisKey := fmt.Sprintf("REDIS:%s", id)
	redisVal := ctr.Redis.Get(ctx, redisKey)
	if redisVal == nil || redisVal.Val() == "" {
		fmt.Println("Key not found in redis for id: ", id)
	} else {
		fmt.Println("Key found in redis")
		if err := json.Unmarshal([]byte(redisVal.Val()), &result); err != nil {
			fmt.Println("Error in unmarshaling redis value for id: ", id)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.Write([]byte(redisVal.Val()))
	}

	ctr.Db.Raw(fmt.Sprintf("Select * from books where id = '%s'", id)).Scan(&result)
	resultB, err := json.Marshal(result)
	if err != nil {
		fmt.Println("Error in getting books", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	redisSet := ctr.Redis.Set(ctx, redisKey, resultB, 10*time.Minute)
	if redisSet == nil || redisSet.Val() == "" {
		fmt.Println("Error in setting redis key against id: ", id)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Write(resultB)
}
