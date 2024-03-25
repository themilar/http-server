package main

import (
	"fmt"
	"net"
	"os"
	"strings"
)

const CRLF = "\r\n"

func main() {
	// You can use print statements as follows for debugging, they'll be visible when running tests.
	fmt.Println("Logs from your program will appear here!")

	ln, err := net.Listen("tcp", "0.0.0.0:4221")
	if err != nil {
		fmt.Println("Failed to bind to port 4221")
		os.Exit(1)
	}
	defer ln.Close()
	conn, err := ln.Accept()
	if err != nil {
		fmt.Println("Error accepting connection: ", err)
		os.Exit(1)
	}
	buf := make([]byte, 1024)
	_, err = conn.Read(buf)
	if err != nil {
		fmt.Println("Error accepting connection: ", err)
	}
	req := string(buf)
	lines := strings.Split(req, CRLF)
	path := strings.Split(lines[0], " ")[1]
	fmt.Println(path)
	var res string
	if path == "/" {
		res = "HTTP/1.1 200 OK\r\n\r\n"
	} else {
		res = "HTTP/1.1 404 Not Found\r\n\r\n"
	}
	fmt.Println(res)
	conn.Write([]byte(res))
	if err != nil {
		fmt.Println("Error accepting connection: ", err)
		os.Exit(1)
	}
}
