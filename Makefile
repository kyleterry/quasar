test:
	go test $(go list ./... | grep -v vendor)

.PHONY: test
