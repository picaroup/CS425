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

func read_ips() []string{
	file, err := os.Open("ips.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	ipAddr := []string{}
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		//fmt.Println(scanner.Text())
		ipAddr = append(ipAddr, scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
	return ipAddr
}

func main() {
	if len(os.Args) <= 2 {
		fmt.Print("Usage: go run client.go grep [options] pattern\n")
		return
	}
	args := os.Args[2:]

	// var ips [1]string
	// ips[0] = "172.22.154.42"
	ipAddr := read_ips()
	fmt.Print(ipAddr)
	wg.Add(len(ipAddr))
	for i := range ipAddr {
		go request(ipAddr[i], args)
	}
	wg.Wait()
}
