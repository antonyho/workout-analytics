BIN=analytics
PWD=$(shell pwd)

.PHONY: all dep dep-update docker-image lint

all: clean dep-download dep-update test build

dep:
	go mod init github.com/antonyho/workout-analytics
	go mod tidy

dep-update: test
	go list -m all
	go mod tidy

dep-download:
	go mod download

clean:
	go clean
	if [ -f ${BIN} ]; then rm ${BIN}; fi

lint:
	docker run -t --rm -v ${PWD}:/app -w /app golangci/golangci-lint:v1.55.2 golangci-lint -E revive run -v

test:
	go test -v ./...

test-coverage:
	go test -cover ./...

build:
	go build -o ${BIN} .

run:
	go run .

docker-image:
	docker build --rm -t twaiv/workout-analytics .