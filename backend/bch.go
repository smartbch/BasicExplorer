package main

import (
	"fmt"
	"net/http"
)

func HandleBchPrice(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "Method is not supported.", http.StatusNotFound)
		return
	}
	_, _ = fmt.Fprintf(w, string(""))
}
