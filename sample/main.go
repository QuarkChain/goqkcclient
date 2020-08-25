package main

import (
	"fmt"
	clt "github.com/QuarkChain/goqkcclient/client"
	"math/big"
	"sync"
)

var (
	client       = clt.NewClient("http://34.222.230.172:38391")
	fullShardKey = uint32(0)
)


var(
	ids=[]int{1,65537,131073,196609,262145,327681,393217,458753}
	heights=[]int{4261560,4228787,4225299,4249469,4260585,4231027,4272980,4190999}
	dist=518400/600
	wg sync.WaitGroup
)

func handle(id int,height int)  {
	//fmt.Println("============fullshardid",id,"from",height,"to",height-dist)

	ff:=make(map[string]bool)
	tt:=make(map[string]bool)
	all:=0
	for h:=height;h>=height-dist;h--{
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
			}
		}
		if h%1000==0{
			fmt.Println("handle h",h,height-dist)
		}
		if h==height-dist/2{
			//fmt.Println("currHeight",h,"from",heights[index],"to",heights[index]-dist,"tx nums",all,"from addr nums",len(ff),"to addr nums",len(tt))
		}
	}
	fmt.Println("all data","from",height,"to",height-dist,"tx nums",all,"from addr nums",len(ff),"to addr nums",len(tt))
wg.Done()
}
func main() {
	wg.Add(8)
	handle(ids[0],heights[0])
	handle(ids[1],heights[1])
	handle(ids[2],heights[2])
	handle(ids[3],heights[3])
	handle(ids[4],heights[4])
	handle(ids[5],heights[5])
	handle(ids[6],heights[6])
	handle(ids[7],heights[7])
	wg.Wait()
	fmt.Println("EEEEEEEEEEE")
}
