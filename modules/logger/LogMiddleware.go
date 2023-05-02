package logger

import "net/http"

type LogMiddleware struct {
	handler http.Handler
}

func (middleware *LogMiddleware) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	Log("Incoming", r.Method, "request on", r.URL.Path)

	middleware.handler.ServeHTTP(w, r)
}

func NewLogMiddleware(handler http.Handler) *LogMiddleware {
	return &LogMiddleware{handler}
}
