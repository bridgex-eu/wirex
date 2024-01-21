package wirex

import "net/http"

type Route struct {
	pattern     string
	handlers    []MethodHandler
	middlewares []Middleware
}

func (r *Route) handler(method string, handler HandlerFunc) *Route {
	for i, h := range r.handlers {
		if h.method == method {
			r.handlers[i].handler = Handler(handler)
			return r
		}
	}

	r.handlers = append(r.handlers, MethodHandler{
		method:  method,
		handler: Handler(handler),
	})
	return r
}

// Options adds an OPTIONS method handler to the Route.
func (r *Route) Options(h HandlerFunc) *Route {
	return r.handler(http.MethodOptions, h)
}

// Get adds a GET method handler to the Route.
func (r *Route) Get(h HandlerFunc) *Route {
	return r.handler(http.MethodGet, h)
}

// Post adds a POST method handler to the Route.
func (r *Route) Post(h HandlerFunc) *Route {
	return r.handler(http.MethodPost, h)
}

// Put adds a PUT method handler to the Route.
func (r *Route) Put(h HandlerFunc) *Route {
	return r.handler(http.MethodPut, h)
}

// Delete adds a DELETE method handler to the Route.
func (r *Route) Delete(h HandlerFunc) *Route {
	return r.handler(http.MethodDelete, h)
}

// Connect adds a CONNECT method handler to the Route.
func (r *Route) Connect(h HandlerFunc) *Route {
	return r.handler(http.MethodConnect, h)
}

// Head adds a HEAD method handler to the Route.
func (r *Route) Head(h HandlerFunc) *Route {
	return r.handler(http.MethodHead, h)
}

// Patch adds a PATCH method handler to the Route.
func (r *Route) Patch(h HandlerFunc) *Route {
	return r.handler(http.MethodPatch, h)
}

// Trace adds a TRACE method handler to the Route.
func (r *Route) Trace(h HandlerFunc) *Route {
	return r.handler(http.MethodTrace, h)
}

// Any handle all requests to the Route.
func (r *Route) Any(h HandlerFunc) *Route {
	return r.Options(h).Get(h).Post(h).Put(h).Delete(h).Connect(h).Head(h).Patch(h).Trace(h)
}