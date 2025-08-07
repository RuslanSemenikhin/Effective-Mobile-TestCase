package main

import (
	"fmt"
	"log"

	"github.com/RuslanSemenikhin/Effective-Mobile-TestCase.git/internal/models"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load("../../../.env")
	if err != nil {
		log.Fatal(".env fail does not exists")
	}

	config, err := models.NewDbConfig()
	if err != nil {
		fmt.Println("===>>>")
	}

	fmt.Println(config)
}
