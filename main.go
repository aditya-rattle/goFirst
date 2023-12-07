package main

import (
	"log"
	"net/http"

	"example.com/goFirst/controller"
	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
)

var Db *gorm.DB

func main() {

	db, err := InitDatabase()
	// redis set up
	// docker
	// kubernetes
	if err != nil {
		log.Fatalf("Panicking Database", err)
	}
	defer db.Close()
	redisClient, err := InitRedis()
	if err != nil {
		log.Fatalf("Panicking Redis", err)
	}
	controller := controller.Controller{
		Db:    db,
		Redis: redisClient,
	}
	r := mux.NewRouter()
	apiV1 := r.PathPrefix("/api/v1").Subrouter()
	apiV1.HandleFunc("/add-book", controller.AddBook).Methods("POST")
	apiV1.HandleFunc("/get-books", controller.GetBooks).Methods("GET")
	apiV1.HandleFunc("/get-book/{id}", controller.GetBookById).Methods("GET")
	log.Fatal(http.ListenAndServe(":9000", r))

}
