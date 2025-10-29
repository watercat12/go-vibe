run:
	air -c .air.toml

worker:
	go run ./cmd/worker

local-db:
# 	docker-compose --env-file ./.env -f docker-compose.yml down
	docker compose --env-file ./.env -f docker-compose.yml up -d

db/migrate:
	go run ./cmd/migrate

gen-swagger:
	swag init -g cmd/api/main.go -o cmd/api/docs

gen-mock:
	mockery

TEST_PATH ?= ./internal/...

unit-test:
	@mkdir -p coverage
	-go test -coverprofile=coverage/coverage.txt.tmp -count=1 $(TEST_PATH)
	@cat coverage/coverage.txt.tmp | grep -v "mock_" > coverage/coverage.txt
	@go tool cover -html=coverage/coverage.txt -o coverage/index-application.html
split-reports:
	bash ./split-html.sh

get-package-change:
	bash ./get-package-change-path.sh