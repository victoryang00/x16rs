package main

import (
	"fmt"
	"github.com/hacash/x16rs"
	"math/big"
	"sync"
)

// 开始挖矿

func startMining(stuff x16rs.MiningPoolStuff, stopCh *chan int, superneve uint32) (bool, []byte, uint64, []byte, []byte, *big.Int, bool) {

	segsize := 4294967294 / superneve

	var stopKind int = 0
	var stopsign *byte = new(byte)
	*stopsign = 0

	var successNonce []byte = nil
	//var totalPower = new(big.Int)
	//var allNonces [][]byte = make([][]byte, 0)
	//var someOneHash = []byte{}
	var mostPowerHash []byte = nil
	var mostPowerNonce []byte = nil
	var mostPower = new(big.Int)

	var group sync.WaitGroup
	group.Add(int(superneve))

	for i := uint32(0); i < superneve; i++ {
		go func(i uint32) {
			segstart := segsize * i
			segend := segstart + segsize
			// 启动挖矿
			success, nonce, reshash := x16rs.MinerNonceHashX16RS(int(stuff.Loopnum), true, stopsign, segstart, segend, stuff.TargetHash, stuff.BlockHeadMeta)
			// 成功
			if success && successNonce == nil {
				fmt.Printf("⬤  h:%d, mining successfully and got rewords! \n", stuff.BlockHeight)
				successNonce = nonce
				*stopCh <- -1 // 写入停止
			}
			//allNonces = append(allNonces, nonce)
			//totalPower = totalPower.Add(totalPower, x16rs.CalculateHashPowerValue(reshash))
			//someOneHash = reshash
			// 判断最大的hash
			if mostPowerHash == nil {
				mostPowerHash = reshash
				mostPowerNonce = nonce
				mostPower = x16rs.CalculateHashPowerValue(reshash)
			} else {
				for i := 0; i < 32; i++ {
					if reshash[i] > mostPowerHash[i] {
						break
					} else if reshash[i] < mostPowerHash[i] {
						mostPowerHash = reshash // 更大
						mostPowerNonce = nonce
						mostPower = x16rs.CalculateHashPowerValue(reshash)
						break
					}
				}
			}
			group.Done()
		}(i)
	}

	go func() {
		stopKind = <-*stopCh // 等待停止
		*stopsign = 1        // 停止其他全部挖矿
	}()

	// 等待全部
	group.Wait()

	// 返回数据
	return successNonce != nil, successNonce, stuff.MiningIndex, mostPowerNonce, mostPowerHash, mostPower, stopKind == 1

}
