package main

import (
	"flag"
	"fmt"
	"net"
	"os"
	"strings"
)

const CRLF = "\r\n"

func main() {
	// You can use print statements as follows for debugging, they'll be visible when running tests.
	fmt.Println("Logs from your program will appear here!")
	dir := flag.String("directory", "", "enter a directory")
	flag.Parse()
	ln, err := net.Listen("tcp", "0.0.0.0:4221")
	if err != nil {
		fmt.Println("Failed to bind to port 4221", err)
		os.Exit(1)
	}
	defer ln.Close()
	for {
		conn, err := ln.Accept()
		if err != nil {
			fmt.Println("Error accepting connection: ", err)
			os.Exit(1)
		}
		go func() {
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
			} else if strings.HasPrefix(path, "/echo/") {
				msg := path[6:]
				res = fmt.Sprintf("HTTP/1.1 200 OK\r\nContent-Type: text/plain\r\nContent-Length: %v\r\n\r\n%v", len(msg), msg)
			} else if path == "/user-agent" {
				msg := strings.Split(lines[2], " ")[1]
				res = fmt.Sprintf("HTTP/1.1 200 OK\r\nContent-Type: text/plain\r\nContent-Length: %v\r\n\r\n%v", len(msg), msg)
			} else if strings.HasPrefix(path, "/files/") && *dir != "" {
				file := path[7:]
				fmt.Println(*dir + file)
				if content, err := os.ReadFile(*dir + file); err == nil {
					content := string(content)
					res = fmt.Sprintf("HTTP/1.1 200 OK\r\nContent-Type: application/octet-stream\r\nContent-Length: %d\r\n\r\n%s", len(content), content)
				} else {
					res = "HTTP/1.1 404 File Not Found\r\n\r\n"
				}

			} else {
				res = "HTTP/1.1 404 Not Found\r\n\r\n"
			}
			fmt.Println(res)
			conn.Write([]byte(res))
			if err != nil {
				fmt.Println("Error accepting connection: ", err)
				os.Exit(1)
			}
		}()
	}

}
