.PHONY: run test local-db lint db/migrate

run:
	air -c ./tools/.air.toml

test:
	go clean -testcache
	@(go run gotest.tools/gotestsum@latest \
	  --format pkgname \
	  -- -cover $$(go list ./... | grep -v -E "(cmd|testutil|tmp|mocks)"))

local-dev:
	docker compose --env-file ./.env -f ./tools/compose/docker-compose.yml down
	docker compose --env-file ./.env -f ./tools/compose/docker-compose.yml up -d

lint:
	golangci-lint run

db/migrate:
	go run ./cmd/migrate

mock-gen:
	docker run -v $(PWD):/src -w /src vektra/mockery:3