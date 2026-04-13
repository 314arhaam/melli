.PHONY: build exec

build:
	@go build -o bin/melli cmd/main.go

exec:
	bin/melli