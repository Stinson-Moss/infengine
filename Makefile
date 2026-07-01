OUT=infengine
FEED_PATH="./feeds.txt"
BIN=bin

build:
	go build -o $(BIN)/$(OUT)

run: build
	./$(BIN)/$(OUT) $(FEED_PATH)