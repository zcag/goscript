BIN := goscript

.PHONY: test build run install clean

test:
	go test ./...

build:
	go build -o $(BIN)

run:
	./$(BIN)

install:
	install -Dm755 $(BIN) $$HOME/bin/$(BIN)

clean:
	rm -f $(BIN)
