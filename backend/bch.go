package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type price string

func HandleBchPrice(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "Method is not supported.", http.StatusNotFound)
		return
	}
	var p price = "400.00"
	out, _ := json.Marshal(p)
	_, _ = fmt.Fprintf(w, string(out))
}
