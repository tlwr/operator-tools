test:
	go test -v ./...

build:
	go build

install: build
	go install
