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

func HandleBlockTxs(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "Method is not supported.", http.StatusNotFound)
		return
	}
	u, _ := url.Parse(r.RequestURI)
	height := path.Base(u.Path)
	out := GetBlockTxs(height)
	_, _ = fmt.Fprintf(w, string(out))
}

type BlockTxs struct {
	Block        string
	Transactions []types.TransactionBasicInfo
}

func GetBlockTxs(height string) []byte {
	bTxs := BlockTxs{Block: height}
	r := types.RequestInfo{
		Jsonrpc: "2.0",
		Method:  "eth_getTransactionsByBlockNumber",
		Params:  []interface{}{height},
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
	var txs []Transaction
	_ = json.Unmarshal(res.Result, &txs)
	fmt.Println(txs)
	bTxs.Transactions = make([]types.TransactionBasicInfo, len(txs))
	for i, tx := range txs {
		bTxs.Transactions[i] = GetTransactionBasicInfoFromTx(tx)
	}
	out, _ := json.MarshalIndent(bTxs, "", "")
	return out
}

func GetTransactionBasicInfoFromTx(tx Transaction) types.TransactionBasicInfo {
	info := types.TransactionBasicInfo{}
	info.From = tx.From
	info.To = tx.To
	info.Value = tx.Value
	info.Hash = tx.Hash
	info.BlockNumber = tx.BlockNumber
	//info.age = ?
	return info
}
