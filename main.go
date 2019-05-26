package main

import (
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/hirokikondo86/API-server/controller"
)

func main() {
	log.Println("Server started on: http://localhost:8001")

	Router := mux.NewRouter()
	Router.HandleFunc("/", controller.Index)
	Router.HandleFunc("/api/v1/show", controller.ShowAll).Methods("GET")
	Router.HandleFunc("/api/v1/show/{id}", controller.Show).Methods("GET")
	Router.HandleFunc("/api/v1/insert", controller.Insert).Methods("POST")
	Router.HandleFunc("/api/v1/update/{id}", controller.Update).Methods("PUT")
	Router.HandleFunc("/api/v1/delete/{id}", controller.Delete).Methods("DELETE")

	server := http.Server{
		Addr:         ":8001",
		WriteTimeout: time.Second * 15,
		ReadTimeout:  time.Second * 15,
		IdleTimeout:  time.Second * 60,
		Handler:      Router,
	}

	log.Fatal(server.ListenAndServe())
}
