package wirex

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
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

func TestRoute(t *testing.T) {
	engine := New()

	setRoutes(&engine.RoutesGroup)

	testServer := httptest.NewServer(engine.Handler())
	defer testServer.Close()

	testRoutes(t, testServer)
}

func TestGroupRoute(t *testing.T) {
	engine := New()

	g := NewRoutesGroup()
	setRoutes(g)

	engine.Group("/", g)

	testServer := httptest.NewServer(engine.Handler())
	defer testServer.Close()

	testRoutes(t, testServer)
}

func setRoutes(g *RoutesGroup) {
	g.Route("/test/").Get(errHandler)
	g.Route("/test/route/").Get(okHandler)
	g.Route("/check").Any(okHandler).Post(errHandler)
}

func testRoutes(t *testing.T, server *httptest.Server) {
	testRoute(t, server, "/test/route/wildcard", http.StatusOK)
	testRoute(t, server, "/test/invalid", http.StatusInternalServerError)
	testRoute(t, server, "/notfound/route", http.StatusNotFound)
	testRoute(t, server, "/check", http.StatusOK)
	testPostRoute(t, server, "/check", http.StatusInternalServerError)
}

type status struct {
	status int
}

func (s status) WriteResponse(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(s.status)
}

func errHandler(r *http.Request) Writer {
	return status{http.StatusInternalServerError}
}

func okHandler(r *http.Request) Writer {
	return status{http.StatusOK}
}

func testRoute(t *testing.T, server *httptest.Server, path string, expectedStatus int) {
	resp, err := http.Get(server.URL + path)
	if err != nil {
		t.Fatalf("Could not make GET request to %s: %v", path, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != expectedStatus {
		t.Errorf("Expected status code %d for path %s, got %d", expectedStatus, path, resp.StatusCode)
	}
}

func testPostRoute(t *testing.T, server *httptest.Server, path string, expectedStatus int) {
	resp, err := http.Post(server.URL+path, MIMEApplicationJSON, strings.NewReader("{}"))
	if err != nil {
		t.Fatalf("Could not make Post request to %s: %v", path, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != expectedStatus {
		t.Errorf("Expected status code %d for path %s, got %d", expectedStatus, path, resp.StatusCode)
	}
}