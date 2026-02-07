package config

import (
	"fmt"
	"log"
	"os"
	"github.com/joho/godotenv"
)

type Config struct {
	PORT int `env:"PORT" validate:"required,min=1,max=65535"`
}

func Load() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("error in loading in env vars!")
	}

	port, ok := os.LookupEnv("PORT")
	if !ok {
		log.Fatal("port variable not defined")
	}
	fmt.Println(port)
}
