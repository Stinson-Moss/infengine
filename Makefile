OUT=infengine
FEED_PATH="./feeds.txt"
BIN=bin

build:
	go build -o $(BIN)/$(OUT)

run: build
	./$(BIN)/$(OUT) $(FEED_PATH)

reset:
	docker-compose down -v
	docker-compose up -d

migrate:
	migrate -path db/postgres/migrations -database "postgres://some_user:1234@localhost:5432/documentdb?sslmode=disable" up