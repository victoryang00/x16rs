package main

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"github.com/xfong/go2opencl/cl"
	"log"
	"math/rand"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// 单个设备执行单例
type minerDeviceExecute struct {
	autoidx           uint64
	deviceIndex       uint32
	device            *cl.Device
	blockHeadBytes    [89]byte
	targetHash        [32]byte
	baseStart         uint32
	groupSize         uint32
	blockHeight       uint32
	blockCoinbaseMsg  [16]byte
	blockCoinbaseAddr [21]byte
	// 挖矿状态
	retry   bool // 重新尝试挖矿
	success bool
	nonce   []byte
}

type GpuMiner struct {
	autoidx uint64

	platform *cl.Platform
	context  *cl.Context
	program  *cl.Program
	devices  []*cl.Device // 所有设备

	// 执行队列
	executeQueueCh chan *minerDeviceExecute

	///////////////////////////////

	// device *cl.Device
	// queue  *cl.CommandQueue
	// kernel   *cl.Kernel

	// data
	blockHeadBytes []byte
	targetHash     []byte

	// config
	openclPath  string
	platName    string // 平台名称
	dvid        int    // 设备id
	groupSize   int    // 组大小
	executeSize int    // 执行单例数量
	rebuild     bool   // 强制重新编译

	// msg
	miningPrevId uint32
	stopMark     map[uint32]bool
}

type MinerResult struct {
	success bool // 是否挖矿成功
	height  uint32
	nonce   []byte
}

func (mr *GpuMiner) ReStartMiner(blockHeight uint32, blkstuff [89]byte, target [32]byte,
	blockCoinbaseMsg [16]byte, blockCoinbaseAddr [21]byte,
) *minerDeviceExecute {
	// 写入停止标记，重新开始下一轮挖矿
	mr.stopMark[mr.miningPrevId] = true
	mid := rand.Uint32()
	mr.miningPrevId = mid
	// 清空待执行队列
	for {
		select {
		case <-mr.executeQueueCh:
		default:
			goto MINERSTART
		}
	}
MINERSTART:

	fmt.Printf("\n\ndo miner height<%d>, target<%s>, block<%s>\n", blockHeight, hex.EncodeToString(target[:]), hex.EncodeToString(blkstuff[:]))

	// 队列是否空了
	miningRetCh := make(chan *minerDeviceExecute)

	// 投喂
	go func(mid uint32, height uint32) {
		for i := 0; ; i++ {
			onmax := uint64(i)*uint64(mr.groupSize) > 4294967290
			if o1, o2 := mr.stopMark[mid]; o1 && o2 || onmax {
				fmt.Println("\nstop miner height", height)
				break // 投喂结束
			}
			// 创建执行单例
			deviceIndex := i % len(mr.devices)
			device := mr.devices[deviceIndex]
			mr.autoidx += 1
			exe := &minerDeviceExecute{
				mr.autoidx,
				uint32(deviceIndex),
				device,
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
			// fmt.Println(" <<<<<<<<<<<<< mr.executeQueueCh <- exe ", exe.baseStart)
			mr.executeQueueCh <- exe
		}
	}(mid, blockHeight)

	// 执行挖矿
	for i := 0; i < mr.executeSize; i++ {
		go func(mid uint32, i int) {
			for {
				if o1, o2 := mr.stopMark[mid]; o1 && o2 {
					return // 挖矿结束
				}
				// 读取一个执行单例
				select {
				case exe := <-mr.executeQueueCh:
					success := mr.executing(exe)
					// fmt.Println("mr.executing(exe) >>>>>>>>>>>>>>>> ", exe.baseStart, success)
					if success {
						// 成功
						miningRetCh <- exe // 返回成功
						return
					}
				case <-time.After(time.Second): // 队列空了
					fmt.Println("  retry reset coinbase to be new block stuff")
					retry := &minerDeviceExecute{}
					retry.retry = true
					miningRetCh <- retry // 返回并重新挖矿
				}
			}
		}(mid, i)
	}

	return <-miningRetCh

	// 删除停止标记， 返回内容
	//delete(mr.stopMark, mid)
}

// 执行
func (mr *GpuMiner) executing(exe *minerDeviceExecute) bool {

	queue, qe1 := mr.context.CreateCommandQueue(exe.device, 0)
	if qe1 != nil {
		panic(qe1)
	}
	//defer queue.Release() // 释放资源

	kernel, ke1 := mr.program.CreateKernel("miner_do_hash_x16rs_v1")
	if ke1 != nil {
		panic(ke1)
	}
	defer kernel.Release()

	iii1 := make([]byte, 32)
	copy(iii1, exe.targetHash[:])
	iii2 := make([]byte, 89)
	copy(iii2, exe.blockHeadBytes[:])
	// fmt.Println(iii1, iii2)
	input_target, ie1 := mr.context.CreateBuffer(cl.MemReadOnly|cl.MemCopyHostPtr, iii1)
	if ie1 != nil {
		panic(ie1)
	}
	//defer k
	input_stuff, ie2 := mr.context.CreateBuffer(cl.MemReadOnly|cl.MemCopyHostPtr, iii2)
	if ie2 != nil {
		panic(ie2)
	}
	defer input_target.Release()
	defer input_stuff.Release()

	local, _ := kernel.WorkGroupSize(exe.device)
	global := int(exe.groupSize)
	d := int(exe.groupSize) % local
	if d != 0 {
		global += local - d
	}

	if exe.autoidx%20 == 0 {
		fmt.Printf(",%d-%d", exe.deviceIndex, exe.baseStart)
	}
	nonce, reshash, success := mr.doGroupWork(kernel, queue, input_target, input_stuff, global, local, exe.baseStart)

	if success {
		noncenum := binary.BigEndian.Uint32(nonce)
		fmt.Printf("\nheight: %d, nonce: %d<%s>[%d,%d,%d,%d], hash: %s, miner success!\n",
			noncenum,
			hex.EncodeToString(nonce),
			nonce[0], nonce[1], nonce[2], nonce[3],
			exe.blockHeight,
			hex.EncodeToString(reshash),
		)
		// 挖矿成功并返回
		exe.success = true
		exe.nonce = nonce
		return true
	}

	return false
}

func (mr *GpuMiner) InitBuildProgram(openclPath string, platName string, dvid int, groupSize int, exeWide int, rebuild bool) error {

	mr.miningPrevId = 0
	mr.stopMark = make(map[uint32]bool)

	mr.openclPath = openclPath
	mr.platName = platName
	mr.dvid = dvid
	mr.groupSize = groupSize
	mr.rebuild = rebuild

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
		fmt.Printf("  - device %d: %s\n", i, dv.Name())
	}
	mr.devices = devices
	if dvid > -1 && dvid < len(devices) {
		mr.devices = []*cl.Device{devices[dvid]}
		fmt.Printf("current use device: %s\n", devices[dvid].Name())
	} else {
		fmt.Printf("current use all %d devices.\n", len(devices))
	}

	// 队列大小
	mr.executeSize = len(mr.devices) * exeWide
	mr.executeQueueCh = make(chan *minerDeviceExecute, mr.executeSize*8)

	if mr.context, err = cl.CreateContext(mr.devices); err != nil {
		panic(err)
	}

	if strings.Compare(mr.openclPath, "") == 0 {
		mr.openclPath = GetCurrentDirectory() + "/opencl"
	}

	fmt.Println("building opencl program from dir " + mr.openclPath + ", please wait...")
	//bderr := mr.program.BuildProgram(nil, "-I "+mr.openclPath) // -I /media/yangjie/500GB/Hacash/src/github.com/hacash/x16rs/opencl
	//if bderr != nil {
	//	panic(bderr)
	//}

	mr.program = mr.buildOrLoadProgram()

	fmt.Println("build complete.")

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

// 启动分组
func (mr *GpuMiner) doGroupWork(kernel *cl.Kernel, queue *cl.CommandQueue, input_target *cl.MemObject, input_stuff *cl.MemObject, global int, local int, base_start uint32) ([]byte, []byte, bool) {

	output_nonce, _ := mr.context.CreateEmptyBuffer(cl.MemReadWrite|cl.MemAllocHostPtr, 4)
	output_hash, _ := mr.context.CreateEmptyBuffer(cl.MemReadWrite|cl.MemAllocHostPtr, 32)
	// 释放资源
	defer output_hash.Release()
	defer output_hash.Release()

	queue.EnqueueWriteBufferByte(output_nonce, true, 0, []byte{0, 0, 0, 0}, nil)
	// set argvs
	kernel.SetArgs(input_target, input_stuff, uint32(base_start), output_nonce, output_hash)

	// run
	queue.EnqueueNDRangeKernel(kernel, nil, []int{global}, []int{local}, nil)
	queue.Finish()

	result_nonce := bytes.Repeat([]byte{0}, 4)
	result_hash := make([]byte, 32)
	// copy get output
	queue.EnqueueReadBufferByte(output_nonce, true, 0, result_nonce, nil)
	queue.EnqueueReadBufferByte(output_hash, true, 0, result_hash, nil)

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

func GetCurrentDirectory() string {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0])) //返回绝对路径  filepath.Dir(os.Args[0])去除最后一个元素的路径
	if err != nil {
		log.Fatal(err)
	}
	return strings.Replace(dir, "\\", "/", -1) //将\替换成/
}
