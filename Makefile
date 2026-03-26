all: build

env:
	go env -w GOPRIVATE=tea.kareha.org

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
