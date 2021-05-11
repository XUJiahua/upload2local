fmt:
	go fmt ./...
server:fmt
	go run main.go server --localtunnel=false
testdata:fmt
	go run main.go testdata -s 1024000
client:testdata
	go run main.go upload --host http://localhost:1234 hello.bin
client-public:testdata
	go run main.go upload -v --host https://new-moth-89.loca.lt hello.bin
#test:
#	curl -F "data=@main.go" localhost:1234
#	diff main.go /tmp/main.go
diff:
	diff hello.bin ./data/hello.bin