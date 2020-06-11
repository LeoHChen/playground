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

type Consensus struct {
	Propose       time.Time
	EnoughPrepare time.Time
	EnoughCommit  time.Time
	StartGrace    time.Time
	EndGrace      time.Time
	OneHundred    time.Time
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
		if strings.Contains(m.Msg, "2/3 Enough commits") {
			consensus.EnoughCommit = t
			continue
		}
		if strings.Contains(m.Msg, "100% Enough commits") {
			consensus.OneHundred = t
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

	fmt.Println("block, consensus, prepare, commit, grace, 100%, finalize")
	for _, k := range blocks {
		c := consensusMap[k]
		var conT, preT, comT, hunT, graT, finT int64
		if !c.Hooray.IsZero() {
			conT = c.Hooray.Sub(c.Propose).Milliseconds()
		}
		preT = c.EnoughPrepare.Sub(c.Propose).Milliseconds()
		comT = c.EnoughCommit.Sub(c.EnoughPrepare).Milliseconds()
		if !c.OneHundred.IsZero() {
			hunT = c.OneHundred.Sub(c.EnoughCommit).Milliseconds()
		}
		if !c.EndGrace.IsZero() {
			graT = c.EndGrace.Sub(c.StartGrace).Milliseconds()
		}
		finT = c.Hooray.Sub(c.EndGrace).Milliseconds()

		over := conT * preT * comT * graT * finT

		if k == 0 || over == 0 {
			continue
		}

		fmt.Printf("%v, %v, %v, %v, %v, %v, %v\n", k, conT, preT, comT, graT, hunT, finT)
	}
}

func main() {
	file := flag.String("logfile", "consensus.log", "timestamp file")

	flag.Parse()

	parseLog(*file)
}
