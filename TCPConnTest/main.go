package main

import (
	"TCP/service"
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"strings"
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

	res, err := http.Get("http://152.136.151.161/output.yml")
	if err != nil {
		fmt.Fprintf(os.Stderr, "fetch: reading %s: %v\n", "http://152.136.151.161/output.yml", err)
		os.Exit(1)
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Fprintf(os.Stderr, "fetch: reading %s: %v\n", "http://152.136.151.161/output.yml", err)
		os.Exit(1)
	}

	res.Body.Close()
	fmt.Printf("%s", string(body))

	output := strings.Fields(string(body))
	fmt.Println(output[1])
	fmt.Println([]byte(output[1]))

	fmt.Println(hex.DecodeString(output[1]))
}
