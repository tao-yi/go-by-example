package main

import (
	"fmt"
	"sync"
	"sync/atomic"
)

func unsafe(v *uint32, wg *sync.WaitGroup) {
	for i := 0; i < 1000; i++ {
		*v++
	}
	wg.Done()
}

func safe(v *uint32, wg *sync.WaitGroup) {
	for i := 0; i < 1000; i++ {
		atomic.AddUint32(v, 1)
	}
	wg.Done()
}

func main() {
	var v1 uint32 = 10
	var wg1 sync.WaitGroup
	wg1.Add(3)

	go unsafe(&v1, &wg1)
	go unsafe(&v1, &wg1)
	go unsafe(&v1, &wg1)
	wg1.Wait()

	var v2 uint32 = 10
	var wg2 sync.WaitGroup
	wg2.Add(3)
	go safe(&v2, &wg2)
	go safe(&v2, &wg2)
	go safe(&v2, &wg2)
	wg2.Wait()

	fmt.Printf("v1:%d, v2:%d\n", v1, v2) // 2084, 3010
}
