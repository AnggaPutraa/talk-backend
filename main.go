package main

import (
	"github.com/AnggaPutraa/talk-backend/app"
	"github.com/AnggaPutraa/talk-backend/configs"
	db "github.com/AnggaPutraa/talk-backend/db/sqlc"
)

func main() {
	configuration := configs.LoadConfig()
	database := configs.OpenConnection(configuration.DBUrl)
	query := db.New(database)
	app.RunServer(configuration, query)
}
