package config

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
)

func ReadOutput() string {
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
	// outputByte, err := hex.DecodeString(output[1])

	// return string(outputByte)
	return output[1]
}

func ReadRunningStatus() string {
	res, err := http.Get("http://152.136.151.161/config.yml")
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
	// outputByte, err := hex.DecodeString(output[1])

	// return string(outputByte)
	return output[1]
}
