package main

import (
	"fmt"
	clt "github.com/QuarkChain/goqkcclient/client"
	"math/big"
)

var (
	client       = clt.NewClient("http://34.222.230.172:38391")
	fullShardKey = uint32(0)
)


var(
	ids=[]int{1,65537,131073,196609,262145,327681,393217,458753}
	heights=[]int{4261560,4228787,4225299,4249469,4260585,4231027,4272980,4190999}
	dist=518
)
func main() {
	for index:=0;index<len(ids);index++{
		id:=ids[index]

		fmt.Println("fullshardid",id,"from",heights[index],"to",heights[index]-dist)

		ff:=make(map[string]bool)
		tt:=make(map[string]bool)
		all:=0
		for h:=heights[index];h>=heights[index]-dist;h--{
			ans,_:=client.GetMinorBlockByHeight(uint32(id),new(big.Int).SetUint64(uint64(h)))
			txs:=ans.Result.(map[string]interface{})["transactions"]
			sv,_:=txs.([]interface{})
			if len(sv)!=0{
				for index:=0;index<len(sv);index++{
					dd:=sv[index].(map[string]interface{})
					from:=dd["from"].(string)
					to:=dd["to"].(string)
					ff[from]=true
					tt[to]=true
					all++
					//fmt.Println("from ",from,to)
				}
			}
			//fmt.Println("hhhhhhhh",h,heights[index]-dist)
			if h%1000==0{
				fmt.Println("handle h",h,heights[index]-dist)
			}
			if h==heights[index]-dist/2{
				fmt.Println("518400/2","tx nums",all,"from addr nums",len(ff),"to addr nums",len(tt))
			}
		}
		fmt.Println("518400","tx nums",all,"from addr nums",len(ff),"to addr nums",len(tt))
	}

}
