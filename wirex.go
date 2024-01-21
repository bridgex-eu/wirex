package wirex

import (
	"errors"
	"log/slog"
	"net"
	"net/http"

	"github.com/go-playground/validator/v10"
)

const EngineContextKey = "wirex-engine"

type Engine struct {
	RoutesGroup
	mux       *http.ServeMux
	Validator *validator.Validate // Used for validating request data in request extractors like from.Json.
	Debug     bool                // A flag indicating if the Engine is in debug mode.

	routesRegistered bool
}

// New creates and returns a new instance of Engine.
func New() *Engine {
	engine := &Engine{
		mux:              http.NewServeMux(),
		Debug:            false,
		Validator:        validator.New(),
		routesRegistered: false,
	}

	engine.With(EngineContextKey, engine)

	return engine
}

// FromRequest extracts the Engine instance from the HTTP request's context.
//
// This method attempts to retrieve the Engine instance stored in the context of the provided HTTP request.
// It uses a predefined EngineContextKey to access the context value. If the Engine is not found or the
// type assertion fails, the method returns an HTTPError indicating an internal server error.
// Otherwise, the method updates the receiver (e) to the Engine instance from the request's context and returns nil.
//
// Usage Example:
//
//		err := e.FromRequest(request)
//		if err != nil {
//	    	// handle error
//		}
func (e *Engine) FromRequest(r *http.Request) HTTPError {
	engine, ok := r.Context().Value(EngineContextKey).(*Engine)
	if !ok {
		return Error(http.StatusInternalServerError, errors.New("WireX engine not found in the request context"))
	}

	e = engine
	return nil
}

func (e *Engine) route(pattern string, handlers []MethodHandler, middlewares []Middleware) {
	for _, handler := range handlers {
		e.mux.Handle(handler.method+" "+pattern, applyMiddlewares(handler.handler, middlewares...))
	}
}

func (e *Engine) registerRoutes() {
	for _, route := range e.routes {
		e.route(route.pattern, route.handlers, route.middlewares)
	}
}

// ServeHTTP implements the http.Handler interface for the Engine.
//
// This method logs the HTTP request method and URL path using slog and then delegates the handling
// of the request to the Engine's internal multiplexer (mux). It allows the Engine to be used directly
// as an http.Handler, making it compatible with standard Go HTTP server functions and tools.
//
// Usage Example:
//
//	http.HandleFunc("/", e.ServeHTTP)
//	http.ListenAndServe(":8080", nil)
func (e *Engine) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	slog.Info("Get request:", r.Method, r.URL.Path)
	e.mux.ServeHTTP(w, r)
}

// Handler returns the http.Handler for the Engine.
//
// This method ensures that all routes are registered before returning the Engine as an http.Handler.
// This is useful for integrating the Engine with other HTTP servers or middleware that require an http.Handler.
//
// Usage Example:
//
//	handler := e.Handler()
//	http.Handle("/", handler)
func (e *Engine) Handler() http.Handler {
	if !e.routesRegistered {
		e.registerRoutes()
	}

	return e
}

// ListenAndServe starts an HTTP server with the specified address using the Engine's handler.
//
// This method logs the server's listening state on the specified address and starts an HTTP server.
// It uses the http.ListenAndServe function along with the Engine's handler. Any errors occurring during
// the server's operation are logged upon exiting the function.
//
// Returns an error if the server fails to start or encounters issues during runtime.
//
// Parameters:
// - addr string: The address for the server to listen and serve.
//
// Usage Example:
//
//	if err := e.ListenAndServe(":8080"); err != nil {
//	    log.Fatal(err)
//	}
func (e *Engine) ListenAndServe(addr string) (err error) {
	slog.Info("Listening and serving HTTP", "addr", addr)
	defer func() { slog.Error(err.Error()) }()

	err = http.ListenAndServe(addr, e.Handler())
	return
}

// ListenAndServeTLS starts an HTTPS server with the specified address, certificate file, and key file.
//
// This method logs the server's listening state and captures any error occurring during operation.
// The server uses http.ListenAndServeTLS with the provided TLS credentials and the Engine's handler.
// In case of an error, it is logged at function exit.
//
// Parameters:
// - addr string: The address for the server to listen and serve.
// - certFile string: The path to the SSL certificate file.
// - keyFile string: The path to the SSL key file.
//
// Usage Example:
//
//	if err := engine.ListenAndServeTLS(":443", "cert.pem", "key.pem"); err != nil {
//	    log.Fatal(err)
//	}
func (engine *Engine) ListenAndServeTLS(addr, certFile, keyFile string) (err error) {
	slog.Info("Listening and serving HTTPS", "addr", addr)
	defer func() { slog.Error(err.Error()) }()

	err = http.ListenAndServeTLS(addr, certFile, keyFile, engine.Handler())
	return
}

// Serve initiates an HTTP server with the specified net.Listener.
//
// This method logs the address bound to the listener and starts an HTTP server using http.Serve and the Engine's handler.
// Errors during the server's operation are logged upon function exit.
//
// Parameters:
//
// - listener net.Listener: The network listener for the server.
//
// Usage Example:
//
//	listener, _ := net.Listen("tcp", ":8080")
//
//	if err := engine.Serve(listener); err != nil {
//	    log.Fatal(err)
//	}
func (engine *Engine) Serve(listener net.Listener) (err error) {
	slog.Info("Listening and serving HTTP on listener what's bind with address", "addr", listener.Addr())
	defer func() { slog.Error(err.Error()) }()

	err = http.Serve(listener, engine.Handler())
	return
}
