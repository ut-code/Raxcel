package main

import (
	"log"

	"github.com/joho/godotenv"
	"github.com/ut-code/Raxcel/server/api"
	"github.com/ut-code/Raxcel/server/db"
)

func main() {

	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}
	db.Migrate()
	e := api.SetupRouter()
	log.Fatal(e.Start(":8080"))
}
