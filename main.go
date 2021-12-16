package main

import (
	"flag"
	"log"
	"math/rand"
	"os"
	"os/signal"
	"runtime"
	"runtime/pprof"
	"runtime/trace"
	"sync"
	"time"
)

func main() {
	var (
		runTime        = flag.Duration("runtime", 10*time.Second, "Run `duration` after target go routine count is reached")
		workDur        = flag.Duration("work", 15*time.Microsecond, "CPU bound work `duration` each cycle")
		cycleDur       = flag.Duration("cycle", 2400*time.Microsecond, "Cycle `duration`")
		gCount         = flag.Int("gcount", runtime.NumCPU(), "Number of `goroutines` to use")
		gStartFreq     = flag.Int("gfreq", 1, "Number of goroutines to start each second until gcount is reached")
		cpuProfilePath = flag.String("cpuprofile", "", "Write CPU profile to `file`")
		tracePath      = flag.String("trace", "", "Write execution trace to `file`")
	)

	flag.Parse()

	sigC := make(chan os.Signal, 1)
	signal.Notify(sigC, os.Interrupt)

	var wg sync.WaitGroup
	done := make(chan struct{})
	stop := make(chan struct{})

	wg.Add(1)
	go func() {
		defer wg.Done()
		select {
		case sig := <-sigC:
			log.Print("got signal ", sig)
		case <-stop:
		}
		close(done)
	}()

	gFreq := time.Second / time.Duration(*gStartFreq)
	jitterCap := int64(gFreq / 2)

	for g := 0; g < *gCount; g++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			ticker := time.NewTicker(*cycleDur)
			defer ticker.Stop()
			for {
				select {
				case <-done:
					return
				case <-ticker.C:
					workUntil(time.Now().Add(*workDur))
				}
			}
		}(g)
		log.Print("goroutine count: ", g+1)
		jitter := time.Duration(rand.Int63n(jitterCap))
		select {
		case <-done:
			g = *gCount // stop loop early
		case <-time.After(gFreq + jitter):
		}
	}

	select {
	case <-done:
	default:
		log.Print("running for ", *runTime)
		runTimer := time.NewTimer(*runTime)
		wg.Add(1)
		go func() {
			wg.Done()
			select {
			case <-runTimer.C:
				log.Print("runTimer fired")
				close(stop)
			}
		}()
	}

	if *cpuProfilePath != "" {
		f, err := os.Create(*cpuProfilePath)
		if err != nil {
			log.Fatal("could not create CPU profile: ", err)
		}
		if err := pprof.StartCPUProfile(f); err != nil {
			log.Fatal("could not start CPU profile: ", err)
		}
		log.Print("profiling")
		defer pprof.StopCPUProfile()
	}

	if *tracePath != "" {
		f, err := os.Create(*tracePath)
		if err != nil {
			log.Fatal("could not create execution trace: ", err)
			os.Exit(1)
		}
		defer f.Close()
		if err := trace.Start(f); err != nil {
			log.Fatal("could not start execution trace: ", err)
		}
		log.Print("tracing")
		defer trace.Stop()
	}

	wg.Wait()
}

func workUntil(deadline time.Time) {
	now := time.Now()
	for now.Before(deadline) {
		now = time.Now()
	}
}