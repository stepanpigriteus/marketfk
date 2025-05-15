package router

import (
	"fmt"
	"net/http"
)

type Route struct {
	Method      string
	Path        string
	Handler     http.HandlerFunc
	Middlewares []func(http.Handler) http.Handler
}

// func RegisterRoutes(mux *http.ServeMux, handlers *handler.AllHandlers) {
// 	routes := []Route{}

// 	for _, route := range routes {
// 		finalHandler := applyMiddleware(route.Handler, route.Middlewares)
// 		mux.Handle(route.Path, methodHandler(route.Method, finalHandler))
// 	}
// }

func applyMiddleware(h http.Handler, middlewares []func(http.Handler) http.Handler) http.Handler {
	for _, m := range middlewares {
		h = m(h)
	}
	return h
}

func methodHandler(method string, handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Printf(">>> methodHandler: got %s request on %s\n", r.Method, r.URL.Path)
		if r.Method != method {
			http.Error(w, `{"error":"method not allowed"}`, http.StatusMethodNotAllowed)
			return
		}
		handler.ServeHTTP(w, r)
	})
}
