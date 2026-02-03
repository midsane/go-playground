package main

import (
	"flag"
	"fmt"
	"github.com/midsane/go-playground/02-cli-app/internal/app"
	"os"
)

func main() {
	formal := flag.Bool("formal", false, "use formal greeting")
	flag.Parse()
	args := flag.Args()
	if len(args) < 1 {
		fmt.Println("usage: greet [--formal] <name>")
		os.Exit(1)
	}

	name := args[0]
	fmt.Println(app.GreetToCli(name, *formal))
}
