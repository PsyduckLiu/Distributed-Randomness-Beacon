package service

import (
	"fmt"
	"net"
	"time"
)

func SendMessage(msg string) error {
	// dial remote TCP port
	conn, err := net.DialTCP("tcp", nil, &net.TCPAddr{IP: net.IPv4(152, 136, 151, 161), Port: 2211})
	if err != nil {
		fmt.Println(time.Now())
		fmt.Println("===>[Fail from sendVRFMsg]Dial tcp err:", err)
		return err
	}

	_, err = conn.Write([]byte(msg))
	if err != nil {
		panic(fmt.Errorf("===>[ERROR from WriteTCP]write to node failed:%s", err))
	}

	return nil
}
