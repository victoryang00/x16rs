package main

import (
	"encoding/hex"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
)

var gpuminer *GpuMiner

func dealHome(response http.ResponseWriter, request *http.Request) {
	response.Write([]byte("gpu miner"))
}

func dealDoMiner(response http.ResponseWriter, request *http.Request) {
	params := parseQueryForm(request)
	// 参数检查
	needs := []string{"height", "blockstuff", "targethash", "coinmsg", "coinaddr"}
	for _, n := range needs {
		if _, ok := params[n]; !ok {
			response.Write([]byte("params must " + n))
			return
		}
	}
	// 参数值
	height, _ := strconv.ParseInt(params["height"], 10, 0)
	stuff, _ := hex.DecodeString(params["blockstuff"])
	var blkstuff [89]byte
	copy(blkstuff[:], stuff)
	targethash, _ := hex.DecodeString(params["targethash"])
	var target [32]byte
	copy(target[:], targethash)
	coinmsg_p, _ := hex.DecodeString(params["coinmsg"])
	var coinmsg [16]byte
	copy(coinmsg[:], coinmsg_p)
	coinaddr_p, _ := hex.DecodeString(params["coinaddr"])
	var coinaddr [21]byte
	copy(coinaddr[:], coinaddr_p)
	// 执行挖矿
	exeRet := gpuminer.ReStartMiner(uint32(height), blkstuff, target, coinmsg, coinaddr)
	//
	//// 进行挖矿
	//var resultCh = make(chan MinerResult, 1)
	//gpuminer.DoMiner(uint32(height), blkstuff, target, resultCh)
	//// 等待状态
	//res := <-resultCh
	//
	if exeRet.retry {
		response.Write([]byte(params["height"] + ".retry")) // 重新生成区块并尝试
	} else if exeRet.success {
		nonce := hex.EncodeToString(exeRet.nonce)
		bts := fmt.Sprintf(".[%d,%d,%d,%d]", exeRet.nonce[0], exeRet.nonce[1], exeRet.nonce[2], exeRet.nonce[3])
		response.Write([]byte(params["height"] + ".success." + nonce + bts))
	} else {
		response.Write([]byte(params["height"] + ".fail"))
	}
}

func parseQueryForm(request *http.Request) map[string]string {
	request.ParseForm()
	params := make(map[string]string)
	for k, v := range request.Form {
		//fmt.Println("key:", k)
		//fmt.Println("val:", strings.Join(v, ""))
		params[k] = strings.Join(v, "")
	}

	return params
}

func RunHttpRpcService(miner *GpuMiner, httpport int) {

	gpuminer = miner

	http.HandleFunc("/", dealHome)           //设置访问的路由
	http.HandleFunc("/dominer", dealDoMiner) //设置访问的路由

	port := strconv.Itoa(httpport)
	fmt.Println("run http rpc service on port: " + port)

	err := http.ListenAndServe(":"+port, nil) //设置监听的端口
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
