package main

import (
	"fmt"
	"log"
	"os"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("error loading .env file", err)
	}
	JWT := os.Getenv("jwt")
	dbUrl := os.Getenv("dburl")
	fmt.Println(JWT, dbUrl)
}
