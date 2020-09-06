SRC=btcrelay.go
BIN=btcrelay

btcrelay: $(SRC)
	go build -o $(BIN) btcrelay.go

fmt: $(SRC)
	go fmt btcrelay.go

test: $(BIN)
	./$(BIN)

clean:
	rm -f $(BIN)
