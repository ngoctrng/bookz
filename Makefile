.PHONY: run test local-db lint db/migrate

run:
	air -c ./tools/.air.toml

worker:
	go run ./cmd/worker/main.go

test:
	go clean -testcache
	@(go run gotest.tools/gotestsum@latest \
	  --format pkgname \
	  -- -cover $$(go list ./... | grep -v -E "(cmd|testutil|tmp|mocks)"))

local-dev:
	docker compose --env-file ./.env -f ./tools/compose/docker-compose.yml up -d

local-clean:
	docker compose --env-file ./.env -f ./tools/compose/docker-compose.yml down

lint:
	golangci-lint run

db/migrate:
	go run ./cmd/migrate

mock-gen:
	docker run -v $(PWD):/src -w /src vektra/mockery:3

diagram:
	docker run -v $(PWD)/docs/diagrams:/work -w /work ghcr.io/plantuml/plantuml -tsvg *.puml

swagger:
	swag init --parseDependency -g httpserver/server.go -d ./internal -o ./api