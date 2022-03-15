package internalhttp

import (
	"net/http"
)

type ResposeWriter struct {
	http.ResponseWriter
	StatusCode int
	Bytes      int
}

func (w *ResposeWriter) WriteHeader(statusCode int) {
	w.StatusCode = statusCode
	w.ResponseWriter.WriteHeader(statusCode)
}

func (w *ResposeWriter) Write(bytes []byte) (int, error) {
	r, err := w.ResponseWriter.Write(bytes)
	w.Bytes += r

	return r, err
}

func loggingMiddleware(next http.Handler, logger Logger) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		wrt := &ResposeWriter{w, 0, 0}
		next.ServeHTTP(wrt, r)
		logger.LogHTTP(r, wrt.StatusCode, wrt.Bytes)
	})
}
