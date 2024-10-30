EXPORTER_VERSION=1.4
PACKAGES_DIR=compiled_packages

all: test build clean

test:
	go fmt ./
	go fix ./
	go vet -v ./
	staticcheck ./ || true
	golangci-lint run
	go mod tidy

build:
	go build -o wanguard_exporter -v

clean:
	rm -f wanguard_exporter

run:
	go run .

compile:
	CGO_ENABLED=0 GOARCH=amd64 GOOS=darwin go build -o ${PACKAGES_DIR}/wanguard_exporter-${EXPORTER_VERSION}-darwin
	CGO_ENABLED=0 GOARCH=amd64 GOOS=linux go build -o ${PACKAGES_DIR}/wanguard_exporter-${EXPORTER_VERSION}-linux
	CGO_ENABLED=0 GOARCH=amd64 GOOS=windows go build -o ${PACKAGES_DIR}/wanguard_exporter-${EXPORTER_VERSION}-windows
