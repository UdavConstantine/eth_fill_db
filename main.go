package main

import (
	"fmt"
	"net/http"
	"bytes"
	"io/ioutil"
	"encoding/json"
	"log"
	"strconv"
	_ "github.com/lib/pq"
	"database/sql"
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

	var jsonStr = []byte(`{"jsonrpc":"2.0","method":"eth_getBlockByNumber","params":["`+ s + `", true],"id":5}`)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	//fmt.Println("response Status:", resp.Status)
	//fmt.Println("response Headers:", resp.Header)
	body, _ := ioutil.ReadAll(resp.Body)
	//fmt.Println("response Body:", string(body))

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


func insertBlock(db *sql.DB, block BlockStruct) {
	_, err := db.Exec(
	`insert into blocks(
		number,
		hash,
		parentHash,
		nonce,
		sha3Uncles,
		logsBloom,
		transactionsRoot,
		stateRoot,
		receiptsRoot,
		miner,
		difficulty,
		totalDifficulty,
		size,
		proofOfAuthorityData,
		gasLimit,
		gasUsed,
		timestamp,
		mixhash) values($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18)`,
	    hextoint(block.Result.Number),
		block.Result.Hash,
		block.Result.ParentHash,
		block.Result.Nonce,
		block.Result.Sha3Uncles,
		block.Result.LogsBloom,
		block.Result.TransactionsRoot,
		block.Result.StateRoot,
		block.Result.ReceiptsRoot,
		block.Result.Miner,
		hextoint(block.Result.Difficulty),
		hextoint(block.Result.TotalDifficulty),
		hextoint(block.Result.Size),
		block.Result.ExtraData,
		hextoint(block.Result.GasLimit),
		hextoint(block.Result.GasUsed),
		hextoint(block.Result.Timestamp),
		block.Result.MixHash)
	if err != nil {
		panic(err)
	}
	for _, tr := range block.Result.Transactions {
		_, err := db.Exec(
			`insert into transactions(
			hash,
			nonce ,
			blockHash,
			blockNumber ,
			transactionIndex,
			"from",
			"to",
			"value",
			gas,
			gasPrice,
			input,
			v,
			r,
			s) values($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14)`,
			tr.Hash,
			hextoint(tr.Nonce),
			tr.BlockHash,
			hextoint(tr.BlockNumber),
			hextoint(tr.TransactionIndex),
			tr.From,
			tr.To,
			hextoint(tr.Value),
			hextoint(tr.Gas),
			hextoint(tr.GasPrice),
			tr.Input,
			tr.V,
			tr.R,
			tr.S)
		if err != nil {
			panic(err)
		}
	}
}


func main() {
	n := getLastBlockNumber()
	fmt.Println(n)
	//block := getBlock(inttohex(n))
	connStr := "user=ethtodb password=ethtodb dbname=eth sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		panic(err)
	}
	defer db.Close()
	//insertBlock(db, block)
	var i uint64
	for i =3024000; i < n; i++ {
		insertBlock(db,  getBlock(inttohex(i)))
	}
}
