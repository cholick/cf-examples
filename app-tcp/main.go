package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
	"strings"
)

func getPort() string {
	port := os.Getenv("PORT")
	if port == "" {
		port = "4000"
	}
	return port
}

func handle(conn net.Conn) {
	reader := bufio.NewReader(conn)
	for {
		m, err := reader.ReadString('\n')
		if err != nil {
			break
		}
		_, err = conn.Write([]byte(strings.ToUpper(m)))
		if err != nil {
			break
		}
	}
	fmt.Printf("%v closed\n", conn.RemoteAddr())
	conn.Close()
}

func main() {
	port := getPort()
	fmt.Printf("Listening on %v\n", port)
	listener, err := net.Listen("tcp", fmt.Sprintf(":%v", port))
	if err != nil {
		log.Fatal(err)
	}
	defer listener.Close()

	for {
		conn, err := listener.Accept()
		if err != nil {
			println(err.Error())
		} else {
			fmt.Printf("%v accepted\n", conn.RemoteAddr())
			go handle(conn)
		}
	}
}
