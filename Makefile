test:
	tmux new-session -d -s eg
	tmux split-window -t "eg:0"
	tmux split-window -t "eg:0.0" -h
	tmux send-keys -t "eg:0.0" "go test -run ^TestTicker  goroutine-switch/ticks -timeout 60s" Enter
	tmux send-keys -t "eg:0.1" "sleep 0.2 && go tool pprof -web http://localhost:6060/debug/pprof/profile?seconds=30" Enter
	tmux attach -t eg
	tmux kill-session -t eg

.PHONY: test
