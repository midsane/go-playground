package main

import (
	"fmt"
	"github.com/midsane/go-playground/02-cli-app/internal/app"
)

func main() {
	/*
		print greeting message depending, get a person's name
	*/
	name := "satmak"
	msg := fmt.Sprintf(Greet(), name)
	fmt.Println(msg)

}
