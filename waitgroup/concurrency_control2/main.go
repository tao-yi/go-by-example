package main

import (
	"fmt"
	"io"
	"net/http"
	"time"
)

var urls []string

func init() {
	urls = []string{
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
}

func fetch(urls <-chan string, results chan<- []byte) {
	for url := range urls {
		fmt.Printf("starting fetching %s\n", url)
		now := time.Now()
		r, _ := http.Get(url)
		defer r.Body.Close()
		bytes, _ := io.ReadAll(r.Body)
		results <- bytes
		fmt.Printf("fetching %s, %f seconds passed\n", url, time.Since(now).Seconds())
	}
}

func main() {
	now := time.Now()
	batch_size := 10
	tasks := make(chan string, batch_size)
	results := make(chan []byte, batch_size)

	for i := 0; i < batch_size; i++ {
		go fetch(tasks, results)
	}

	go func() {
		for _, url := range urls {
			tasks <- url
		}
		close(tasks)
	}()

	for i := 0; i < len(urls); i++ {
		<-results
	}
	fmt.Printf("total %f seconds passed", time.Since(now).Seconds())
}
