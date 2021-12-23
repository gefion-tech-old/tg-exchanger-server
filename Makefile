SOURCES=./cmd/main
SERVICE=server

.PHONY: init

init:
	mkdir logs

.PHONY: run

run:
	go build -o $(SERVICE) -v $(SOURCES)
	clear		
	./$(SERVICE)


.PHONY: build

build:
	go build -o $(SERVICE) -v $(SOURCES)
	clear


count:
	find . -name tests -prune -o -type f -name '*.go' | xargs wc -l

.DEFAULT_GOAL := run