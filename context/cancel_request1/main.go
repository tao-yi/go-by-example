package main

import (
	"fmt"
	"net/http"
	"time"
)

func slowGetRequest() <-chan string {
	ch := make(chan string)
	go func() {
		time.Sleep(time.Second * 3)
		ch <- "hello"
	}()
	return ch
}

/*
当客户端取消请求时，会通过req.Context发送一个cancel信号
我们可以通过<-ctx.Done()拿到这个channel中的信号
*/
func main() {
	http.HandleFunc("/hello", func(w http.ResponseWriter, req *http.Request) {
		ctx := req.Context()
		fmt.Println("server: hello handler started")
		defer fmt.Println("server: hello handler ended")
		select {
		case msg := <-slowGetRequest():
			fmt.Fprint(w, msg)
		case <-ctx.Done():
			err := ctx.Err()
			fmt.Println("server:", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	})
	err := http.ListenAndServe(":8090", nil)
	if err != nil {
		panic(err)
	}
}
