package main

import (
	"fmt"
	"io"
	"net"
	"time"
)

func Accept() {
	port := 2211
	s, err := net.ListenTCP("tcp4", &net.TCPAddr{
		Port: port,
	})
	if err != nil {
		panic(fmt.Errorf("===>[ERROR from NewSimpleP2pLib]Listen TCP err:%s", err))
	}
	fmt.Printf("===>[P2P]Node is waiting at:%s\n", s.Addr().String())

	for {
		conn, err := s.AcceptTCP()

		if err != nil {
			fmt.Printf("===>[P2P]P2p network accept err:%s\n", err)
			if err == io.EOF {
				fmt.Printf("Remote Error%s", conn.RemoteAddr().String())
			}
			continue
		}

		go waitData(conn)
	}
}

func waitData(conn *net.TCPConn) {
	buf := make([]byte, 8192)
	for {
		n, err := conn.Read(buf)

		if err != nil {
			fmt.Printf("===>[P2P]P2p network accept err:%s\n", err)
			if err == io.EOF {
				fmt.Printf("Remote Error%s", conn.RemoteAddr().String())
			}
			continue
		}

		// handle a message
		fmt.Println("read from", conn.RemoteAddr().String(), time.Now())
		fmt.Println("Message contains:")
		fmt.Println(buf[:n])
	}
}
