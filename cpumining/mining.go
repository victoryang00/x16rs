package cpumining

type CPUMining struct {
	supervene int
}

func NewCPUMining(supervene int) *CPUMining {
	return &CPUMining{
		supervene: supervene,
	}
}

// 初始化
func (c *CPUMining) Init() error {
	return nil
}

// 关闭算力统计
func (c *CPUMining) CloseUploadHashrate() {

}

// 并发数
func (c *CPUMining) GetSuperveneWide() int {
	return c.supervene
}
