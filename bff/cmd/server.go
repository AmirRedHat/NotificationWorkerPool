package main

import (
	"log"
	"net/http"
	"task01/bff/handler"
)

func main() {

	notificationHandler := handler.NewNotificationHandler()

	mux := http.NewServeMux()

	notificationHandler.Mount(mux)

	log.Println("Listening on :8080")
	log.Fatal(http.ListenAndServe(":8080", mux))
}
