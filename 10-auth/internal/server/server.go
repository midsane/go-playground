package server

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"time"
)

type Server struct {
	addr string
	mux  *http.ServeMux
}

func NewServer(addr string) *Server {
	return &Server{
		addr: addr,
		mux:  &http.ServeMux{},
	}
}

func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		oldTime := time.Now()
		defer func() {
			fmt.Printf("%s, %s ", r.Method, time.Since(oldTime))
		}()

		next.ServeHTTP(w, r)
	})

}

func RecoveryMiddleWare(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				log.Println("panic:", err)
			}
		}()
		next.ServeHTTP(w, r)
	})
}

func (s *Server) routes() {
	s.mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "welcome to mid auth")
	})

	s.mux.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "welcome to mid auth login route")
	})

	protectedProfile := Chain(
		http.HandlerFunc(ProfileHandler),
		JWTMiddleware,
	)

	s.mux.Handle("/profile", protectedProfile)
}

func writeJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

func parseJSON(r *http.Request, dst interface{}) error {
	if r.Header.Get("Content-Type") != "application/json" {
		return errors.New("content type must be application/json")
	}
	return json.NewDecoder(r.Body).Decode(dst)
}

func Chain(h http.Handler, middlewares ...func(http.Handler) http.Handler) http.Handler{
	for i := len(middlewares)-1; i>=0; i--{
		h = middlewares[i](h)
	}
	return h
}

func Start(addr string) {
	srv := NewServer(addr)
	srv.routes()

	fmt.Println("server is running on port ", addr)
	finalHandler := Chain(
		srv.mux,
		RecoveryMiddleWare,
		LoggingMiddleware,
	)
	log.Fatal(http.ListenAndServe(addr, finalHandler))
}
