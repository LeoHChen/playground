package main

import (
	"bytes"
//	"encoding/hex"
	"crypto/sha256"
	"encoding/json"
	"flag"
	"fmt"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/btcsuite/btcd/wire"
	"github.com/cbergoon/merkletree"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"
	//	"github.com/btcsuite/btcutil"
)

type Tx struct {
	Hash []byte `json:"hash"`
}

// Reverse returns its argument string reversed rune-wise left to right.
func (t Tx) Reverse() Tx {
	r := t.Hash
//	n := make([]rune, l)
	l := len(r)
   n := make([]byte, l)

	for i := 0; i < l-1; i = i + 2 {
		n[l-i-2], n[l-i-1] = r[i], r[i+1]
	}
//	t.Hash = string(n)
   t.Hash = n
	return t
}

func (t Tx) CalculateHash() ([]byte, error) {
	s1 := sha256.Sum256([]byte(t.Hash))
	s2 := sha256.Sum256(s1[:])

	return s2[:], nil
}

func (t Tx) Equals(other merkletree.Content) (bool, error) {
	return string(t.Hash) == string(other.(Tx).Hash), nil
}

// block structure returned by https://blockchain.info/rawblock/
type BlockRaw struct {
	Hash       string `json:"hash"`
	PrevBlock  string `json:"prev_block"`
	MerkleRoot string `json:"mrkl_root"`
	Ver        int32  `json:"ver"`
	Height     uint64 `json:"height:`
	Bits       uint32 `json:"bits"`
	Nonce      uint32 `json:"nonce"`
	Fees       uint32 `json:"fee"`
	NumTxs     uint32 `json:"n_tx"`
	Size       uint32 `json:"size"`
	Weight     uint32 `json:"weight"`
	Time       int64  `json:"time"`
	Main       bool   `json:"main_chain"`
	Index      uint32 `json:"block_index"`
	Txs        []Tx   `json:"tx"`
}

func getBlockRaw(ep string, hash string) ([]byte, error) {
	url := fmt.Sprintf("%s/%v", ep, hash)

	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	data, _ := ioutil.ReadAll(resp.Body)
	return data, err
}

func createHeader(block BlockRaw) ([]byte, error) {
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
		Timestamp:  time.Unix(block.Time, 0),
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

func getShaSha(data []byte) [32]byte {
	s1 := sha256.Sum256(data)
	return sha256.Sum256(s1[:])
}

var (
	buf    bytes.Buffer
	logger = log.New(&buf, "logger: ", log.Lshortfile)
)

func main() {
	endpoint := flag.String("endpoint", "https://blockchain.info/rawblock/", "end point of the query")
	blockhash := flag.String("blockhash", "0000000000000000000fa9aab97fac410916c4b5e7e1edd6e401b73e812f8510", "block hash")
	//	output := flag.String("output", "", "output file of the block header")
	debug := flag.Bool("debug", false, "debug output")

	flag.Parse()

	if *debug {
		logger.Printf("endpoint:%s", *endpoint)
		logger.Printf("blockhash:%v", *blockhash)
	}

	data, err := getBlockRaw(*endpoint, *blockhash)
	if err != nil {
		logger.Fatal("getBlock failed: ", err)
	}

	var block BlockRaw
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

		var list []merkletree.Content
		for _, tx := range block.Txs {
			list = append(list, tx.Reverse())
		}

		t, err := merkletree.NewTree(list)

		if err != nil {
			log.Fatal(err)
		}
		/*
			vt, err := t.VerifyTree()
			if err != nil {
				log.Fatal(err)
			}

			vc, err := t.VerifyContent(list[0])
			if err != nil {
				log.Fatal(err)
			}

			log.Println("Verify Content: ", vc)
			log.Println("Verify Tree: ", vt)
			log.Println("The Tree: ", t)
		*/
		log.Printf("The MerkleRoot: %x", t.MerkleRoot())

		header, err := createHeader(block)

		if err != nil {
			logger.Printf("create header error: %v\n", err)
		} else {
			logger.Printf("header: %x\n", header)
			sha := getShaSha(header)
			logger.Printf("sha: %x\n", sha)
		}
	}

	fmt.Print(&buf)
}
