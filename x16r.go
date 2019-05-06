package x16r

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

func HashX16RS(data []byte) []byte {
	var res [32]C.char
	var cstr = C.CString(string(data))
	defer C.free(unsafe.Pointer(cstr))
	C.x16rs_hash(cstr, &res[0])
	return []byte(C.GoStringN(&res[0], 32))
}
func HashX16RS_Optimize(data []byte) []byte {
	var res [32]C.char
	var cstr = C.CString(string(data))
	defer C.free(unsafe.Pointer(cstr))
	C.x16rs_hash_optimize(cstr, &res[0])
	return []byte(C.GoStringN(&res[0], 32))
}

func MinerNonceHashX16RS(stopmark *byte, tarhashvalue []byte, blockheadmeta []byte) []byte {
	var nonce [4]C.char
	var tarhash = C.CString(string(tarhashvalue))
	var stuff = C.CString(string(blockheadmeta))
	defer C.free(unsafe.Pointer(tarhash))
	defer C.free(unsafe.Pointer(stuff))
	//fmt.Println("C.miner_x16rs_hash_v1")
	C.miner_x16rs_hash_v1((*C.char)((unsafe.Pointer)(stopmark)), tarhash, stuff, &nonce[0])
	//fmt.Println("C.miner_x16rs_hash_v1 finish")
	return []byte(C.GoStringN(&nonce[0], 4))
}

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

func Diamond(blockhash []byte, nonce []byte, address []byte) string {
	stuff := new(bytes.Buffer)
	stuff.Write(blockhash)
	stuff.Write(nonce)
	stuff.Write(address)
	ssshash := sha3.Sum256(stuff.Bytes())
	//fmt.Println(ssshash)
	reshash := HashX16RS(ssshash[:])
	//fmt.Println(reshash)
	diamond_str := DiamondHash(reshash)
	return diamond_str
}

// 判断是否为钻石
func IsDiamondHashResultString(diamondStr string) (string, bool) {
	if len(diamondStr) != 16 {
		return "", false
	}
	diamond_prefixs := []byte(diamondStr)[0:10]
	if bytes.Compare(diamond_prefixs, bytes.Repeat(diamond_hash_base_stuff[0:1], 10)) != 0 {
		return "", false
	}
	diamond_value := []byte(diamondStr)[10:]
	for _, a := range diamond_value {
		if a == diamond_hash_base_stuff[0] {
			return "", false
		}
		if bytes.IndexByte(diamond_hash_base_stuff, a) == -1 {
			return "", false
		}
	}
	// 检查成功
	return string(diamond_value), true
}

func MinerHacashDiamond(stopmark *byte, blockhash []byte, address []byte) ([]byte, string) {
	var nonce [8]C.char
	var diamond [16]C.char
	var tarhash = C.CString(string(blockhash))
	var taraddr = C.CString(string(address))
	defer C.free(unsafe.Pointer(tarhash))
	defer C.free(unsafe.Pointer(taraddr))
	C.miner_diamond_hash((*C.char)((unsafe.Pointer)(stopmark)), tarhash, taraddr, &nonce[0], &diamond[0])
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
	nonce, diamond := MinerHacashDiamond(stopmark, blockhash, address)
	fmt.Println("miner finish nonce is", binary.BigEndian.Uint64(nonce), "bytes", nonce, "diamond is", diamond)
	// 验证钻石算法是否正确
	diamond_str := Diamond(blockhash, nonce, address)
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
