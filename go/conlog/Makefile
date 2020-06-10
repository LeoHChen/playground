BIN:=analyze

bin: analyze.go
	go build -o $(BIN) analyze.go

clean:
	goimports -w analyze.go
	rm -f $(BIN)

DIR=.
analysis: $(DIR)/consensus.*.log
	for file in $^ ; do \
		./jj.sh $${file}; \
	done

all: bin

.PHONY: bin clean
