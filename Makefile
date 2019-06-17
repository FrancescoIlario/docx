run:
	go run main.go

deps:
	go get -u github.com/shurcooL/vfsgen

build:
	go generate
	go build main.go

.PHONY: run build deps
