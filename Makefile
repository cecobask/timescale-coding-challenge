.PHONY: *

start:
	docker compose up -d

stop:
	docker compose down --volumes
	rm -rf dbdata

build:
	go build -o build/ts main.go

lint:
	golangci-lint run

lint-fix:
	golangci-lint run --fix
