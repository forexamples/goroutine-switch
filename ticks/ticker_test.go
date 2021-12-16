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

func startHTTPServer() {
	http.HandleFunc("/", pprof.Index)
	err := http.ListenAndServe(":6060", nil)
	if err != nil {
		log.Printf("ListenAndServe err:%#v", err)
	}
}

func createTickers(count int) []TickerI {
	list := make([]TickerI, count)
	for i := 0; i < count; i++ {
		list[i] = NewTicker(time.Millisecond * 33)
	}
	return list
}

func TestTicker(t *testing.T) {
	m := createTickers(10000)
	resultChan := make(chan int, 1)

	// var wg sync.WaitGroup
	for i := 0; i < len(m); i++ {
		// wg.Add(1)
		localI := i
		ti := m[i]
		ti.Tick(func() {
			ti.Stop()
			// wg.Done()
			resultChan <- localI
		})
	}

	for i := 0; i < len(m); i++ {
		res := <-resultChan
		t.Logf("received: %d", res)
	}

	t.Log("全部关闭,等待10秒")
	time.Sleep(10 * time.Second)
}
