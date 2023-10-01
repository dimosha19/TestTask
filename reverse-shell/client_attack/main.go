package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
	"sync"
	"time"
)

func sendToServer(conn net.Conn) {
	for {
		reader := bufio.NewReader(os.Stdin)
		fmt.Println(">>>")

		text, err := reader.ReadString('\n')
		if err != nil {
			log.Fatal(err)
		}
		_, err = conn.Write([]byte(text + "\n"))
		if err != nil {
			fmt.Println("connection error")
			return
		}
	}
}

func receiveFromServer(conn net.Conn) {
	for {
		message := bufio.NewScanner(conn)

		for message.Scan() {
			fmt.Println(message.Text())
		}
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

	wg := sync.WaitGroup{}
	wg.Add(2)

	go func() {
		defer wg.Done()
		receiveFromServer(conn)
	}()

	go func() {
		defer wg.Done()
		sendToServer(conn)
	}()

	wg.Wait()
}
