package main

import "github.com/midsane/go-playground/10-auth/internal/server"

const PORT = ":8080"

/*
lets start a basic net/http server and do jwt authentication
*/

func main() {
	server.Start(PORT)
}
