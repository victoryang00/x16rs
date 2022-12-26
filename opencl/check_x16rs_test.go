package opencl

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/hacash/x16rs"
	"github.com/xfong/go2opencl/cl"
	"golang.org/x/crypto/sha3"
)

func Test1(t *testing.T) {

	input_hash, _ := hex.DecodeString("899b7aed0bc68793e479245683c32d15cd92a4a1babb0f9ffd08cfb51339950d")

	for k := 0; k < 255; k++ {
		for i := 0; i < 255; i++ {
			input_hash[28] = uint8(i)
			input_hash[29] = uint8(k)
			check_step_by_input_hash(1, input_hash)
			fmt.Println("--------------------------------------------", k, i)
		}
	}

}

func Test111(t *testing.T) {

	input_hash, _ := hex.DecodeString("899b7aed0bc68793e479245683c32d15cd92a4a1babb0f9ffd08cfb51339950d")
	input_hash[28] = uint8(0) // 错误 3、4、5
	check_step_by_input_hash(1, input_hash)

}

func Test2(t *testing.T) {

	input_stuff_89, _ := hex.DecodeString("01000000008500605eeb7806eec8d3aec85b91dd17a999400294f3566320b782964f9e7b46fd1dfc3712157348605b995344856356406b785fa0457ec77a403f099ef67b8a2b07a645f2670000000100000000fafb58bb0000")

	check_step_pre_hash(input_stuff_89)

}

func Test3(t *testing.T) {

	input_stuff_89, _ := hex.DecodeString("01000000008500605eeb7806eec8d3aec85b91dd17a999400294f3566320b782964f9e7b46fd1dfc3712157348605b995344856356406b785fa0457ec77a403f099ef67b8a2b07a645f2670000000100000077fafb58bb0000")

	mindemhash := check_step_pre_hash(input_stuff_89)

	x16rsrepeat := x16rs.HashRepeatForBlockHeight(133)
	endhash := check_step_by_input_hash(x16rsrepeat, mindemhash)

	fmt.Println("hex.EncodeToString( endhash ) =======================")
	fmt.Println(hex.EncodeToString(endhash))
}

func check_step_pre_hash(input_stuff_89 []byte) []byte {

	//input_hash, _ := hex.DecodeString("4906f613be6708dca0ed8222368acc477036919485059c01a0735092474fe485")

	device, _, kernel, context := buildOrLoadProgram("/home/victoryang00/Documents/x16rs/opencl", 0, 0, false, "check_x16rs_prehash")

	queue, _ := context.CreateCommandQueue(device, 0)

	input_hash_param, _ := context.CreateEmptyBuffer(cl.MemReadOnly, 89)
	output_hash_param, _ := context.CreateEmptyBuffer(cl.MemReadWrite|cl.MemAllocHostPtr, 32)

	queue.EnqueueWriteBufferByte(input_hash_param, true, 0, input_stuff_89, nil)

	// set argvs
	kernel.SetArgs(input_hash_param, output_hash_param)

	// run
	queue.EnqueueNDRangeKernel(kernel, []int{0}, []int{1}, []int{1}, nil)
	queue.Finish()

	result_hash := make([]byte, 32)
	queue.EnqueueReadBufferByte(output_hash_param, true, 0, result_hash, nil)
	queue.Flush()

	//
	fmt.Println("check_step_pre_hash result_hash_1:", hex.EncodeToString(result_hash))

	hashbase := sha3.Sum256(input_stuff_89)
	result_hash_2 := hashbase[:]

	fmt.Println("check_step_pre_hash result_hash_2:", hex.EncodeToString(result_hash_2))

	if bytes.Compare(result_hash, result_hash_2) != 0 {
		panic("check_step_pre_hash bytes.Compare(result_hash, result_hash_2) != 0")
	}

	return result_hash_2

}

func check_step_by_input_hash(x16rsrepeat int, input_hash []byte) []byte {

	//input_hash, _ := hex.DecodeString("4906f613be6708dca0ed8222368acc477036919485059c01a0735092474fe485")

	device, _, kernel, context := buildOrLoadProgram("/home/victoryang00/Documents/x16rs/opencl", 0, 0, false, "check_x16rs_step")

	queue, _ := context.CreateCommandQueue(device, 0)

	input_hash_param, _ := context.CreateEmptyBuffer(cl.MemReadOnly, 32)
	output_hash_param, _ := context.CreateEmptyBuffer(cl.MemReadWrite|cl.MemAllocHostPtr, 32)

	queue.EnqueueWriteBufferByte(input_hash_param, true, 0, input_hash, nil)

	// set argvs
	kernel.SetArgs(uint32(x16rsrepeat), input_hash_param, output_hash_param)

	// run
	queue.EnqueueNDRangeKernel(kernel, []int{0}, []int{1}, []int{1}, nil)
	queue.Finish()

	result_hash := make([]byte, 32)
	queue.EnqueueReadBufferByte(output_hash_param, true, 0, result_hash, nil)
	queue.Flush()

	//
	fmt.Println("check_step_by_input_hash result_hash_1:", hex.EncodeToString(result_hash))

	result_hash_2 := x16rs.HashX16RS(x16rsrepeat, input_hash)

	fmt.Println("check_step_by_input_hash result_hash_2:", hex.EncodeToString(result_hash_2))

	if bytes.Compare(result_hash, result_hash_2) != 0 {
		panic("check_step_by_input_hash  bytes.Compare(result_hash, result_hash_2) != 0")
	}

	return result_hash_2

}

func buildOrLoadProgram(cldir string, platform_id int, device_id int, rebuild bool, kernelname string) (*cl.Device, *cl.Program, *cl.Kernel, *cl.Context) {

	platforms, _ := cl.GetPlatforms()

	if len(platforms) == 0 {
		fmt.Println("not find any platforms.")
		return nil, nil, nil, nil
	}
	for i, pt := range platforms {
		fmt.Printf("- platform %d: %s\n", i, pt.Name())
	}

	platform := platforms[platform_id]

	fmt.Printf("current use platform: %s\n", platform.Name())

	devices, _ := platform.GetDevices(cl.DeviceTypeAll)

	if len(devices) == 0 {
		fmt.Println("not find any devices.")
		return nil, nil, nil, nil
	}

	for i, dv := range devices {
		fmt.Printf("- device %d: %s\n", i, dv.Name())
	}

	device := devices[device_id]
	fmt.Printf("current use device: %s\n", device.Name())

	// context
	context, _ := cl.CreateContext([]*cl.Device{device})

	var program *cl.Program

	binfilestuff := platform.Name() // + "_" + mr.devices[0].Name()
	binfilename := strings.Replace(binfilestuff, " ", "_", -1)
	binfilepath := cldir + "/" + binfilename + ".objcache"
	fmt.Println("binfilename:", binfilepath)
	binstat, staterr := os.Stat(binfilepath)
	fmt.Println("Stat:", binstat, staterr)
	if rebuild || staterr != nil {
		fmt.Print("Create opencl program with source: " + cldir + ", Please wait...")
		buildok := false
		go func() { // 打印
			for {
				time.Sleep(time.Second * 3)
				if buildok {
					break
				}
				fmt.Print(".")
			}
		}()
		program, _ = context.CreateProgramWithSource([]string{` #include "x16rs_main.cl" `})
		bderr := program.BuildProgram(nil, "-I "+cldir) // -I /media/yangjie/500GB/Hacash/src/github.com/hacash/x16rs/opencl
		if bderr != nil {
			panic(bderr)
		}
		buildok = true // build 完成
		fmt.Println("\nBuild complete get binaries...")
		//fmt.Println("program.GetBinarySizes_2()")
		sizes, _ := program.GetBinarySizes_2(1)
		//fmt.Println(sizes)
		fmt.Println("GetBinarySizes_2", sizes[0])
		fmt.Println("program.GetBinaries_2()")
		bins, _ := program.GetBinaries_2([]int{sizes[0]})
		fmt.Println("bins[0].size", len(bins[0]))
		f, e := os.OpenFile(binfilepath, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0777)
		if e != nil {
			panic(e)
		}
		fmt.Println("f.Write(wbin) "+binfilepath, sizes[0])
		var berr error
		_, berr = f.Write(bins[0])
		if berr != nil {
			panic(berr)
		}
		berr = f.Close()
		if berr != nil {
			panic(berr)
		}

	} else {
		fmt.Printf("Load binary program file from \"%s\"\n", binfilepath)
		file, _ := os.OpenFile(binfilepath, os.O_RDONLY, 0777)
		bin := make([]byte, 0)
		fmt.Println("file.Read(bin) size", binstat.Size())
		var berr error
		bin, berr = ioutil.ReadAll(file)
		if berr != nil {
			panic(berr)
		}
		if int64(len(bin)) != binstat.Size() {
			panic("int64(len(bin)) != binstat.Size()")
		}
		berr = file.Close()
		if berr != nil {
			panic(berr)
		}
		//fmt.Println(bin)
		// 仅仅支持同一个平台的同一种设备
		bins := make([][]byte, len(devices))
		sizes := make([]int, len(devices))
		for k, _ := range devices {
			bins[k] = bin
			sizes[k] = int(binstat.Size())
		}
		fmt.Println("Create program with binary...")
		program, berr = context.CreateProgramWithBinary(devices, sizes, bins)
		if berr != nil {
			panic(berr)
		}
		program.BuildProgram(devices, "")
		//fmt.Println("context.CreateProgramWithBinary")
	}
	fmt.Println("GPU miner program create complete successfully.")

	fmt.Println("build complete create kernel call...")
	kernel, _ := program.CreateKernel(kernelname)

	// 返回
	return device, program, kernel, context
}

func Test4(t *testing.T) {
	fmt.Println(os.Args)
	//cldir := flag.String("cldir", "./", "Opencl source file absolute path")
	//item_wide := flag.Int("iw", 64, "Number of concurrent processing at a time")
	//flag.Parse()

	BuildProgram("/home/victoryang00/Documents/x16rs/opencl", "", 0, false, "miner_do_hash_x16rs_v1")

}

func CreateWorkKernel(item_wide int) {

}

func BuildProgram(cldir string, plat string, dvid int, rebuild bool, kernelname string) {

	// init
	platids := 0

	platforms, _ := cl.GetPlatforms()

	if len(platforms) == 0 {
		fmt.Println("not find any platforms.")
		return
	}
	for i, pt := range platforms {
		fmt.Printf("- platform %d: %s\n", i, pt.Name())
		if strings.Compare(plat, "") != 0 && strings.Contains(pt.Name(), plat) {
			platids = i
		}
	}

	platform := platforms[platids]

	fmt.Printf("current use platform: %s\n", platform.Name())

	devices, _ := platform.GetDevices(cl.DeviceTypeAll)

	if len(devices) == 0 {
		fmt.Println("not find any devices.")
		return
	}

	for i, dv := range devices {
		fmt.Printf("- device %d: %s\n", i, dv.Name())
	}

	device := devices[dvid]
	fmt.Printf("current use device: %s\n", device.Name())

	context, _ := cl.CreateContext([]*cl.Device{device})
	queue, _ := context.CreateCommandQueue(device, 0)
	program, _ := context.CreateProgramWithSource([]string{` #include "x16rs_main.cl" `})

	fmt.Println("building opencl program from dir " + cldir + ", please wait...")
	bderr := program.BuildProgram(nil, "-I "+cldir) // -I /media/yangjie/500GB/Hacash/src/github.com/hacash/x16rs/opencl
	if bderr != nil {
		panic(bderr)
	}

	/*
		binfilestuff := platform.Name() + "_" + device.Name()
		binfilename := strings.Replace(binfilestuff, " ", "_", -1)
		binfilepath := cldir + "/" + binfilename + ".objcache"
		binstat, staterr := os.Stat(binfilepath)
		if rebuild || staterr != nil {
			bderr := program.BuildProgram(nil, "-I " + cldir) // -I /media/yangjie/500GB/Hacash/src/github.com/hacash/x16rs/opencl
			if bderr != nil {
				panic(bderr)
			}
			fmt.Println("program.GetBinarySizes_2()")
			sizes, _ := program.GetBinarySizes_2( 1 )
			fmt.Println(sizes)
			fmt.Println(sizes[0])
			fmt.Println("program.GetBinaries_2()")
			bins, _ := program.GetBinaries_2( []int{sizes[0]} )
			fmt.Println(bins[0])

			f, e := os.OpenFile(binfilepath, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0777)
			if e != nil {
				panic(e)
			}
			fmt.Println("f.Write(wbin) "+binfilepath, sizes[0])
			f.Write(wbin)

		}else{
			fmt.Println("os.Open(binfilepath) ", binfilepath)
			file, _ := os.OpenFile(binfilepath, os.O_RDWR, 0777)
			bin := make([]byte, binstat.Size())
			fmt.Println("file.Read(bin) size", binstat.Size())
			file.Read(bin)
			fmt.Println(bin)
			var berr error
			program, berr = context.CreateProgramWithBinary_2([]*cl.Device{device}, []int{int(binstat.Size())}, [][]byte{bin})
			if berr != nil {
				panic(berr)
			}
			fmt.Println("context.CreateProgramWithBinary")
		}
	*/

	fmt.Println("build complete create kernel call...")
	kernel, _ := program.CreateKernel("miner_do_hash_x16rs_v1")

	t1 := time.Now() // get current time
	fmt.Println(t1.Format("15:04:05"))
	/////////////////////////////////////////
	group_size := 4096 * 4
	var i uint32
	var nonce []byte = nil
	for i = 0; i < 1000000; i++ {
		nonce = doGroupMiner(device, context, queue, kernel, i*uint32(group_size), group_size)
		fmt.Println(i, uint32(group_size)*(i+1), nonce)
		if nonce != nil {
			break
		}
	}
	fmt.Println("i=", i, "nonce=", nonce, "sum=", uint32(group_size)*(i+1))
	////////////////////////////////////////
	fmt.Println(time.Now().Format("15:04:05"))
	elapsed := time.Since(t1)
	fmt.Println("App elapsed: ", elapsed)

}

func doGroupMiner(
	device *cl.Device,
	context *cl.Context,
	queue *cl.CommandQueue,
	kernel *cl.Kernel,
	base_start uint32,
	group_size int) []byte {

	input_target, _ := context.CreateEmptyBuffer(cl.MemReadOnly, 32)
	input_stuff, _ := context.CreateEmptyBuffer(cl.MemReadOnly, 89)
	output_nonce, _ := context.CreateEmptyBuffer(cl.MemReadWrite|cl.MemAllocHostPtr, 4)
	output_hash, _ := context.CreateEmptyBuffer(cl.MemReadWrite|cl.MemAllocHostPtr, 32)
	// copy set input
	target := bytes.Repeat([]byte{255}, 32)
	target[0] = 0
	target[1] = 0
	target[2] = 0
	target[3] = 10
	target[4] = 198
	//target[5] = 100
	queue.EnqueueWriteBufferByte(input_target, true, 0, target, nil)
	blockbytes, _ := hex.DecodeString("010000000000005c57b08c0000000000000000000000000000000000000000000000000000000000000000ad557702fc70afaf70a855e7b8a4400159643cb5a7fc8a89ba2bce6f818a9b0100000001098b344500000000000000000c1aaa4e6007cc58cfb932052ac0ec25ca356183f80101686172646572746f646f62657474657200")
	// block hash: [0 0 0 10 198 94 211 152 131 143 206 7 61 245 177 81 50 218 67 111 126 41 147 53 63 211 102 43 248 178 207 145]
	// block hash: 0000000ac65ed398838fce073df5b15132da436f7e2993353fd3662bf8b2cf91
	// target nonce: [0 112 219 121]
	//fmt.Println(blockbytes[79:83])
	queue.EnqueueWriteBufferByte(input_stuff, true, 0, blockbytes[0:89], nil)
	//queue.EnqueueWriteBufferByte(input_basestart, true, 0, []byte{0,112,219,120}, nil
	queue.EnqueueWriteBufferByte(output_nonce, true, 0, []byte{0, 0, 0, 0}, nil)

	// set argvs
	kernel.SetArgs(input_target, input_stuff, uint32(base_start), output_nonce, output_hash)
	//fmt.Println("kernel.SetArgs")

	local, _ := kernel.WorkGroupSize(device)
	//fmt.Println("Work group size: %d", local)
	//size, _ := kernel.PreferredWorkGroupSizeMultiple(nil)
	//fmt.Println("Preferred Work Group Size Multiple: %d", size)
	global := group_size
	d := group_size % local
	if d != 0 {
		global += local - d
	}

	// run
	queue.EnqueueNDRangeKernel(kernel, nil, []int{global}, []int{local}, nil)

	if err := queue.Finish(); err != nil {
		fmt.Printf("Finish NDRangeKernel failed: %+v", err)
		return nil
	}

	result_nonce := bytes.Repeat([]byte{0}, 4)
	result_hash := make([]byte, 32)
	// copy get output
	_, e1 := queue.EnqueueReadBufferByte(output_nonce, true, 0, result_nonce, nil)
	if e1 != nil {
		fmt.Printf("EnqueueReadBufferByte failed, output_nonce: %+v", e1)
		return nil
	}
	_, e2 := queue.EnqueueReadBufferByte(output_hash, true, 0, result_hash, nil)
	if e2 != nil {
		fmt.Printf("EnqueueReadBufferByte failed, output_hash: %+v", e2)
		return nil
	}
	//
	//if err := queue.Finish(); err != nil {
	//	fmt.Println("Finish ReadBufferByte failed: %+v", err)
	//	return nil
	//}

	//fmt.Println("result_nonce: ", result_nonce )
	//
	//
	//if result_nonce[0] != 0 {
	//	fmt.Println("result_nonce:", result_nonce)
	//	fmt.Println("result_hash:", result_hash)
	//}
	//fmt.Println(results)

	//fmt.Println(result_nonce)
	nonce := binary.BigEndian.Uint32(result_nonce)
	if nonce > 0 {
		// check results
		fmt.Println("==========================", nonce, result_nonce)
		fmt.Println("output_hash", result_hash, hex.EncodeToString(result_hash))
		// return
		return result_nonce
	}
	return nil

}
