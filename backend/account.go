package main

import (
	"bytes"
	"encoding/binary"
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

type accountInfo struct {
	Balance      string
	Transactions []types.TransactionBasicInfo
}

var (
	AccountInfoKey = []byte{0x01}
)

func buildAccountInfoKey(address string, page int) []byte {
	b := [2]byte{}
	binary.BigEndian.PutUint16(b[:], uint16(page))
	return append(append(AccountInfoKey, address...), b[:]...)
}

func HandleAccount(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "Method is not supported.", http.StatusNotFound)
		return
	}
	u, _ := url.Parse(r.RequestURI)
	address := path.Base(u.Path)
	fmt.Println(address)
	values, _ := url.ParseQuery(u.RawQuery)
	from := values.Get("from")
	to := values.Get("to")
	pageNumber := values.Get("page")
	out := GetAccountInfo(address, from, to, pageNumber)
	_, _ = fmt.Fprintf(w, string(out))
}

func GetAccountInfo(address string, from string, to string, page string) []byte {
	var info accountInfo
	//get balance
	info.Balance = getBalance(address)
	//get transactions
	info.Transactions = getTransactionsByPage(address, from, to, page)
	out, _ := json.MarshalIndent(info, "", "    ")
	return out
}

func getBalance(address string) string {
	r := types.RequestInfo{
		Jsonrpc: "2.0",
		Method:  "eth_getBalance",
		Params:  []interface{}{address, "latest"},
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
	balance := ""
	_ = json.Unmarshal(res.Result, &balance)
	fmt.Println(balance)
	return balance
}

func getTransactionsByPage(address, from, to, page string) (out []types.TransactionBasicInfo) {
	p, err := strconv.Atoi(page)
	if err != nil {
		return nil
	}
	if p == 0 {
		p = 1
	}
	if p == 1 {
		//always get newest info from node when page = 1
		txs := getTransactions(address, from, to)
		if len(txs) <= 25 {
			out = txs
		} else {
			out = txs[:25]
		}
		_ = cacheAndSplitTxs(address, txs, 0)
	} else {
		key := buildAccountInfoKey(address, p)
		value, err := server.Cache.Get(key)
		if err == nil {
			_ = json.Unmarshal(value, &out)
		} else {
			txs := getTransactions(address, from, to)
			out = cacheAndSplitTxs(address, txs, p)
		}
	}
	return out
}

/*
1. check if txs and cache is same
2. split txs into pages with 15 entries per page
3. cache pages 5min
4. return specific page txs
page number is 1, 2, 3, 4...
*/
func cacheAndSplitTxs(address string, txs []types.TransactionBasicInfo, page int) []types.TransactionBasicInfo {
	if len(txs) == 0 {
		return nil
	}
	setCache := true
	value, err := server.Cache.Get(buildAccountInfoKey(address, 1))
	if err == nil {
		b, _ := json.Marshal(txs[0])
		if bytes.Equal(b, value) {
			setCache = false
		}
	}
	if setCache {
		for i := 0; i <= (len(txs)-1)/25; i++ {
			var v []byte
			if i == (len(txs)-1)/25 {
				v, _ = json.Marshal(txs)
			} else {
				v, _ = json.Marshal(txs[:25])
				txs = txs[25:]
			}
			err := server.Cache.Set(buildAccountInfoKey(address, i), v, 5*60)
			if err != nil {
				break
			}
		}
	}
	var out []types.TransactionBasicInfo
	if page != 0 {
		value, err := server.Cache.Get(buildAccountInfoKey(address, page))
		if err != nil {
			_ = json.Unmarshal(value, &out)
		}
	}
	return out
}

func getTransactions(address, from, to string) (out []types.TransactionBasicInfo) {
	r := types.RequestInfo{}
	if from == "true" && to == "true" {
		r = types.RequestInfo{
			Jsonrpc: "2.0",
			Method:  "moe_queryTxByAddr",
			Params:  []interface{}{address, "0x1", "latest"},
			Id:      1,
		}
	} else if to == "true" {
		r = types.RequestInfo{
			Jsonrpc: "2.0",
			Method:  "moe_queryTxByDst",
			Params:  []interface{}{address, "0x1", "latest"},
			Id:      1,
		}
	} else {
		r = types.RequestInfo{
			Jsonrpc: "2.0",
			Method:  "moe_queryTxBySrc",
			Params:  []interface{}{address, "0x1", "latest"},
			Id:      1,
		}
	}
	b, _ := json.Marshal(r)
	resp, _ := http.Post(config.NodeUrl, "application/json", bytes.NewReader(b))
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	res := types.ResponseInfo{}
	fmt.Println(string(body))
	err := json.Unmarshal(body, &res)
	if err != nil {
		panic(err)
	}
	fmt.Println(string(res.Result))
	var allTxs []types.TransactionBasicInfo
	_ = json.Unmarshal(res.Result, &allTxs)
	return allTxs
}
