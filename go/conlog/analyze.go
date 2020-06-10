package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
)

type Message struct {
	Time  string `json:time`
	Msg   string `json:msg`
	Block int64  `json:block`
}

func parseLog(logfile string) {
	file, err := os.Open(logfile)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)

	var blockNum int64
	for scanner.Scan() {
		line := strings.NewReader(scanner.Text())
		dec := json.NewDecoder(line)
		var m Message
		err := dec.Decode(&m)
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}
		if m.Block == 0 {
			m.Block = blockNum
		} else {
			blockNum = m.Block
		}
		fmt.Printf("%v:%s => %s\n", m.Block, m.Time, m.Msg)
	}
}

func main() {
	file := flag.String("logfile", "consensus.log", "timestamp file")

	flag.Parse()

	parseLog(*file)
}
