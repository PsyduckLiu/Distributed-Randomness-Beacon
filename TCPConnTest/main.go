package main

import (
	"TCP/service"
	"os"
	"strconv"
	"time"
)

func main() {
	if len(os.Args) < 2 {
		panic("[Command line arguments]Usage: input id")
	}

	role, _ := strconv.Atoi(os.Args[1])

	if role == 0 {
		for i := 0; i < 10; i++ {
			time.Sleep(500 * time.Millisecond)
			service.SendMessage("hello" + string(rune(i)))
		}
	}

	if role == 1 {
		service.Accept()
	}

}
