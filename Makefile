test:
	go test ./...

# used in CI/CD, add more targets if needed
ci:
	make test
