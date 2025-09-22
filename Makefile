run:
	air -c .air.toml

local-db:
	docker-compose --env-file ./.env -f docker-compose.yml down
	docker-compose --env-file ./.env -f docker-compose.yml up -d

db/migrate:
	go run ./cmd/migrate