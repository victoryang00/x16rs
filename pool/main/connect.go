package main

import (
	"bufio"
	"bytes"
	"encoding/hex"
	"fmt"
	"github.com/hacash/x16rs"
	"io"
	"math/rand"
	"net"
	"time"
)

var prevStopCh chan int = nil

var connectStatus bool = false
var connectObj net.Conn = nil

// 矿池

func StartPoolWorker(conf *Config) {

	go connect(conf)

	// 重连
	go func() {
		time.Sleep(time.Second * 11)
		for {
			time.Sleep(time.Second * 7)
			if connectStatus == false {
				timeout := rand.Intn(30)
				time.Sleep(time.Second * time.Duration(timeout))
				go connect(conf) // 重新连接
			}
		}
	}()

	// 定时发送心跳活跃消息
	go func() {
		for {
			time.Sleep(time.Minute * 5)
			if connectObj != nil && connectStatus {
				// type=11 心跳包
				x16rs.MiningPoolWriteTcpMsgBytes(connectObj, 11, []byte{})
			}
		}
	}()

}

func connect(conf *Config) {

	var tcpAddr *net.TCPAddr
	tcpAddr, _ = net.ResolveTCPAddr("tcp", conf.Pool) // "127.0.0.1:3339"

	conn, err := net.DialTCP("tcp", nil, tcpAddr)
	connectObj = conn

	if err != nil {
		fmt.Println("client connect error ! " + err.Error())
		return
	}

	defer func() {
		conn.Close()
		connectObj = nil
	}()

	fmt.Printf("connected pool: %s, cpu superneve: %d, reward address: %s\n", conf.Pool, conf.Supervene, conf.Reward)
	fmt.Println(conn.LocalAddr().String() + " : client connected !")

	onMessageReceived(conn, conf)

}

func onMessageReceived(conn *net.TCPConn, conf *Config) {

	// 发送注册
	x16rs.MiningPoolWriteTcpMsgBytes(conn, 0, []byte(conf.Reward)) // "1MzNY1oA3kfgYi75zquj3SRUPYztzXHzK9"

	//b := []byte(conn.LocalAddr().String() + " Say hello to Server... \n")
	//conn.Write(b)

	//bufcache := bytes.NewBuffer([]byte{})

	reader := bufio.NewReader(conn)

	for {
		msgbytes, err := x16rs.MiningPoolReadTcpMsgBytes(reader)
		//fmt.Println("ReadString")
		//fmt.Println([]byte(msgbytes))
		if err != nil {
			fmt.Println("reader message error:", err)
			if err != nil || err == io.EOF {
				connectStatus = false
				return // 连接失败
			}
		}

		// 连接状态
		connectStatus = true

		//fmt.Println(len(msgbytes), msgbytes)
		//time.Sleep(time.Millisecond * 1100)

		// 错误消息
		if msgbytes[0] == 255 { // type=255
			panic(string(msgbytes[1:]))
		}

		// 开始挖矿消息
		if msgbytes[0] == 1 { // type=1
			// 单次挖矿停止
			if prevStopCh != nil {
				prevStopCh <- 1 // 停止上一次挖矿
			}
			var stopCh chan int = make(chan int, 2)
			prevStopCh = stopCh
			go func() {
				var stuff x16rs.MiningPoolStuff
				stuff.Parse(msgbytes, 1)

				superneve := conf.Supervene // 多线程

				fmt.Printf("mining block height: %d, ", stuff.BlockHeight)
				mlok, successNonce, _, mostPowerNonce, mostPowerHash, totalPower, isbreak := startMining(stuff, &stopCh, superneve)
				fmt.Printf("work %d result last hash: %s, total power: %s \n",
					stuff.BlockHeight,
					hex.EncodeToString(mostPowerHash[0:16])+"...",
					totalPower.String())

				if mlok == true {

					//// 传递算力统计
					//bhm := stuff.BlockHeadMeta
					//bhm[79] = successNonce[0]
					//bhm[80] = successNonce[1]
					//bhm[81] = successNonce[2]
					//bhm[82] = successNonce[3]
					//fmt.Println("stuff.BlockHeadMeta", hex.EncodeToString(stuff.BlockHeadMeta))
					//hashbase := sha3.Sum256(stuff.BlockHeadMeta)
					//reshash := x16rs.HashX16RS_Optimize(int(stuff.Loopnum), hashbase[:])
					//fmt.Println("mining reshash", hex.EncodeToString(reshash))

					// 传递挖矿成功消息回去
					success := &x16rs.MiningSuccess{
						stuff.BlockHeight,
						stuff.MiningIndex,
						successNonce,
					}
					// type=2
					x16rs.MiningPoolWriteTcpMsgBytes(conn, 2, success.Serialize())

				} else {

					//fmt.Println("<<<<<<<<<<<<<<<<<< bhm := stuff.BlockHeadMeta")
					//fmt.Println(stuff.Loopnum, mostPowerHash)

					// 传递算力统计
					buf := bytes.NewBuffer([]byte{stuff.Loopnum})
					buf.Write(mostPowerNonce)
					if isbreak {
						buf.Write([]byte{0})
					} else {
						buf.Write([]byte{1}) // 遍历完所有数字，请求下一次挖矿
					}

					//buf.Write(stuff.BlockHeadMeta)
					//bhm[79] = mostPowerNonce[0]
					//bhm[80] = mostPowerNonce[1]
					//bhm[81] = mostPowerNonce[2]
					//bhm[82] = mostPowerNonce[3]
					//for i := 0; i < len(allPowerNonces); i++ {
					//	for k:=0; k<1; k++ {
					//		buf.Write(allPowerNonces[i])
					//	}
					//}
					//fmt.Println("stuff.BlockHeadMeta", stuff.BlockHeadMeta)
					//hashbase := sha3.Sum256(stuff.BlockHeadMeta)
					//reshash := x16rs.HashX16RS_Optimize(int(stuff.Loopnum), hashbase[:])
					//fmt.Println("reshash", reshash)
					// type=3
					x16rs.MiningPoolWriteTcpMsgBytes(conn, 3, buf.Bytes())
				}
			}()
		}

		//time.Sleep(time.Second * 2)

		//fmt.Println("writing...")

		//b := []byte(conn.LocalAddr().String() + " write data to Server... \n")
		//_, err = conn.Write(b)

	}
}
