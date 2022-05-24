package cpumining

import (
	"github.com/hacash/mint/difficulty"
	"github.com/hacash/x16rs"
	"sync"
	"sync/atomic"
)

// execute a mining job
func (c *CPUMining) DoMining(blockHeight uint64, reporthashrate bool, stopmark *byte, tarhashvalue []byte, blockheadmeta_list [][]byte) (bool, int, []byte, []byte) {

	startNonce := uint32(0)
	endNonce := uint32(4294967295)

	supervene := len(blockheadmeta_list)

	// success mark
	var smm uint32 = 0
	successMiningMark := &smm // failed mark

	// return value
	var rSuccess bool = false
	var rBlkhmi int = 0
	var rNonce []byte = nil
	var rPowerHash []byte = nil

	// Concurrent group
	var checkLock sync.Mutex // Synchronization check
	var syncWait = sync.WaitGroup{}
	syncWait.Add(supervene)

	// mining
	for i := 0; i < supervene; i++ {
		go func(i int, nextstop *byte) {
			defer func() {
				checkLock.Unlock()
				syncWait.Done()
			}()
			blockheadmeta := blockheadmeta_list[i]
			// start mining
			_, success, nonce, endhash := x16rs.MinerNonceHashX16RS(blockHeight, reporthashrate, stopmark, startNonce, endNonce, tarhashvalue, blockheadmeta)
			checkLock.Lock() // Serial lock
			if success && atomic.CompareAndSwapUint32(successMiningMark, 0, 1) {
				*nextstop = 1 // Stop all mining
				rSuccess = true
				rBlkhmi = i
				rNonce = nonce
				rPowerHash = endhash
			} else if reporthashrate {
				// Compare the calculation force
				if rPowerHash == nil || difficulty.CheckHashDifficultySatisfy(endhash, rPowerHash) {
					rBlkhmi = i
					rNonce = nonce
					rPowerHash = endhash
				}
			}

		}(i, stopmark)
	}

	syncWait.Wait()

	// return
	return rSuccess, rBlkhmi, rNonce, rPowerHash

}
