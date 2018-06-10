package main

import (
	"fmt"
	"net/http"
	"bytes"
	"io/ioutil"
	"encoding/json"
	"log"
	"strconv"
)


func hextoint(h string) uint64 {
	var n uint64
	n, err := strconv.ParseUint(h[2:], 16, 64)
	if err != nil {
		log.Fatal("error")
	}
	return n
}


func inttohex(i uint64) string{
	s:= strconv.FormatUint(i, 16)
	return "0x" + s
}


type BlockStruct struct {
	Jsonrpc string `json:"jsonrpc"`
	ID      int    `json:"id"`
	Result  struct {
		Difficulty      string `json:"difficulty"`
		ExtraData       string `json:"extraData"`
		GasLimit        string `json:"gasLimit"`
		GasUsed         string `json:"gasUsed"`
		Hash            string `json:"hash"`
		LogsBloom       string `json:"logsBloom"`
		Miner           string `json:"miner"`
		MixHash         string `json:"mixHash"`
		Nonce           string `json:"nonce"`
		Number          string `json:"number"`
		ParentHash      string `json:"parentHash"`
		ReceiptsRoot    string `json:"receiptsRoot"`
		Sha3Uncles      string `json:"sha3Uncles"`
		Size            string `json:"size"`
		StateRoot       string `json:"stateRoot"`
		Timestamp       string `json:"timestamp"`
		TotalDifficulty string `json:"totalDifficulty"`
		Transactions    []struct {
			BlockHash        string `json:"blockHash"`
			BlockNumber      string `json:"blockNumber"`
			From             string `json:"from"`
			Gas              string `json:"gas"`
			GasPrice         string `json:"gasPrice"`
			Hash             string `json:"hash"`
			Input            string `json:"input"`
			Nonce            string `json:"nonce"`
			To               string `json:"to"`
			TransactionIndex string `json:"transactionIndex"`
			Value            string `json:"value"`
			V                string `json:"v"`
			R                string `json:"r"`
			S                string `json:"s"`
		} `json:"transactions"`
		TransactionsRoot string        `json:"transactionsRoot"`
		Uncles           []interface{} `json:"uncles"`
	} `json:"result"`
}


func getBlock(s string) BlockStruct {
	url := "https://sidechain-dev.sonm.com/"
	fmt.Println("URL:>", url)

	var jsonStr = []byte(`{"jsonrpc":"2.0","method":"eth_getBlockByNumber","params":["`+ s + `", true],"id":5}`)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	fmt.Println("response Status:", resp.Status)
	fmt.Println("response Headers:", resp.Header)
	body, _ := ioutil.ReadAll(resp.Body)
	fmt.Println("response Body:", string(body))

	var block = BlockStruct{}
	err1 := json.Unmarshal(body, &block)
	if err1 != nil {
		log.Fatal("error")
	}

	return block
}

type BlockNumber struct {
	Jsonrpc string `json:"jsonrpc"`
	ID      int    `json:"id"`
	Result  string `json:"result"`
}

func getLastBlockNumber() uint64{
	url := "https://sidechain-dev.sonm.com/"
	fmt.Println("URL:>", url)

	var jsonStr = []byte(`{"jsonrpc":"2.0","method":"eth_blockNumber","params":[],"id":5}`)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)

	var blocknumber = BlockNumber{}
	err1 := json.Unmarshal(body, &blocknumber)
	if err1 != nil {
		log.Fatal("error")
	}

	return hextoint(blocknumber.Result)
}

func main() {
	block := getBlock("0x2e2488")

	fmt.Println(block.Result.Number)
	n := hextoint(block.Result.Number)
	fmt.Println(n)
	fmt.Println(inttohex(n))
	fmt.Println(getLastBlockNumber())
}
