package x16rs

/*
#cgo LDFLAGS: -L. -lx16rs_hash
#include <stdlib.h>
#include "x16rs.h"
*/
import "C"

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"golang.org/x/crypto/sha3"
	"unsafe"
)

/*

mkdir -p build && cd build
rm -f ../libx16rs_hash.a  && rm -rf * && cmake ../ && make && mv -f ./libx16rs_hash.a ../

*/

//
func Sum(data []byte) []byte {
	var res [32]C.char
	var cstr = C.CString(string(data))
	defer C.free(unsafe.Pointer(cstr))
	C.x16r_hash(cstr, &res[0])
	return []byte(C.GoStringN(&res[0], 32))
}

//
func Sha3_256(data []byte) []byte {
	var res [32]C.char
	var cstr = C.CString(string(data))
	var cint = C.int(len(data))
	// fmt.Println(len(data))
	defer C.free(unsafe.Pointer(cstr))
	C.sha3_256(cstr, cint, &res[0])
	return []byte(C.GoStringN(&res[0], 32))
}

func HashX16RS_Old(loopnum int, data []byte) []byte {
	var res [32]C.char
	var lpnm = C.int(loopnum)
	var cstr = C.CString(string(data))
	defer C.free(unsafe.Pointer(cstr))
	C.x16rs_hash(lpnm, cstr, &res[0])
	return []byte(C.GoStringN(&res[0], 32))
}

func HashX16RS_Optimize(loopnum int, data []byte) []byte {
	var res [32]C.char
	var lpnm = C.int(loopnum)
	var cstr = C.CString(string(data))
	defer C.free(unsafe.Pointer(cstr))
	C.x16rs_hash__development(lpnm, cstr, &res[0])
	return []byte(C.GoStringN(&res[0], 32))
}

func CalculateBlockHash(blockHeight uint64, stuff []byte) []byte {
	loopnum := int(blockHeight/50000 + 1)
	if loopnum > 16 {
		loopnum = 16 // 8年时间上升到16次
	}
	hashbase := sha3.Sum256(stuff)
	return HashX16RS_Optimize(loopnum, hashbase[:])
}

// stopkind:  0.自然循环完毕后停止   1.外部信号强制停止   2.挖出成功停止
func MinerNonceHashX16RS(blockHeight uint64, retmaxhash bool, stopmark *byte, hashstart uint32, hashend uint32, tarhashvalue []byte, blockheadmeta []byte) (byte, bool, []byte, []byte) {
	loopnum := int(blockHeight/50000 + 1)
	if loopnum > 16 {
		loopnum = 16 // 8年时间上升到16次
	}
	var retmaxhashsig int = 0
	if retmaxhash {
		retmaxhashsig = 1
	}
	var stopkind [1]C.char
	var success [1]C.char
	var nonce [4]C.char
	var reshash [32]C.char
	var hsstart = C.uint(hashstart) // uint32(1)
	var hsend = C.uint(hashend)     // uint32(4294967294)
	var tarhash = C.CString(string(tarhashvalue))
	var stuff = C.CString(string(blockheadmeta))
	defer C.free(unsafe.Pointer(tarhash))
	defer C.free(unsafe.Pointer(stuff))
	//fmt.Println("C.miner_x16rs_hash_v1") //
	C.miner_x16rs_hash(C.int(loopnum), C.int(retmaxhashsig), (*C.char)((unsafe.Pointer)(stopmark)), hsstart, hsend, tarhash, stuff, &stopkind[0], &success[0], &nonce[0], &reshash[0])
	//fmt.Println("C.miner_x16rs_hash_v1 finish")
	return byte(stopkind[0]), success[0] == 1, []byte(C.GoStringN(&nonce[0], 4)), []byte(C.GoStringN(&reshash[0], 32))
}

/////////////////////////////////////////////////////////

var diamond_hash_base_stuff = []byte("0WTYUIAHXVMEKBSZN")

func DiamondHash(reshash []byte) string {
	diamond_str := make([]byte, 16)
	p := 13
	for i := 0; i < 16; i++ {
		num := p * int(reshash[i*2]) * int(reshash[i*2+1])
		p = num % 17
		diamond_str[i] = diamond_hash_base_stuff[p]
		if p == 0 {
			p = 13
		}
	}
	return string(diamond_str)
}

func Diamond(diamondNumber uint32, blockhash []byte, nonce []byte, address []byte, extendMessage []byte) ([]byte, string) {
	loopnum := diamondNumber/8192 + 1 // 每 8192 颗钻石（约140天小半年）调整一下哈希次数
	if loopnum > 16 {
		loopnum = 16 // 最高16次 x16rs 哈希
	}
	stuff := new(bytes.Buffer)
	stuff.Write(blockhash)
	stuff.Write(nonce)
	stuff.Write(address)
	stuff.Write(extendMessage)
	//fmt.Println(stuff.Bytes())
	ssshash := sha3.Sum256(stuff.Bytes())
	//fmt.Println(ssshash)
	reshash := HashX16RS_Optimize(int(loopnum), ssshash[:])
	//fmt.Println(reshash)
	diamond_str := DiamondHash(reshash)
	//fmt.Println(diamond_str)
	return reshash, diamond_str
}

// 判断是否为钻石
func IsDiamondHashResultString(diamondStr string) (string, bool) {
	if len(diamondStr) != 16 {
		return "", false
	}
	prefixlen := 10 // 前导0的数量
	diamond_prefixs := []byte(diamondStr)[0:prefixlen]
	if bytes.Compare(diamond_prefixs, bytes.Repeat(diamond_hash_base_stuff[0:1], prefixlen)) != 0 {
		return "", false
	}
	diamond_value := []byte(diamondStr)[prefixlen:]
	isdmd := IsDiamondValueString(string(diamond_value))
	if !isdmd {
		return "", false
	}
	// 检查成功
	return string(diamond_value[10-prefixlen:]), true
}

// 判断是否为钻石
func IsDiamondValueString(diamondStr string) bool {
	if len(diamondStr) != 6 {
		return false
	}
	for _, a := range diamondStr {
		// drop 0
		if bytes.IndexByte(diamond_hash_base_stuff[1:], byte(a)) == -1 {
			return false
		}
	}
	// 检查成功
	return true
}

// 检查钻石难度值，是否满足要求
func CheckDiamondDifficulty(dNumber uint32, dBytes []byte) bool {
	// 每 3277 颗钻石调整一下难度 3277 = 16^6 / 256 / 20
	// 难度最高时hash前20位为0，而不是32位都为0。
	diffnum := dNumber / 3277
	for _, bt := range dBytes {
		if diffnum < 255 {
			if uint32(bt)+diffnum > 255 {
				return false // 难度检查失败
			} else {
				return true
			}
		} else if diffnum >= 255 {
			if uint8(bt) != 0 {
				return false // 难度检查失败
			}
			// 下一轮检查
			diffnum -= 255
		}
	}
	return false
}

// 钻石挖矿
func MinerHacashDiamond(hash_start uint32, hash_end uint32, diamondnumber int, stopmark *byte, blockhash []byte, address []byte, extendmsg []byte) ([]byte, string) {
	var nonce [8]C.char
	var diamond [16]C.char
	var hsstart = C.uint(hash_start)
	var hsend = C.uint(hash_end)
	var dmnb = C.int(diamondnumber)
	var tarhash = C.CString(string(blockhash))
	var taraddr = C.CString(string(address))
	var tarextmsg = C.CString(string(extendmsg))
	defer C.free(unsafe.Pointer(tarhash))
	defer C.free(unsafe.Pointer(taraddr))
	defer C.free(unsafe.Pointer(tarextmsg))
	C.miner_diamond_hash(hsstart, hsend, dmnb, (*C.char)((unsafe.Pointer)(stopmark)), tarhash, taraddr, tarextmsg, &nonce[0], &diamond[0])
	return []byte(C.GoStringN(&nonce[0], 8)), C.GoStringN(&diamond[0], 16)
}

func TestPrintX16RS(stuff32 []byte) [][]byte {
	var res [32 * 16]C.char
	var cstr = C.CString(string(stuff32))
	defer C.free(unsafe.Pointer(cstr))
	C.test_print_x16rs(cstr, &res[0])
	var bytes = []byte(C.GoStringN(&res[0], 32*16))
	var resbytes [][]byte
	for i := 0; i < 16; i++ {
		resbytes = append(resbytes, bytes[32*i:32*i+32])
	}
	return resbytes
}

////////////////////////  GPU OpenCL  ////////////////////////////

func OpenCLMinerNonceHashX16RS(stopmark *byte, tarhashvalue []byte, blockheadmeta []byte) []byte {
	return nil
}

func main() {
	//fmt.Println(Sum([]byte("test")))
	//fmt.Println(HashX16RS([]byte("test")))

	blockhash, _ := hex.DecodeString("000009d0d0af1c87d65c581310bd7ae803b23c69754be16df02a7b156c03c87")
	address, _ := hex.DecodeString("000c1aaa4e6007cc58cfb932052ac0ec25ca356183") // 1271438866CSDpJUqrnchoJAiGGBFSQhjd

	var stopmark *byte = new(byte)
	*stopmark = 0
	nonce, diamond := MinerHacashDiamond(1, 4200009999, 1, stopmark, blockhash, address, []byte{})
	fmt.Println("miner finish nonce is", binary.BigEndian.Uint64(nonce), "bytes", nonce, "diamond is", diamond)
	// 验证钻石算法是否正确
	_, diamond_str := Diamond(1, blockhash, nonce, address, []byte{})
	fmt.Println("diamond_str is", diamond_str)

}

/*

//#include "x16r.h"

//#include "sha3/sph_blake.h"
//#include "sha3/sph_bmw.h"
//#include "sha3/sph_groestl.h"
//#include "sha3/sph_jh.h"
//#include "sha3/sph_keccak.h"
//#include "sha3/sph_skein.h"
//#include "sha3/sph_luffa.h"
//#include "sha3/sph_cubehash.h"
//#include "sha3/sph_shavite.h"
//#include "sha3/sph_simd.h"
//#include "sha3/sph_echo.h"
//#include "sha3/sph_hamsi.h"
//#include "sha3/sph_fugue.h"
//#include "sha3/sph_shabal.h"
//#include "sha3/sph_whirlpool.h"
//#include "sha3/sph_sha2.h"

*/
