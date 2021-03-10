package config

import "os"

var NodeUrl = "http://localhost:8545"
var ServerDir = os.ExpandEnv("$HOME/.moeing-server")
var (
	DataDir   = ServerDir + "/data"
	ConfigDir = ServerDir + "/config"
)

var (
	/*
		token address => {
			symbol,decimals
		}
		not store totalSupply as it may not const
	*/
	KeyERC20InfoByAddr = []byte("0x01")
	/*
		tokenSymbol => [address]
		store erc20 token addresses with same symbol
	*/
	KeyERC20sBySymbol = []byte("0x02")
	//address => alias
	//like 0x6f259637dcD74C767781E37Bc6133cd6A68aa161 => Huobi Token
	//like 0xE93381fB4c4F14bDa253907b18faD305D799241a => Huobi 10
	KeyAccountAlias = []byte("0x10")
)

func BuildERC20InfoKey(address string) []byte {
	return append(KeyERC20InfoByAddr, address...)
}

func BuildERC20sKey(symbol string) []byte {
	return append(KeyERC20sBySymbol, symbol...)
}

func BuildAccAliasKey(address string) []byte {
	return append(KeyAccountAlias, address...)
}

const (
	UrlGetBlock          = "/v1/block/"
	UrlGetAccount        = "/v1/account/"
	UrlGetBlockTxs       = "/v1/txs/"
	UrlGetTx             = "/v1/tx/"
	UrlGetToken          = "/v1/erc20/"
	UrlGetTokenAddresses = "/v1/erc20s/"
	UrlGetBchPrice       = "/v1/bch_price"
)
