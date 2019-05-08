package main

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"github.com/xfong/go2opencl/cl"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

// 单个设备执行单例
type minerDeviceExecute struct {
	autoidx     uint64
	deviceIndex uint32
	device      *cl.Device
	//kernel    *cl.Kernel
	//queue     *cl.CommandQueue
	blockHeadBytes    [89]byte
	targetHash        [32]byte
	baseStart         uint32
	groupSize         uint32
	blockHeight       uint32
	blockCoinbaseMsg  [16]byte
	blockCoinbaseAddr [21]byte

	// 状态数据

	// 挖矿状态
	retry   bool // 重新尝试挖矿
	success bool
	nonce   []byte
}

// 执行任务
type minerExecuteWork struct {
	autoidx uint64

	blockHeight       uint32
	blkstuff          [89]byte
	target            [32]byte
	blockCoinbaseMsg  [16]byte
	blockCoinbaseAddr [21]byte

	// id标记
	mid uint64
	// 结果通知
	resultCh chan *minerDeviceExecute
	// 执行队列
	executeQueueChList []chan *minerDeviceExecute
	// 等待全部关闭
	wg sync.WaitGroup
	// 全部停止标记
	stopMarkCh chan bool
}

type GpuMiner struct {
	autoidx uint64

	platform *cl.Platform
	context  *cl.Context
	program  *cl.Program
	devices  []*cl.Device // 所有设备
	//kernel  *cl.Kernel // 核心函数
	//queues   []*cl.CommandQueue // 所有队列

	// 执行队列
	//executeQueueChList []chan *minerDeviceExecute

	// 任务列表
	minerExecuteWorkListCh chan *minerExecuteWork

	// 执行任务
	//////////////////

	// device *cl.Device
	// queue  *cl.CommandQueue
	// kernel   *cl.Kernel

	// data
	blockHeadBytes []byte
	targetHash     []byte

	// config
	openclPath   string
	platName     string // 平台名称
	dvid         int    // 设备id
	groupSize    int    // 组大小
	loopNum      int    // 循环次数
	exeWide      int
	executeSize  int // 执行单例数量
	printNumBase int
	rebuild      bool // 强制重新编译

	workLoopLock sync.Mutex

	// msg
	//miningPrevId uint32
	//stopMarkCh   chan uint32
}

type MinerResult struct {
	success bool // 是否挖矿成功
	height  uint32
	nonce   []byte
}

/*
func (mr *GpuMiner) stopAll() {
	if mr.miningPrevId > 0 {
		// 清空待执行队列
		mr.stopMarkCh <- mr.miningPrevId // 投喂线程

		//for range mr.executeQueueChList {
		//	mr.stopMarkCh <- mr.miningPrevId // 执行线程
		//}
	}
}
*/

func (mr *GpuMiner) minerWorkLoop() {
	for {
		work := <-mr.minerExecuteWorkListCh
		// 判断当前执行id
		if work.mid < mr.autoidx {
			continue // 永远只执行最后一个任务
		}
		// 执行
		fmt.Printf("\n\ndo miner work<%d>, height<%d>, target<%s>, block<%s>\n", work.mid, work.blockHeight, hex.EncodeToString(work.target[:]), hex.EncodeToString(work.blkstuff[:]))
		// 投喂执行单元
		work.wg.Add(1)
		go func(work *minerExecuteWork) {
			defer func() {
				fmt.Printf(" , work %d exit ", work.mid)
				work.wg.Done()
			}()
			// 持续投喂
			for i := 0; ; i++ {
				// 判断停止
				if work.mid < mr.autoidx {
					return // 工作过期
				}
				baseItemNum := uint32(mr.groupSize * mr.loopNum)
				baseStart := uint32(i) * baseItemNum
				// 已尝试全部
				onmax := baseStart+baseItemNum > 4294967290
				if onmax {
					// 等待停止
					fmt.Print(" , #retry reset coinbase to be new block stuff# ")
					retry := &minerDeviceExecute{}
					retry.retry = true
					select {
					case work.resultCh <- retry:
					case <-time.After(time.Second * 3):
					}
					return
				}
				// 创建执行单例
				deviceIndex := i % len(mr.devices)
				device := mr.devices[deviceIndex]
				//queue := mr.queues[deviceIndex]
				work.autoidx += 1
				chindex := (mr.exeWide * deviceIndex) + int(work.autoidx%uint64(mr.exeWide))
				exe := &minerDeviceExecute{
					work.autoidx,
					uint32(deviceIndex),
					device,
					//mr.kernel,
					//queue,
					work.blkstuff,
					work.target,
					baseStart,
					uint32(mr.groupSize),
					work.blockHeight,
					work.blockCoinbaseMsg,
					work.blockCoinbaseAddr,
					false,
					false,
					nil,
				}
				//fmt.Println("<<<<<<<<<<<<<<<< ")
				select {
				case <-time.After(time.Second * 3):
					return // 三秒没有执行则退出
				case work.executeQueueChList[chindex] <- exe:
				}
			}
		}(work)
		// 执行队列
		// 执行挖矿
		work.wg.Add(mr.executeSize)
		for i := 0; i < mr.executeSize; i++ {
			go func(work *minerExecuteWork, i int) {
				defer func() {
					fmt.Printf(" , item %d-%d quit ", work.mid, i)
					work.wg.Done()
				}()
				// 运行创建执行单元
				iii1 := make([]byte, 32)
				copy(iii1, work.target[:])
				iii2 := make([]byte, 89)
				copy(iii2, work.blkstuff[:])
				// |cl.MemCopyHostPtr
				input_target, _ := mr.context.CreateBuffer(cl.MemReadOnly|cl.MemCopyHostPtr, iii1)
				//defer k
				input_stuff, _ := mr.context.CreateBuffer(cl.MemReadOnly|cl.MemCopyHostPtr, iii2)
				defer input_target.Release()
				defer input_stuff.Release()

				kernel, ke1 := mr.program.CreateKernel("miner_do_hash_x16rs_v1")
				if ke1 != nil {
					panic(ke1)
				}
				defer kernel.Release()

				device := mr.devices[i/mr.exeWide]
				queue, qe1 := mr.context.CreateCommandQueue(device, 0)
				if qe1 != nil {
					panic(qe1)
				}
				defer queue.Release()

				fmt.Printf(" , execute %d-%d start ", work.mid, i)

				for {
					// 读取一个执行单例
					// 判断停止
					if work.mid < mr.autoidx {
						return // 工作过期
					}
					// time.Sleep(time.Millisecond * 10)
					select {
					case exe := <-work.executeQueueChList[i]:
						//fmt.Println("mr.executing(exe) >>>>>>>>>>>>>>>> ", exe.baseStart)
						success := mr.executing(exe, kernel, queue, input_stuff, input_target)
						if success {
							// 成功返回
							select {
							case work.resultCh <- exe:
							case <-time.After(time.Second * 7):
							}
							return
						}
						//if exe.autoidx > 8 { panic("") }
					case <-time.After(time.Second * 7):
						// 等待超时失败
						return

					}
				}
			}(work, i)
		}
		// 等待结束
		// go func(work *minerExecuteWork) {
		work.wg.Wait()
		fmt.Printf(" , finish close execute work <%d>. ", work.mid)
		// }(work)
	}
}

//
func (mr *GpuMiner) ReStartMiner(blockHeight uint32, blkstuff [89]byte, target [32]byte,
	blockCoinbaseMsg [16]byte, blockCoinbaseAddr [21]byte,
) *minerDeviceExecute {

	mr.workLoopLock.Lock()

	// 停止之前的工作，启动一个队列任务并等待
	mr.autoidx += 1

	executeQueueChList := make([]chan *minerDeviceExecute, mr.executeSize)
	for i := 0; i < mr.executeSize; i++ {
		executeQueueChList[i] = make(chan *minerDeviceExecute, mr.executeSize*8)
	}

	// 创建挖矿工作
	minerWork := &minerExecuteWork{
		0,
		blockHeight,
		blkstuff,
		target,
		blockCoinbaseMsg,
		blockCoinbaseAddr,
		mr.autoidx,
		make(chan *minerDeviceExecute, 1),
		executeQueueChList,
		sync.WaitGroup{},
		make(chan bool, mr.executeSize+1),
	}

	// 放入工作池
	mr.minerExecuteWorkListCh <- minerWork

	mr.workLoopLock.Unlock()

	// 等待执行结果
	resexe := <-minerWork.resultCh
	return resexe
}

/*
func (mr *GpuMiner) ReStartMiner_old(blockHeight uint32, blkstuff [89]byte, target [32]byte,
	blockCoinbaseMsg [16]byte, blockCoinbaseAddr [21]byte,
) *minerDeviceExecute {
	// 写入停止标记，重新开始下一轮挖矿

	mr.stopAll()

	for i := range mr.executeQueueChList {
		for {
			select {
			case <-mr.executeQueueChList[i]:
			default:
				goto CLEAR_EXE_QUEUE_ONE
			}
		}
	CLEAR_EXE_QUEUE_ONE:
	}

	mid := rand.Uint32()
	mr.miningPrevId = mid
	mr.autoidx = 1

	fmt.Printf("\n\ndo miner height<%d>, target<%s>, block<%s>\n", blockHeight, hex.EncodeToString(target[:]), hex.EncodeToString(blkstuff[:]))

	// 队列是否空了
	miningRetCh := make(chan *minerDeviceExecute, mr.executeSize+1)

	// 投喂
	go func(mid uint32, height uint32) {
		for i := 0; ; i++ {
			select {
			case <-mr.stopMarkCh:
				fmt.Print(" , stop miner for height: ", height)
				return // 投喂结束
			default:
			}
			onmax := uint64(i)*uint64(mr.groupSize) > 4294967290
			if onmax {
				// 等待停止
				time.Sleep(time.Second * 3)
				fmt.Print(" , retry reset coinbase to be new block stuff")
				go func() {
					mr.stopAll()
					mr.miningPrevId = 0 // 防止死锁
					retry := &minerDeviceExecute{}
					retry.retry = true
					miningRetCh <- retry // 返回并重新挖矿
				}()
				continue
				//<- mr.stopMarkCh
				//return
			}
			// 创建执行单例
			deviceIndex := i % len(mr.devices)
			device := mr.devices[deviceIndex]
			//queue := mr.queues[deviceIndex]
			mr.autoidx += 1
			chindex := (mr.exeWide * deviceIndex) + int(mr.autoidx%uint64(mr.exeWide))
			exe := &minerDeviceExecute{
				mr.autoidx,
				uint32(deviceIndex),
				device,
				//mr.kernel,
				//queue,
				blkstuff,
				target,
				uint32(i) * uint32(mr.groupSize),
				uint32(mr.groupSize),
				blockHeight,
				blockCoinbaseMsg,
				blockCoinbaseAddr,
				false,
				false,
				nil,
			}
			//fmt.Println(" <<<<<<<<<<<<< mr.executeQueueCh <- exe ", exe.baseStart)
			select {
			case <-mr.stopMarkCh:
				fmt.Print(" , stop miner for height: ", height)
				return // 投喂结束
			case mr.executeQueueChList[chindex] <- exe:
			}

		}
	}(mid, blockHeight)

	// 执行挖矿
	for i := 0; i < mr.executeSize; i++ {
		go func(mid uint32, i int) {
			iii1 := make([]byte, 32)
			copy(iii1, target[:])
			iii2 := make([]byte, 89)
			copy(iii2, blkstuff[:])
			// |cl.MemCopyHostPtr
			input_target, _ := mr.context.CreateBuffer(cl.MemReadOnly|cl.MemCopyHostPtr, iii1)
			//defer k
			input_stuff, _ := mr.context.CreateBuffer(cl.MemReadOnly|cl.MemCopyHostPtr, iii2)
			defer input_target.Release()
			defer input_stuff.Release()

			kernel, ke1 := mr.program.CreateKernel("miner_do_hash_x16rs_v1")
			if ke1 != nil {
				panic(ke1)
			}
			defer kernel.Release()

			device := mr.devices[i/mr.exeWide]
			queue, qe1 := mr.context.CreateCommandQueue(device, 0)
			if qe1 != nil {
				panic(qe1)
			}
			defer queue.Release()

			for {
				// 读取一个执行单例
				select {
				case <-mr.stopMarkCh:
					fmt.Print(" , stop miner queue ", i)
					return // 挖矿结束
				case exe := <-mr.executeQueueChList[i]:
					//fmt.Println("mr.executing(exe) >>>>>>>>>>>>>>>> ", exe.baseStart)
					success := mr.executing(exe, kernel, queue, input_stuff, input_target)
					if success {
						go func() {
							// 成功
							mr.stopAll()       // 停止所有挖矿
							miningRetCh <- exe // 返回成功
						}()
					}
				case <-time.After(time.Second * 1): // 等待队列
				}
			}
		}(mid, i)
	}

	return <-miningRetCh

	// 删除停止标记， 返回内容
	//delete(mr.stopMark, mid)
}


*/

// 执行
func (mr *GpuMiner) executing(exe *minerDeviceExecute, kernel *cl.Kernel, queue *cl.CommandQueue, input_stuff *cl.MemObject, input_target *cl.MemObject) bool {

	//local, _ := kernel.WorkGroupSize(exe.device)
	//global := int(exe.groupSize)
	//d := int(exe.groupSize) % local
	//if d != 0 {
	//	global += local - d
	//}

	/////////test/////////
	//global = mr.groupSize
	//local = mr.groupSize
	///////test end///////

	if exe.autoidx%uint64(mr.printNumBase) == 0 {
		fmt.Printf(",%d-%d", exe.deviceIndex, exe.baseStart)
	}
	nonce, reshash, success := mr.doGroupWork(exe, kernel, queue, input_target, input_stuff, mr.groupSize, mr.groupSize, exe.baseStart)

	if success {
		noncenum := binary.BigEndian.Uint32(nonce)
		fmt.Printf("\nheight: %d, nonce: %d<%s>[%d,%d,%d,%d], hash: %s, miner success!\n",
			exe.blockHeight,
			noncenum,
			hex.EncodeToString(nonce),
			nonce[0], nonce[1], nonce[2], nonce[3],
			hex.EncodeToString(reshash),
		)
		// 挖矿成功并返回
		exe.success = true
		exe.nonce = nonce
		return true
	}

	return false
}

// 启动分组
func (mr *GpuMiner) doGroupWork(exe *minerDeviceExecute, kernel *cl.Kernel, queue *cl.CommandQueue, input_target *cl.MemObject, input_stuff *cl.MemObject, global int, local int, base_start uint32) ([]byte, []byte, bool) {

	// 参数
	// |cl.MemAllocHostPtr
	output_nonce, _ := mr.context.CreateEmptyBuffer(cl.MemReadWrite|cl.MemAllocHostPtr, 4)
	output_hash, _ := mr.context.CreateEmptyBuffer(cl.MemReadWrite|cl.MemAllocHostPtr, 32)
	defer output_nonce.Release()
	defer output_hash.Release()

	queue.EnqueueWriteBufferByte(output_nonce, true, 0, []byte{0, 0, 0, 0}, nil)
	// set argvs
	kernel.SetArgs(input_target, input_stuff, uint32(base_start), uint32(mr.loopNum), output_nonce, output_hash)
	defer kernel.Retain()

	// run
	queue.EnqueueNDRangeKernel(kernel, nil, []int{global}, []int{local}, nil)
	queue.Finish()

	result_nonce := bytes.Repeat([]byte{0}, 4)
	result_hash := make([]byte, 32)
	// copy get output
	queue.EnqueueReadBufferByte(output_nonce, true, 0, result_nonce, nil)
	queue.EnqueueReadBufferByte(output_hash, true, 0, result_hash, nil)
	queue.Flush()
	defer queue.Retain()
	//defer queue.Release()

	//fmt.Println(result_nonce)
	nonce := binary.BigEndian.Uint32(result_nonce)
	if nonce > 0 {
		// check results
		// fmt.Println("==========================", nonce, result_nonce)
		// fmt.Println("output_hash", result_hash, hex.EncodeToString(result_hash))
		// return
		return result_nonce, result_hash, true
	}
	return nil, nil, false

}

func (mr *GpuMiner) InitBuildProgram(openclPath string, platName string, dvid int, groupSize int, loopNum int, exeWide int, printNumBase int, rebuild bool) error {

	//mr.miningPrevId = 0
	mr.autoidx = 0

	mr.openclPath = openclPath
	mr.platName = platName
	mr.dvid = dvid
	mr.groupSize = groupSize
	mr.loopNum = loopNum
	mr.exeWide = exeWide
	mr.printNumBase = printNumBase
	mr.rebuild = rebuild

	// 工作池
	mr.minerExecuteWorkListCh = make(chan *minerExecuteWork, 16)

	// 执行工作
	go mr.minerWorkLoop()

	var err error

	// init
	platids := 0

	platforms, _ := cl.GetPlatforms()

	if len(platforms) == 0 {
		return fmt.Errorf("not find any platforms.")
	}
	for i, pt := range platforms {
		fmt.Printf("  - platform %d: %s\n", i, pt.Name())
		if strings.Compare(mr.platName, "") != 0 && strings.Contains(pt.Name(), mr.platName) {
			platids = i
		}
	}

	mr.platform = platforms[platids]

	fmt.Printf("current use platform: %s\n", mr.platform.Name())

	devices, _ := mr.platform.GetDevices(cl.DeviceTypeAll)

	if len(devices) == 0 {
		return fmt.Errorf("not find any devices.")
	}

	for i, dv := range devices {
		fmt.Printf("  - device %d: %s, (max_work_group_size: %d)\n", i, dv.Name(), dv.MaxWorkGroupSize())
	}
	mr.devices = devices
	if dvid > -1 && dvid < len(devices) {
		mr.devices = []*cl.Device{devices[dvid]}
		fmt.Printf("current use device %d: %s\n", dvid, devices[dvid].Name())
	} else {
		fmt.Printf("current use all %d devices.\n", len(devices))
	}

	// 队列大小
	mr.executeSize = len(mr.devices) * mr.exeWide

	//mr.executeQueueChList = make([]chan *minerDeviceExecute, mr.executeSize)
	//for i := 0; i < mr.executeSize; i++ {
	//	mr.executeQueueChList[i] = make(chan *minerDeviceExecute, mr.executeSize*8)
	//}

	// 停止标记
	//mr.stopMarkCh = make(chan uint32) // , mr.executeSize + 1

	if mr.context, err = cl.CreateContext(mr.devices); err != nil {
		panic(err)
	}

	if strings.Compare(mr.openclPath, "") == 0 {
		mr.openclPath = GetCurrentDirectory() + "/opencl"
	}

	fmt.Println("create building opencl program from dir " + mr.openclPath + ", please wait...")
	//bderr := mr.program.BuildProgram(nil, "-I "+mr.openclPath) // -I /media/yangjie/500GB/Hacash/src/github.com/hacash/x16rs/opencl
	//if bderr != nil {
	//	panic(bderr)
	//}

	mr.program = mr.buildOrLoadProgram()

	fmt.Println("create program complete.")
	//
	//kernel, ke1 := mr.program.CreateKernel("miner_do_hash_x16rs_v1")
	//if ke1 != nil {
	//	panic(ke1)
	//}
	//mr.kernel = kernel
	//
	//mr.queues = make([]*cl.CommandQueue, len(mr.devices))
	//for k, d := range mr.devices {
	//	queue, qe1 := mr.context.CreateCommandQueue(d, 0)
	//	if qe1 != nil {
	//		panic(qe1)
	//	}
	//	mr.queues[k] = queue
	//}

	return nil
}

func (mr *GpuMiner) buildOrLoadProgram() *cl.Program {

	var program *cl.Program

	binfilestuff := mr.platform.Name() // + "_" + mr.devices[0].Name()
	binfilename := strings.Replace(binfilestuff, " ", "_", -1)
	binfilepath := mr.openclPath + "/" + binfilename + ".objcache"
	binstat, staterr := os.Stat(binfilepath)
	if mr.rebuild || staterr != nil {
		program, _ = mr.context.CreateProgramWithSource([]string{` #include "x16rs_main.cl" `})
		bderr := program.BuildProgram(nil, "-I "+mr.openclPath) // -I /media/yangjie/500GB/Hacash/src/github.com/hacash/x16rs/opencl
		if bderr != nil {
			panic(bderr)
		}
		fmt.Println("build complete get binaries...")
		//fmt.Println("program.GetBinarySizes_2()")
		sizes, _ := program.GetBinarySizes_2(1)
		//fmt.Println(sizes)
		//fmt.Println(sizes[0])
		//fmt.Println("program.GetBinaries_2()")
		bins, _ := program.GetBinaries_2([]int{sizes[0]})
		//fmt.Println(bins[0])
		f, e := os.OpenFile(binfilepath, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0777)
		if e != nil {
			panic(e)
		}
		//fmt.Println("f.Write(wbin) "+binfilepath, sizes[0])
		f.Write(bins[0])

	} else {
		fmt.Printf("load binary program file from \"%s\"\n", binfilepath)
		file, _ := os.OpenFile(binfilepath, os.O_RDWR, 0777)
		bin := make([]byte, binstat.Size())
		//fmt.Println("file.Read(bin) size", binstat.Size())
		file.Read(bin)
		//fmt.Println(bin)
		// 仅仅支持同一个平台的同一种设备
		bins := make([][]byte, len(mr.devices))
		sizes := make([]int, len(mr.devices))
		for k, _ := range mr.devices {
			bins[k] = bin
			sizes[k] = int(binstat.Size())
		}
		fmt.Println("create program with binary...")
		var berr error
		program, berr = mr.context.CreateProgramWithBinary_2(mr.devices, sizes, bins)
		if berr != nil {
			panic(berr)
		}
		program.BuildProgram(mr.devices, "")
		//fmt.Println("context.CreateProgramWithBinary")
	}
	// 返回
	return program
}

/*



// 开始、重新开始挖矿
func (mr *GpuMiner) DoMiner(blockHeight uint32, blkstuff [89]byte, target [32]byte, resultCh chan MinerResult) {

	// 写入停止标记，重新开始下一轮挖矿
	mr.stopMark[mr.miningPrevId] = true
	mid := rand.Uint32()
	mr.miningPrevId = mid

	// 开始挖矿
	go func(mid uint32, height uint32, blkstuff [89]byte, target [32]byte, resultCh chan MinerResult) {

		fmt.Printf("\n\ndo miner height<%d>, target<%s>, block<%s>\n", height, hex.EncodeToString(target[:]), hex.EncodeToString(blkstuff[:]))

		input_target, _ := mr.context.CreateEmptyBuffer(cl.MemReadOnly, 32)
		input_stuff, _ := mr.context.CreateEmptyBuffer(cl.MemReadOnly, 89)

		// 参数
		mr.queue.EnqueueWriteBufferByte(input_target, true, 0, target[:], nil)
		mr.queue.EnqueueWriteBufferByte(input_stuff, true, 0, blkstuff[:], nil)

		local, _ := mr.kernel.WorkGroupSize(mr.device)
		global := mr.groupSize
		d := mr.groupSize % local
		if d != 0 {
			global += local - d
		}

		// 循环挖矿
		var i uint32
		var nonce []byte = nil
		var reshash []byte = nil
		var success bool = false
		for i = 0; ; i++ {
			if uint64(i) * uint64(mr.groupSize) > 4294967290 {
				success = false
				break
			}
			base_start := i * uint32(mr.groupSize)
			if i%10 == 0 {
				fmt.Printf(",%d", base_start)
			}
			nonce, reshash, success = mr.doGroupWork(input_target, input_stuff, global, local, base_start)
			if success {
				noncenum := binary.BigEndian.Uint32(nonce)
				fmt.Printf("\nheight: %d, nonce: %d<%s>[%d,%d,%d,%d], hash: %s, miner success!\n",
					noncenum,
					hex.EncodeToString(nonce),
					nonce[0], nonce[1], nonce[2], nonce[3],
					height,
					hex.EncodeToString(reshash),
				)
				break // 成功
			}
			// 结束判断
			if stop, ok := mr.stopMark[mid]; ok && stop {
				delete(mr.stopMark, mid) // 删除结束标记
				fmt.Printf("\nstop miner height: %d\n", height)
				success = false
				nonce = nil
				break
			}
			// 继续下一组挖矿
		}

		fmt.Printf("height: %d miner finish with %t\n", height, success)

		// 挖矿状态，返回
		resultCh <- MinerResult{
			success,
			height,
			nonce,
		}

	}(mid, blockHeight, blkstuff, target, resultCh)

}

*/

func GetCurrentDirectory() string {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0])) //返回绝对路径  filepath.Dir(os.Args[0])去除最后一个元素的路径
	if err != nil {
		log.Fatal(err)
	}
	return strings.Replace(dir, "\\", "/", -1) //将\替换成/
}
