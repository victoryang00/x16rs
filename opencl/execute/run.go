package execute

import "sync"

// 开始采矿
func (mr *GpuMiner) DoMining(blockHeight uint64, stopmark *byte, tarhashvalue []byte, blockheadmeta []byte) (byte, bool, []byte, []byte) {

	deviceNum := len(mr.devices)
	groupsize := 256
	overallstep := groupsize * deviceNum

	var syncWait = sync.WaitGroup{}
	syncWait.Add(deviceNum)

}
