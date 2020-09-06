package main

import (
	"bytes"
	// "encoding/binary"
	"crypto/sha256"
	"encoding/json"
	"flag"
	"fmt"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/btcsuite/btcd/wire"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"
	//	"github.com/btcsuite/btcutil"
)

// block structure returned by https://api.blockcypher.com/v1/btc/main/blocks/
// document https://www.blockcypher.com/dev/bitcoin/#block-height-endpoint
type Block struct {
	Hash       string    `json:"hash"`
	Ver        int32     `json:"ver"`
	Height     uint64    `json:"height:`
	Chain      string    `json:"chain"`
	Total      uint64    `json:"total"`
	Fees       uint32    `json:"fees"`
	Size       uint32    `json:"size"`
	Time       time.Time `json:"time"`
	PrevBlock  string    `json:"prev_block"`
	MerkleRoot string    `json:"mrkl_root"`
	Bits       uint32    `json:"bits"`
	Nonce      uint32    `json:"nonce"`
	NumTxs     uint32    `json:"n_tx"`
}

func getBlock(ep string, height uint64) ([]byte, error) {
	url := fmt.Sprintf("%s/%v?txstart=1&limit=1", ep, height)

	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	data, _ := ioutil.ReadAll(resp.Body)
	return data, err
}

func createHeader(block Block) ([]byte, error) {
	prevBlock, err := chainhash.NewHashFromStr(block.PrevBlock)
	if err != nil {
		return nil, err
	}

	merkleRoot, err := chainhash.NewHashFromStr(block.MerkleRoot)
	if err != nil {
		return nil, err
	}

	var header = wire.BlockHeader{
		Version:    block.Ver,
		PrevBlock:  *prevBlock,
		MerkleRoot: *merkleRoot,
		Timestamp:  block.Time,
		Bits:       block.Bits,
		Nonce:      block.Nonce,
	}

	buf := new(bytes.Buffer)
	err = header.Serialize(buf)

	return buf.Bytes(), err
}

func writeToFile(filename string, data []byte) {
	file, err := os.OpenFile(
		filename, os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0666,
	)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	written, err := file.Write(data)
	if err != nil {
		log.Fatal(err)
	}
	logger.Printf("written: %v\n", written)
}

func getShaSha(data []byte) {
	s1 := sha256.Sum256(data)
	fmt.Printf("s1:%x\n", s1)
	s2 := sha256.Sum256(s1[:])
	fmt.Printf("s2:%x\n", s2)
}

var (
	buf    bytes.Buffer
	logger = log.New(&buf, "logger: ", log.Lshortfile)
)

func main() {
	endpoint := flag.String("endpoint", "https://api.blockcypher.com/v1/btc/main/blocks/", "end point of the query")
	blockheight := flag.Uint64("blockheight", 646880, "block height")
	output := flag.String("output", "", "output file of the block header")
	debug := flag.Bool("debug", true, "debug output")

	flag.Parse()

	if *debug {
		logger.Printf("endpoint:%s", *endpoint)
		logger.Printf("blockheight:%v", *blockheight)
	}

	data, err := getBlock(*endpoint, *blockheight)
	if err != nil {
		logger.Fatal("getBlock failed: ", err)
	}

	var block Block
	err = json.Unmarshal(data, &block)
	if err != nil {
		logger.Printf("error decoding response: %v", err)
		if e, ok := err.(*json.SyntaxError); ok {
			logger.Printf("syntax error at byte offset %d", e.Offset)
		}
		logger.Printf("response: %q", data)
	} else {
		if *debug {
			pretty, err := json.MarshalIndent(block, "", "  ")
			if err != nil {
				log.Fatal("failed to generate json", err)
			}
			fmt.Printf("block: %s\n", string(pretty))
		}

		header, err := createHeader(block)

		if err != nil {
			logger.Printf("create header error: %v\n", err)
		} else {
			logger.Printf("header: %x\n", header)
			getShaSha(header)

			if len(*output) > 0 {
				writeToFile(*output, header)
			}
		}
	}

	fmt.Print(&buf)
}
