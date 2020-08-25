package main

import (
	"fmt"
	clt "github.com/QuarkChain/goqkcclient/client"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"math/big"
)

var (
	client       = clt.NewClient("http://34.222.230.172:38391")
	fullShardKey = uint32(0)
)


var(
	ids=[]int{0}
	heights=[]int{4261327}
)
func main() {
	for index:=0;index<len(ids);index++{
		id:=ids[index]

		fmt.Println("fullshardid",id,"from",heights[index],"to",heights[index]-518400)

		ff:=make(map[string]bool)
		all:=0

		for h:=heights[index];h>=heights[index]-5184;h--{
			ans,_:=client.GetMinorBlockByHeight(uint32(id),new(big.Int).SetUint64(uint64(h)))
			//fmt./Println("hhhh",h,err)
			txs:=ans.Result.(map[string]interface{})["transactions"]
			//fmt.Println("txsss",txs)
			sv,_:=txs.([]interface{})
			if len(sv)!=0{
				for index:=0;index<len(sv);index++{
					dd:=sv[index].(map[string]interface{})
					from:=dd["from"].(string)
					to:=dd["to"].(string)
					ff[from]=true
					ff[to]=true
					all++
					fmt.Println("from ",from,to)
				}

			}
		}
		fmt.Println("all tx",all)
		for k,v:=range ff{
			fmt.Println("addr",k)
		}
	}

}

//获取余额
func getBalance(addr *clt.QkcAddress) {
	//address := common.HexToAddress(ctx.FormValue("address"))
	balance, err := client.GetBalance(addr)
	if err != nil {
		fmt.Println(err.Error())
	}
	fmt.Println(balance)
}

//获取区块和交易内容
func getBlock(ctx map[string]string) map[string]interface{} {
	heightStr := ctx["height"]
	height := new(big.Int).SetBytes(common.FromHex(heightStr))
	result, err := client.GetRootBlockByHeight(height)
	if err != nil {
		fmt.Println(err.Error())
	}

	headers := result.Result.(map[string]interface{})["minorBlockHeaders"]
	headerIdList := make([]string, 0)
	txList := make([]interface{}, 0)
	for _, h := range headers.([]interface{}) {
		info := h.(map[string]interface{})
		id := (info["id"]).(string)
		headerIdList = append(headerIdList, id)
	}
	fmt.Println("headerIdList len", len(headerIdList))
	for _, headerId := range headerIdList {
		mresult, err := client.GetMinorBlockById(headerId)
		if err != nil {
			fmt.Println(err.Error())
		}
		txs := mresult.Result.(map[string]interface{})["transactions"]
		for _, tx := range txs.([]interface{}) {
			txList = append(txList, tx)
		}
	}
	result.Result.(map[string]interface{})["transactions"] = txList
	fmt.Println("txList len", len(txList))
	fmt.Println(result.Result)
	return result.Result.(map[string]interface{})
}

//获取交易回执
func getReceipt(ctx map[string]string) {
	txid, err := clt.ByteToTransactionId(common.FromHex(ctx["txid"]))
	if err != nil {
		fmt.Println(err.Error())
	}
	result, err := client.GetTransactionReceipt(txid)
	if err != nil {
		fmt.Println("getTransactionReceipt error: ", err.Error())
	}
	fmt.Println(result.Result)
}

func getHeight(ctx map[string]string) uint64 {
	height, err := client.GetRootBlockHeight()
	if err != nil {
		fmt.Println(err.Error())
	}
	fmt.Println(height)
	return height
}

func getTransaction(ctx map[string]string) {
	txid, err := clt.ByteToTransactionId(common.FromHex(ctx["txid"]))
	if err != nil {
		fmt.Println(err.Error())
	}
	result, err := client.GetTransactionById(txid)
	if err != nil {
		fmt.Println(err.Error())
	}
	fmt.Println("txid", result.Result.(map[string]interface{})["id"])
	fmt.Println(result.Result)
}

func sent(ctx map[string]string) string {
	from := common.HexToAddress(ctx["from"])
	to := common.HexToAddress(ctx["to"])
	amount, _ := new(big.Int).SetString(ctx["amount"], 10)
	gasPrice, _ := new(big.Int).SetString(ctx["price"], 10)
	privateKey := ctx["privateKey"]
	prvkey, _ := crypto.ToECDSA(common.FromHex(privateKey))
	fromFullShardKey := fullShardKey
	if _, ok := ctx["fromFullShardKey"]; ok {
		fromFullShardKey = uint32(new(big.Int).SetBytes(common.FromHex(ctx["fromFullShardKey"])).Uint64())
	}
	toFullShardKey := fullShardKey
	if _, ok := ctx["toFullShardKey"]; ok {
		toFullShardKey = uint32(new(big.Int).SetBytes(common.FromHex(ctx["toFullShardKey"])).Uint64())
	}
	tx, err := client.CreateTransaction(&clt.QkcAddress{Recipient: from, FullShardKey: fromFullShardKey}, &clt.QkcAddress{Recipient: to, FullShardKey: toFullShardKey}, amount, uint64(30000), gasPrice)
	if err != nil {
		fmt.Println(err.Error())
	}
	tx, err = clt.SignTx(tx, prvkey)
	if err != nil {
		fmt.Println(err.Error())
	}
	txid, err := client.SendTransaction(tx)
	if err != nil {
		fmt.Println("SendTransaction error: ", err.Error())
	}

	fmt.Println(common.Bytes2Hex(txid))
	return common.Bytes2Hex(txid)
}
