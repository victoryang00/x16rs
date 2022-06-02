package cpumining

type CPUMining struct {
	supervene int
}

func NewCPUMining(supervene int) *CPUMining {
	return &CPUMining{
		supervene: supervene,
	}
}

// init
func (c *CPUMining) Init() error {
	return nil
}

// Turn off force statistics
func (c *CPUMining) CloseUploadHashrate() {

}

// Concurrent number
func (c *CPUMining) GetSuperveneWide() int {
	return c.supervene
}
