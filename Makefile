test:
	go test $(shell go list ./... | grep -v tinygo)
