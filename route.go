package wirex

import "net/http"

type Route struct {
	pattern     string
	handlers    []MethodHandler
	middlewares []Middleware
}

// Options adds an OPTIONS method handler to the Route.
func (r *Route) Options(h HandlerFunc) *Route {
	r.handlers = append(r.handlers, MethodHandler{
		method:  http.MethodOptions,
		handler: Handler(h),
	})

	return r
}

// Get adds a GET method handler to the Route.
func (r *Route) Get(h HandlerFunc) *Route {
	r.handlers = append(r.handlers, MethodHandler{
		method:  http.MethodGet,
		handler: Handler(h),
	})

	return r
}

// Post adds a POST method handler to the Route.
func (r *Route) Post(h HandlerFunc) *Route {
	r.handlers = append(r.handlers, MethodHandler{
		method:  http.MethodPost,
		handler: Handler(h),
	})

	return r
}

// Put adds a PUT method handler to the Route.
func (r *Route) Put(h HandlerFunc) *Route {
	r.handlers = append(r.handlers, MethodHandler{
		method:  http.MethodPut,
		handler: Handler(h),
	})

	return r
}

// Delete adds a DELETE method handler to the Route.
func (r *Route) Delete(h HandlerFunc) *Route {
	r.handlers = append(r.handlers, MethodHandler{
		method:  http.MethodDelete,
		handler: Handler(h),
	})

	return r
}

// Connect adds a CONNECT method handler to the Route.
func (r *Route) Connect(h HandlerFunc) *Route {
	r.handlers = append(r.handlers, MethodHandler{
		method:  http.MethodConnect,
		handler: Handler(h),
	})

	return r
}

// Head adds a HEAD method handler to the Route.
func (r *Route) Head(h HandlerFunc) *Route {
	r.handlers = append(r.handlers, MethodHandler{
		method:  http.MethodHead,
		handler: Handler(h),
	})

	return r
}

// Patch adds a PATCH method handler to the Route.
func (r *Route) Patch(h HandlerFunc) *Route {
	r.handlers = append(r.handlers, MethodHandler{
		method:  http.MethodPatch,
		handler: Handler(h),
	})

	return r
}

// Trace adds a TRACE method handler to the Route.
func (r *Route) Trace(h HandlerFunc) *Route {
	r.handlers = append(r.handlers, MethodHandler{
		method:  http.MethodTrace,
		handler: Handler(h),
	})

	return r
}
