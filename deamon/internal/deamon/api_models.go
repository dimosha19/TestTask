package deamon

type PID int

type proc struct {
	User    string   `json:"User"`
	Pid     PID      `json:"Pid"`
	Cpu     string   `json:"Cpu"`
	Mem     string   `json:"Mem"`
	Vsz     string   `json:"Vsz"`
	Rss     string   `json:"Rss"`
	Tty     string   `json:"Tty"`
	Stat    string   `json:"Stat"`
	Start   string   `json:"Start"`
	Time    string   `json:"Time"`
	Command []string `json:"Command"`
}

type ProcessResponse struct {
	Error   *string `json:"error,omitempty"`
	Process []proc  `json:"Process,omitempty"`
}
