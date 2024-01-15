package wirex

import (
	"log/slog"
	"net/http"
)

// In WireX, each object returned by a handler function must adhere to the Writer interface.
// Each Writer contains WriteResponse method, which is automatically invoked once a handler returns an object. This architecture enables a highly flexible and tailored approach to crafting and delivering HTTP responses. Each implementation of Writer can uniquely shape how responses are constructed and transmitted to clients.
//
// WireX also offers a suite of pre-built writers, including write.String, write.Json, write.Blob, and write.Redirect, catering to a variety of response formats.
type Writer interface {
	WriteResponse(w http.ResponseWriter, r *http.Request)
}

// The FromRequest module provides a way to extract and process information from HTTP requests.
//
// Struct including this module must implement the `from_request` method, which is responsible
// for parsing the HTTP request object. It should return an HTTP error if any issues arise during
// the extraction process. This module is particularly useful for controllers or middleware where
// structured and reusable data extraction from requests is required.
type FromRequest interface {
	FromRequest(r *http.Request) HTTPError
}

type ResponseHeaderWriter interface {
	Header() http.Header
}

type HeaderWriter interface {
	WriteHeader(w ResponseHeaderWriter, r *http.Request) HTTPError
}

type HandlerFunc func(*http.Request) Writer

type MethodHandler struct {
	method  string
	handler http.Handler
}

func Handler(h HandlerFunc) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		wr := h(r)

		if err, ok := wr.(error); ok {
			slog.Error("handler error", "error", err)
		}

		wr.WriteResponse(w, r)
	})
}
