package server

import (
	"net/http"
	"strings"

	log "github.com/sirupsen/logrus"
)

const bearer = "bearer"

// Authorizer defines a type that can be used to authorize to the API
type Authorizer interface {
	AuthenticateRead(http.Handler) http.Handler
	AuthenticateWrite(http.Handler) http.Handler
	AuthenticateCreate(http.Handler) http.Handler
}

// NewAuthorizer creates a new instance of a default authorizer
func NewAuthorizer(readToken string, writeToken string, allowPublicCreateGenerated bool) *DefaultAuthorizer {
	return &DefaultAuthorizer{
		readToken:                  readToken,
		writeToken:                 writeToken,
		allowPublicCreateGenerated: allowPublicCreateGenerated,
	}
}

// DefaultAuthorizer is the default authorizing middleware
type DefaultAuthorizer struct {
	readToken                  string
	writeToken                 string
	allowPublicCreateGenerated bool
}

// AuthenticateRead Authenticates for read permissions
func (h *DefaultAuthorizer) AuthenticateRead(next http.Handler) http.Handler {
	return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		if h.readToken != "" {
			token := h.getToken(req)
			if token != h.readToken {
				log.Debugf("Failed to authorize read %s", req.RemoteAddr)
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
				log.Debugf("Failed to authorize write %s", req.RemoteAddr)
				res.WriteHeader(http.StatusUnauthorized)
				return
			}
		}

		next.ServeHTTP(res, req)
	})
}

// AuthenticateCreate Authenticates for creating a new entry
// where it can be possible to only authenticate when creating custom tiny URLs
func (h *DefaultAuthorizer) AuthenticateCreate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		if h.writeToken != "" {
			if !h.allowPublicCreateGenerated {
				h.AuthenticateWrite(next).ServeHTTP(res, req)
				return
			}

			// If id form field has value check the token
			token := h.getToken(req)
			err := req.ParseForm()
			if err != nil {
				log.Error("Failed to parse form data: %s", err)
				res.WriteHeader(http.StatusInternalServerError)
				return
			}
			reqID := req.Form.Get("id")
			if reqID != "" {
				if token != h.writeToken {
					log.Debugf("Failed to authorize create %s", req.RemoteAddr)
					res.WriteHeader(http.StatusUnauthorized)
					return
				}
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
