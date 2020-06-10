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
	"time"
)

type Message struct {
	Time  string `json:time`
	Msg   string `json:msg`
	Block int64  `json:block`
}

type Stage byte

const (
	Propose Stage = iota
	EnoughPrepare
	EnoughCommit
	StartGrace
	EndGrace
	Hooray
)

type Consensus struct {
	Timestamp time.Time
	Stage     Stage
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

		var consensus Consensus

		t, err := time.Parse(time.RFC3339, m.Time)
		if err != nil {
			log.Fatal("Invalid timestamp")
		}
		consensus.Timestamp = t
		if strings.Contains(m.Msg, "PROPOSING") {
			consensus.Stage = Propose
		}
		if strings.Contains(m.Msg, "Enough Prepare") {
			consensus.Stage = EnoughPrepare
		}
		if strings.Contains(m.Msg, "Enough commits") {
			consensus.Stage = EnoughCommit
		}
		if strings.Contains(m.Msg, "Starting Grace") {
			consensus.Stage = StartGrace
		}
		if strings.Contains(m.Msg, "Commit Grace") {
			consensus.Stage = EndGrace
		}
		if strings.Contains(m.Msg, "HOORAY") {
			consensus.Stage = Hooray
		}

		fmt.Printf("%v:%s => %s\n", m.Block, m.Time, m.Msg)
	}
}

func main() {
	file := flag.String("logfile", "consensus.log", "timestamp file")

	flag.Parse()

	parseLog(*file)
}
