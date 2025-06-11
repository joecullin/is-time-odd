# This is the first I've touched make in a long time.
#
# I guess JS's package.json has conditioned me to look for build commands in a standard place,
# rather than some random collection of shell scripts.
#
# I'm sure they could be cleaner, but they definitely saved me time during dev & test.
#
###############################################################################

# For building the server app: use the Makefile in ../server

# Use this Makefile for:
# - Commands for working on and building the client app.
# - Copying the built app (all versions, plus info file ) into server dir.
# - Creating "dist" zip files for each platform with client, server, and data.

format:
	gofmt -w *.go

clean:
	rm -rf bin/* buildtmp/* ummmm/* ../server/server_data/* ../dist/*
	mkdir -p bin
	mkdir -p buildtmp/even
	mkdir -p buildtmp/odd
	mkdir -p buildtmp/app
	mkdir -p ../server/server_data
	mkdir -p ../dist

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

	# copy odd & even builds to server's data dir:
	cp bin/* ../server/server_data

	# update server's appData.json file:
	node ../utils/generateDataFile.js

	# build "un-initialized" versions of app:
	cp *.go buildtmp/app
	perl -pi -e 's/PLACEHOLDER_VERSION_STRING/1.0/g' buildtmp/app/*.go
	perl -pi -e 's/PLACEHOLDER_TIME_VALUE/UNCERTAIN/g' buildtmp/app/*.go
	GOOS=darwin GOARCH=amd64 go build -o bin/is_time_odd_mac buildtmp/app/*.go
	GOOS=linux GOARCH=amd64 go build -o bin/is_time_odd_linux buildtmp/app/*.go
	GOOS=windows GOARCH=amd64 go build -o bin/is_time_odd.exe buildtmp/app/*.go

	@echo

dist:
	# Create a zip file per platform, with just the apps & data for that platform

	mkdir -p ../dist

	mkdir -p buildtmp/dist/mac
	mkdir -p buildtmp/dist/linux
	mkdir -p buildtmp/dist/windows

	cp ../server/bin/server_mac buildtmp/dist/mac/server
	cp ../server/bin/server.exe buildtmp/dist/windows/server.exe
	cp ../server/bin/server_linux buildtmp/dist/linux/server

	cp bin/is_time_odd_mac buildtmp/dist/mac/is_time_odd
	cp bin/is_time_odd.exe buildtmp/dist/windows/is_time_odd.exe
	cp bin/is_time_odd_linux buildtmp/dist/linux/is_time_odd

	cp -R ../server/server_data buildtmp/dist/mac
	cp -R ../server/server_data buildtmp/dist/windows
	cp -R ../server/server_data buildtmp/dist/linux

	# prune the "other OS" files from server_data dir, to make demo zip smaller
	pushd buildtmp/dist/mac/server_data && rm *.exe *linux && popd
	pushd buildtmp/dist/windows/server_data && rm *mac *linux && popd
	pushd buildtmp/dist/linux/server_data && rm *mac *.exe && popd

	# pushd buildtmp/dist/mac && find . -type f
	pushd buildtmp/dist/mac && zip -r ../../../../dist/ito_mac.zip * && popd
	pushd buildtmp/dist/windows && zip -r ../../../../dist/ito_windows.zip * && popd
	pushd buildtmp/dist/linux && zip -r ../../../../dist/ito_linux.zip * && popd

	@echo "done creating dist files!"
