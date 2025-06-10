format:
	gofmt -w ./*.go

build:
	GOOS=darwin GOARCH=amd64 go build -o bin/server_mac *.go
	GOOS=linux GOARCH=amd64 go build -o bin/server_linux *.go
	GOOS=windows GOARCH=amd64 go build -o bin/server.exe *.go

	@echo