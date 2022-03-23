clean:
	go clean
darwin:
	go clean
	GOOS=darwin GOARCH=amd64 go build -ldflags '-s -w'
all:
	make darwin