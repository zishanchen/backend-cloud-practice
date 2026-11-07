package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"sync/atomic"
	"time"
)

type contextKey string

const requestIDKey contextKey = "requestID"

var requestCounter atomic.Uint64

func newRequestID() string {
	id := requestCounter.Add(1)
	return fmt.Sprintf("req-%d", id)
}

func requestIDFromContext(ctx context.Context) string {
	id, ok := ctx.Value(requestIDKey).(string)
	if !ok {
		return ""
	}
	return id
}

func requestIDMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			requestID := r.Header.Get("X-Request-ID")

			if requestID == "" {
				requestID = newRequestID()
			}

			w.Header().Set(
				"X-Request-ID",
				requestID,
			)

			ctx := context.WithValue(
				r.Context(),
				requestIDKey,
				requestID,
			)

			next.ServeHTTP(
				w,
				r.WithContext(ctx),
			)
		},
	)
}

type responseRecorder struct {
	http.ResponseWriter
	status int
}

func (rec *responseRecorder) WriteHeader(status int) {
	rec.status = status
	rec.ResponseWriter.WriteHeader(status)
}

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()

			rec := &responseRecorder{
				ResponseWriter: w,
				status:         http.StatusOK,
			}

			next.ServeHTTP(rec, r)

			log.Printf(
				"request_id=%s method=%s path=%s status=%d duration=%s",
				requestIDFromContext(r.Context()),
				r.Method,
				r.URL.Path,
				rec.status,
				time.Since(start),
			)
		},
	)
}
