package x16rs

/**

go build -o poolworker github.com/hacash/x16rs/pool/main/ && ./poolworker

**/

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"math/big"
	"net"
	"strings"
)

////////////////////////////////////////////////

type MiningSuccess struct {
	BlockHeight uint64
	MiningIndex uint64 // 挖矿标号
	Nonce       []byte // len:4
}

func (mp *MiningSuccess) Parse(stuff []byte, seek int) {
	if len(stuff) - seek < 8+8+4 {
		panic("stuff length not enough.")
	}

	mp.BlockHeight = binary.BigEndian.Uint64(stuff[seek : seek + 8])
	seek += 8
	mp.MiningIndex = binary.BigEndian.Uint64(stuff[seek : seek + 8])
	seek += 8
	mp.Nonce = stuff[seek : seek + 4]
}

func (mp *MiningSuccess) Serialize() []byte {
	if len(mp.Nonce) != 4 {
		panic("stuff length error.")
	}

	s1 := make([]byte, 8)
	binary.BigEndian.PutUint64(s1, mp.BlockHeight)
	s2 := make([]byte, 8)
	binary.BigEndian.PutUint64(s2, mp.MiningIndex)
	buf := bytes.NewBuffer(s1)
	buf.Write(s2)
	buf.Write(mp.Nonce)

	return buf.Bytes()
}

////////////////////////////////////////////////
type MiningPoolStuff struct {
	BlockHeight   uint64
	MiningIndex   uint64 // 挖矿标号
	Loopnum       uint8
	TargetHash    []byte // len:32
	BlockHeadMeta []byte // len:89
}

func (mp *MiningPoolStuff) Parse(stuff []byte, seek int) {
	if len(stuff) - seek < 8+8+1+32+89 {
		panic("stuff length not enough.")
	}

	mp.BlockHeight = binary.BigEndian.Uint64(stuff[seek : seek + 8])
	seek += 8
	mp.MiningIndex = binary.BigEndian.Uint64(stuff[seek : seek + 8])
	seek += 8
	mp.Loopnum = stuff[seek]
	seek += 1
	mp.TargetHash = stuff[seek : seek+32]
	seek += 32
	mp.BlockHeadMeta = stuff[seek : seek+89]
}

func (mp *MiningPoolStuff) Serialize() []byte {
	if len(mp.TargetHash) != 32 || len(mp.BlockHeadMeta) != 89 {
		panic("stuff length error.")
	}

	s1 := make([]byte, 8)
	binary.BigEndian.PutUint64(s1, mp.BlockHeight)
	s2 := make([]byte, 8)
	binary.BigEndian.PutUint64(s2, mp.MiningIndex)
	buf := bytes.NewBuffer(s1)
	buf.Write(s2)
	buf.Write([]byte{mp.Loopnum})
	buf.Write(mp.TargetHash)
	buf.Write(mp.BlockHeadMeta)

	return buf.Bytes()
}

// 发送tcp数据
func MiningPoolWriteTcpMsgBytes(conn net.Conn, typeid uint8, stuff []byte) error {
	// 信息发送给客户端
	buf := bytes.NewBuffer([]byte{typeid})
	buf.Write(stuff)
	msg := buf.Bytes()
	conn.Write([]byte(hex.EncodeToString(msg) + "\n"))

	return nil
}

// 读取tcp数据
func MiningPoolReadTcpMsgBytes(reader *bufio.Reader) ([]byte, error) {
	msgstr, err := reader.ReadString('\n')
	if err != nil {
		return nil, err
	}

	msgstr = strings.TrimRight(msgstr, "\n")
	msgbytes, err2 := hex.DecodeString(msgstr)
	if err2 != nil {
		return nil, fmt.Errorf("error tcp data from pool.")
	}

	return msgbytes, nil
}

// 计算算力分值
func CalculateHashPowerValue(hash []byte) *big.Int {
	value := []byte{}
	for i := 0; i < len(hash); i++ {
		if hash[i] == 0 {
			value = append(value, 0)
		} else {
			value = append(value, 0)
			value[0] = 255 - hash[i]
			return big.NewInt(0).SetBytes(value)
		}
	}
	return big.NewInt(0)
}
