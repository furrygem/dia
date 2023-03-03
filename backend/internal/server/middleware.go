package server

import (
	"net/http"

	"github.com/furrygem/dia/internal/logging"
)

func loggingMiddleware(next http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		lrw := newLoggingResponseWriter(w)
		next.ServeHTTP(lrw, r)
		e := logging.GetLogger().WithField("remote_addr", r.RemoteAddr)
		e = e.WithField("method", r.Method)
		e = e.WithField("url", r.URL)
		e = e.WithField("status", lrw.status)
		e.Info()
	}
}

type loggingResponseWriter struct {
	http.ResponseWriter
	status int
}

func newLoggingResponseWriter(w http.ResponseWriter) *loggingResponseWriter {
	return &loggingResponseWriter{
		status:         200,
		ResponseWriter: w,
	}
}

func (lrw *loggingResponseWriter) WriteHeader(status int) {
	lrw.status = status
	lrw.ResponseWriter.WriteHeader(status)
}
