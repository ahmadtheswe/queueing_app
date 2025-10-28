package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/ahmadtheswe/queueing_app/routers"
)

func main() {
	fmt.Println("Hello, World!")

	routers.SetupRoutes()

	log.Fatal(http.ListenAndServe(":8080", nil))
}
