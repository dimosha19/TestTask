package deamon

import (
	"bytes"
	"log"
	"net/http"
	"os/exec"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

func GetProcList(c *gin.Context) {
	var out bytes.Buffer

	cmd := exec.Command("ps", "aux")
	cmd.Stdout = &out

	err := cmd.Run()
	if err != nil {
		log.Fatal(err)
	}

	c.IndentedJSON(
		http.StatusOK,
		ProcessResponse{
			Process: splitByRow(out.String()),
		},
	)
}

func splitByRow(str string) []proc {
	var res []proc
	procList := strings.Split(str, "\n")

	for i, proc := range procList {
		if i > 0 && proc != "" {
			convertedProc := convertProcToStruct(proc)
			res = append(res, convertedProc)
		}
	}

	return res
}

func convertProcToStruct(process string) proc {
	procProperties := strings.Fields(process)

	pid, err := strconv.Atoi(procProperties[1])
	if err != nil {
		log.Fatal(err)
	}

	return proc{
		User:    procProperties[0],
		Pid:     PID(pid),
		Cpu:     procProperties[2],
		Mem:     procProperties[3],
		Vsz:     procProperties[4],
		Rss:     procProperties[5],
		Tty:     procProperties[6],
		Stat:    procProperties[7],
		Start:   procProperties[8],
		Time:    procProperties[9],
		Command: procProperties[10:],
	}
}
