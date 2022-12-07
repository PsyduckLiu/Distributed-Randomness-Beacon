package service

import (
	"fmt"
	"io"
	"net"
	"time"
)

func Accept() {
	port := 2211
	s, err := net.ListenTCP("tcp4", &net.TCPAddr{
		// IP:   net.IPv4(152, 136, 151, 161),
		Port: port,
	})
	if err != nil {
		panic(fmt.Errorf("===>[ERROR from NewSimpleP2pLib]Listen TCP err:%s", err))
	}
	fmt.Printf("===>[P2P]Node is waiting at:%s\n", s.Addr().String())

	buf := make([]byte, 8192)
	for {
		conn, err := s.AcceptTCP()

		if err != nil {
			fmt.Printf("===>[P2P]P2p network accept err:%s\n", err)
			if err == io.EOF {
				fmt.Printf("Remote Error%s", conn.RemoteAddr().String())
			}
			continue
		}
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
		fmt.Println(string(buf[:n]))
		fmt.Printf("===>[P2P]P2p network accept success:%s\n", conn.RemoteAddr().String())
	}
}
