package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"path"

	"backend/config"
	"backend/types"
)

type Transaction struct {
	Hash             string
	Nonce            string
	BlockNumber      string
	TransactionIndex string
	From             string
	To               string
	Value            string
	Gas              string
	GasPrice         string
	Input            string
}

type TransactionInfo struct {
	Hash             string
	Nonce            string
	BlockNumber      string
	TransactionIndex string
	From             string
	To               string
	Value            string
	Gas              string
	GasPrice         string
	Input            string
	Logs             []Log
}

type Log struct {
	Address         string
	TransactionHash string
	Topics          []string
	Data            string
}

func HandleTransactionInfo(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "Method is not supported.", http.StatusNotFound)
		return
	}
	u, _ := url.Parse(r.RequestURI)
	hash := path.Base(u.Path)
	out := GetTransactionInfo(hash)
	_, _ = fmt.Fprintf(w, string(out))
}

func GetTransactionInfo(hash string) []byte {
	t := GetTransaction(hash)
	info := TransactionInfo{
		Hash:             t.Hash,
		Nonce:            t.Nonce,
		BlockNumber:      t.BlockNumber,
		TransactionIndex: t.TransactionIndex,
		From:             t.From,
		To:               t.To,
		Value:            t.Value,
		Gas:              t.Gas,
		GasPrice:         t.GasPrice,
		Input:            t.Input,
	}
	r := GetTransactionReceipt(hash)
	info.Logs = r.Logs
	out, _ := json.MarshalIndent(info, "", "    ")
	return out
}

func GetTransaction(hash string) Transaction {
	r := types.RequestInfo{
		Jsonrpc: "2.0",
		Method:  "eth_getTransactionByHash",
		Params:  []interface{}{hash},
		Id:      1,
	}
	b, _ := json.Marshal(r)
	resp, _ := http.Post(config.NodeUrl, "application/json", bytes.NewReader(b))
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	res := types.ResponseInfo{}
	err := json.Unmarshal(body, &res)
	if err != nil {
		panic(err)
	}
	t := Transaction{}
	_ = json.Unmarshal(res.Result, &t)
	fmt.Println(t)
	return t
}

type Receipt struct {
	GasUsed         string
	ContractAddress string
	Logs            []Log
	Status          string
}

func GetTransactionReceipt(hash string) Receipt {
	r := types.RequestInfo{
		Jsonrpc: "2.0",
		Method:  "eth_getTransactionReceipt",
		Params:  []interface{}{hash},
		Id:      1,
	}
	b, _ := json.Marshal(r)
	resp, _ := http.Post(config.NodeUrl, "application/json", bytes.NewReader(b))
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	res := types.ResponseInfo{}
	err := json.Unmarshal(body, &res)
	if err != nil {
		panic(err)
	}
	receipt := Receipt{}
	_ = json.Unmarshal(res.Result, &receipt)
	fmt.Println(receipt)
	return receipt
}
