# Moeing explorer server API

#### 1. 根据高度查询block

```
GET /v1/block/{block_height}
```

**响应**

```
{
    "Number": "0x13",
    "Hash": "0x9c9c50705e431ad294e1fb6917ec8b0c00c28b282f776b20f1ad0e927f8ce448",
    "ParentHash": "0x89c82372fdb9aea8b3e5041c0fb1611985abed026f90e269046dfc282a37f2fa",
    "Miner": "0xd0b1a15d2759af1fe569e4d4074ae13c91087f31",
    "Size": "0x234",
    "GasLimit": "0x0",
    "GasUsed": "0x0",
    "Timestamp": "0x604f0ffd",
    "TransactionsCount": "0x1", //区块中的交易总数
    "BlockReward": "", //修改节点app代码计算出reward，添加到getBlock的字段里。
    "StateRoot": "0x8dce8194d72c0acc8ad78073ec7e7d54e39c6906cff457e8fc106711d933cd21",
    "TransactionsRoot": "0x0000000000000000000000000000000000000000000000000000000000000000",
    "Transactions": [
        "0x880c564aeaacdf1a877f9398a423f79a4831c0ec461e9b6d75f628bbc5b5ff3e"
    ]
}
```

#### 2. 根据区块的高度查询交易列表

```
GET /v1/txs/{block_height}
```

**响应**

```
{
    "Block": "0x6",
    "Transactions": [
        {
            "Hash": "0x09938944c3a7ed0179b2953883c5006e4ca15af8fcd4508761cc35f4acc48114",
            "BlockNumber": "0x6",
            "From": "0x09f236e4067f5fca5872d0c09f92ce653377ae41",
            "To": "0xc7bbd3373c6d9f582102c332be91e8dcdd087e35",
            "Age": "", //暂时不支持
            "Value": "0x0"
        }
    ]
}
```

#### 3. 根据地址查询账户信息

```
GET /api/v1/account/{address} //显示from和to地址为该地址的全部交易，不分页
GET /api/v1/account/{address}?page=1 //显示from和to地址为该地址的全部交易的第一页，每页25笔交易
GET /api/v1/account/{address}?from=true&page=1 //显示from为address的交易的第一页，每页25个交易
GET /api/v1/account/{address}?to=true&page=1 //只显示to为address的交易的第一页，每页25个交易
```

交易列表里的交易按区块高度从高到低排列。

**响应**

````
{
    "Balance": "0xde0b6b3a7640100",
    "Transactions": [
        {
            "Hash": "0x880c564aeaacdf1a877f9398a423f79a4831c0ec461e9b6d75f628bbc5b5ff3e",
            "BlockNumber": "0x13",
            "From": "0xab5d62788e207646fa60eb3eebdc4358c7f5686c",
            "To": "0x3e144eb45c5ff912b2b29b2823fa674c972e9ec0",
            "Age": "", //暂时不提供
            "Value": "0x100"
        }
    ]
}
````

#### 4. 根据交易HASH查询

```
GET /api/v1/tx/{hash}
```

**响应**

```
{
		"hash": "0x9fc76417374aa880d4449a1f7f31ec597f00b1f6f3dd2d66f4c9c6c445836d8b",
    "nonce": "0x2",
    "blockNumber": "0x3",
    "transactionIndex": "0x1",
    "from": "0xa94f5374fce5edbc8e2a8697c15331677e6ebf0b",
    "to": "0x6295ee1b4f6dd65047762f924ecd367c17eabf8f",
    "value": "0x123450000000000000",
    "gas": "0x314159",
    "gasPrice": "0x2000000000000",
    "gasUsed": "0x3000",
    "input": "0x57cb2fc4"
    "logs": [
    		{
    			"address":"0x56badd9B06bBaA1dF336A5A9524A90592a5Db962",
    			"topics":[
    				"0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef",
    				"0x56badd9b06bbaa1df336a5a9524a90592a5db962",
    				"0x18bdcf05a8093051472a09dcd3bfcaa0cdd546c0",
    			],
    			"data":"0x000000000000000000000000000000000000000c9f2c9cd04674edea40000000",
    		}
    ]
}
```

浏览器中显示的transaction fee可以用gasUsed乘以gasPrice乘以bch的实时价格计算出来。前端维护bch的最新价格

#### 5. 根据erc20的token名字查询token信息

```
GET /v1/erc20/{token_address} //erc20 token的全量交易记录，最多1000条
GET /v1/erc20/{token_address}?page=1 //erc20转账相关的交易，一页25条记录
GET /v1/erc20/{token_address}?address=0x49Fd1607a0b93334F090eBaF42C72BaBb38a0f76&page=1 //按地址筛选，address作为erc20 token的发送地址或者接收地址
```

**请求**

```
http://localhost:8080/v1/erc20/0xc7bbd3373c6d9f582102c332be91e8dcdd087e35?page=1
```

**响应**

```
{
    "Symbol": "4f50540000000000000000000000000000000000000000000000000000000000",
    "Decimals": "0x0000000000000000000000000000000000000000000000000000000000000001",
    "MaxSupply": "0x0000000000000000000000000000000000000000000000000000000000002710",
    "ContractAddress": "0xc7bbd3373c6d9f582102c332be91e8dcdd087e35",
    "Transactions": [
        {
            "Hash": "0x09938944c3a7ed0179b2953883c5006e4ca15af8fcd4508761cc35f4acc48114",
            "BlockNumber": "0x6",
            "From": "0x09f236e4067f5fca5872d0c09f92ce653377ae41",
            "To": "0xc7bbd3373c6d9f582102c332be91e8dcdd087e35",
            "Age": "", //暂时不提供
            "Value": "0x0" //暂时设为0
        }
    ]
}
```

#### 6.根据symbol查询erc20 token列表

```
GET /v1/erc20s/{token_symbol} //返回具有相同symbol的erc20token合约地址列表
```

多个erc20 token可能有相同的symbol，浏览器用户可以输入一个symbol，浏览器返回一个已知token的下拉列表，用户选择相应的token来查询详细信息。

**响应**

```
["0xc7bbd3373c6d9f582102c332be91e8dcdd087e35"]
```

#### 7. 查询最新BCH价格

```
GET /v1/bch_price
```

**响应**

```
"600.01" //美金
```

