package main

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"time"
)

func slowGetRequest(ctx context.Context, millisecond int) <-chan string {
	ch := make(chan string)
	go func() {
		ctx, cancel := context.WithTimeout(ctx, time.Millisecond*time.Duration(millisecond))
		defer cancel()
		req, err := http.NewRequestWithContext(ctx, http.MethodGet, "http://www.google.com", nil)
		if err != nil {
			fmt.Println("server err:", err)
			ch <- err.Error()
			return
		}
		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			// 如果客户端取消request，那么这里error会弹出context canceled
			if errors.Is(err, context.Canceled) {
				fmt.Println("client request canceled")
			} else {
				fmt.Println("server err:", err)
			}
			ch <- err.Error()
			return
		}
		defer resp.Body.Close()
		bytes, err := io.ReadAll(resp.Body)
		if err != nil {
			fmt.Println("server err:", err)
			ch <- err.Error()
			return
		}
		ch <- string(bytes)
	}()
	return ch
}

/*
当客户端取消请求时，会通过req.Context发送一个cancel信号
我们可以通过<-ctx.Done()拿到这个channel中的信号
*/
func main() {
	http.HandleFunc("/hello", func(w http.ResponseWriter, req *http.Request) {
		timeout := req.URL.Query().Get("timeout")
		millisecond, err := strconv.Atoi(timeout)
		if err != nil {
			millisecond = 10 // 10 seconds
		}
		ctx := req.Context()
		fmt.Println("server: hello handler started")
		defer fmt.Println("server: hello handler ended")
		msg := <-slowGetRequest(ctx, millisecond)
		fmt.Fprint(w, msg)
	})
	err := http.ListenAndServe(":8090", nil)
	if err != nil {
		panic(err)
	}
}
