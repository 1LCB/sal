package sal

import (
	"log"
	"net/http"
	"time"
)

type CustomResponseWriter struct {
	status_code int
	http.ResponseWriter
}

func (m *CustomResponseWriter) Write(p []byte) (int, error) {
	if m.status_code == 0 {
		m.status_code = http.StatusOK
	}
	return m.ResponseWriter.Write(p)
}

func (m *CustomResponseWriter) WriteHeader(statusCode int) {
	m.status_code = statusCode
	m.ResponseWriter.WriteHeader(statusCode)
}

func LoggerMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		m := &CustomResponseWriter{
			ResponseWriter: w,
		}

		start := time.Now()
		next.ServeHTTP(m, r)
		duration := time.Since(start).Milliseconds()

		log.Printf(
			"| %s %s %s [ %d %s ] (%dms)\n",
			r.RemoteAddr,
			r.Method,
			r.URL.String(),
			m.status_code,
			StatusCodes[m.status_code],
			duration,
		)
	}
}
