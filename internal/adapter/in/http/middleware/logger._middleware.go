package middleware

import (
	"fmt"
	"net/http"
	"time"

	"marketfuck/pkg/logger"
)

func LoggerMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		// Логиру ем информацию о запросе
		logger.NewSlogAdapter().Info(fmt.Sprintf("Started %s %s", r.Method, r.URL.Path))

		// Обрабатываем запрос
		next.ServeHTTP(w, r)

		// Логируем информацию о времени выполнения запроса
		duration := time.Since(start)
		logger.NewSlogAdapter().Info(fmt.Sprintf("Completed %s %s in %v", r.Method, r.URL.Path, duration))
	})
}
