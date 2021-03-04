package execute

import "github.com/xfong/go2opencl/cl"

// GPU挖矿统筹
type GpuMiner struct {
	platform *cl.Platform
	context  *cl.Context
	program  *cl.Program
	devices  []*cl.Device // 所有设备

	deviceworkers []*GpuMinerDeviceWorker

	// config
	openclPath string
	rebuild    bool   // 强制重新编译
	platName   string // 选择的平台

}

// 初始化
func NewGpuMiner(
	openclPath string,
	platName string,
	rebuild bool,
) *GpuMiner {

	miner := &GpuMiner{
		openclPath: openclPath,
		platName:   platName,
		rebuild:    rebuild,
	}

	// 初始化
	miner.Init()

	// 编译源码
	miner.program = miner.buildOrLoadProgram()

	// 初始化执行环境
	devlen := len(miner.devices)
	miner.deviceworkers = make([]*GpuMinerDeviceWorker, devlen)
	for i := 0; i < devlen; i++ {
		miner.deviceworkers[i] = miner.createWorkContext(i)
	}

	// 创建成功
	return miner
}
