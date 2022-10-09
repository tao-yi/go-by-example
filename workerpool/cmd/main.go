package main

import (
	"fmt"
	"go-by-example/workerpool"
	"io"
	"net/http"
	"time"
)

func main() {
	urls := []string{
		"https://jsonplaceholder.typicode.com/todos/1",
		"https://jsonplaceholder.typicode.com/todos/2",
		"https://jsonplaceholder.typicode.com/todos/3",
		"https://jsonplaceholder.typicode.com/todos/4",
		"https://jsonplaceholder.typicode.com/todos/5",
		"https://jsonplaceholder.typicode.com/todos/6",
		"https://jsonplaceholder.typicode.com/todos/7",
		"https://jsonplaceholder.typicode.com/todos/8",
		"https://jsonplaceholder.typicode.com/todos/9",
		"https://jsonplaceholder.typicode.com/todos/10",
	}
	start := time.Now()
	defer func() {
		fmt.Printf("fetching all cost %f seconds\n", time.Since(start).Seconds())
	}()

	pool := workerpool.New(1)
	for _, url := range urls {
		url := url
		pool.Submit(func() *workerpool.Result {
			start := time.Now()
			defer func() {
				fmt.Printf("fetching %s cost %f seconds\n", url, time.Since(start).Seconds())
			}()
			resp, err := http.Get(url)
			if err != nil {
				return &workerpool.Result{Err: err, Data: nil}
			}
			defer resp.Body.Close()
			body, err := io.ReadAll(resp.Body)
			if err != nil {
				return &workerpool.Result{Err: err, Data: nil}
			}
			return &workerpool.Result{Err: nil, Data: string(body)}
		})
	}

	results := pool.Await()
	for _, result := range results {
		fmt.Println(result.Data)
	}
}
