package service

import (
	"fmt"
	"net"
	"time"
)

func SendMessage(msg string) error {
	fmt.Println("message is ",msg)
	// dial remote TCP port
	conn, err := net.DialTCP("tcp", nil, &net.TCPAddr{IP: net.IPv4(152, 136, 151, 161), Port: 2211})
	if err != nil {
		fmt.Println(time.Now())
		fmt.Println("===>[Fail from sendVRFMsg]Dial tcp err:", err)
		return err
	}

	fmt.Println(time.Now())
	fmt.Println("===>Dial tcp success")

	_, err = conn.Write([]byte(msg))
	if err != nil {
		panic(fmt.Errorf("===>[ERROR from WriteTCP]write to node failed:%s", err))
	}

	fmt.Println(time.Now())
	fmt.Println("===>Send message success")

	return nil
}
