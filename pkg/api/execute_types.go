package api

type ExecuteRequest struct {
	Code string `json:"code"`
}

type ExecuteResponse struct {
	Stdout   string `json:"stdout"`
	Stderr   string `json:"stderr"`
	ExitCode int    `json:"exitCode"`
	CpuTime  int64  `json:"cpuTime"`
	Memory   int64  `json:"memory"`
}
