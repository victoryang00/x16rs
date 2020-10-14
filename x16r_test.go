package x16rs

import (
	"bytes"
	bmr "crypto/rand"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"github.com/hacash/core/fields"
	"github.com/xfong/go2opencl/cl"
	"golang.org/x/crypto/sha3"
	"log"
	"math/rand"
	"os"
	"testing"
	"time"
)

func TestNewDiamond(t *testing.T) {
	prevhash, _ := hex.DecodeString("000000000ae89f8a1c93fbffe9baad198fda287f13e39c37294eb7a7b617bd70")
	extmsg, _ := hex.DecodeString("7d8367f6e46e9ffee311b9f3a38519d674b52407fd0aa287442715fe2f0c4db0")
	nonce, _ := hex.DecodeString("00000000c9babb01")
	addr, _ := fields.CheckReadableAddress("1KcXiRhMgGcvgxZGLBkLvogKLNKNXfKjEr")

	dmdhash, dmastr := Diamond(20001, prevhash, nonce, *addr, extmsg)

	fmt.Println(dmdhash, dmastr)

	fmt.Println(Diamond(20001, prevhash, nonce, *addr, []byte{}))

}

func TestX16R(t *testing.T) {
	// name+year+name+year+10001
	data, _ := hex.DecodeString("514eb391138bc40330d54c1d8ba0c2bff5b055602ba01fa7f9b3f466a042d08f")
	hash, _ := hex.DecodeString("f3bfada6cf5bb8c898fe81e37195287520b1ee08d97672b821bbe6f1ba4492ce")
	res := Sum(data)
	if !bytes.Equal(res, hash) {
		t.Error("hash", hex.EncodeToString(res))
	}
}

func TestX16RS(t *testing.T) {

	loopnum := 1

	data, _ := hex.DecodeString("514eb391138bc40330d54c1d8ba0c2bff5b055602ba01fa7f9b3f466a042d08f")
	hash, _ := hex.DecodeString("57cef097f9a7cc0c45bcac6325b5b6e58199c8197763734cac6664e8d2b8e63e")
	for i := 0; i < 1; i++ {
		res1 := HashX16RS_Optimize(loopnum, data)
		fmt.Println(hex.EncodeToString(res1))
		res2 := HashX16RS_Optimize(loopnum, data)
		fmt.Println(hex.EncodeToString(res2))
		//time.Sleep(time.Duration(100) * time.Millisecond)
	}
	res1 := HashX16RS_Optimize(loopnum, data)
	fmt.Println(hex.EncodeToString(res1))
	res2 := HashX16RS_Optimize(loopnum, data)
	fmt.Println(hex.EncodeToString(res2))
	//fmt.Println(data)
	//fmt.Println(hash)
	//fmt.Println(res)
	//fmt.Println(hex.EncodeToString(res))

	sha3results := sha3.Sum256(bytes.Repeat([]byte{1, 2, 3, 4}, 8))
	fmt.Println("sha3results", sha3results)

	if !bytes.Equal(res1, hash) {
		t.Error("hash", hex.EncodeToString(res1))
	}

}

func TestX16RS_LOOP(t *testing.T) {

	data1 := bytes.Repeat([]byte{1, 2, 3, 4, 5, 6, 7, 8}, 4)
	for i := 0; i < 10000*450; i++ { // 0000*450
		//fmt.Println(token)
		data1[4] = uint8(i % 255)
		HashX16RS_Optimize(1, data1)
		HashX16RS_Optimize(1, data1)
		//res := data
		//if bytes.Compare(res1, res2) != 0 {
		//	t.Error("hash1", hex.EncodeToString(res1), "hash2", hex.EncodeToString(res2))
		//}
	}

}

func TestX16RS_miner(t *testing.T) {
	var tarhash, _ = hex.DecodeString("00000007b37f53178a3e353d6dd319db6b62a88b5f8be80fb4e56b5f8a066fa3")
	var signstuff, _ = hex.DecodeString("00000007b37f53178a3e353d6dd319db6b62a88b5f8be80fb4e56b5f8a066fa3")
	var stopmark *byte = new(byte)
	*stopmark = 0
	go func() {
		fmt.Println("wait to stop (5s)")
		time.Sleep(time.Second * 5)
		fmt.Println("set stop mark !")
		*stopmark = 1 // 通知停止
	}()
	_, _, nonce, _ := MinerNonceHashX16RS(1, false, stopmark, 1, 4294967294, tarhash, signstuff)
	fmt.Println("miner finish nonce is", binary.BigEndian.Uint32(nonce), "bytes", nonce)

}

func TestX16RS_miner_do(t *testing.T) {

	blkbts, _ := hex.DecodeString("010000003f37005c90a5b80000000d0d0af1c87d65c581310bd7ae803b23c69754be16df02a7b156c03c87aadd0ada0615668c7bf3658efeab80ef2a6be1e884a2844d52afdb88fa82f5c6000000010070db79e48fffa400000000ff89de02003bea1b64e8d5659d314c078ad37551f801012020202020202020202020202020202000")
	blockheadmeta := blkbts[0:89]
	targetdiffhash, _ := hex.DecodeString("000009d0d0af1c87d65c581310bd7ae803b23c69754be16df02a7b156c03c87")

	//fmt.Println(blockheadmeta)
	//fmt.Println(len(targetdiffhash))

	var stopmark *byte = new(byte)
	*stopmark = 0

	_, _, nonce, _ := MinerNonceHashX16RS(1, false, stopmark, 1, 4294967294, targetdiffhash, blockheadmeta)
	fmt.Println("miner finish nonce is", binary.BigEndian.Uint32(nonce), "bytes", nonce)

}

func TestSha3_256(t *testing.T) {
	stuff := []byte("12345678901234567890123456789012")
	checkResult, _ := hex.DecodeString("dcb35cb4900cd08e524b8609b1df612e2e9d1fbfeedaa9d58a00fc0984f4a387")
	checkResult2 := sha3.Sum256(stuff)

	result := Sha3_256(stuff)

	fmt.Println(result)
	fmt.Println(checkResult)
	fmt.Println(checkResult2)

	if !bytes.Equal(checkResult, result) {
		t.Error("hash", hex.EncodeToString(result))
	}
}

func TestX16RS_num(t *testing.T) {
	data, _ := hex.DecodeString("f3bfada6cf5bb8c898fe81e37195287520b1ee08d97672b821bbe6f1ba4492ce")
	hash1 := HashX16RS_Optimize(1, data)
	fmt.Println(hash1)
	hash2 := HashX16RS_Optimize(2, data)
	fmt.Println(hash2)
	hash3 := HashX16RS_Optimize(3, data)
	fmt.Println(hash3)

}

func Test_diamond_miner_do(t *testing.T) {

	//blkbts, _ := hex.DecodeString("010000003f37005c90a5b80000000d0d0af1c87d65c581310bd7ae803b23c69754be16df02a7b156c03c87aadd0ada0615668c7bf3658efeab80ef2a6be1e884a2844d52afdb88fa82f5c6000000010070db79e48fffa400000000ff89de02003bea1b64e8d5659d314c078ad37551f801012020202020202020202020202020202000")
	//blockheadmeta := blkbts[0:89]
	blockhash, _ := hex.DecodeString("000000077790ba2fcdeaef4a4299d9b667135bac577ce204dee8388f1b97f7e6")
	//address, _ := hex.DecodeString("000c1aaa4e6007cc58cfb932052ac0ec25ca356183") // 1271438866CSDpJUqrnchoJAiGGBFSQhjd
	address, _ := fields.CheckReadableAddress("1MzNY1oA3kfgYi75zquj3SRUPYztzXHzK9")

	//fmt.Println(blockheadmeta)
	//fmt.Println(len(targetdiffhash))

	var stopmark *byte = new(byte)
	*stopmark = 0

	go func() {
		time.Sleep(time.Second)
		//*stopmark = 1
	}()

	nonce, diamond := MinerHacashDiamond(1, 4200008888, 1, stopmark, blockhash, *address, []byte{})
	fmt.Println("miner finish nonce is", binary.BigEndian.Uint64(nonce), "bytes", nonce, "hex", hex.EncodeToString(nonce), "diamond is", diamond)

	// 验证钻石算法是否正确
	_, diamond_str := Diamond(1, blockhash, nonce, *address, []byte{})
	fmt.Println("diamond_str is", diamond_str)

	if !bytes.Equal([]byte(diamond), []byte(diamond_str)) {
		t.Error("diamond: ", diamond, "but get", diamond_str)
	}

}

func Test_print_x16rs(t *testing.T) {

	data := bytes.Repeat([]byte{12, 52, 5, 230, 151, 150, 139, 223, 254, 37, 62, 187, 3, 34, 169, 36, 48, 200, 23, 127, 166, 146, 160, 123, 134, 36, 215, 137, 113, 139, 34, 241}, 1)
	fmt.Println(data)
	resultBytes := TestPrintX16RS(data)
	for i := 0; i < 16; i++ {
		fmt.Println(i, resultBytes[i])
	}

}

func Test_print_testX16RS(t *testing.T) {

	data := bytes.Repeat([]byte{12, 52, 5, 230, 151, 150, 139, 223, 254, 37, 62, 187, 3, 34, 169, 36, 48, 200, 23, 127, 166, 146, 160, 123, 134, 36, 215, 137, 113, 139, 34, 240}, 1)
	fmt.Println(data)
	resultBytes := HashX16RS_Optimize(1, data)
	fmt.Println(resultBytes)

}

// 234,214,164,90,45,197,130,255,13,248,176,44,151,46,87,41,204,138,20,15,157,191,112, 255,107,107,118,6,83,243,227,192
// 12,52,5,230,151,150,139,223,254,37,62,187,3,34,169,36,48,200,23,127,166,146,160,123,134,36,215,137,113,139,34,241
// 12 52 5 230 151 150 139 223 254 37 62 187 3 34 169 36 48 200 23 127 166 146 160 123 134 36 215 137 113 139 34 241
// 190, 201, 237, 69, 96, 107, 53, 61, 164, 23, 100, 251, 210, 169,203,189,199,200,184,172,187,60,210,209,109,96,122,78,2,172,220,201
//////////////////// OpenCL /////////////////////
// 108 220 63 239 43 104 233 103 219 79 119 139 26 152 146 61 47 77 229 77 11 14 13 202 42 188 120 72 225 240 38 167

//按字节读取，将整个文件读取到缓冲区buffer
func ReadFileBytes(filename string) []byte {
	file, err := os.Open(filename)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	fileinfo, err := file.Stat()
	if err != nil {
		log.Fatal(err)
	}
	fileSize := fileinfo.Size()
	buffer := make([]byte, fileSize)
	bytesread, err := file.Read(buffer)
	if err != nil {
		log.Fatal(err, bytesread)
	}
	//fmt.Println("bytes read:", bytesread)
	//fmt.Println("bytestream to string:", string(buffer))
	return buffer
}

func Test_OpenCL(t *testing.T) {

	// source
	kernelSource := ReadFileBytes("./opencl/x16rs_main.cl")

	// input data
	var data [2]float32
	for i := 0; i < len(data); i++ {
		data[i] = rand.Float32()
	}
	// init
	platforms, _ := cl.GetPlatforms()
	platform := platforms[0]
	devices, _ := platform.GetDevices(cl.DeviceTypeAll)
	device := devices[0]
	context, _ := cl.CreateContext([]*cl.Device{device})
	queue, _ := context.CreateCommandQueue(device, 0)
	program, _ := context.CreateProgramWithSource([]string{string(kernelSource)})
	program.BuildProgram(nil, "-I /media/yangjie/500GB/Hacash/src/github.com/hacash/x16rs/opencl") // -I /media/yangjie/500GB/Hacash/src/github.com/hacash/x16rs/opencl
	kernel, _ := program.CreateKernel("hash_sha3")
	// input and output
	input, _ := context.CreateEmptyBuffer(cl.MemReadOnly, 89)
	output, _ := context.CreateEmptyBuffer(cl.MemReadOnly, 32)
	// copy set input
	blockbytes, _ := hex.DecodeString("010000000000005c57b08c0000000000000000000000000000000000000000000000000000000000000000ad557702fc70afaf70a855e7b8a4400159643cb5a7fc8a89ba2bce6f818a9b0100000001098b344500000000000000000c1aaa4e6007cc58cfb932052ac0ec25ca356183f80101686172646572746f646f62657474657200")
	input_stuff := blockbytes[0:89]
	fmt.Println(Sha3_256(input_stuff))
	fmt.Println(sha3.Sum256(input_stuff))
	queue.EnqueueWriteBufferByte(input, true, 0, input_stuff, nil)
	// set argvs
	kernel.SetArgs(input, uint32(89), output)

	// run prepare
	local, _ := kernel.WorkGroupSize(device)
	fmt.Printf("Work group size: %d\n", local)
	size, _ := kernel.PreferredWorkGroupSizeMultiple(nil)
	fmt.Printf("Preferred Work Group Size Multiple: %d\n", size)
	global := len(data)
	d := len(data) % local
	if d != 0 {
		global += local - d
	}
	// run
	queue.EnqueueNDRangeKernel(kernel, nil, []int{1}, []int{1}, nil)
	queue.Finish()
	results := make([]byte, 32)
	// copy get output
	queue.EnqueueReadBufferByte(output, true, 0, results, nil)
	fmt.Println(results)

	// check results

	fmt.Println("==========================")

}

func Test_OpenCL_2(t *testing.T) {

	// source
	//kernelSource := ReadFileBytes("./opencl/x16rs_main.cl")

}

// 检查钻石难度值
func Test_Diamond_CheckDiamondDifficulty(t *testing.T) {

	// fmt.Println(CheckDiamondDifficulty(3277*700+3, []byte{0, 0, 69, 255, 255, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}))

	//
	//for c:=0; c<=32*256; c++ {
	//	tarhx := getdiffhashtarget(c)
	//	fmt.Println(c, tarhx)
	//}

	// 循环计算出难度目标
	for dn := uint32(1); dn < 2000000; dn += 3277 {
		for c := 0; c <= 32*256; c++ {
			tarhx := getdiffhashtarget(c)
			if CheckDiamondDifficulty(dn, tarhx) {
				fmt.Println(dn, tarhx)
				break
			}
		}
	}

	// source
	//kernelSource := ReadFileBytes("./opencl/x16rs_main.cl")

}

func getdiffhashtarget(subnum int) []byte {
	tarhash := bytes.Repeat([]byte{255}, 32)
	for i := 0; i < 32; i++ {
		if subnum < 255 {
			tarhash[i] -= uint8(subnum)
			break
		} else {
			tarhash[i] = 0
			subnum -= 255
		}
	}
	return tarhash
}

// 检查钻石哈希分布
func Test_Diamond_HashMap(t *testing.T) {

	bts := bytes.Repeat([]byte{0}, 32)

	for i := uint64(0); i < 30000; i++ {
		bmr.Read(bts)
		diamond := DiamondHash(bts)
		fmt.Print(diamond, " ")
	}
}
