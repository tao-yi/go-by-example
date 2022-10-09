package main

import "net/http"

func helloHandler(message string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, err := w.Write([]byte(message))
		if err != nil {
			panic(err)
		}
	})
}

func helloMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// middleware logic
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
