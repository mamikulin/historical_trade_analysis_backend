package main

import (
	"archpath/internal/api"
	"log"
)

func main() {
	log.Println("Application start")
	api.StartServer()
	log.Println("Application terminated")

}
