package server

import (
	"net/http"

	"github.com/chrisvdg/gotiny/business"
)

const apiSpecFile = "specs/api.yaml"

// Handlers represents the handlers needed for the API
type Handlers interface {
	APISpec(http.ResponseWriter, *http.Request)
	List(http.ResponseWriter, *http.Request)
	CreateTinyURL(http.ResponseWriter, *http.Request)
	FollowURL(http.ResponseWriter, *http.Request)
	UpdateTinyURL(http.ResponseWriter, *http.Request)
	ExpandURL(http.ResponseWriter, *http.Request)
}

// DefaultHandlers implements Handlers with the default implementation
type DefaultHandlers struct {
	b *business.Logic
}

// APISpec Shows API spec
func (h *DefaultHandlers) APISpec(res http.ResponseWriter, req *http.Request) {
	res.Header().Add("Content-Type", "application/x-yaml")
	http.ServeFile(res, req, apiSpecFile)
}

// List Lists all shorthand IDs
func (h *DefaultHandlers) List(res http.ResponseWriter, req *http.Request) {

}

// CreateTinyURL Create a new URL shorthand
func (h *DefaultHandlers) CreateTinyURL(res http.ResponseWriter, req *http.Request) {
}

// FollowURL Get redirected to long URL
func (h *DefaultHandlers) FollowURL(res http.ResponseWriter, req *http.Request) {

}

// UpdateTinyURL Update URL shorthand
func (h *DefaultHandlers) UpdateTinyURL(res http.ResponseWriter, req *http.Request) {
}

// ExpandURL Get info of the shorthand ID
func (h *DefaultHandlers) ExpandURL(res http.ResponseWriter, req *http.Request) {
}
