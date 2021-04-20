package types

import (
	"encoding/json"
)

type RequestInfo struct {
	Jsonrpc string        `json:"jsonrpc"`
	Method  string        `json:"method"`
	Params  []interface{} `json:"params"`
	Id      uint          `json:"id"`
}

type ErrorInfo struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

type ResponseInfo struct {
	Jsonrpc string          `json:"jsonrpc"`
	Id      uint            `json:"id"`
	Result  json.RawMessage `json:"result"`
	Error   ErrorInfo       `json:"error"`
}

type Block struct {
	Number            string   `json:"number"`
	Hash              string   `json:"hash"`
	ParentHash        string   `json:"parentHash"`
	Miner             string   `json:"miner"`
	Size              string   `json:"size"`
	GasLimit          uint     `json:"gasLimit"`
	GasUsed           string   `json:"gasUsed"`
	Timestamp         string   `json:"timestamp"`
	TransactionsCount string   `json:"transactionsCount"`
	BlockReward       string   `json:"blockReward"`
	StateRoot         string   `json:"stateRoot"`
	TransactionsRoot  string   `json:"transactionsRoot"`
	Transactions      []string `json:"transactions"`
}

type TransactionBasicInfo struct {
	Hash        string `json:"hash"`
	BlockNumber string `json:"blockNumber"`
	From        string `json:"from"`
	To          string `json:"to"`
	Age         string `json:"age"`
	Value       string `json:"value"`
}
