package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
)

func main() {
	resp, err := http.Get("http://localhost:8080")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer resp.Body.Close()
	fmt.Println("resp: ", resp.Status)
	body, err := io.ReadAll(resp.Body)

	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(string(body))
	}
}
