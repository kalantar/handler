IMG ?= handler:latest

# Before pushing, `make` or `make all` will run through most of the build tests
all: lint vet test coverage

lint:
	golint ./...

vet:
	go vet ./...

test:
	go test ./... -coverprofile=coverage.out

coverage:
	@echo "test coverage: $(shell go tool cover -func coverage.out | grep total | awk '{print substr($$3, 1, length($$3)-1)}')"

docker-build:
	docker build . -t ${IMG}

docker-push:
	docker push ${IMG}
