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

func HandleBlock(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "Method is not supported.", http.StatusNotFound)
		return
	}

	w.Header().Add("Access-Control-Allow-Origin", "*")

	u, _ := url.Parse(r.RequestURI)
	number := path.Base(u.Path)
	_, _ = fmt.Fprintf(w, string(GetBlock(number)))
}

func GetBlock(number string) []byte {
	r := types.RequestInfo{
		Jsonrpc: "2.0",
		Method:  "eth_getBlockByNumber",
		Params:  []interface{}{number, true},
		Id:      1,
	}
	b, _ := json.Marshal(r)
	resp, _ := http.Post(config.NodeUrl, "application/json", bytes.NewReader(b))
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	fmt.Println(string(body))
	res := types.ResponseInfo{}
	block := types.Block{}
	err := json.Unmarshal(body, &res)
	if err != nil {
		panic(err)
	}
	err = json.Unmarshal(res.Result, &block)
	block.TransactionsCount = "0x" + strconv.FormatInt(int64(len(block.Transactions)), 16)
	out, _ := json.MarshalIndent(block, "", "    ")
	return out
}
