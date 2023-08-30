build:
	make server

server:
	go build -o bin/server cmd/server/main.go

clean:
	rm -rf bin

test:
	go test ./...

# used in CI/CD, add more targets if needed
ci:
	make build test
