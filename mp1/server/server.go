package main

import (
	"bufio"
	"bytes"
	"fmt"
	"log"
	"net"
	"os/exec"
	"strings"
)

/*
  This function handles the read pattern and send data back stuff.
*/
func HandleConn(conn net.Conn) {
	pattern, err := bufio.NewReader(conn).ReadString('\n')
	if err != nil {
		log.Fatal(err)
	}
	pattern = strings.TrimSuffix(pattern, "\n")
	input := "cat ../machine.i.log | grep -n " + pattern
	cmd := exec.Command("bash", "-c", input)
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		log.Fatal(err)
	}
	err = cmd.Start()
	if err != nil {
		log.Fatal(err)
	}
	buf := new(bytes.Buffer)
	buf.ReadFrom(stdout)
	output := buf.String()
	dataSend := "machine.i.log\n" + output
	fmt.Print("data sent\n")
	fmt.Fprintf(conn, dataSend)
	conn.Close()
}

/*
  The main function for server. In this function, the requests are handled.
*/
func main() {
	se, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Fatal(err)
	}
	for {
		//infinite loop
		conn, err := se.Accept()
		if err != nil {
			log.Fatal(err)
		}
		go HandleConn(conn)
	}
}
