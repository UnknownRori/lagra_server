package main

import (
	"os"

	echoKoishi "github.com/UnknownRori/lagra_server/src/echo"

	"github.com/charmbracelet/log"
	"github.com/joho/godotenv"
)

func main() {
	log.Info("Loading .env")
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	appPort := os.Getenv("APP_PORT")

	log.Info("Initialize Server")
	server := echoKoishi.CreateServer()

	server.Start(appPort)
}
