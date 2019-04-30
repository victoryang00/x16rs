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

	if _, ok := params["height"]; !ok {
		response.Write([]byte("params must height"))
		return
	}
	if _, ok := params["blockstuff"]; !ok {
		response.Write([]byte("params must blockstuff"))
		return
	}
	if _, ok := params["targethash"]; !ok {
		response.Write([]byte("params must targethash"))
		return
	}

	height, _ := strconv.ParseInt(params["height"], 10, 0)
	stuff, _ := hex.DecodeString(params["blockstuff"])
	var blkstuff [89]byte
	copy(blkstuff[:], stuff)
	targethash, _ := hex.DecodeString(params["targethash"])
	var target [32]byte
	copy(target[:], targethash)

	// 进行挖矿
	var resultCh = make(chan MinerResult, 1)
	gpuminer.DoMiner(uint32(height), blkstuff, target, resultCh)

	// 等待状态
	res := <-resultCh
	if res.success {
		nonce := hex.EncodeToString(res.nonce)
		bts := fmt.Sprintf(".[%d,%d,%d,%d]", res.nonce[0], res.nonce[1], res.nonce[2], res.nonce[3])
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
