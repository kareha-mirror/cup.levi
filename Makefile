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

tidy:
	grep -v '^.tea.kareha.org' go.mod > go.mod.clipped
	mv go.mod.clipped go.mod
	GOPRIVATE=tea.kareha.org go mod tidy

windows:
	GOOS=windows go build -o levi-windows ./cmd/levi
