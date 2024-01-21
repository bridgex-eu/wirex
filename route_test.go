package wirex

import (
	"net/http"
	"testing"
)

var tests = []string{
	http.MethodOptions,
	http.MethodGet,
	http.MethodPost,
	http.MethodPut,
	http.MethodDelete,
	http.MethodConnect,
	http.MethodHead,
	http.MethodPatch,
	http.MethodTrace,
}

func TestMethods(t *testing.T) {
	for _, test := range tests {
		t.Run(test, func(t *testing.T) {
			r := &Route{}

			mockHandler := func(*http.Request) Writer { return nil }

			switch test {
			case http.MethodOptions:
				r.Options(mockHandler)
			case http.MethodGet:
				r.Get(mockHandler)
			case http.MethodPost:
				r.Post(mockHandler)
			case http.MethodPut:
				r.Put(mockHandler)
			case http.MethodDelete:
				r.Delete(mockHandler)
			case http.MethodConnect:
				r.Connect(mockHandler)
			case http.MethodHead:
				r.Head(mockHandler)
			case http.MethodPatch:
				r.Patch(mockHandler)
			case http.MethodTrace:
				r.Trace(mockHandler)
			}

			if len(r.handlers) != 1 {
				t.Errorf("expected 1 handler for method %s, got %d", test, len(r.handlers))
			}

			if r.handlers[0].method != test {
				t.Errorf("expected %s method, got %s", test, r.handlers[0].method)
			}

		})
	}
}

func TesnAny(t *testing.T) {
	r := &Route{}

	mockHandler := func(*http.Request) Writer { return nil }

	r.Any(mockHandler)

	if len(r.handlers) != len(tests) {
		t.Errorf("expected %d handler for method Any, got %d", len(tests), len(r.handlers))
	}
}