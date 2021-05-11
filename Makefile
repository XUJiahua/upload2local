fmt:
	go fmt ./...
server:fmt
	go run main.go server
testdata:fmt
	go run main.go testdata -s 102400000
client:testdata
	go run main.go upload --host http://localhost:1234 hello.bin
#test:
#	curl -F "data=@main.go" localhost:1234
#	diff main.go /tmp/main.go
