all:
	mkdir -p dist
	go build -o dist/ ./cmd/...

clean:
	rm -rf ./dist

run:
	go run main.go
