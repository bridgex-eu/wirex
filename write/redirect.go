package write

import (
	"net/http"

	"github.com/bridgex-eu/wirex"
)

// RedirectData represents the URL needed for an HTTP redirect.
type RedirectData struct {
	Url string
}

var _ wirex.Writer = &RedirectData{}

// WriteResponse writes the HTTP redirect response.
func (rd *RedirectData) WriteResponse(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Location", rd.Url)
	w.WriteHeader(http.StatusFound) // or use 302 directly
}

// Redirect with the provided URL and 302(Found) status code.
func Redirect(url string, other ...wirex.HeaderWriter) wirex.Writer {
	return WithHeader(&RedirectData{Url: url}, other...)
}
