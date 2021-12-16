test:
	tmux new-session -d -s eg
	tmux split-window -t "eg:0"
	# tmux split-window -t "eg:0.0" -h
	tmux send-keys -t "eg:0.0" "go test -run ^TestTicker  goroutine-switch/ticks -timeout 60s -cpuprofile cpu.prof -memprofile mem.prof" Enter
	tmux send-keys -t "eg:0.1" "go tool pprof -web cpu.prof" Enter
	# tmux send-keys -t "eg:0.1" "sleep 0.1 && go tool pprof -web http://localhost:6060/debug/pprof/profile?seconds=30" Enter
	tmux attach -t eg
	tmux kill-session -t eg


test-and-show:
	rm -fr *.prof
	go test -run ^TestTicker  goroutine-switch/ticks -timeout 60s -cpuprofile cpu.prof -memprofile mem.prof \
		&& go tool pprof -web cpu.prof \
		&& go tool pprof -web mem.prof

.PHONY: test
