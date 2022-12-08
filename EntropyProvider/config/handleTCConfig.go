package config

import (
	"fmt"
	"io/ioutil"
	"math/big"
	"net/http"
	"os"
	"strconv"
	"strings"
)

// get groupLength L from config file
func GetL() int {
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

	output := strings.Fields(string(body))
	l, err := strconv.Atoi(output[3])
	if err != nil {
		panic(fmt.Errorf("===>[ERROR from GroupLength]get GroupLength:%s", err))
	}

	return l
}

// get g from config file
func GetG() *big.Int {
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

	output := strings.Fields(string(body))
	g := new(big.Int)
	g.SetString(output[1], 10)

	return g
}

// get N from config file
func GetN() *big.Int {
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

	output := strings.Fields(string(body))
	N := new(big.Int)
	N.SetString(output[30], 10)

	return N
}

// get primes from config file
func GetPrimes() []*big.Int {
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

	output := strings.Fields(string(body))

	var primes []*big.Int
	prime0 := new(big.Int)
	prime0.SetString(strings.Trim(output[32], "\""), 10)
	primes = append(primes, prime0)
	prime1 := new(big.Int)
	prime1.SetString(strings.Trim(output[34], "\""), 10)
	primes = append(primes, prime1)

	return primes
}

// get N from config file
func GetMArray() []*big.Int {
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

	output := strings.Fields(string(body))

	t, err := strconv.Atoi(output[107])
	if err != nil {
		panic(fmt.Errorf("===>[ERROR from TimeParameter]get TimeParameter:%s", err))
	}

	var mArray []*big.Int
	for i := 0; i < t+2; i++ {
		mBigint, _ := new(big.Int).SetString(output[6+2*i], 10)
		mArray = append(mArray, mBigint)
	}

	return mArray
}
