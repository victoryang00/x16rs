package execute

import (
	"fmt"
	"github.com/xfong/go2opencl/cl"
	"os"
	"strings"
	"time"
)

func (mr *GpuMiner) buildOrLoadProgram() *cl.Program {

	var program *cl.Program

	binfilestuff := mr.platform.Name() // + "_" + mr.devices[0].Name()
	binfilename := strings.Replace(binfilestuff, " ", "_", -1)
	binfilepath := mr.openclPath + "/" + binfilename + ".objcache"
	binstat, staterr := os.Stat(binfilepath)
	if mr.rebuild || staterr != nil {
		fmt.Print("Create opencl program with source: " + mr.openclPath + ", Please wait...")
		buildok := false
		go func() { // 打印
			for {
				time.Sleep(time.Second * 3)
				if buildok {
					break
				}
				fmt.Print(".")
			}
		}()
		program, _ = mr.context.CreateProgramWithSource([]string{` #include "x16rs_main.cl" `})
		bderr := program.BuildProgram(nil, "-I "+mr.openclPath) // -I /media/yangjie/500GB/Hacash/src/github.com/hacash/x16rs/opencl
		if bderr != nil {
			panic(bderr)
		}
		buildok = true // build 完成
		fmt.Println("\nBuild complete get binaries...")
		//fmt.Println("program.GetBinarySizes_2()")
		sizes, _ := program.GetBinarySizes_2(1)
		//fmt.Println(sizes)
		//fmt.Println(sizes[0])
		//fmt.Println("program.GetBinaries_2()")
		bins, _ := program.GetBinaries_2([]int{sizes[0]})
		//fmt.Println(bins[0])
		f, e := os.OpenFile(binfilepath, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0777)
		if e != nil {
			panic(e)
		}
		//fmt.Println("f.Write(wbin) "+binfilepath, sizes[0])
		f.Write(bins[0])

	} else {
		fmt.Printf("Load binary program file from \"%s\"\n", binfilepath)
		file, _ := os.OpenFile(binfilepath, os.O_RDWR, 0777)
		bin := make([]byte, binstat.Size())
		//fmt.Println("file.Read(bin) size", binstat.Size())
		file.Read(bin)
		//fmt.Println(bin)
		// 仅仅支持同一个平台的同一种设备
		bins := make([][]byte, len(mr.devices))
		sizes := make([]int, len(mr.devices))
		for k, _ := range mr.devices {
			bins[k] = bin
			sizes[k] = int(binstat.Size())
		}
		fmt.Println("Create program with binary...")
		var berr error
		program, berr = mr.context.CreateProgramWithBinary_2(mr.devices, sizes, bins)
		if berr != nil {
			panic(berr)
		}
		program.BuildProgram(mr.devices, "")
		//fmt.Println("context.CreateProgramWithBinary")
	}
	fmt.Println("GPU miner program create complete successfully.")

	// 返回
	return program
}
