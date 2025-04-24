package main

import (
	"log"
	"logistic-app/internal/adapters/db"
	"logistic-app/internal/adapters/http"
	"logistic-app/internal/app/service"
)

func main() {
	repo, err := db.NewPostgresDB()
	if err != nil {
		log.Fatal("could not connect to postgres: ", err)
	}
	defer repo.Close()

	logSer := service.NewLogisticService(repo)
	server := http.NewServer(logSer)

	server.Run()
}
