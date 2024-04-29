.PHONY: *

validate:
	@test -f .env || (echo "The .env file is required!\nPlease, create it by using the .env.example file as guidance..." && exit 1)

start: validate
	docker compose up --build --detach

stop:
	docker compose down --volumes
	rm -rf dbdata

build:
	go build -o build/ts main.go

lint:
	golangci-lint run

lint-fix:
	golangci-lint run --fix
