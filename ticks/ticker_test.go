package ticks

import (
	"log"
	"net/http"
	"net/http/pprof"
	"sync"
	"testing"
	"time"
)

func TestOne(t *testing.T) {
	wg := sync.WaitGroup{}
	wg.Add(1)
	num := 2
	for i := 0; i < num; i++ {
		ti := NewTicker(time.Millisecond * 33)
		ti.Tick(func(i int) func() {
			return func() {
				if i == num-1 {
					wg.Done()
				}
			}
		}(i))
	}

	wg.Wait()
}

func startHttpServer() {
	http.HandleFunc("/", pprof.Index)
	err := http.ListenAndServe(":6060", nil)
	if err != nil {
		log.Printf("ListenAndServe err:%#v", err)
	}
}

func TestTicker(t *testing.T) {
	m := map[int]TickerI{}
	go startHttpServer()

	num := 10000
	resultChan := make(chan int, 1)

	go func() {
		for i := 0; i < num; i++ {
			ti := NewTicker(time.Millisecond * 33)
			m[i] = ti
		}

		for i := 0; i < num; i++ {
			ti := m[i]
			ti.Tick(func(i int) func() {
				return func() {
					ti := m[i]
					ti.Stop()
					resultChan <- i
				}
			}(i))
		}
	}()

	for i := 0; i < num; i++ {
		res := <-resultChan
		t.Log(res)
	}

	t.Log("全部关闭")
	select {}
}
