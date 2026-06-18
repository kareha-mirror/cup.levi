all: build

build:
	go build -o levi ./cmd/levi

clean:
	rm -f levi

run:
	go run ./cmd/levi

fmt:
	go fmt ./...

test:
	go test ./...

windows:
	GOOS=windows go build -o levi-windows ./cmd/levi
