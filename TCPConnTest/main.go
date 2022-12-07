package main

import (
	"os"
	"strconv"
	"time"
)

func main() {
	if len(os.Args) < 2 {
		panic("[Command line arguments]Usage: input id")
	}

	service, _ := strconv.Atoi(os.Args[1])

	if service == 0 {
		for i := 0; i < 10; i++ {
			time.Sleep(10 * time.Millisecond)
			SendMessage("hello" + string(rune(i)))
		}
	}

	if service == 1 {

	}

}
