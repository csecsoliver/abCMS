.PHONY: generate build test run

generate:
	go tool templ generate

build: generate
	go build -o abcms ./cmd/abcms

test: generate
	go test -v ./...

run: generate
	go run ./cmd/abcms
