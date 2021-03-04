package execute

import (
	"fmt"
	"github.com/xfong/go2opencl/cl"
	"strings"
)

func (mr *GpuMiner) Init() error {

	var e error = nil
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
	mr.devices = devices
	if mr.context, e = cl.CreateContext(mr.devices); e != nil {
		return e
	}

	if strings.Compare(mr.openclPath, "") == 0 {
		mr.openclPath = GetCurrentDirectory() + "/opencl"
	}

	// 初始化成功
	return nil
}
