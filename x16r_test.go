package x16r

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"math/rand"
	"testing"
)

func TestX16R(t *testing.T) {
	data, _ := hex.DecodeString("000000302aefb1655612597af8b82943fbe79efe672d44a226ae9a548aab00000000000002aeda49909cbedc98b1101bce8beae89a342ddf2946705a528fb1bc95e18bfce7eae15bfbc6001b43f2c36a")
	hash, _ := hex.DecodeString("01b800b258ed908d8ed4dffd87b87105e0b6e11a5e7f465b741d000000000000")
	res := Sum(data)
	if !bytes.Equal(res, hash) {
		t.Error("hash", hex.EncodeToString(res))
	}
}

func TestX16RS(t *testing.T) {
	data, _ := hex.DecodeString("000000302aefb1655612597af8b82943fbe79efe672d44a226ae9a548aab00000000000002aeda49909cbedc98b1101bce8beae89a342ddf2946705a528fb1bc95e18bfce7eae15bfbc6001b43f2c36a")
	hash, _ := hex.DecodeString("e49e48653c3f456924bbeaa3fb0df1c0c8d511577e693bdab8067222262e0ac1")
	res := HashX16RS(data)
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
	for i := 0; i < 10000*10; i++ {
		rand.Read(data)
		//fmt.Println(token)
		res := HashX16RS(data)
		if bytes.HasPrefix(res, []byte{0}) {
			fmt.Println(hex.EncodeToString(res))
		}
	}

}
