run-app:
	@echo "Running app"
	@go run cmd/main.go

run-app-logs:
	@echo "Running app export archive logs"
	@go run cmd/main.go > ~/tmp/app/customers-ms.log 2>&1

unit:
	@go test ./tests/unit/... --tags=unit -v

unit-verbose:
	ginkgo -r --race --tags=unit --randomize-all --randomize-suites --fail-on-pending

unit-cover:
	@go test ./tests/unit/... -coverpkg ./internal/... --tags=unit

unit-cover-verbose:
	@go test ./tests/unit/... -coverpkg ./internal/... --tags=unit -v

unit-report:
	mkdir -p "tests/unit/coverage" \
	&& go test ./tests/unit/... -coverprofile=tests/unit/coverage/cover.out -coverpkg ./internal/... --tags=unit \
	&& go tool cover -html=tests/unit/coverage/cover.out -o tests/unit/coverage/cover.html \
	&& go tool cover -func=tests/unit/coverage/cover.out -o tests/unit/coverage/cover.functions.html

up:
	@docker compose -f docker-compose.yml up -d --build

down:
	@docker compose -f docker-compose.yml down -v

pgpool-logs:
	@echo "Running pgpool logs"
	@docker logs -f pgpool

tag:
	scripts/bump_version.sh

integration:
	@go test ./tests/integration/... --tags=integration -v -count=1

diff:
	atlas migrate diff create_users_table \
		--env local \
		--to file://atlas/ \
		--config file://atlas/atlas.hcl \
		--dev-url "docker://postgres/17.4"



.PHONY: run-app,
		unit,
		unit-cover,
		unit-report,
		integration,
		docker-dev,
		up,
		down,
		tag,
