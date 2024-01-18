package main

import (
	"CRUD_API/database"
	"CRUD_API/pkg/handler"
	"CRUD_API/pkg/repository"
	"CRUD_API/pkg/service"
	"github.com/joho/godotenv"
	"log"
)

/*http://localhost:8080/swagger/index.html#/*/

// @title CRUD API
// @version 1.0
// @description Welcome

// @host localhost:8080
// @BasePath /

// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file:%s", err)
	}
	db, err := database.ConnectToDb()
	if err != nil {
		log.Fatalf("err in connecttodb:%s", err)
	}
	//handler.InitRedis()
	repos := repository.NewRepository(db)
	services := service.NewService(repos)
	handlers := handler.NewHandler(services)

	r := handlers.InitRoutes()

	err = r.Run(":8080")
	if err != nil {
		log.Fatalf("error in router.run %s", err.Error())
	}

}
