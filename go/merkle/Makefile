SRC=merkle.go
BIN=merkle

$(BIN): $(SRC)
	go build -o $(BIN) $(SRC)

fmt: $(SRC)
	go fmt $(SRC)

test: $(BIN)
	./$(BIN)

clean:
	rm -f $(BIN)
