fmt:
	go fmt ./...
run:fmt
	go run main.go
test:
	curl -F "data=@main.go" localhost:8090
	diff main.go /tmp/main.go
