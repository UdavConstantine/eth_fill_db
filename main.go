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
	"math/big"
	"time"
)


func hextointdb(h string) string {
	var n big.Int
	n.SetString(h[2:], 16)
	return n.String()
}

func hextoint(h string) uint64 {
	n, err := strconv.ParseUint(h[2:], 16, 64)
	if err != nil {
		fmt.Println("Ошибка при преобразовании hex в uin64", err)
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


func getBlock(s string, sint uint64) (BlockStruct, error) {
	url := "https://sidechain-dev.sonm.com/"
	//url := "https://mainnet.infura.io/Ol4LW5vVUUUV0SrUxkzv"
	var jsonStr = []byte(`{"jsonrpc":"2.0","method":"eth_getBlockByNumber","params":["`+ s + `", true],"id":`+ strconv.FormatUint(sint, 10) +`}`)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		//log.Fatalln("Ошибка получения ответа от сервера:", err)
		log.Println("Ошибка получения ответа от сервера:", err)
		return BlockStruct{}, err
	}
	defer resp.Body.Close()

	if resp.Status != "200 OK" {
		//log.Fatalln("Статус не 200: ", resp.Status, ", Req: ", string(jsonStr))
		log.Println("Статус не 200: ", resp.Status, ", Req: ", string(jsonStr))
		return BlockStruct{}, err
	}
	//fmt.Println("response Status:", resp.Status)
	//fmt.Println("response Headers:", resp.Header)
	body, _ := ioutil.ReadAll(resp.Body)
	//fmt.Println("response Body:", string(body))

	var block = BlockStruct{}
	err1 := json.Unmarshal(body, &block)
	if err1 != nil {
		//log.Fatalln("Ошибка преобразования в json", string(body))
		log.Println("Ошибка преобразования в json", string(body))
		return BlockStruct{}, err1
	}

	return block, nil
}


type BlockNumber struct {
	Jsonrpc string `json:"jsonrpc"`
	ID      int    `json:"id"`
	Result  string `json:"result"`
}


func getLastBlockNumber() uint64{
	url := "https://sidechain-dev.sonm.com/"
	//url := "https://mainnet.infura.io/Ol4LW5vVUUUV0SrUxkzv"
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
	    hextointdb(block.Result.Number),
		block.Result.Hash,
		block.Result.ParentHash,
		block.Result.Nonce,
		block.Result.Sha3Uncles,
		block.Result.LogsBloom,
		block.Result.TransactionsRoot,
		block.Result.StateRoot,
		block.Result.ReceiptsRoot,
		block.Result.Miner,
		hextointdb(block.Result.Difficulty),
		hextointdb(block.Result.TotalDifficulty),
		hextointdb(block.Result.Size),
		block.Result.ExtraData,
		hextointdb(block.Result.GasLimit),
		hextointdb(block.Result.GasUsed),
		hextointdb(block.Result.Timestamp),
		block.Result.MixHash)
	if err != nil {
		log.Fatalln("Ошибка вставки в blocks:", err)
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
			hextointdb(tr.Nonce),
			tr.BlockHash,
			hextointdb(tr.BlockNumber),
			hextointdb(tr.TransactionIndex),
			tr.From,
			tr.To,
			hextointdb(tr.Value),
			hextointdb(tr.Gas),
			hextointdb(tr.GasPrice),
			tr.Input,
			tr.V,
			tr.R,
			tr.S)
		if err != nil {
			log.Fatalln("Ошибка вставки в transactions:", err)
		}
	}
}


func processBlock(db *sql.DB, i uint64){

	for {
		block, err := getBlock(inttohex(i), i)
		if err != nil || len(block.Result.Number) == 0 {
			time.Sleep(time.Millisecond * 100)
		} else {
			insertBlock(db,  block)
			break
		}
	}
}

const maxt = 200

func main() {
	n := getLastBlockNumber()
	//n = 1000
	fmt.Println(n)
	connStr := "user=ethtodb password=ethtodb dbname=eth sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	var i uint64 = 0
	sem := make(chan int, maxt)
	for i = 0; i < n; i++{
		sem <- 1
		go func(i uint64) {
			fmt.Println(i)
			processBlock(db, i)
			<-sem
		}(i)
	}

}
