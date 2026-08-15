run:
	CONFIG_PATH=./config/config.yaml go run ./cmd/server
swagger update:
	swag init -g cmd/server/main.go --parseInternal --parseDependency
docker-build:
	docker build -t media-processing-platform .
docker-run:
	docker run --rm \
		-p 8000:8000 \
		-e CONFIG_PATH=/app/config/config.yaml \
		-v "$(PWD)/config/config.yaml:/app/config/config.yaml:ro" \
		media-processing-platform
local-postgresql:
	sudo service postgresql start
docker-compose-build:
	docker compose -f deployments/docker-compose.yml build
docker-compose-up:
	docker compose -f deployments/docker-compose.yml up
docker-compose-down:
	docker compose -f deployments/docker-compose.yml down