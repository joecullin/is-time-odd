format:
	gofmt -w *.go

clean:
	rm -rf bin/* buildtmp/* ummmm/*
	mkdir -p bin
	mkdir -p buildtmp/even
	mkdir -p buildtmp/odd
	mkdir -p buildtmp/app

build:
	make clean

	# build "even" versions of app:
	cp *.go buildtmp/even
	perl -pi -e 's/PLACEHOLDER_VERSION_STRING/1.2/g' buildtmp/even/*.go
	perl -pi -e 's/PLACEHOLDER_TIME_VALUE/NOT odd/g' buildtmp/even/*.go
	GOOS=darwin GOARCH=amd64 go build -o bin/app_even_mac buildtmp/even/*.go
	GOOS=linux GOARCH=amd64 go build -o bin/app_even_linux buildtmp/even/*.go
	GOOS=windows GOARCH=amd64 go build -o bin/app_even.exe buildtmp/even/*.go

	# build "odd" versions of app:
	cp *.go buildtmp/odd
	perl -pi -e 's/PLACEHOLDER_VERSION_STRING/1.3/g' buildtmp/odd/*.go
	perl -pi -e 's/PLACEHOLDER_TIME_VALUE/ODD/g' buildtmp/odd/*.go
	GOOS=darwin GOARCH=amd64 go build -o bin/app_odd_mac buildtmp/odd/*.go
	GOOS=linux GOARCH=amd64 go build -o bin/app_odd_linux buildtmp/odd/*.go
	GOOS=windows GOARCH=amd64 go build -o bin/app_odd.exe buildtmp/odd/*.go

	# build "un-initialized" versions of app:
	cp *.go buildtmp/app
	perl -pi -e 's/PLACEHOLDER_VERSION_STRING/1.0/g' buildtmp/app/*.go
	perl -pi -e 's/PLACEHOLDER_TIME_VALUE/UNCERTAIN/g' buildtmp/app/*.go
	GOOS=darwin GOARCH=amd64 go build -o bin/is_time_odd_mac buildtmp/app/*.go
	GOOS=linux GOARCH=amd64 go build -o bin/is_time_odd_linux buildtmp/app/*.go
	GOOS=windows GOARCH=amd64 go build -o bin/is_time_odd.exe buildtmp/app/*.go

	#TODO (after progress on server) - copy builds to server dir, create json file for server

	@echo