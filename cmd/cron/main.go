package main

import (
	"log"
	"logistic-app/internal/adapters/cron"
	"logistic-app/internal/adapters/db"
	"logistic-app/internal/app/service"
)

func main() {
	repo, err := db.NewPostgresDB()
	if err != nil {
		log.Fatal("could not connect to postgres: ", err)
	}
	defer repo.Close()

	logSer := service.NewLogisticService(repo)
	scheduler := cron.NewScheduler(logSer)

	scheduler.Run()
}
