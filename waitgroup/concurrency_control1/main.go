package main

import (
	"fmt"
	"io"
	"net/http"
	"sync"
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

func fetch(url string, res chan<- string) {
	fmt.Printf("starting fetching %s\n", url)
	now := time.Now()
	r, err := http.Get(url)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	defer r.Body.Close()
	_, err = io.ReadAll(r.Body)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	fmt.Printf("fetching %s, %f seconds passed\n", url, time.Since(now).Seconds())
}

/*
using sliding window to send http request in fixed batch size
*/
func main() {
	now := time.Now()
	var wg sync.WaitGroup
	batch_size := 3
	res := make(chan string, batch_size)
	start := 0
	end := 0
	for end < len(urls) {
		wg.Add(1)
		go func(url string) {
			fetch(url, res)
			wg.Done()
		}(urls[end])
		window_size := end - start + 1
		if window_size == batch_size || end+1 == len(urls) {
			wg.Wait()
			start = end + 1
		}
		end = end + 1
	}
	fmt.Printf("total %f seconds passed", time.Since(now).Seconds())
}
