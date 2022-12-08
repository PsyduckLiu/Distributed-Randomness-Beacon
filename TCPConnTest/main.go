package main

import (
	"TCP/service"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
)

func main() {
	if len(os.Args) < 2 {
		panic("[Command line arguments]Usage: input id")
	}

	role, _ := strconv.Atoi(os.Args[1])

	if role == 0 {
		fmt.Println("I'm sender")
		for i := 0; i < 10; i++ {
			fmt.Println("Round", i)
			// time.Sleep(500 * time.Millisecond)
			service.SendMessage("hello")
		}
	}

	if role == 1 {
		service.Accept()
	}

	res, err := http.Get("http://152.136.151.161/TC.yml")
	if err != nil {
		fmt.Fprintf(os.Stderr, "fetch: reading %s: %v\n", "http://152.136.151.161/TC.yml", err)
		os.Exit(1)
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Fprintf(os.Stderr, "fetch: reading %s: %v\n", "http://152.136.151.161/TC.yml", err)
		os.Exit(1)
	}

	res.Body.Close()
	fmt.Println(body)

	// fmt.Println(N)
	// outputByte, err := hex.DecodeString(output[1])
	// fmt.Println(string(outputByte))
}
