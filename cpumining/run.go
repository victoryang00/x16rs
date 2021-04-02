package cpumining

import (
	"github.com/hacash/mint/difficulty"
	"github.com/hacash/x16rs"
	"sync"
	"sync/atomic"
)

// 执行一次挖矿
func (c *CPUMining) DoMining(blockHeight uint64, reporthashrate bool, stopmark *byte, tarhashvalue []byte, blockheadmeta_list [][]byte) (bool, int, []byte, []byte) {

	startNonce := uint32(0)
	endNonce := uint32(4294967295)

	supervene := len(blockheadmeta_list)

	// 成功
	var smm uint32 = 0
	successMiningMark := &smm // 未成功标记

	// 返回值
	var rSuccess bool = false
	var rBlkhmi int = 0
	var rNonce []byte = nil
	var rPowerHash []byte = nil

	// 并发 group
	var checkLock sync.Mutex // 同步检查
	var syncWait = sync.WaitGroup{}
	syncWait.Add(supervene)

	// 挖矿
	for i := 0; i < supervene; i++ {
		go func(i int, nextstop *byte) {
			defer func() {
				checkLock.Unlock()
				syncWait.Done()
			}()
			blockheadmeta := blockheadmeta_list[i]
			// 开始挖款
			_, success, nonce, endhash := x16rs.MinerNonceHashX16RS(blockHeight, reporthashrate, stopmark, startNonce, endNonce, tarhashvalue, blockheadmeta)
			checkLock.Lock() // 串行锁
			if success && atomic.CompareAndSwapUint32(successMiningMark, 0, 1) {
				*nextstop = 1 // 停止所有挖矿
				rSuccess = true
				rBlkhmi = i
				rNonce = nonce
				rPowerHash = endhash
			} else if reporthashrate {
				// 比较算力大小
				if rPowerHash == nil || difficulty.CheckHashDifficultySatisfy(endhash, rPowerHash) {
					rBlkhmi = i
					rNonce = nonce
					rPowerHash = endhash
				}
			}

		}(i, stopmark)
	}

	syncWait.Wait()

	// 返回
	return rSuccess, rBlkhmi, rNonce, rPowerHash

}
