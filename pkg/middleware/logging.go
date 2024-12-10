package middleware

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/fatih/color"
)

var (
	methodColor   = color.New(color.FgCyan)
	uriColor      = color.New(color.FgMagenta)
	durationColor = color.New(color.FgYellow)
	successColor  = color.New(color.FgGreen)
	errorColor    = color.New(color.FgRed)
)

type responseWriter struct {
	http.ResponseWriter
	status int
	body   []byte
}

func (rw *responseWriter) WriteHeader(status int) {
	rw.status = status
	rw.ResponseWriter.WriteHeader(status)
}

func (rw *responseWriter) Write(b []byte) (int, error) {
	rw.body = b
	return rw.ResponseWriter.Write(b)
}

func NewColorLogger() *log.Logger {
	return log.New(os.Stdout, "", log.Ldate|log.Ltime)
}

func LoggingMiddleware(logger *log.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			crw := &responseWriter{
				ResponseWriter: w,
				status:         http.StatusOK,
			}
			start := time.Now()
			defer func() {
				if err := recover(); err != nil {
					crw.status = http.StatusInternalServerError
					errorColor.Fprintf(os.Stderr, "[PANIC] %v\n", err)
					logger.Printf(
						"%s %s\t%s\t%s\t%s",
						errorColor.Sprintf("[ERROR]"),
						methodColor.Sprintf("%s", r.Method),
						uriColor.Sprintf("%s", r.RequestURI),
						errorColor.Sprintf("Status: %d", crw.status),
						durationColor.Sprintf("Duration: %v", time.Since(start)),
					)
				}
			}()
			next.ServeHTTP(crw, r)
			var statusColorFunc func(format string, a ...interface{}) string
			if crw.status >= 200 && crw.status < 300 {
				statusColorFunc = successColor.Sprintf
			} else if crw.status >= 400 {
				statusColorFunc = errorColor.Sprintf
			} else {
				statusColorFunc = color.New(color.FgYellow).Sprintf
			}
			logger.Printf(
				"%s %s\t%s\t%s\t%s",
				statusColorFunc("[%d]", crw.status),
				methodColor.Sprintf("%s", r.Method),
				uriColor.Sprintf("%s", r.RequestURI),
				statusColorFunc("Status: %d", crw.status),
				durationColor.Sprintf("Duration: %v", time.Since(start)),
			)
		})
	}
}
