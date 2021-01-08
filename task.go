package ci

type Step struct {
	Name        string
	Image       string
	Cmd         []string
	Args        []string
	CpuLimit    int
	MemoryLimit int
	//BuildID???
	//Host??
}

type stepExec struct {
	Name     string
	Image    string
	Cmd      []string
	Args     []string
	CpuLimit int
	Memory   int
	BuildId  int64
}
