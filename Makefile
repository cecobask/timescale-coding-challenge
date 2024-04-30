.PHONY: *

validate:
	@test -f .env || (echo "The .env file is required!\nPlease, create it by using the .env.example file as guidance..." && exit 1)

start: validate
	@docker compose --profile=app up --build --detach --remove-orphans

start-db: validate
	@docker compose up --detach --remove-orphans

stop:
	@docker compose --profile=app down --volumes
	@rm -rf dbdata

build:
	@go build -o build/ts main.go

lint:
	@golangci-lint run

lint-fix:
	@golangci-lint run --fix

print-results:
	@docker compose --profile=app logs ts --no-log-prefix