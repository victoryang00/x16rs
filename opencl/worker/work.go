package worker

import (
	"bytes"
	"encoding/binary"
	"fmt"
)

// 启动分组
func (mr *GpuMiner) doGroupWork(ctx *GpuMinerDeviceWorkerContext, global int, local int, x16rsrepeat uint32, base_start uint32) (bool, []byte, []byte) {

	// time.Sleep(time.Millisecond * 300)

	var e error

	// 重置
	_, e = ctx.queue.EnqueueWriteBufferByte(ctx.output_nonce, true, 0, []byte{0, 0, 0, 0}, nil)
	if e != nil {
		panic(e)
	}
	// set argvs
	e = ctx.kernel.SetArgs(ctx.input_target, ctx.input_stuff, x16rsrepeat, uint32(base_start), uint32(mr.itemLoop), ctx.output_nonce, ctx.output_hash)
	if e != nil {
		panic(e)
	}
	// run
	//fmt.Println("EnqueueNDRangeKernel")
	_, e = ctx.queue.EnqueueNDRangeKernel(ctx.kernel, []int{0}, []int{global}, []int{local}, nil)
	if e != nil {
		fmt.Println("EnqueueNDRangeKernel ERROR:")
		panic(e)
	}
	//fmt.Println("EnqueueNDRangeKernel END!!!")
	//fmt.Println("ctx.queue.Finish() start")
	e = ctx.queue.Finish()
	if e != nil {
		panic(e)
	}
	//fmt.Println("ctx.queue.Finish() end")

	result_nonce := bytes.Repeat([]byte{0}, 4)
	result_hash := make([]byte, 32)
	// copy get output
	//fmt.Println("EnqueueReadBufferByte output_nonce start")
	_, e = ctx.queue.EnqueueReadBufferByte(ctx.output_nonce, true, 0, result_nonce, nil)
	if e != nil {
		panic(e)
	}
	//fmt.Println("EnqueueReadBufferByte output_nonce end")
	//fmt.Println("EnqueueReadBufferByte output_hash start")
	_, e = ctx.queue.EnqueueReadBufferByte(ctx.output_hash, true, 0, result_hash, nil)
	if e != nil {
		panic(e)
	}
	//fmt.Println("EnqueueReadBufferByte output_hash end")
	//fmt.Println("ctx.queue.Finish() start")
	e = ctx.queue.Finish()
	if e != nil {
		panic(e)
	}
	//fmt.Println("ctx.queue.Finish() end")

	// check results
	//fmt.Println("==========================", result_nonce, hex.EncodeToString(result_nonce))
	//fmt.Println("output_hash", result_hash, hex.EncodeToString(result_hash))
	//fmt.Println(result_nonce)
	nonce := binary.BigEndian.Uint32(result_nonce)
	if nonce > 0 {
		// check results
		// fmt.Println("==========================", nonce, result_nonce)
		// fmt.Println("output_hash", result_hash, hex.EncodeToString(result_hash))
		// return
		return true, result_nonce, result_hash
	}
	return false, nil, nil

}
