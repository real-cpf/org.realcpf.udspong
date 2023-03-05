package main

import (
	"fmt"
	"io"
	"log"
	"net"
	"os"

	"github.com/panjf2000/gnet/v2"
	"reacpf.org/udspong/pongserver"
)

var udspath = "/tmp/uds.sock"

func main() {

	s := new(pongserver.PongServer)
	log.Fatal(gnet.Run(s, "unix:///tmp/uds.sock", gnet.WithMulticore(true)))
}

func gonative() {
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
