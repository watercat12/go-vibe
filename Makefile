run:
	air -c .air.toml

worker:
	go run ./cmd/worker

local-db:
# 	docker-compose --env-file ./.env -f docker-compose.yml down
	docker-compose --env-file ./.env -f docker-compose.yml up -d

db/migrate:
	go run ./cmd/migrate

unit-test:
	@mkdir coverage || true
	-go test -coverprofile=coverage/coverage.txt.tmp -count=1 ./internal/adapters/handler/...
	@cat coverage/coverage.txt.tmp | grep -v "mock_" > coverage/coverage.txt
	@go tool cover -func=coverage/coverage.txt
	@go tool cover -html=coverage/coverage.txt -o coverage/index-application.html