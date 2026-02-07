package main

import (
	"fmt"

	"github.com/midsane/go-playground/03-config-management/internal/config"
)

func main() {
	config.Load()
	fmt.Println("all env variables are loaded!")
}
