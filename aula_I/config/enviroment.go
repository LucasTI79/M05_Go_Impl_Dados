package config

import (
	"log"

	"github.com/joho/godotenv"
)

func initEnviroment() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("error ao carregar o arquivo .env")
	}
}
