all: fmt vet lint

fmt:
	go fmt .
	go fmt ./fleet
	go fmt ./cmd/nadm

vet:
	go vet .
	go vet ./fleet
	go vet ./cmd/nadm

lint:
	golint .
	golint ./fleet
	golint ./cmd/nadm

install:
	gp run "go install github.com/shiftr-io/nadm/cmd/nadm"
	gp run "cp ./bin/nadm /usr/local/bin" -r
