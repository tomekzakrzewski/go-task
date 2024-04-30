build:
	@go build -o bin/task

run: build
	@./bin/task