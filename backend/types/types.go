package types

import (
	"encoding/json"
)

type RequestInfo struct {
	Jsonrpc string
	Method  string
	Params  []interface{}
	Id      uint
}

type ErrorInfo struct {
	Code    uint
	Message string
}

type ResponseInfo struct {
	Jsonrpc string
	Id      uint
	Result  json.RawMessage
	Error   ErrorInfo
}

type Block struct {
	Number            string
	Hash              string
	ParentHash        string
	Miner             string
	Size              string
	GasLimit          uint
	GasUsed           string
	Timestamp         string
	TransactionsCount string
	BlockReward       string
	StateRoot         string
	TransactionsRoot  string
	Transactions      []string
}

type TransactionBasicInfo struct {
	Hash        string
	BlockNumber string
	From        string
	To          string
	Age         string
	Value       string
}
