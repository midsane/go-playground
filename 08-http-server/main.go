// most minimal http server
package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
)

//what function can handle http request -> any function that implements Handler interface
/*
type handler interface {
	ServeHTTP(ResponseWriter, *Request)
}
*/

func helloHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	fmt.Fprintf(w, "{\"message\": \"Hello, World!\"}")
}

func main() {
	//now lets register a handler function to :8080 port to always listen to
	//that handler function can be made by wrapping a function with HandlerFunction

	/*
		we can only register function that satisfies Handler interface to a route.
		now HandlerFunc in named type of this function signature:
			func(w http.ResponseWriter, r *http.Request)
		and there is also a method associated to it i.e
		ServeHTTP(w http.ResponseWriter, r *http.Request)

		so this makes the wrapped function satisy Handler interface
	*/
	http.HandleFunc("/", helloHandler)
	/*
		internally, this http.HandleFunc("/", helloHandler)
		is converted to var h http.Handler = http.HandlerFunc(helloHandler)
		HandlerFunc adapter makes sure helloHandler satisfies Handler interface
	*/

	fmt.Println("server is running on port 8080")
	//below make it listent on port 8080
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal("err occurred !,", err)
		os.Exit(1)
	}

}
