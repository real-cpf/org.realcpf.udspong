package main

import (
	"fmt"
	"io"
	"net"
	"os"
)

var udspath = "/tmp/uds.sock"

func main() {
	os.Remove(udspath)
	l, err := net.Listen("unix", udspath)
	if err != nil {
		fmt.Printf("unix listen error %s", err)
	}
	fmt.Println("listen..." + udspath)
	defer l.Close()
	for {
		conn, _ := l.Accept()
		fmt.Println("conn ", conn.LocalAddr().String())

		go func(c net.Conn) {
			for {
				buf := make([]byte, 32)
				_, err := c.Read(buf)
				if err != nil && err != io.EOF {
					fmt.Println(err)
					c.Close()
					break
				}
				if err == io.EOF {
					break
				}
				fmt.Println("rev", string(buf))
				c.Write([]byte("hello client"))
				c.Close()
			}
		}(conn)

	}
}
