package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"strings"
	"sync"
)

type Package struct {
	Type  string            `json:"type"`
	Data  map[string]string `json:"data"`
	Key   string            `json:"key"`
	Value string            `json:"value"`
}

var mutex = &sync.Mutex{}

func main() {
	target := flag.String("d", "", "target peer to dial")
	flag.Parse()

	if *target == "" {
		log.Fatal("Provide ip:port for connection using flag -d")
	}

	conn, err := net.Dial("tcp", *target)
	if err != nil {
		log.Fatal("Cannot connect")
	}

	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Println("Text to send: ")
		sendData, err := reader.ReadString('\n')
		if err != nil {
			log.Fatal(err)
		}
		bufio.NewReader(conn).ReadString('\n')

		sendData = strings.Replace(sendData, "\n", "", -1)
		s := strings.Split(sendData, ":")

		mutex.Lock()
		switch s[0] {
		case "add":
			if len(s) < 3 {
				fmt.Print("format is add:key:value")

			} else {
				kv := map[string]string{s[1]: s[2]}
				p := Package{
					Type:  "add",
					Data:  kv,
					Key:   s[1],
					Value: s[2],
				}

				// send to socket
				fmt.Fprintf(conn, p.Marshal())
			}
		case "get":
			if len(s) < 2 {
				fmt.Print("format is add:key:value")

			} else {
				p := Package{
					Type: "get",
					Key:  s[1],
				}

				fmt.Fprintf(conn, p.Marshal())
				message, _ := bufio.NewReader(conn).ReadString('\n')
				fmt.Print("Message from server: " + message)
			}
		}

		mutex.Unlock()
	}
}

func (p Package) Marshal() string {
	bytes, _ := json.Marshal(p)
	return string(bytes) + "\n"
}
