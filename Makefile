GOOSE := go run github.com/pressly/goose/v3/cmd/goose@v3.27.3

.PHONY: migrate-up migrate-down migrate-status migrate-reset migrate-create

migrate-up:
	$(GOOSE) up

migrate-down:
	$(GOOSE) down

migrate-status:
	$(GOOSE) status

migrate-create:
	$(GOOSE) create $(name) sql
