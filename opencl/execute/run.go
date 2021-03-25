package execute

import (
	"fmt"
	"sync"
)

// 开始采矿
func (g *GpuMiner) DoMining(blockHeight uint64, stopmark *byte, tarhashvalue []byte, blockheadmeta []byte) (byte, bool, []byte, []byte) {

	deviceNum := len(g.devices)
	groupsize := 256
	overallstep := groupsize * deviceNum
	fmt.Println(overallstep)

	var syncWait = sync.WaitGroup{}
	syncWait.Add(deviceNum)

	return 0, false, nil, nil

}
