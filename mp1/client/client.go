package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
	"strings"
	"sync"
)

var wg sync.WaitGroup

func request(ipAddr string, args []string) {
	defer wg.Done()
	pattern := strings.Join(args, " ")
	conn, err := net.Dial("tcp", ipAddr+":8080")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("connected to %s\n", conn.RemoteAddr().String())
	fmt.Fprintf(conn, pattern+"\n")
	reader := bufio.NewReader(conn)
	machineName, err := reader.ReadString('\n')
	if err != nil {
		log.Fatal(err)
	}
	machineName = strings.TrimSuffix(machineName, "\n")
	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			break
		}
		fmt.Printf("%s %s", machineName, line)
	}
}

func main() {
	if len(os.Args) <= 2 {
		fmt.Print("Usage: go run client.go grep [options] pattern\n")
		return
	}
	args := os.Args[2:]

	var ips [1]string
	ips[0] = "172.22.154.42"

	wg.Add(len(ips))
	for i := range ips {
		go request(ips[i], args)
	}
	wg.Wait()
}
