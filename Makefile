.PHONY: test vet lint

test:
	go test -v -race -coverprofile=coverage.txt -covermode=atomic ./...

vet:
	go vet ./...

lint:
	golint -set_exit_status ./...
