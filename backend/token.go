package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"path"
	"strconv"

	"backend/config"
	"backend/types"
)

const (
	TotalSupplyAbiData = "0x18160ddd"
	DecimalsAbiData    = "0x313ce567"
	SymbolAbiData      = "0x95d89b41"
)

func HandleToken(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "Method is not supported.", http.StatusNotFound)
		return
	}
	u, _ := url.Parse(r.RequestURI)
	tokenAddress := path.Base(u.Path)
	values, _ := url.ParseQuery(u.RawQuery)
	pageNumber := values.Get("page")
	accountAddress := values.Get("address")
	out := GetTokenInfo(tokenAddress, accountAddress, pageNumber)
	_, _ = fmt.Fprintf(w, string(out))
}

func HandleTokens(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "Method is not supported.", http.StatusNotFound)
		return
	}
	u, _ := url.Parse(r.RequestURI)
	symbol := path.Base(u.Path)
	out := GetTokensBySymbol(symbol)
	_, _ = fmt.Fprintf(w, string(out))
}

func GetTokensBySymbol(symbol string) []byte {
	key := config.BuildERC20sKey(symbol)
	value, err := server.DB.Get(key, nil)
	if err != nil || value == nil {
		return []byte("")
	} else {
		return value
	}
}

type TokenInfo struct {
	Symbol          string
	Decimals        string
	MaxSupply       string
	ContractAddress string
	Transactions    []types.TransactionBasicInfo
}

type simpleTokenInfo struct {
	Symbol   string
	Decimals string
}

func GetTokenInfo(tokenAddress, accountAddress, page string) []byte {
	t := TokenInfo{ContractAddress: tokenAddress}
	key := config.BuildERC20InfoKey(tokenAddress)
	value, err := server.DB.Get(key, nil)
	if err != nil {
		t.Decimals = GetErc20Info(tokenAddress, DecimalsAbiData)
		t.Symbol = GetErc20Info(tokenAddress, SymbolAbiData)
		if t.Symbol == "" {
			return nil
		} else {
			t.Symbol = t.Symbol[2+64*2:]
			tmp := t.Symbol
			result := ""
			for i := 0; i < len(t.Symbol)/2; i++ {
				d, _ := strconv.ParseInt(tmp[:2], 16, 8)
				if d == 0 {
					break
				}
				result += string(rune(d))
				tmp = tmp[2:]
			}
			t.Symbol = result //fast string abi decode
		}
		s := simpleTokenInfo{
			Symbol:   t.Symbol,
			Decimals: t.Decimals,
		}
		b, _ := json.Marshal(s)
		_ = server.DB.Put(key, b, nil)
	} else {
		s := simpleTokenInfo{}
		_ = json.Unmarshal(value, &s)
		t.Decimals = s.Decimals
		t.Symbol = s.Symbol
	}
	t.MaxSupply = GetErc20Info(tokenAddress, TotalSupplyAbiData)
	insertErc20AddrToSymbolTable(t.Symbol, tokenAddress)
	var logs []Log
	if accountAddress != "" {
		logs = GetErc20HistoryTransferLogsByAccount(tokenAddress, accountAddress)
	} else {
		logs = GetErc20HistoryTransferLogs(tokenAddress)
	}
	t.Transactions = GetErc20TransferTxsFromLogs(logs, page)
	out, _ := json.MarshalIndent(t, "", "    ")
	return out
}

func insertErc20AddrToSymbolTable(symbol, tokenAddress string) {
	key := config.BuildERC20sKey(symbol)
	var table []string
	store := true
	value, err := server.DB.Get(key, nil)
	if err != nil {
		_ = json.Unmarshal(value, &table)
		for _, a := range table {
			if a == tokenAddress {
				store = false
			}
		}
	}
	if store {
		table = append(table, tokenAddress)
		b, _ := json.Marshal(table)
		_ = server.DB.Put(key, b, nil)
	}
}

type CallObject struct {
	From     string //optional
	To       string
	gas      string //optional
	gasPrice string //optional
	value    string //optional
	Data     string //optional
}

func GetErc20Info(address string, data string) string {
	r := types.RequestInfo{
		Jsonrpc: "2.0",
		Method:  "eth_call",
		Params: []interface{}{CallObject{
			From: address, //todo: from is not need
			To:   address,
			Data: data,
		}, "latest"},
		Id: 1,
	}
	b, _ := json.Marshal(r)
	//fmt.Println(string(b))
	resp, _ := http.Post(config.NodeUrl, "application/json", bytes.NewReader(b))
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	res := types.ResponseInfo{}
	err := json.Unmarshal(body, &res)
	//fmt.Println(string(body))
	if err != nil {
		panic(err)
	}
	var t string
	_ = json.Unmarshal(res.Result, &t)
	return t
}

type FilterOption struct {
	FromBlock string     //set 1
	toBlock   string     //not set mean latest
	Address   string     //erc20 token address
	Topics    [][]string //not set, matches any topic list
	blockHash string     //not set
}

func GetErc20HistoryTransferLogs(address string) []Log {
	r := types.RequestInfo{
		Jsonrpc: "2.0",
		Method:  "eth_getLogs",
		Params: []interface{}{
			FilterOption{
				FromBlock: "0x1",
				Address:   address,
				Topics: [][]string{
					//Transfer(address indexed from,address indexed to,uint value)
					{"0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef"},
				},
			},
		},
		Id: 1,
	}
	b, _ := json.Marshal(r)
	resp, _ := http.Post(config.NodeUrl, "application/json", bytes.NewReader(b))
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	res := types.ResponseInfo{}
	err := json.Unmarshal(body, &res)
	//.Println("GET_LOGS:" + string(body))
	if err != nil {
		panic(err)
	}
	var logs []Log
	_ = json.Unmarshal(res.Result, &logs)
	//fmt.Println("LOGS:")
	//fmt.Println(logs)
	return logs
}

//http://localhost:8080/v1/erc20/0xb4589F19e4dA21de2450d79544f082FB94167DcD?address=0x0c60b56403637dc9059fff3603a58db3d5d76d38&page=1
func GetErc20HistoryTransferLogsByAccount(tokenAddr, accountAddr string) []Log {
	r := types.RequestInfo{
		Jsonrpc: "2.0",
		Method:  "moe_queryLogs",
		Params: []interface{}{
			tokenAddr,
			[]string{
				//Transfer(address indexed from,address indexed to,uint value)
				"0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef",
				"0x000000000000000000000000" + accountAddr[2:],
			},
			"0x1",
			"latest",
		},
		Id: 1,
	}
	b, _ := json.Marshal(r)
	resp, _ := http.Post(config.NodeUrl, "application/json", bytes.NewReader(b))
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	res := types.ResponseInfo{}
	err := json.Unmarshal(body, &res)
	//fmt.Println(string(body))
	if err != nil {
		panic(err)
	}
	var logs []Log
	_ = json.Unmarshal(res.Result, &logs)
	//fmt.Println(logs)
	return logs
}

//return all erc20 events, tx may repeat,not cache
func GetErc20TransferTxsFromLogs(logs []Log, page string) []types.TransactionBasicInfo {
	l := len(logs)
	p, err := strconv.Atoi(page)
	if err != nil {
		p = 0
	}
	if p < 0 {
		p = 0
	}
	if p*25 > l+25 {
		return nil
	}
	end := p * 25
	if end > l {
		end = l
	}
	if p == 0 {
		p = 1
		end = l
	}
	infos := make([]types.TransactionBasicInfo, len(logs))
	for i, log := range logs[p*25-25 : end] {
		t := GetTransaction(log.TransactionHash)
		infos[i] = GetTransactionBasicInfoFromTx(t)
	}
	return infos
}
