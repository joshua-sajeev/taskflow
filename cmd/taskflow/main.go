package main

import (
	"taskflow/internal/db"
	repositories "taskflow/internal/domain/repository"
	"taskflow/internal/handlers"
	"taskflow/internal/routes"

	"github.com/sirupsen/logrus"
)

func main() {
	database, err := db.InitDB()
	if err != nil {
		logrus.Warn("Couldn'r Initialise Database")
	}

	jobRepo := repositories.NewGormJobRepository(database)
	jobHandler := handlers.NewJobHandler(jobRepo)
	router := routes.SetupRouter()
	routes.RegisterJobRoutes(router, jobHandler)
	router.Run(":8080")
}
