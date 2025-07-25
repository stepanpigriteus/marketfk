package router

import (
	"fmt"
	"marketfuck/internal/adapter/in/http/handler"
	"marketfuck/internal/adapter/in/http/middleware"
	"net/http"
)

type Route struct {
	Method      string
	Path        string
	Handler     http.HandlerFunc
	Middlewares []func(http.Handler) http.Handler
}

func RegisterRoutes(mux *http.ServeMux, handlers *handler.AllHandlers) {
	routes := []Route{
		{Method: http.MethodGet, Path: "/prices/latest/", Handler: handlers.Price.HandleGetLatestPrice, Middlewares: []func(http.Handler) http.Handler{middleware.LoggerMiddleware}},
		{Method: http.MethodGet, Path: "/prices/latest/{exchange}/{symbol}", Handler: handlers.Price.HandleGetLatestPriceByExchange, Middlewares: []func(http.Handler) http.Handler{middleware.LoggerMiddleware}},
		{Method: http.MethodGet, Path: "/prices/highest/{symbol}", Handler: handlers.Price.HandleGetHighestPriceInPeriod, Middlewares: []func(http.Handler) http.Handler{middleware.LoggerMiddleware}},
		{Method: http.MethodGet, Path: "/prices/highest/{exchange}/{symbol}", Handler: handlers.Price.HandleGetHighestPriceByExchange, Middlewares: []func(http.Handler) http.Handler{middleware.LoggerMiddleware}},
		{Method: http.MethodGet, Path: "/prices/lowest/{symbol}", Handler: handlers.Price.HandleGetLowestPrice, Middlewares: []func(http.Handler) http.Handler{middleware.LoggerMiddleware}},
		{Method: http.MethodGet, Path: "/prices/lowest/{exchange}/{symbol}", Handler: handlers.Price.HandleGetLowestPriceByExchange, Middlewares: []func(http.Handler) http.Handler{middleware.LoggerMiddleware}},
		{Method: http.MethodGet, Path: "/prices/average/{symbol}", Handler: handlers.Price.HandleGetAveragePrice, Middlewares: []func(http.Handler) http.Handler{middleware.LoggerMiddleware}},
		{Method: http.MethodGet, Path: "/prices/average/{exchange}/{symbol}", Handler: handlers.Price.HandleGetAveragePriceByExchange, Middlewares: []func(http.Handler) http.Handler{middleware.LoggerMiddleware}},
		{Method: http.MethodGet, Path: "/mode/test", Handler: handlers.Mode.HandleSwitchToTestMode, Middlewares: []func(http.Handler) http.Handler{middleware.LoggerMiddleware}},
		{Method: http.MethodGet, Path: "/mode/live", Handler: handlers.Mode.HandleSwitchToLiveMode, Middlewares: []func(http.Handler) http.Handler{middleware.LoggerMiddleware}},
		{Method: http.MethodGet, Path: "/health", Handler: handlers.Health.HandleHealthCheck, Middlewares: []func(http.Handler) http.Handler{middleware.LoggerMiddleware}},
	}

	for _, route := range routes {
		finalHandler := applyMiddleware(route.Handler, route.Middlewares)
		mux.Handle(route.Path, methodHandler(route.Method, finalHandler))
	}
}

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
