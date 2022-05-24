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
	"strconv"
	"unsafe"
)

/*
mkdir -p build && cd build
rm -f ../libx16rs_hash.a  && rm -rf * && cmake ../ && make && mv -f ./libx16rs_hash.a ../
*/

func Sum(data []byte) []byte {
	var res [32]C.char
	var cstr = C.CString(string(data))

	defer C.free(unsafe.Pointer(cstr))
	C.x16r_hash(cstr, &res[0])

	return []byte(C.GoStringN(&res[0], 32))
}

func Sha3_256(data []byte) []byte {
	var res [32]C.char
	var cstr = C.CString(string(data))
	var cint = C.int(len(data))

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

func HashX16RS(loopnum int, data []byte) []byte {
	var res [32]C.char
	var lpnm = C.int(loopnum)
	var cstr = C.CString(string(data))

	defer C.free(unsafe.Pointer(cstr))
	C.x16rs_hash(lpnm, cstr, &res[0])

	return []byte(C.GoStringN(&res[0], 32))
}

func HashRepeatForBlockHeight(blockHeight uint64) int {
	repeat := int(blockHeight/50000 + 1)
	if repeat > 16 {
		repeat = 16 // Up to 16 rounds in 8 years
	}
	return repeat
}

func CalculateBlockHash(blockHeight uint64, stuff []byte) []byte {
	repeat := HashRepeatForBlockHeight(blockHeight)
	hashbase := sha3.Sum256(stuff)
	return HashX16RS(repeat, hashbase[:])
}

// stopkind: 0.Stop after natural circulation 1.External signal forced stop 2.Excavation successfully stopped
func MinerNonceHashX16RS(blockHeight uint64, retmaxhash bool, stopmark *byte, hashstart uint32, hashend uint32, tarhashvalue []byte, blockheadmeta []byte) (byte, bool, []byte, []byte) {
	repeat := HashRepeatForBlockHeight(blockHeight)

	var retmaxhashsig int = 0
	if retmaxhash {
		retmaxhashsig = 1
	}

	var stopkind [1]C.char
	var success [1]C.char
	var nonce [4]C.char
	var reshash [32]C.char
	var hsstart = C.uint32_t(hashstart) // uint32(1)
	var hsend = C.uint32_t(hashend)     // uint32(4294967294)
	var tarhash = C.CString(string(tarhashvalue))
	var stuff = C.CString(string(blockheadmeta))
	defer C.free(unsafe.Pointer(tarhash))
	defer C.free(unsafe.Pointer(stuff))

	C.miner_x16rs_hash(C.int(repeat), C.int(retmaxhashsig), (*C.char)((unsafe.Pointer)(stopmark)), hsstart, hsend, tarhash, stuff, &stopkind[0], &success[0], &nonce[0], &reshash[0])

	return byte(stopkind[0]), success[0] == 1, []byte(C.GoStringN(&nonce[0], 4)), []byte(C.GoStringN(&reshash[0], 32))
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

func HashRepeatForDiamondNumber(diamondNumber uint32) int {
	repeat := int(diamondNumber/8192 + 1) // Adjust the hashing times every 8192 diamonds (about 140 days and half a year)
	if repeat > 16 {
		repeat = 16 // atmost 16 round due to x16rs algorithm
	}
	return repeat
}

func Diamond(diamondNumber uint32, prevblockhash []byte, nonce []byte, address []byte, extendMessage []byte) ([]byte, []byte, string) {
	repeat := HashRepeatForDiamondNumber(diamondNumber)
	stuff := new(bytes.Buffer)
	stuff.Write(prevblockhash)
	stuff.Write(nonce)
	stuff.Write(address)
	stuff.Write(extendMessage)

	/* get ssshash by sha3 algrotithm */
	ssshash := sha3.Sum256(stuff.Bytes())
	/* get diamond hash value by HashX16RS algorithm */
	reshash := HashX16RS(repeat, ssshash[:])
	/* get diamond name by DiamondHash function */
	diamond_str := DiamondHash(reshash)

	return ssshash[:], reshash, diamond_str
}

// to check if a string is a valid diamond
func IsDiamondHashResultString(diamondStr string) (string, bool) {
	if len(diamondStr) != 16 {
		return "", false
	}

	prefixlen := 10 // Number of leading zeros
	diamond_prefixs := []byte(diamondStr)[0:prefixlen]
	if bytes.Compare(diamond_prefixs, bytes.Repeat(diamond_hash_base_stuff[0:1], prefixlen)) != 0 {
		return "", false
	}

	diamond_value := []byte(diamondStr)[prefixlen:]
	isdmd := IsDiamondValueString(string(diamond_value))
	if !isdmd {
		return "", false
	}

	// to check success
	return string(diamond_value[10-prefixlen:]), true
}

// to check if a string is a diamond
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

	// check success
	return true
}

// Judge whether it is a legal diamond name or number
func IsDiamondNameOrNumber(diamondStr string) bool {
	// check number is valid or not, diamond number can't more than 16777216
	if dianum, e := strconv.Atoi(diamondStr); e == nil && dianum > 0 && dianum < 16777216 {
		return true
	}

	// to check diamondStr is valid diamond or not
	return IsDiamondValueString(diamondStr)
}

// Check whether the diamond difficulty value meets the requirements
/*
 0 128 2.00
 1 132 1.94
 2 136 1.88
 3 140 1.83
 4 144 1.78
 5 148 1.73
 6 152 1.68
 7 156 1.64
 8 160 1.60
 9 164 1.56
10 168 1.52
11 172 1.49
12 176 1.45
13 180 1.42
14 184 1.39
15 188 1.36
16 192 1.33
17 196 1.31
18 200 1.28
19 204 1.25
20 208 1.23
21 212 1.21
22 216 1.19
23 220 1.16
24 224 1.14
25 228 1.12
26 232 1.10
27 236 1.08
28 240 1.07
29 244 1.05
30 248 1.03
31 252 1.02

[128 132 136 140 144 148 152 156 160 164 168 172 176 180 184 188 192 196 200 204 208 212 216 220 224 228 232 236 240 244 248 252]
*/
func CheckDiamondDifficulty(dNumber uint32, sha3hash, dBytes []byte) bool {
	var DiaMooreDiffBits = []byte{ // difficulty requirements
		128, 132, 136, 140, 144, 148, 152, 156, // step +4
		160, 164, 168, 172, 176, 180, 184, 188,
		192, 196, 200, 204, 208, 212, 216, 220,
		224, 228, 232, 236, 240, 244, 248, 252,
	}

	// Referring to Moore's law, the excavation difficulty of every 42000 diamonds will double in about 2 years,
	// and the difficulty increment will tend to decrease to zero in 64 years
	shnum := dNumber / 42000
	if shnum > 32 {
		shnum = 32 // Up to 64 years
	}

	for i := 0; i < int(shnum); i++ {
		if sha3hash[i] >= DiaMooreDiffBits[i] {
			return false // Check failed, difficulty value does not meet requirements
		}
	}

	// Every 3277 diamonds is about 56 days. Adjust the difficulty 3277 = 16 ^ 6 / 256 / 20
	// When the difficulty is the highest, the first 20 bits of the hash are 0, not all 32 bits are 0.
	diffnum := dNumber / 3277
	for _, bt := range dBytes {
		if diffnum < 255 {
			if uint32(bt)+diffnum > 255 {
				return false // Difficulty check failed
			} else {
				return true
			}
		} else if diffnum >= 255 {
			if uint8(bt) != 0 {
				return false // Difficulty check failed
			}
			// to do next round check
			diffnum -= 255
		}
	}

	return false
}

// to mint diamond
func MinerHacashDiamond(hash_start uint32, hash_end uint32, diamondnumber int, stopmark *byte, blockhash []byte, address []byte, extendmsg []byte) ([]byte, string) {
	var nonce [8]C.char
	var diamond [16]C.char
	var hsstart = C.uint32_t(hash_start)
	var hsend = C.uint32_t(hash_end)
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
	blockhash, _ := hex.DecodeString("000009d0d0af1c87d65c581310bd7ae803b23c69754be16df02a7b156c03c87")
	address, _ := hex.DecodeString("000c1aaa4e6007cc58cfb932052ac0ec25ca356183") // 1271438866CSDpJUqrnchoJAiGGBFSQhjd

	var stopmark *byte = new(byte)
	*stopmark = 0
	nonce, diamond := MinerHacashDiamond(1, 4200009999, 1, stopmark, blockhash, address, []byte{})
	fmt.Println("miner finish nonce is", binary.BigEndian.Uint64(nonce), "bytes", nonce, "diamond is", diamond)

	// Verify whether the diamond algorithm is correct
	_, _, diamond_str := Diamond(1, blockhash, nonce, address, []byte{})
	fmt.Println("diamond_str is", diamond_str)
}
