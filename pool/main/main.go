package main

import (
	"fmt"
	"os"
	"os/signal"
	"regexp"
	"strconv"
	"strings"
)

// 矿池
type Config struct {
	Pool      string // 127.0.0.1:3339 // 矿池地址
	Reward    string // 奖励地址
	Supervene uint32 // 多线程
}

func main() {

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, os.Kill)

	fmt.Println("mining pool start ...")

	conf := ReadParseConfig()
	StartPoolWorker(conf)

	s := <-c
	fmt.Println("Got signal:", s)

}

// 读取并解析配置文件
func ReadParseConfig() *Config {
	file, err := os.Open("./poolworker.config.yml")
	if err != nil {
		panic("cannot find the config file './config.yml'")
	}
	stat, _ := file.Stat()
	content := make([]byte, stat.Size())
	file.Read(content)
	contentstr := string(content)
	// 处理注释
	rex0, _ := regexp.Compile(`\n`)
	contentstr = rex0.ReplaceAllString(contentstr, "\n;")
	rex1, _ := regexp.Compile(`#[^\n]*\n`)
	contentstr = rex1.ReplaceAllString(contentstr, ";")
	rex2, _ := regexp.Compile(`\s+`)
	contentstr = rex2.ReplaceAllString(contentstr, "")
	rex3, _ := regexp.Compile(`;+`)
	contentstr = rex3.ReplaceAllString(contentstr, ";")
	contentstr = strings.Trim(contentstr, ";")
	// 解析值
	params := make(map[string]string)
	keyvals := strings.Split(contentstr, ";")
	for i := 0; i < len(keyvals); i++ {
		ps := strings.SplitN(keyvals[i], ":", 2)
		if len(ps) == 2 {
			params[ps[0]] = ps[1]
		}
	}
	// 判断赋值
	var conf Config
	pool, hs1 := params["pool"]
	if !hs1 {
		panic("config file key 'pool' is must set.")
	}
	conf.Pool = pool
	reward, hs2 := params["reward"]
	if !hs2 {
		panic("config file key 'reward' is must set.")
	}
	conf.Reward = reward
	supervene, hs3 := params["supervene"]
	if !hs3 {
		panic("config file key 'supervene' is must set.")
	}
	sunm, e3 := strconv.ParseUint(supervene, 10, 0)
	if e3 != nil {
		panic("config file key 'supervene' is error.")
	}
	if sunm <= 0 || sunm > 256 {
		panic("config file key 'supervene' value must between 1 and 256.")
	}
	conf.Supervene = uint32(sunm)

	//fmt.Println(conf)

	return &conf
}
