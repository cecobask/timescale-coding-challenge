.PHONY: *

start:
	docker compose up -d

stop:
	docker compose down --volumes
	rm -rf dbdata
