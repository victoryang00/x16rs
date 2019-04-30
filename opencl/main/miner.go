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
)

// 单个设备执行单例
type minerDeviceExecute struct {
	device         *cl.Device
	blockHeadBytes [89]byte
	targetHash     [32]byte
}

type GpuMiner struct {
	platform *cl.Platform
	context  *cl.Context
	program  *cl.Program
	kernel   *cl.Kernel
	devices  []*cl.Device // 所有设备

	executes map[int]*cl.Device // 正在执行的设备

	device *cl.Device
	queue  *cl.CommandQueue

	// data
	blockHeadBytes []byte
	targetHash     []byte

	// config
	openclPath string
	platName   string // 平台名称
	dvid       int    // 设备id
	groupSize  int

	// msg
	miningPrevId uint32
	stopMark     map[uint32]bool
}

type MinerResult struct {
	success bool // 是否挖矿成功
	height  uint32
	nonce   []byte
}

func (mr *GpuMiner) InitBuildProgram(openclPath string, platName string, dvid int, groupSize int) error {

	mr.miningPrevId = 0
	mr.stopMark = make(map[uint32]bool)

	mr.openclPath = openclPath
	mr.platName = platName
	mr.dvid = dvid
	mr.groupSize = groupSize

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

	mr.device = devices[mr.dvid]
	fmt.Printf("current use device: %s\n", mr.device.Name())

	mr.context, _ = cl.CreateContext(devices)
	mr.program, _ = mr.context.CreateProgramWithSource([]string{` #include "x16rs_main.cl" `})

	if strings.Compare(mr.openclPath, "") == 0 {
		mr.openclPath = GetCurrentDirectory() + "/opencl"
	}

	fmt.Println("building opencl program from dir " + mr.openclPath + ", please wait...")
	bderr := mr.program.BuildProgram(nil, "-I "+mr.openclPath) // -I /media/yangjie/500GB/Hacash/src/github.com/hacash/x16rs/opencl
	if bderr != nil {
		panic(bderr)
	}

	fmt.Println("build complete.")
	mr.kernel, _ = mr.program.CreateKernel("miner_do_hash_x16rs_v1")

	return nil
}

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
		var success bool = false
		for i = 0; ; i++ {
			base_start := i * uint32(mr.groupSize)
			if i%50 == 0 {
				fmt.Printf(",%d", base_start)
			}
			nonce, success = mr.doGroupWork(input_target, input_stuff, global, local, base_start)
			if success {
				noncenum := binary.BigEndian.Uint32(nonce)
				fmt.Printf("\nnonce %d<%s>[%d,%d,%d,%d] height<%d> miner success!\n",
					noncenum,
					hex.EncodeToString(nonce),
					nonce[0], nonce[1], nonce[2], nonce[3],
					height,
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

// 启动分组
func (mr *GpuMiner) doGroupWork(input_target *cl.MemObject, input_stuff *cl.MemObject, global int, local int, base_start uint32) ([]byte, bool) {

	mr.queue, _ = mr.context.CreateCommandQueue(mr.device, 0)

	output_nonce, _ := mr.context.CreateEmptyBuffer(cl.MemReadWrite|cl.MemAllocHostPtr, 4)
	output_hash, _ := mr.context.CreateEmptyBuffer(cl.MemReadWrite|cl.MemAllocHostPtr, 32)
	mr.queue.EnqueueWriteBufferByte(output_nonce, true, 0, []byte{0, 0, 0, 0}, nil)
	// set argvs
	mr.kernel.SetArgs(input_target, input_stuff, uint32(base_start), output_nonce, output_hash)

	// run
	mr.queue.EnqueueNDRangeKernel(mr.kernel, nil, []int{global}, []int{local}, nil)
	mr.queue.Finish()

	result_nonce := bytes.Repeat([]byte{0}, 4)
	result_hash := make([]byte, 32)
	// copy get output
	mr.queue.EnqueueReadBufferByte(output_nonce, true, 0, result_nonce, nil)
	mr.queue.EnqueueReadBufferByte(output_hash, true, 0, result_hash, nil)

	//fmt.Println(result_nonce)
	nonce := binary.BigEndian.Uint32(result_nonce)
	if nonce > 0 {
		// check results
		// fmt.Println("==========================", nonce, result_nonce)
		// fmt.Println("output_hash", result_hash, hex.EncodeToString(result_hash))
		// return
		return result_nonce, true
	}
	return nil, false
}

func GetCurrentDirectory() string {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0])) //返回绝对路径  filepath.Dir(os.Args[0])去除最后一个元素的路径
	if err != nil {
		log.Fatal(err)
	}
	return strings.Replace(dir, "\\", "/", -1) //将\替换成/
}
