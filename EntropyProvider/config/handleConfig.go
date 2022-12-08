package config

import (
	"encoding/hex"
	"entropyNode/util"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"strings"
)

type Output struct {
	PreviousOutput string `mapstructure:"previousOutput"`
}

type NodeConfig struct {
	Ip string
	Pk string
}

// get difficulty from config file
func GetDifficulty() int {
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

	output := strings.Fields(string(body))
	difficulty, err := strconv.Atoi(output[1])
	if err != nil {
		panic(fmt.Errorf("===>[ERROR from GetDifficulty]get difficulty:%s", err))
	}

	return difficulty
}

// get curve from config file
func GetCurve() string {
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

	output := strings.Fields(string(body))
	curve, err := hex.DecodeString(output[3])
	if err != nil {
		panic(fmt.Errorf("===>[ERROR from GetCurve]get Curve:%s", err))
	}

	// fmt.Println(string(outputByte))
	return string(curve)
}

// get previous output from config file
func GetPreviousOutput() string {
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
	output := strings.Fields(string(body))

	return output[1]
}

// get consensus nodes from config file
func GetConsensusNode() []NodeConfig {
	var nodeConfig []NodeConfig

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

	output := strings.Fields(string(body))

	for i := 0; i < util.TotalNodeNum; i++ {
		var node NodeConfig
		node.Ip = output[i*4+7]
		node.Pk = output[i*4+9]
		nodeConfig = append(nodeConfig, node)
	}

	return nodeConfig
}

// read configurations from config file
func ReadConfig() {
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

	output := strings.Fields(string(body))

	fmt.Printf("\nReading Configuration:\n")
	fmt.Printf("Running:%v\n", output[23])
	fmt.Printf("Version:%s\n", output[27])
}
