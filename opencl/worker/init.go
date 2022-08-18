package worker

import (
	"fmt"
	"github.com/xfong/go2opencl/cl"
	"os"
	"strings"
)

func (mr *GpuMiner) Init() error {

	var e error = nil

	// opencl file prepare
	if strings.Compare(mr.openclPath, "") == 0 {
		tardir := GetCurrentDirectory() + "/opencl/"
		if _, err := os.Stat(tardir); err != nil {
			fmt.Println("Create opencl dir and render files...")
			files := getRenderCreateAllOpenclFiles() // 输出所有文件
			err := writeClFiles(tardir, files)
			if err != nil {
				fmt.Println(e)
				os.Exit(0) // 致命错误
			}
			fmt.Println("all file ok.")
		} else {
			fmt.Println("Opencl dir already here.")
		}
		mr.openclPath = tardir
	}

	// start
	platforms, e := cl.GetPlatforms()
	if e != nil {
		return e
	}

	if len(platforms) == 0 {
		return fmt.Errorf("not find any platforms.")
	}

	chooseplatids := 0
	for i, pt := range platforms {
		fmt.Printf("  - platform %d: %s\n", i, pt.Name())
		if strings.Compare(mr.platName, "") != 0 && strings.Contains(pt.Name(), mr.platName) {
			chooseplatids = i
		}
	}

	mr.platform = platforms[chooseplatids]
	fmt.Printf("current use platform: %s\n", mr.platform.Name())

	devices, _ := mr.platform.GetDevices(cl.DeviceTypeAll)

	if len(devices) == 0 {
		return fmt.Errorf("not find any devices.")
	}

	for i, dv := range devices {
		fmt.Printf("  - device %d: %s, (max_work_group_size: %d)\n", i, dv.Name(), dv.MaxWorkGroupSize())
	}

	// 是否单设备编译
	if mr.useOneDeviceBuild {
		fmt.Println("Only use single device to build and run.")
		mr.devices = []*cl.Device{devices[0]} // 使用单台设备
	} else {
		mr.devices = devices
	}

	if mr.context, e = cl.CreateContext(mr.devices); e != nil {
		return e
	}

	// 编译源码
	mr.program = mr.buildOrLoadProgram()

	// 初始化执行环境
	devlen := len(mr.devices)
	mr.deviceworkers = make([]*GpuMinerDeviceWorkerContext, devlen)
	for i := 0; i < devlen; i++ {
		mr.deviceworkers[i] = mr.createWorkContext(i)
	}

	// 初始化成功
	return nil
}

// 写入 opencl 文件
func writeClFiles(tardir string, files map[string]string) error {

	e := os.MkdirAll(tardir, os.ModePerm)
	if e != nil {
		return e
	}
	for name, content := range files {
		fmt.Print(name + " ")
		f, e := os.OpenFile(tardir+name, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0777)
		if e != nil {
			return e
		}
		//fmt.Println(e)
		_, e = f.Write([]byte(content))
		if e != nil {
			return e
		}
		e = f.Close()
		if e != nil {
			return e
		}
	}
	// 成功
	return nil
}
