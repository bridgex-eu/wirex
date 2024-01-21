package wirex

import (
	"net/url"
)

type RoutesGroup struct {
	routes []*Route
}

func NewRoutesGroup() *RoutesGroup {
	return &RoutesGroup{}
}

func (g *RoutesGroup) Route(pattern string, handlers ...MethodHandler) *Route {
	route := Route{pattern: pattern, handlers: handlers}
	g.routes = append(g.routes, &route)

	return &route
}

func (g *RoutesGroup) Use(middleware ...Middleware) {
	for _, route := range g.routes {
		route.middlewares = append(route.middlewares, middleware...)
	}
}

func (g *RoutesGroup) With(key string, val any) {
	g.Use(with(key, val))
}

func (g *RoutesGroup) Group(pattern string, group *RoutesGroup, middlewares ...Middleware) {
	for _, route := range group.routes {
		pattern, err := url.JoinPath(pattern, route.pattern)
		if err != nil {
			panic(err)
		}

		route.pattern = pattern
		route.middlewares = append(route.middlewares, middlewares...)
		g.routes = append(g.routes, route)
	}
}
