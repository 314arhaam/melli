.PHONY: build exec

build-api:
	@go build -o bin/melli cmd/main.go

exec:
	bin/melli