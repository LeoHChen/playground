BIN:=analyze

bin: analyze.go
	go build -o $(BIN) analyze.go

clean:
	goimports -w analyze.go
	rm -f $(BIN)

all: bin

.PHONY: bin clean
