# SERVER TEST

http://localhost:8080/v1/account/0xAb5D62788E207646fA60EB3eEbDC4358C7F5686c?page=1&from=false&to=true

```
{
    "Balance": "0xde0b6b3a763adf8",
    "Transactions": [
        {
            "Hash": "0x330969a2e04e92329d88e136920b89793f9ff7c9264e08fb703d534c77c0b4c2",
            "BlockNumber": "0x50",
            "From": "0x3e144eb45c5ff912b2b29b2823fa674c972e9ec0",
            "To": "0xab5d62788e207646fa60eb3eebdc4358c7f5686c",
            "Age": "",
            "Value": "0x100"
        }
    ]
}
```

http://localhost:8080/v1/block/0x51

```
{
    "Number": "0x51",
    "Hash": "0x982bf834edfc9efebc5db754de440f4276cef8e2d9e41350e6fa14fa3d2564c5",
    "ParentHash": "0x86b41c5139dedbe70bd25b3fca1886e8497c1983ed0b39fa28b818c0b0c413a7",
    "Miner": "0x0000000000000000000000000000000000000000",
    "Size": "0x1d2",
    "GasLimit": 0,
    "GasUsed": "0x0",
    "Timestamp": "0x6040a25d",
    "TransactionsCount": "0x1",
    "BlockReward": "",
    "StateRoot": "0xd19df6b17b6b4483080baa19256930f849cfd902c426308f5da12e454ba7354b",
    "TransactionsRoot": "0xe3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855",
    "Transactions": [
        "0x330969a2e04e92329d88e136920b89793f9ff7c9264e08fb703d534c77c0b4c2"
    ]
}
```

http://localhost:8080/v1/tx/0x330969a2e04e92329d88e136920b89793f9ff7c9264e08fb703d534c77c0b4c2

```
{
    "Hash": "0x330969a2e04e92329d88e136920b89793f9ff7c9264e08fb703d534c77c0b4c2",
    "Nonce": "0x0",
    "BlockNumber": "0x50",
    "TransactionIndex": "0x0",
    "From": "0x3e144eb45c5ff912b2b29b2823fa674c972e9ec0",
    "To": "0xab5d62788e207646fa60eb3eebdc4358c7f5686c",
    "Value": "0x100",
    "Gas": "0x100000",
    "GasPrice": "0x1",
    "Input": "0x",
    "Logs": []
}
```

