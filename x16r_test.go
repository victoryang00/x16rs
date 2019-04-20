package x16r

import (
	"bytes"
	crypto_rand "crypto/rand"
	"math/rand"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"github.com/xfong/go2opencl/cl"
	"golang.org/x/crypto/sha3"
	"log"
	"os"
	"testing"
	"time"
)

func TestX16R(t *testing.T) {
	data, _ := hex.DecodeString("514eb391138bc40330d54c1d8ba0c2bff5b055602ba01fa7f9b3f466a042d08f")
	hash, _ := hex.DecodeString("f3bfada6cf5bb8c898fe81e37195287520b1ee08d97672b821bbe6f1ba4492ce")
	res := Sum(data)
	if !bytes.Equal(res, hash) {
		t.Error("hash", hex.EncodeToString(res))
	}
}

func TestX16RS(t *testing.T) {

	data, _ := hex.DecodeString("514eb391138bc40330d54c1d8ba0c2bff5b055602ba01fa7f9b3f466a042d08f")
	hash, _ := hex.DecodeString("57cef097f9a7cc0c45bcac6325b5b6e58199c8197763734cac6664e8d2b8e63e")
	for i := 0; i < 1000; i++ {
		res := HashX16RS(data)
		fmt.Println(hex.EncodeToString(res))
		//time.Sleep(time.Duration(100) * time.Millisecond)
	}

	res := HashX16RS(data)
	fmt.Println(hex.EncodeToString(res))
	//fmt.Println(data)
	//fmt.Println(hash)
	//fmt.Println(res)
	//fmt.Println(hex.EncodeToString(res))
	if !bytes.Equal(res, hash) {
		t.Error("hash", hex.EncodeToString(res))
	}

}

func TestX16RS_LOOP(t *testing.T) {

	data := make([]byte, 32)
	for i := 0; i < 10000*10000*100; i++ {
		crypto_rand.Read(data)
		//fmt.Println(token)
		res := HashX16RS(data)
		//res := data
		if bytes.HasPrefix(res, []byte{0, 0}) {
			fmt.Println(hex.EncodeToString(res))
		}
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
	nonce := MinerNonceHashX16RS(stopmark, tarhash, signstuff)
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

	nonce := MinerNonceHashX16RS(stopmark, targetdiffhash, blockheadmeta)
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

func Test_diamond_miner_do(t *testing.T) {

	//blkbts, _ := hex.DecodeString("010000003f37005c90a5b80000000d0d0af1c87d65c581310bd7ae803b23c69754be16df02a7b156c03c87aadd0ada0615668c7bf3658efeab80ef2a6be1e884a2844d52afdb88fa82f5c6000000010070db79e48fffa400000000ff89de02003bea1b64e8d5659d314c078ad37551f801012020202020202020202020202020202000")
	//blockheadmeta := blkbts[0:89]
	blockhash, _ := hex.DecodeString("000009d0d0af1c87d65c581310bd7ae803b23c69754be16df02a7b156c03c8788")
	address, _ := hex.DecodeString("000c1aaa4e6007cc58cfb932052ac0ec25ca356183") // 1271438866CSDpJUqrnchoJAiGGBFSQhjd

	//fmt.Println(blockheadmeta)
	//fmt.Println(len(targetdiffhash))

	var stopmark *byte = new(byte)
	*stopmark = 0

	nonce, diamond := MinerHacashDiamond(stopmark, blockhash, address)
	fmt.Println("miner finish nonce is", binary.BigEndian.Uint64(nonce), "bytes", nonce, "diamond is", diamond)

	// 验证钻石算法是否正确
	diamond_str := Diamond(blockhash, nonce, address)
	fmt.Println("diamond_str is", diamond_str)

	if !bytes.Equal([]byte(diamond), []byte(diamond_str)) {
		t.Error("diamond: ", diamond, "but get", diamond_str)
	}

}

func Test_print_x16rs(t *testing.T)  {

	data := bytes.Repeat([]byte{9,2,3,4,5,6,7,8}, 4)
	fmt.Println(data)
	resultBytes := TestPrintX16RS(data)
	for i:=0; i<16; i++ {
		fmt.Println(i, resultBytes[i])
	}

}


//////////////////// OpenCL /////////////////////


//按字节读取，将整个文件读取到缓冲区buffer
func ReadFileBytes( filename string ) []byte {
	file,err := os.Open(filename)
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
	var data [100]float32
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
	kernel, _ := program.CreateKernel("test_hash_x16rs")
	// input and output
	input, _ := context.CreateEmptyBuffer(cl.MemReadOnly, 32)
	output, _ := context.CreateEmptyBuffer(cl.MemReadOnly, 32)
	// copy set input
	queue.EnqueueWriteBufferByte(input, true, 0, bytes.Repeat([]byte{9,2,3,4,5,6,7,8}, 4), nil)
	// set argvs
	kernel.SetArgs(input, output)

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
	queue.EnqueueNDRangeKernel(kernel, nil, []int{global}, []int{local}, nil)
	queue.Finish()
	results := make([]byte, 32)
	// copy get output
	queue.EnqueueReadBufferByte(output, true, 0, results, nil)
	fmt.Println(results)

	// check results

	fmt.Println("==========================")





}

