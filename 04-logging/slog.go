package main

import (
	"log"
	"log/slog"
	"os"

	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()

	//standard logging
	log.Println("user created")
    
	slog.Info("hello world", "port", os.Getenv("PORT"))
}
