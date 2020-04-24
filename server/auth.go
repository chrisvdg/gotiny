package server

import (
	"net/http"
	"strings"

	"github.com/chrisvdg/gotiny/business"
	log "github.com/sirupsen/logrus"
)

const bearer = "bearer"

// Authorizer defines a type that can be used to authorize to the API
type Authorizer interface {
	AuthenticateRead(http.Handler) http.Handler
	AuthenticateWrite(http.Handler) http.Handler
}

// NewAuthorizer creates a new instance of a default authorizer
func NewAuthorizer(b *business.Logic, readToken string, writeToken string) *DefaultAuthorizer {
	return &DefaultAuthorizer{
		readToken:  readToken,
		writeToken: writeToken,
		b:          b,
	}
}

// DefaultAuthorizer is the default authorizing middleware
type DefaultAuthorizer struct {
	readToken  string
	writeToken string
	b          *business.Logic
}

// AuthenticateRead Authenticates for read permissions
func (h *DefaultAuthorizer) AuthenticateRead(next http.Handler) http.Handler {
	return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		if h.readToken != "" {
			token := h.getToken(req)
			if token != h.readToken {
				log.Debugf("Failed to authorize %s", req.RemoteAddr)
				res.WriteHeader(http.StatusUnauthorized)
				return
			}
		}

		next.ServeHTTP(res, req)
	})
}

// AuthenticateWrite Authenticates for write permissions
func (h *DefaultAuthorizer) AuthenticateWrite(next http.Handler) http.Handler {
	return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		if h.writeToken != "" {
			token := h.getToken(req)
			if token != h.writeToken {
				log.Debugf("Failed to authorize %s", req.RemoteAddr)
				res.WriteHeader(http.StatusUnauthorized)
				return
			}
		}

		next.ServeHTTP(res, req)
	})
}

// getToken fetches the Bearer token in the Authorization header
func (h *DefaultAuthorizer) getToken(req *http.Request) string {
	token := ""
	reqToken := req.Header.Get("Authorization")
	if strings.HasPrefix(strings.ToLower(reqToken), bearer) {
		token = strings.TrimSpace(reqToken[len(bearer):])
	}

	return token
}
