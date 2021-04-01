package worker

import "github.com/xfong/go2opencl/cl"

// GPU挖矿统筹
type GpuMiner struct {
	platform *cl.Platform
	context  *cl.Context
	program  *cl.Program
	devices  []*cl.Device // 所有设备

	deviceworkers []*GpuMinerDeviceWorkerContext

	// config
	openclPath        string
	rebuild           bool   // 强制重新编译
	platName          string // 选择的平台
	groupNum          int    // 同时执行组数量
	groupSize         int    // 组大小
	itemLoop          int    // 单次执行循环次数
	emptyFuncTest     bool   // 空函数编译测试
	useOneDeviceBuild bool   // 使用单个设备编译

}

// 初始化
func NewGpuMiner(
	openclPath string,
	platName string,
	groupSize int, // 组宽度
	groupNum int, // 同时执行组数量: 1 ~ 64
	itemLoop int, // 建议 20 ～ 100
	useOneDeviceBuild bool, // 使用一个设备去编译
	rebuild bool,
	emptyFuncTest bool,
) *GpuMiner {

	miner := &GpuMiner{
		openclPath:        openclPath,
		platName:          platName,
		rebuild:           rebuild,
		emptyFuncTest:     emptyFuncTest,
		useOneDeviceBuild: useOneDeviceBuild,
		groupSize:         groupSize,
		groupNum:          groupNum,
		itemLoop:          itemLoop,
	}

	// 创建成功
	return miner
}
