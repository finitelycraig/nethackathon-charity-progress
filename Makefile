GOFILES := $(shell find . -type f -name *.go)

charityprogress: $(GOFILES)
	@go mod tidy
	@go build -o build/charityprogress cmd/progress/main.go

charityprogress-ssh: $(GOFILES)
	@go mod tidy
	@go build -o build/charityprogress-ssh cmd/serve/main.go

build-remote: build/charityprogress-ssh-linux-amd64

build/charityprogress-ssh-linux-amd64: $(GOFILES)
	@env GOOS=linux GOARCH=amd64 go build -o build/charityprogress-ssh-linux-amd64 cmd/serve/main.go

.PHONY=run
run:
	@go mod tidy
	@go run cmd/progress/main.go

.PHONY=serve
serve:
	@go mod tidy
	@go run cmd/serve/main.go

.PHONY=fmt
fmt:
	@go fmt cmd/progress/main.go
	@go fmt db/*.go
	@go fmt data/*.go
	@go fmt internal/*.go
