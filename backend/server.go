package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"runtime/debug"
	"syscall"

	"github.com/coocood/freecache"
	"github.com/syndtr/goleveldb/leveldb"

	"backend/config"
)

type Server struct {
	Address string
	Cache   *freecache.Cache
	DB      *leveldb.DB
}

var server *Server

func NewServer(address string, cacheSize int) *Server {
	s := Server{
		Address: address,
		Cache:   freecache.NewCache(cacheSize),
	}
	debug.SetGCPercent(20)
	return &s
}

func (s *Server) init() {
	//http://localhost:8080/v1/block/0x50
	http.HandleFunc(config.UrlGetBlock, HandleBlock)
	//http://localhost:8080/v1/account/0xAb5D62788E207646fA60EB3eEbDC4358C7F5686c?page=1&from=false&to=true
	http.HandleFunc(config.UrlGetAccount, HandleAccount)
	//
	http.HandleFunc(config.UrlGetBlockTxs, HandleBlockTxs)
	//http://localhost:8080/v1/tx/0x330969a2e04e92329d88e136920b89793f9ff7c9264e08fb703d534c77c0b4c2
	http.HandleFunc(config.UrlGetTx, HandleTransactionInfo)
	//http://localhost:8080/v1/erc20/0xC7BBd3373c6D9f582102c332bE91e8dCDd087e35
	http.HandleFunc(config.UrlGetToken, HandleToken)
	http.HandleFunc(config.UrlGetTokenAddresses, HandleTokens)
	http.HandleFunc(config.UrlGetBchPrice, HandleBchPrice)
}

func (s *Server) start() {
	db, err := leveldb.OpenFile(config.DataDir, nil)
	if err != nil {
		panic(err)
	}
	s.DB = db
	fmt.Printf("Starting server at port 8080\n")
	if err := http.ListenAndServe(s.Address, nil); err != nil {
		log.Fatal(err)
	}
}

func (s *Server) stop() {
	_ = s.DB.Close()
}

func TrapSignal(cleanupFunc func()) {
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		sig := <-sigs
		if cleanupFunc != nil {
			cleanupFunc()
		}
		exitCode := 128
		switch sig {
		case syscall.SIGINT:
			exitCode += int(syscall.SIGINT)
		case syscall.SIGTERM:
			exitCode += int(syscall.SIGTERM)
		}
		os.Exit(exitCode)
	}()
}

func main() {
	server = NewServer(":8080", 100*1024*1024)
	server.init()
	server.start()
	TrapSignal(func() {
		server.stop()
		//server.Logger.Info("exiting...")
	})
	select {}
}
