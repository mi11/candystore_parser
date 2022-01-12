API_PORT=4000

## build: builds all binaries
build: clean build_candystore_parser
	@printf "All binaries built!\n"

## clean: cleans all binaries and runs go clean
clean:
	@echo "Cleaning..."
	@- rm -f dist/*
	@go clean
	@echo "Cleaned!"

## build_candystore_parser: builds candystore parser
build_candystore_parser:
	@echo "Building candystore parser..."
	@go build -o dist/candystore_parser ./cmd/candystore_parser
	@echo "Candystore parser built!"

## start: starts candystore parser
start: start_candystore_parser

## start_candystore_parser: starts candystore parser
start_candystore_parser: build_candystore_parser
	@echo "Starting the candystore parser..."
	@env ./dist/candystore_parser -port=${API_PORT} &
	@echo "Candystore parser is running!"

## stop: stops candystore parser
stop: stop_candystore_parser
	@echo "All applications stopped"

## stop_candystore_parser: stops candystore parser
stop_candystore_parser:
	@echo "Stopping the candystore parser..."
	@-pkill -SIGTERM -f "candystore_parser -port=${API_PORT}"
	@echo "Stopped the candystore parser"