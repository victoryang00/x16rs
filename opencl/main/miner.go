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

	workContext *WorkContext

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

	workContexts []*WorkContext

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

type WorkContext struct {
	device       *cl.Device
	kernel       *cl.Kernel
	queue        *cl.CommandQueue
	input_stuff  *cl.MemObject
	input_target *cl.MemObject
	output_nonce *cl.MemObject
	output_hash  *cl.MemObject
}

func (w *WorkContext) Retain() {
	w.kernel.Retain()
	w.queue.Retain()
	w.input_stuff.Retain()
	w.input_target.Retain()
	w.output_nonce.Retain()
	w.output_hash.Retain()
}

func (w *WorkContext) ReInit(stuff_buf []byte, target_buf []byte) {
	w.queue.EnqueueWriteBufferByte(w.input_stuff, true, 0, stuff_buf, nil)
	w.queue.EnqueueWriteBufferByte(w.input_target, true, 0, target_buf, nil)

}

func (mr *GpuMiner) ReInitWorkContext(stuff_buf []byte, target_buf []byte) {
	for _, ctx := range mr.workContexts {
		ctx.ReInit(stuff_buf, target_buf)
	}
}

// chua
func (mr *GpuMiner) createWorkContext(queue_index int) *WorkContext {

	// 运行创建执行单元
	//input_target_buf := make([]byte, 32)
	//copy(input_target_buf, work.target[:])
	//input_stuff_buf := make([]byte, 89)
	//copy(input_stuff_buf, work.blkstuff[:])
	// |cl.MemCopyHostPtr
	input_target, _ := mr.context.CreateEmptyBuffer(cl.MemReadOnly, 32)
	input_stuff, _ := mr.context.CreateEmptyBuffer(cl.MemReadOnly, 89)
	//defer input_target.Release()
	//defer input_stuff.Release()

	// 参数
	// |cl.MemAllocHostPtr
	output_nonce, _ := mr.context.CreateEmptyBuffer(cl.MemReadWrite|cl.MemAllocHostPtr, 4)
	output_hash, _ := mr.context.CreateEmptyBuffer(cl.MemReadWrite|cl.MemAllocHostPtr, 32)
	//defer output_nonce.Release()
	//defer output_hash.Release()

	kernel, ke1 := mr.program.CreateKernel("miner_do_hash_x16rs_v1")
	if ke1 != nil {
		panic(ke1)
	}
	//defer kernel.Release()

	device := mr.devices[queue_index/mr.exeWide]
	queue, qe1 := mr.context.CreateCommandQueue(device, 0)
	if qe1 != nil {
		panic(qe1)
	}
	//defer queue.Release()

	return &WorkContext{
		device,
		kernel,
		queue,
		input_stuff,
		input_target,
		output_nonce,
		output_hash,
	}

}

////////////////////////////////////////////

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
		make(chan *minerDeviceExecute, mr.executeSize),
		executeQueueChList,
		sync.WaitGroup{},
	}

	// 放入工作池
	mr.minerExecuteWorkListCh <- minerWork

	mr.workLoopLock.Unlock()

	// 等待执行结果
	resexe := <-minerWork.resultCh
	return resexe
}

func (mr *GpuMiner) minerWorkLoop() {
	for {
		work := <-mr.minerExecuteWorkListCh
		// 判断当前执行id
		if work.mid < mr.autoidx {
			continue // 永远只执行最后一个任务
		}
		// 执行
		fmt.Printf("\ndo miner work<%d>, height<%d>, target<%s>, block<%s>\n", work.mid, work.blockHeight, hex.EncodeToString(work.target[:]), hex.EncodeToString(work.blkstuff[:]))

		// 初始化状态
		input_stuff_buf := make([]byte, 89)
		copy(input_stuff_buf, work.blkstuff[:])
		input_target_buf := make([]byte, 32)
		copy(input_target_buf, work.target[:])
		mr.ReInitWorkContext(input_stuff_buf, input_target_buf)

		// 投喂执行单元
		work.wg.Add(1)
		go func(work *minerExecuteWork) {
			defer func() {
				fmt.Printf(", queue <%d> exit ", work.mid)
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
					fmt.Print(", ## retry reset coinbase to be new block stuff ## ")
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
					nil,
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
					fmt.Printf(", item <%d>%d quit ", work.mid, i)
					work.wg.Done()
				}()
				// 载入环境
				var work_ctx = mr.workContexts[i]

				fmt.Printf(", execute <%d>%d start ", work.mid, i)

				for {
					// 读取一个执行单例
					// 判断停止
					if work.mid < mr.autoidx {
						return // 工作过期
					}
					// time.Sleep(time.Millisecond * 10)
					select {
					case exe := <-work.executeQueueChList[i]:
						exe.workContext = work_ctx
						//fmt.Println("mr.executing(exe) >>>>>>>>>>>>>>>> ", exe.baseStart)
						success := mr.executing(exe)
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
		fmt.Printf(", work <%d> finish closed.\n", work.mid)
		// }(work)
	}
}

// 执行
func (mr *GpuMiner) executing(exe *minerDeviceExecute) bool {

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
	nonce, reshash, success := mr.doGroupWork(exe, mr.groupSize, mr.groupSize, exe.baseStart)

	if success {
		noncenum := binary.BigEndian.Uint32(nonce)
		fmt.Printf("\n\n⬤  ㄜ height: %d, nonce: %d<%s>[%d,%d,%d,%d], hash: %s, miner success!\n\n",
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
func (mr *GpuMiner) doGroupWork(exe *minerDeviceExecute, global int, local int, base_start uint32) ([]byte, []byte, bool) {

	ctx := exe.workContext
	// 结束后复用
	defer ctx.Retain()

	// 重置
	ctx.queue.EnqueueWriteBufferByte(ctx.output_nonce, true, 0, []byte{0, 0, 0, 0}, nil)
	// set argvs
	ctx.kernel.SetArgs(ctx.input_target, ctx.input_stuff, uint32(base_start), uint32(mr.loopNum), ctx.output_nonce, ctx.output_hash)

	// run
	ctx.queue.EnqueueNDRangeKernel(ctx.kernel, nil, []int{global}, []int{local}, nil)
	ctx.queue.Finish()

	result_nonce := bytes.Repeat([]byte{0}, 4)
	result_hash := make([]byte, 32)
	// copy get output
	ctx.queue.EnqueueReadBufferByte(ctx.output_nonce, true, 0, result_nonce, nil)
	ctx.queue.EnqueueReadBufferByte(ctx.output_hash, true, 0, result_hash, nil)
	ctx.queue.Flush()

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

	fmt.Println("program complete.")
	//

	// 队列宽度
	mr.executeSize = len(mr.devices) * mr.exeWide
	mr.workContexts = make([]*WorkContext, mr.executeSize)
	for i := 0; i < mr.executeSize; i++ {
		mr.workContexts[i] = mr.createWorkContext(i)
	}

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

func GetCurrentDirectory() string {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0])) //返回绝对路径  filepath.Dir(os.Args[0])去除最后一个元素的路径
	if err != nil {
		log.Fatal(err)
	}
	return strings.Replace(dir, "\\", "/", -1) //将\替换成/
}
