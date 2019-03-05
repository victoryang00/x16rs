package x16r

import (
	"bytes"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"testing"
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
	for i:=0; i<1000; i++ {
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
		rand.Read(data)
		//fmt.Println(token)
		res := HashX16RS(data)
		//res := data
		if bytes.HasPrefix(res, []byte{0, 0}) {
			fmt.Println(hex.EncodeToString(res))
		}
	}

}
