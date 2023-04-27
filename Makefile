build:
	go build -o app cmd/main.go 

run:
	./app

build-run: build run