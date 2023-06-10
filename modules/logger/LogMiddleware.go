package logger

import (
	"net/http"
)

type LogMiddleware struct {
	handler http.Handler
	logger  *LoggerService
}

func (middleware *LogMiddleware) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	middleware.logger.Log("Incoming", r.Method, "request on", r.URL.Path)

	middleware.handler.ServeHTTP(w, r)
}

func NewLogMiddleware(handler http.Handler, logger *LoggerService) *LogMiddleware {
	return &LogMiddleware{
		handler,
		logger,
	}
}
