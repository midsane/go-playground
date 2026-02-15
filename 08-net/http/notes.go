package main

import (
	"fmt"
	"log"
	"net/http"
)

/*
	ServeMux ->struct type of the server -> contains fields like r/w mutext, index, tree .
	defaultServeMux -> zero value of construction of ServeMux
	DefaultServeMux -> pointer of defaultServeMux


	Handler ->
	interface having the function -> ServeHTTP(w http.ResponseWriter, r *http.Request)
	something like this:
	type Handler interface {
		ServeHTTP(w http.ResponseWriter, r *http.Request)
	}

	HandlerFunc ->
	adapter function -> that wraps a function with signature func(w http.ResponseWriter, r *http.Request)
	and associated with it ServeHTTP(w http.ResponseWriter, r *http.Request) fuction to make it satisfy handler
	interface.
	something like this:
	func HandlerFunc func(w http.ResponseWriter, r *http.Request)
	func(f HandlerFunc) ServeHTTP(w http.ResponseWriter, r *http.Request){
		f(w,r)
	}

	HandleFunc -> accepts two params pattern and Handler , a simple function with just params as (w http.ResponseWrite, r *http.Request)
	will need a wrapping around HanlderFunc

	Handle -> accepts two params patter and any func with signature func(w http.ResponseWriter, r *http.Request), since internally it wraps
	that function was HandlerFunc an then pass it back HandleFunc.
	something like this:
	func Handle(p Pattern, f func(w http.ResponseWriter, r *http.Request)){
		funcSatisfyingHandler := HandlerFunc(f)
		HandleFunc(p, funcSatisfyingHandler)
	}

	http.ListenAndServe() -> asks for two params -> port and a Handler of which ServeHTTP will be invoked to find
	matching route.
	so in the handler we can pass any type whose method sets satisfies all methods of interface Handler and add some logging logic as well,
	otherwise the default Handler is assigned in the Handler of the New instance of ServeMux created.

	ex of loggingMiddleware:
	func LoggingMiddleWare(next ServeMux) http.Handler{
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request){
			log.Printf("%s %s %s", r.Method, r.RequestURI, r.RemoteAddr)
			next.ServeHTTP(w,r)
		})
	}
	now when ServeHTTP function associated to the func returned by LoggingMiddleware is called -> it will inturn called that exact function,
	causing to log and then proced to ServeHTTP of handler of ServerInstance.
	--> can be used like this
		http.ListenAndServe(":8080", loggingMiddleware(http.NewServeMux()))

	ServeHTTP -> it looks for a suitable handler by matching with request.url.path, if it does find out ->
	invokes ServeHTTP of that handler otherwise writes to response -> something like not found.

	A new go routine is spawned for every new connection.

*/

/*
concept of method set in golang:
method sets basically determines which methods are associated with a given type:
a method can have two types of receiver type ->
1. value receiver type func (t T) MethodName()
2. ptr receiver type func(t *T) MethodName()

for a value of a given type -> it can invoke methods of type1 only, while the for ptr of
a given type -> it can invoke both type1 and typ2 methods.



*/

func LoggingMiddleWare(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s %s %s", r.Method, r.RequestURI, r.RemoteAddr)
		next.ServeHTTP(w, r)
	})
}

/*
notice it's gives a warning if i try to pass next http.ServeMux here, which will
disappear if i pass next *http.ServeMux or http.Handler
that's because ServeHTTP method receiver type is *ServeMux.
so value ServeMux dont implement Handler, like *ServeMux implements.
(this is done because ServeMux has fields like mutext which should be shared via ptr, not copies)
and http.Handler ofc implments Handler.
*/

func Test() {
	srv := http.ServeMux{}
	srv.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "hello word")
	})

	fmt.Println("server is running on port 8080")
	http.ListenAndServe(":8080", LoggingMiddleWare(&srv))
}
