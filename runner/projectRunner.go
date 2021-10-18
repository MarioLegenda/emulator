package runner

type ProjectRunResult struct {
	Success bool `json:"success"`
	Result string `json:"result"`
	Timeout int `json:"timeout"`
}

