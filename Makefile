.PHONY=run
run:
	@go mod tidy
	@go run cmd/progress/main.go

.PHONY=fmt
fmt:
	@go fmt cmd/progress/main.go
	@go fmt db/*.go
	@go fmt data/*.go
