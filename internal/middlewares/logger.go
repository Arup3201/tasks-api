package middlewares

import (
	"net/http"
	"os"
	"time"

	"github.com/sirupsen/logrus"
)

type HTTPLogger struct {
	*logrus.Logger
}

func NewLogger(level, format string) *HTTPLogger {
	logger := logrus.New()

	switch level {
	case "debug":
		logger.SetLevel(logrus.DebugLevel)
	case "info":
		logger.SetLevel(logrus.InfoLevel)
	case "warn":
		logger.SetLevel(logrus.WarnLevel)
	case "error":
		logger.SetLevel(logrus.ErrorLevel)
	default:
		logger.SetLevel(logrus.InfoLevel)
	}

	if format == "json" {
		logger.SetFormatter(&logrus.JSONFormatter{
			TimestampFormat: "2006-01-02T15:04:05.000Z07:00",
		})
	} else {
		logger.SetFormatter(&logrus.TextFormatter{})
	}

	logger.SetOutput(os.Stdout)

	return &HTTPLogger{logger}
}

func (logger *HTTPLogger) RequireLogging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		var ww = &responseWriter{w, http.StatusOK}

		next.ServeHTTP(ww, r)

		logger.WithFields(logrus.Fields{
			"method":      r.Method,
			"url":         r.URL.String(),
			"status_code": ww.statusCode,
			"duration":    time.Since(start).String(),
			"remote_addr": r.RemoteAddr,
			"user_agent":  r.UserAgent(),
		}).Info("HTTP Request")
	})
}

type responseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (w *responseWriter) WriteHeader(code int) {
	w.ResponseWriter.WriteHeader(code)
	w.statusCode = code
}
