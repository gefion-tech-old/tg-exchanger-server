SOURCES=./cmd/main
SERVICE=server

.PHONY: run

run:
	go build -o $(SERVICE) -v $(SOURCES)
	clear		
	./$(SERVICE)


.PHONY: build

build:
	go build -o $(SERVICE) -v $(SOURCES)
	clear	


.PHONY: test

test:
	go test -v -race -timeout 30s ./...


count:
	find . -name tests -prune -o -type f -name '*.go' | xargs wc -l

.DEFAULT_GOAL := run