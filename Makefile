run:
	go run main.go

deps:
	go get -u github.com/shurcooL/vfsgen

build:
	go generate
	go build -o main/main main/main.go

.PHONY: run build deps
