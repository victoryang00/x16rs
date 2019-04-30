package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
)

func main() {

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, os.Kill)

	///////////////////////////////////////////////////

	fmt.Println(os.Args)
	cldir := flag.String("oc", "", "Opencl source file absolute path")
	group_size := flag.Int("gs", 512, "Number of concurrent processing at a time, group size")
	plat_name := flag.String("pn", "", "Platform name your choise")
	dv_id := flag.Int("di", 0, "Device idx your choise")
	http_port := flag.Int("port", 3330, "Http Api listen port")
	flag.Parse()

	// 启动 gpu miner
	var miner GpuMiner
	miner.InitBuildProgram(*cldir, *plat_name, *dv_id, *group_size)

	// 启动 http 监听
	go RunHttpRpcService(&miner, *http_port)

	///////////////////////////////////////////////////

	s := <-c
	fmt.Println("Got signal:", s)

}
