package main

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"os"
	"strings"
)

func main() {
	// You can use print statements as follows for debugging, they'll be visible when running tests.
	fmt.Println("Logs from your program will appear here!")

	l, err := net.Listen("tcp", "0.0.0.0:6379")
	if err != nil {
		fmt.Println("Failed to bind to port 6379")
		os.Exit(1)
	}

	storage := NewStorage()

	for {
		conn, err := l.Accept()
		if err != nil {
			fmt.Println("Error accepting connection: ", err.Error())
			os.Exit(1)
		}

		go handleConnection(conn, storage)
	}
}

func handleConnection(conn net.Conn, storage *Storage) {
	defer conn.Close()

	for {
		resp, err := DecodeRESP(bufio.NewReader(conn))
		if err != nil {
			if err == io.EOF {
				break
			} else {
				fmt.Println("error decoding RESP: ", err.Error())
				return
			}
		}

		arr := resp.Array()
		if arr == nil || len(arr) == 0 {
			fmt.Println("no command found")
			return
		}

		command := arr[0].String()
		args := arr[1:]

		switch strings.ToLower(command) {
		case "ping":
			conn.Write([]byte("+PONG\r\n"))
		case "echo":
			if len(args) > 0 {
				conn.Write([]byte(fmt.Sprintf("$%d\r\n%s\r\n", len(args[0].String()), args[0].String())))
			} else {
				conn.Write([]byte("-ERR wrong number of arguments for command '" + command + "'\r\n"))
			}
		case "set":
			if len(args) == 2 {
				storage.Set(args[0].String(), args[1].String())
				conn.Write([]byte("+OK\r\n"))
			} else {
				conn.Write([]byte("-ERR wrong number of arguments for command '" + command + "'\r\n"))
			}
		case "get":
			if len(args) == 1 {
				value, err := storage.Get(args[0].String())
				if err != nil {
					conn.Write([]byte("-ERR cannot get key '" + args[0].String() + "'\r\n"))
				} else {
					conn.Write([]byte(fmt.Sprintf("+%s\r\n", value)))
				}
			} else {
				conn.Write([]byte("-ERR wrong number of arguments for command '" + command + "'\r\n"))
			}
		default:
			conn.Write([]byte("-ERR unknown command '" + command + "'\r\n"))
		}
	}
}
