package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"sort"
	"strings"
	"time"
)

type Message struct {
	Time  string `json:time`
	Msg   string `json:msg`
	Block int    `json:block`
}

/*
type Stage byte

const (
	Propose Stage = iota
	EnoughPrepare
	EnoughCommit
	StartGrace
	EndGrace
	Hooray
)
*/

type Consensus struct {
	Propose       time.Time
	EnoughPrepare time.Time
	EnoughCommit  time.Time
	StartGrace    time.Time
	EndGrace      time.Time
	Hooray        time.Time
}

func parseLog(logfile string) {
	file, err := os.Open(logfile)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)

	var blockNum int
	var consensusMap = make(map[int]*Consensus)

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

		var consensus *Consensus
		var ok bool
		consensus, ok = consensusMap[m.Block]
		if !ok {
			consensusMap[m.Block] = new(Consensus)
		}
		consensus, _ = consensusMap[m.Block]

		t, err := time.Parse(time.RFC3339, m.Time)
		if err != nil {
			log.Fatal("Invalid timestamp")
		}
		if strings.Contains(m.Msg, "PROPOSING") {
			consensus.Propose = t
			continue
		}
		if strings.Contains(m.Msg, "Enough Prepare") {
			consensus.EnoughPrepare = t
			continue
		}
		if strings.Contains(m.Msg, "Enough commits") {
			consensus.EnoughCommit = t
			continue
		}
		if strings.Contains(m.Msg, "Starting Grace") {
			consensus.StartGrace = t
			continue
		}
		if strings.Contains(m.Msg, "Commit Grace") {
			consensus.EndGrace = t
			continue
		}
		if strings.Contains(m.Msg, "HOORAY") {
			consensus.Hooray = t
			continue
		}
	}

	var blocks []int
	for k := range consensusMap {
		blocks = append(blocks, k)
	}
	sort.Ints(blocks)

	fmt.Println("block, consensus, prepare, commit, grace, finalize")
	for _, k := range blocks {
		c := consensusMap[k]
		fmt.Printf("%v, %v, %v, %v, %v, %v\n", k, c.Hooray.Sub(c.Propose), c.EnoughPrepare.Sub(c.Propose), c.EnoughCommit.Sub(c.EnoughPrepare), c.EndGrace.Sub(c.StartGrace), c.Hooray.Sub(c.EndGrace))
	}
}

func main() {
	file := flag.String("logfile", "consensus.log", "timestamp file")

	flag.Parse()

	parseLog(*file)
}
