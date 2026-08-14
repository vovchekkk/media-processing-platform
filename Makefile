run:
	CONFIG_PATH=./config/config.yaml go run ./cmd/server
swag update:
	swag init -g cmd/server/main.go --parseInternal --parseDependency