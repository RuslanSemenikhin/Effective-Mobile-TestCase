package main

import (
	"fmt"
	"log"

	"github.com/RuslanSemenikhin/Effective-Mobile-TestCase.git/grpc"
	"github.com/RuslanSemenikhin/Effective-Mobile-TestCase.git/internal/models"
	"github.com/joho/godotenv"
)

func main() {
	// err := godotenv.Load("../../../.env")
	err := godotenv.Load("/app/.env")
	if err != nil {
		log.Fatal(".env fail does not exists")
	}

	config, err := models.NewConfig()
	if err != nil {
		log.Fatal(err)
	}

	grpc.Start(config.Port)

	fmt.Println(config)
}
