package server

import (
	"context"
	"fmt"
	"net/http"

	"github.com/chrisvdg/gotiny/business"
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
)

const defaultBackendFile string = "./backend.json"

// New creates a new server instance
func New(c *Config) (*Server, error) {
	return &Server{cfg: c}, nil
}

// Server represents a server instance
type Server struct {
	cfg *Config
}

// AddAPIRoutes adds the API routes with default handlers
func (s *Server) AddAPIRoutes(r *mux.Router, b *business.Logic) error {
	h := &DefaultHandlers{b: b}

	auth := &DefaultAuthorizer{}
	s.AddAPIRoutesAndHandlers(r, h, auth)

	return nil
}

// AddAPIRoutesAndHandlers adds the API routes with provided handlers and authorizers
func (s *Server) AddAPIRoutesAndHandlers(r *mux.Router, handlers Handlers, auth Authorizer) error {
	if r == nil {
		return fmt.Errorf("Router is nil")
	}
	listHandler := http.HandlerFunc(handlers.List)
	createHandler := http.HandlerFunc(handlers.CreateTinyURL)
	updateHandler := http.HandlerFunc(handlers.UpdateTinyURL)
	expandHandler := http.HandlerFunc(handlers.ExpandURL)
	deleteHandler := http.HandlerFunc(handlers.RemoveTinyURL)

	r.HandleFunc("/api", handlers.APISpec).Methods("GET")
	r.Handle("/api/tiny", auth.AuthenticateRead(listHandler)).Methods("GET")
	r.Handle("/api/tiny", auth.AuthenticateCreate(createHandler)).Methods("POST")
	r.HandleFunc("/api/tiny/{id}", handlers.FollowURL).Methods("GET")
	r.Handle("/api/tiny/{id}", auth.AuthenticateWrite(updateHandler)).Methods("POST")
	r.Handle("/api/tiny/{id}", auth.AuthenticateWrite(deleteHandler)).Methods("DELETE")
	r.Handle("/api/tiny/{id}/expand", auth.AuthenticateRead(expandHandler)).Methods("GET")

	return nil
}

// ListenAndServeFileBackedAPI sets the API routes only with a file backend
// and listens for requests and serves them
func (s *Server) ListenAndServeFileBackedAPI() error {
	file := s.cfg.FileBackendPath
	if file == "" {
		file = defaultBackendFile
	}
	b, err := business.NewFileBackedLogic(file, s.cfg.PrettyJSON, 5)
	if err != nil {
		return err
	}
	h, err := NewDefaultHandlers(b)
	if err != nil {
		return err
	}

	s.ListenAndServeAPI(h)

	return nil
}

// ListenAndServeAPI sets the API routes only with provided backend
// and listens for requests and serves them
func (s *Server) ListenAndServeAPI(handlers Handlers) error {
	r := mux.NewRouter()
	auth := NewAuthorizer(s.cfg.ReadAuthToken, s.cfg.WriteAuthToken, s.cfg.AllowPublicCreateGenerated)
	err := s.AddAPIRoutesAndHandlers(r, handlers, auth)
	if err != nil {
		return err
	}
	s.ListenAndServe(r)

	return nil
}

// ListenAndServe listens for requests and serves them
func (s *Server) ListenAndServe(handler http.Handler) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	tlsEnabled := s.cfg.TLS.CertFile != "" && s.cfg.TLS.KeyFile != ""

	if !s.cfg.TLSOnly {
		go listenAndServe(ctx, cancel, s.cfg.ListenAddr, handler)
	}

	if tlsEnabled {
		go listenAndServeTLS(ctx, cancel, s.cfg.TLSListenAddr, s.cfg.TLS, handler)
	}

	<-ctx.Done()
}

// listenAndServe serves a plain http webserver
func listenAndServe(ctx context.Context, cancel func(), addr string, handler http.Handler) {
	defer cancel()
	log.Infof("http server listening on: localhost%s\n", addr)
	log.Print(http.ListenAndServe(addr, handler))
}

// listenAndServeTLS serves a tls webserver
func listenAndServeTLS(ctx context.Context, cancel func(), addr string, tls *TLSConfig, handler http.Handler) {
	defer cancel()
	log.Infof("https server listening on: localhost%s\n", addr)
	log.Print(http.ListenAndServeTLS(addr, tls.CertFile, tls.KeyFile, handler))
}
