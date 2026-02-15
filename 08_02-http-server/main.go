package main

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"
)

type User struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}


type userStore struct {
	mu    sync.Mutex
	users map[int]User
	next  int
}

func newUserStore() *userStore {
	return &userStore{
		users: make(map[int]User),
		next:  1,
	}
}

func (s *userStore) Create(u User) User {
	s.mu.Lock()
	defer s.mu.Unlock()

	u.ID = s.next
	s.next++
	s.users[u.ID] = u
	return u
}

func (s *userStore) GetAll() []User {
	s.mu.Lock()
	defer s.mu.Unlock()

	out := make([]User, 0, len(s.users))
	for _, u := range s.users {
		out = append(out, u)
	}
	return out
}

func (s *userStore) Get(id int) (User, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()

	u, ok := s.users[id]
	return u, ok
}

func (s *userStore) Update(id int, u User) (User, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, ok := s.users[id]; !ok {
		return User{}, errors.New("not found")
	}
	u.ID = id
	s.users[id] = u
	return u, nil
}

func (s *userStore) Delete(id int) bool {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, ok := s.users[id]; !ok {
		return false
	}
	delete(s.users, id)
	return true
}


func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(v)
}

func readJSON(r *http.Request, dst any) error {
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()
	return dec.Decode(dst)
}

func parseID(path string) (int, error) {
	parts := strings.Split(strings.Trim(path, "/"), "/")
	if len(parts) != 2 {
		return 0, errors.New("invalid path")
	}
	return strconv.Atoi(parts[1])
}

type server struct {
	store *userStore
}

func (s *server) users(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		users := s.store.GetAll()
		writeJSON(w, http.StatusOK, users)

	case http.MethodPost:
		var u User
		if err := readJSON(r, &u); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		if u.Name == "" || u.Email == "" {
			http.Error(w, "name and email required", http.StatusBadRequest)
			return
		}
		created := s.store.Create(u)
		writeJSON(w, http.StatusCreated, created)

	default:
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	}
}

func (s *server) userByID(w http.ResponseWriter, r *http.Request) {
	id, err := parseID(r.URL.Path)
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	switch r.Method {
	case http.MethodGet:
		u, ok := s.store.Get(id)
		if !ok {
			http.NotFound(w, r)
			return
		}
		writeJSON(w, http.StatusOK, u)

	case http.MethodPut:
		var u User
		if err := readJSON(r, &u); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		updated, err := s.store.Update(id, u)
		if err != nil {
			http.NotFound(w, r)
			return
		}
		writeJSON(w, http.StatusOK, updated)

	case http.MethodDelete:
		if !s.store.Delete(id) {
			http.NotFound(w, r)
			return
		}
		w.WriteHeader(http.StatusNoContent)

	default:
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	}
}


func logging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)
		log.Println(r.Method, r.URL.Path, time.Since(start))
	})
}

func main() {
	store := newUserStore()
	srv := &server{store: store}

	mux := http.NewServeMux()
	mux.HandleFunc("/users", srv.users)
	mux.HandleFunc("/users/", srv.userByID)

	handler := logging(mux)

	log.Println("listening on :8080")
	log.Fatal(http.ListenAndServe(":8080", handler))
}
