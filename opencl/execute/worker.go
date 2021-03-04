package execute

import "github.com/xfong/go2opencl/cl"

// 单个设备执行环境
type GpuMinerDeviceWorker struct {
	device       *cl.Device
	kernel       *cl.Kernel
	queue        *cl.CommandQueue
	input_stuff  *cl.MemObject
	input_target *cl.MemObject
	output_nonce *cl.MemObject
	output_hash  *cl.MemObject
}

func (w *GpuMinerDeviceWorker) Retain() {
	w.kernel.Retain()
	w.queue.Retain()
	w.input_stuff.Retain()
	w.input_target.Retain()
	w.output_nonce.Retain()
	w.output_hash.Retain()
}

func (w *GpuMinerDeviceWorker) ReInit(stuff_buf []byte, target_buf []byte) {
	w.queue.EnqueueWriteBufferByte(w.input_stuff, true, 0, stuff_buf, nil)
	w.queue.EnqueueWriteBufferByte(w.input_target, true, 0, target_buf, nil)

}

func (mr *GpuMiner) ReInitWorkContext(stuff_buf []byte, target_buf []byte) {
	for _, ctx := range mr.deviceworkers {
		ctx.ReInit(stuff_buf, target_buf)
	}
}

// chua
func (mr *GpuMiner) createWorkContext(devidx int) *GpuMinerDeviceWorker {

	// 运行创建执行单元
	//input_target_buf := make([]byte, 32)
	//copy(input_target_buf, work.target[:])
	//input_stuff_buf := make([]byte, 89)
	//copy(input_stuff_buf, work.blkstuff[:])
	// |cl.MemCopyHostPtr
	input_target, _ := mr.context.CreateEmptyBuffer(cl.MemReadOnly, 32)
	input_stuff, _ := mr.context.CreateEmptyBuffer(cl.MemReadOnly, 89)
	//defer input_target.Release()
	//defer input_stuff.Release()

	// 参数
	// |cl.MemAllocHostPtr
	output_nonce, _ := mr.context.CreateEmptyBuffer(cl.MemReadWrite|cl.MemAllocHostPtr, 4)
	output_hash, _ := mr.context.CreateEmptyBuffer(cl.MemReadWrite|cl.MemAllocHostPtr, 32)
	//defer output_nonce.Release()
	//defer output_hash.Release()

	kernel, ke1 := mr.program.CreateKernel("miner_do_hash_x16rs_v2")
	if ke1 != nil {
		panic(ke1)
	}
	//defer kernel.Release()

	device := mr.devices[devidx]
	queue, qe1 := mr.context.CreateCommandQueue(device, 0)
	if qe1 != nil {
		panic(qe1)
	}
	//defer queue.Release()

	return &GpuMinerDeviceWorker{
		device,
		kernel,
		queue,
		input_stuff,
		input_target,
		output_nonce,
		output_hash,
	}

}
