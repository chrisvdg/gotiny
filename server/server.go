package server

import (
	"fmt"
	"net/http"

	"github.com/chrisvdg/gotiny/business"
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
)

// New creates a new server instance
func New(c *Config) (*Server, error) {
	return &Server{cfg: c}, nil
}

// Server represents a server instance
type Server struct {
	cfg *Config
}

// AddAPIRoutes adds the API routes with default handlers
func (s *Server) AddAPIRoutes(r *mux.Router) error {
	b, err := business.NewLogic()
	if err != nil {
		return err
	}
	h := &DefaultHandlers{b: b}

	auth := &DefaultAuthorizer{}
	s.AddAPIRoutesAndHandlers(r, h, auth)

	return nil
}

// AddAPIRoutesAndHandlers adds the API routes with provided handlers and authorizers
func (s *Server) AddAPIRoutesAndHandlers(r *mux.Router, handler Handlers, auth Authorizer) {
	if r == nil {
		r = mux.NewRouter()
	}
	listHandler := http.HandlerFunc(handler.List)
	createHandler := http.HandlerFunc(handler.CreateTinyURL)
	updateHandler := http.HandlerFunc(handler.UpdateTinyURL)
	expandHandler := http.HandlerFunc(handler.ExpandURL)
	deleteHandler := http.HandlerFunc(handler.RemoveTinyURL)

	r.HandleFunc("/api", handler.APISpec).Methods("GET")
	r.Handle("/api/tiny", auth.AuthenticateRead(listHandler)).Methods("GET")
	r.Handle("/api/tiny", auth.AuthenticateWrite(createHandler)).Methods("POST")
	r.HandleFunc("/api/tiny/{id}", handler.FollowURL).Methods("GET")
	r.Handle("/api/tiny/{id}", auth.AuthenticateWrite(updateHandler)).Methods("POST")
	r.Handle("/api/tiny/{id}", auth.AuthenticateWrite(deleteHandler)).Methods("DELETE")
	r.Handle("/api/tiny/{id}/expand", auth.AuthenticateRead(expandHandler)).Methods("GET")
}

// ListenAndServe listens for requests and serves them
func (s *Server) ListenAndServe(addr string, handler http.Handler) error {
	return nil
}

// listenAndServe serves a plain http webserver
func listenAndServe(cancel func(), addr string, dir string, handler http.Handler) {
	defer cancel()
	fmt.Printf("Now serving plain http on: localhost%s\n", addr)
	log.Fatal(http.ListenAndServe(addr, handler))
}

// listenAndServeTLS serves a tls webserver
func listenAndServeTLS(cancel func(), addr string, dir string, cert string, key string, handler http.Handler) {
	defer cancel()
	fmt.Printf("Now serving tls on: localhost%s\n", addr)
	log.Fatal(http.ListenAndServeTLS(addr, cert, key, handler))
}
