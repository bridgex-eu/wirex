package wirex

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestNewEngine checks if a new Engine is correctly initialized.
func TestNewEngine(t *testing.T) {
	engine := New()

	assert.NotNil(t, engine.mux)
	assert.NotNil(t, engine.Validator)
}

// TestEngineFromRequest verifies the Engine retrieval from HTTP request context.
func TestEngineFromRequest(t *testing.T) {
	engine := New()
	ctx := context.WithValue(context.Background(), EngineContextKey, engine)
	reqWithEngine, _ := http.NewRequestWithContext(ctx, "GET", "/", nil)

	// Test successful retrieval from context
	err := engine.FromRequest(reqWithEngine)
	assert.Equal(t, nil, err)

	// Test failure when engine is not in context
	reqWithoutEngine, _ := http.NewRequest("GET", "/", nil)
	err = engine.FromRequest(reqWithoutEngine)
	assert.NotEqual(t, nil, err)
}

// TestEngineRoute tests the route function.
func TestEngineRoute(t *testing.T) {
	engine := New()

	// Define a mock middleware
	middleware := func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// You can perform middleware logic here. For testing, we'll just add a header.
			w.Header().Set("X-Middleware-Test", "true")
			next.ServeHTTP(w, r)
		})
	}

	// Define the test handler
	testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	// Add the route with middleware
	engine.route("/test", []MethodHandler{{method: "GET", handler: testHandler}}, []Middleware{middleware})

	// Make a request to the test route
	req, _ := http.NewRequest("GET", "/test", nil)
	rr := httptest.NewRecorder()
	engine.mux.ServeHTTP(rr, req)

	// Assert that the status code is correct
	assert.Equal(t, http.StatusOK, rr.Code)

	// Assert that the middleware was executed
	assert.Equal(t, "true", rr.Header().Get("X-Middleware-Test"))
}

// TestEngineRegisterRoutes tests the registerRoutes function.
func TestEngineRegisterRoutes(t *testing.T) {
	engine := New()

	// Mock handlers for testing
	handler1 := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})
	handler2 := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	})

	// Mock routes
	engine.routes = []*Route{
		{
			pattern: "/ok",
			handlers: []MethodHandler{
				{method: "GET", handler: handler1},
			},
			middlewares: nil,
		},
		{
			pattern: "/notfound",
			handlers: []MethodHandler{
				{method: "GET", handler: handler2},
			},
			middlewares: nil,
		},
	}

	// Register mock routes
	engine.registerRoutes()

	// Test if routes are correctly registered
	rr1 := httptest.NewRecorder()
	req1, _ := http.NewRequest("GET", "/ok", nil)
	engine.mux.ServeHTTP(rr1, req1)
	assert.Equal(t, http.StatusOK, rr1.Code)

	rr2 := httptest.NewRecorder()
	req2, _ := http.NewRequest("GET", "/notfound", nil)
	engine.mux.ServeHTTP(rr2, req2)
	assert.Equal(t, http.StatusNotFound, rr2.Code)
}

// TestListenAndServe tests the ListenAndServe method of the Engine.
func TestListenAndServe(t *testing.T) {
	engine := New()

	// Mock handler
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("ok"))
	})
	engine.mux.Handle("/test", handler)

	// Create a test server (not using ListenAndServe directly)
	testServer := httptest.NewServer(engine.Handler())
	defer testServer.Close()

	// Send a request to the test server
	resp, err := http.Get(testServer.URL + "/test")
	if err != nil {
		t.Fatalf("Could not make GET request: %v", err)
	}
	defer resp.Body.Close()

	// Check if the status code is as expected
	assert.Equal(t, http.StatusOK, resp.StatusCode)
}
