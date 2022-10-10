package main

import (
	"context"
	"fmt"
	"net/http"
)

type RequestID string

const RequestIDKey RequestID = "request_id"

func helloHandler(message string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestID := r.Context().Value(RequestIDKey)
		fmt.Fprintf(w, "request[%s]: %s", requestID, message)
	})
}

func helloMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// middleware logic
		// extract request id
		requestID := r.URL.Query().Get(string(RequestIDKey))
		ctx := r.Context()
		r = r.WithContext(context.WithValue(ctx, RequestIDKey, requestID))
		next.ServeHTTP(w, r)
	})
}

func main() {
	http.Handle("/hello", helloMiddleware(helloHandler("hello")))
	err := http.ListenAndServe(":8090", nil)
	if err != nil {
		panic(err)
	}
}
