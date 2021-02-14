package main

import (
	"fmt"

	clt "github.com/QuarkChain/goqkcclient/client"
	"github.com/QuarkChain/goquarkchain/qkcdb"
	ethcom "github.com/ethereum/go-ethereum/common"
)

var (
	client       = clt.NewClient("http://34.222.230.172:38391")
	fullShardKey = uint32(0)
)

func main() {
	paths := make([]string, 8)
	paths[0] = "/home/gocode/src/github.com/QuarkChain/goquarkchain/cmd/cluster/qkc-data/mainnet/S0/shard-1/db"
	paths[1] = "/home/gocode/src/github.com/QuarkChain/goquarkchain/cmd/cluster/qkc-data/mainnet/S1/shard-65537/db"
	/*	paths[2] = "/home/gocode/src/github.com/QuarkChain/goquarkchain/cmd/cluster/qkc-data/mainnet/S2/shard-131073/db"
		paths[3] = "/home/gocode/src/github.com/QuarkChain/goquarkchain/cmd/cluster/qkc-data/mainnet/S3/shard-196609/db"
		paths[4] = "/home/gocode/src/github.com/QuarkChain/goquarkchain/cmd/cluster/qkc-data/mainnet/S0/shard-262145/db"
		paths[5] = "/home/gocode/src/github.com/QuarkChain/goquarkchain/cmd/cluster/qkc-data/mainnet/S1/shard-327681/db"
		paths[6] = "/home/gocode/src/github.com/QuarkChain/goquarkchain/cmd/cluster/qkc-data/mainnet/S2/shard-393217/db"
		paths[7] = "/home/gocode/src/github.com/QuarkChain/goquarkchain/cmd/cluster/qkc-data/mainnet/S3/shard-458753/db"*/

	for idx, path := range paths {
		m := GetAccounts(idx, path)
		for acc, _ := range m {
			m[acc] = getBalance(acc, uint32(idx))
		}
	}
}

func getBalance(address ethcom.Address, shardid uint32) uint64 {
	context := make(map[string]string)
	addr := clt.QkcAddress{Recipient: address, FullShardKey: shardid}
	context["address"] = addr.Recipient.Hex()
	context["fromFullShardKey"] = addr.FullShardKeyToHex()
	balance, err := client.GetBalance(&addr)
	if err != nil {
		fmt.Println(err.Error())
	}
	return balance.Uint64()
}

func GetAccounts(idx int, path string) map[ethcom.Address]uint64 {
	fmt.Println("--------", idx, "---------")
	accounts := make(map[ethcom.Address]uint64)
	db, err := qkcdb.NewDatabase(path, false, false)
	if err != nil {
		panic(err)
	}
	it := db.NewIterator()
	it.Seek([]byte{})
	for it.Valid() {
		size := len(it.Key())
		if size == 43 {
			value, err := db.Get(it.Key())
			if err != nil {
				panic(err)
			}
			if len(value) == 20 {
				accounts[value] = 0
			}
		}
		it.Next()
	}
}
