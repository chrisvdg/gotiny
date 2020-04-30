package server

import (
	"fmt"
	"net/http"

	"github.com/chrisvdg/gotiny/business"
	"github.com/gorilla/mux"
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
	RemoveTinyURL(http.ResponseWriter, *http.Request)
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

// List Lists all tiny URL entries
func (h *DefaultHandlers) List(res http.ResponseWriter, req *http.Request) {
	data, err := h.b.List()
	if err != nil {
		writeError(res, err)
		return
	}

	writeJSONResp(res, data)
}

// CreateTinyURL Create a new tiny URL entry
func (h *DefaultHandlers) CreateTinyURL(res http.ResponseWriter, req *http.Request) {
	err := req.ParseForm()
	if err != nil {
		writeError(res, fmt.Errorf("Failed to parse form data: %s", err))
		return
	}
	id := req.Form.Get("id")
	url := req.Form.Get("url")
	data, err := h.b.Create(id, url)
	if err != nil {
		writeErrorWithValidationCheck(res, err)
	}

	writeJSONResp(res, data)
}

// FollowURL Get redirected to full URL
func (h *DefaultHandlers) FollowURL(res http.ResponseWriter, req *http.Request) {
	id := mux.Vars(req)["id"]
	url, err := h.b.GetURL(id)
	if err != nil {
		writeError(res, err)
		return
	}

	http.Redirect(res, req, url, 301)
}

// UpdateTinyURL Update a tiny URL entry
func (h *DefaultHandlers) UpdateTinyURL(res http.ResponseWriter, req *http.Request) {
	id := mux.Vars(req)["id"]
	err := req.ParseForm()
	if err != nil {
		writeError(res, fmt.Errorf("Failed to parse form data: %s", err))
		return
	}
	url := req.Form.Get("url")

	err = h.b.Update(id, url)
	if err != nil {
		writeErrorWithValidationCheck(res, err)
		return
	}

	res.WriteHeader(http.StatusNoContent)
}

// ExpandURL Get info for the tiny URL ID entry
func (h *DefaultHandlers) ExpandURL(res http.ResponseWriter, req *http.Request) {
	id := mux.Vars(req)["id"]

	data, err := h.b.Get(id)
	if err != nil {
		writeError(res, err)
		return
	}

	writeJSONResp(res, data)
}

// RemoveTinyURL Remove a tiny URL entry
func (h *DefaultHandlers) RemoveTinyURL(res http.ResponseWriter, req *http.Request) {
	id := mux.Vars(req)["id"]

	err := h.b.Delete(id)
	if err != nil {
		writeError(res, err)
		return
	}

	res.WriteHeader(http.StatusNoContent)
}

func writeJSONResp(res http.ResponseWriter, data []byte) {
	res.WriteHeader(http.StatusOK)
	res.Header().Set("Content-Type", "application/json")
	res.Write(data)
}

// writeError writes error to the response writer
// Don't forget to return after calling this function in the handler
func writeError(res http.ResponseWriter, err error) {
	if err == business.ErrTinyURLNotFound {
		res.WriteHeader(http.StatusNotFound)
	} else {
		res.WriteHeader(http.StatusInternalServerError)
		res.Write([]byte(err.Error()))
	}
}

// writeError writes error to the response writer, also checks for validation error
// Don't forget to return after calling this function in the handler
func writeErrorWithValidationCheck(res http.ResponseWriter, err error) {
	if business.IsValidationError(err) {
		res.WriteHeader(http.StatusBadRequest)
		res.Write([]byte(err.Error()))
	} else {
		writeError(res, err)
	}
}
