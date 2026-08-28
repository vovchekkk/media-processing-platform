swagger update:
	swag init -g server/cmd/main.go --parseInternal --parseDependency
local-postgresql:
	sudo service postgresql start
docker-compose-build:
	docker compose --env-file .env -f deployments/docker-compose.yml build
docker-compose-up:
	docker compose --env-file .env -f deployments/docker-compose.yml up
docker-compose-down:
	docker compose --env-file .env -f deployments/docker-compose.yml down
docker-compose-test:
	docker compose --env-file .env --profile test -f deployments/docker-compose.yml run --rm tests
proto:
	protoc -I=. \
           --go_out=. --go_opt=paths=source_relative \
           pkg/proto/task.proto