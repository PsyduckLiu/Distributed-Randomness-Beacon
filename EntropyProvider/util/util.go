package util

import (
	"bytes"
	"crypto/rand"
	"crypto/sha256"
	"io/ioutil"
	"net/http"
	"os"
	"strings"

	"github.com/algorand/go-algorand/crypto"
	"github.com/algorand/go-algorand/protocol"

	"fmt"
)

// to be modified
const MaxFaultyNode = 1
const TotalNodeNum = 3*MaxFaultyNode + 1

type MessageHashable struct {
	Data []byte
}

func (s MessageHashable) ToBeHashed() (protocol.HashID, []byte) {
	return "msg", s.Data
}

// generate random message string for VRF
func RandString() (b MessageHashable) {
	d := make([]byte, 32)
	_, err := rand.Read(d)
	if err != nil {
		panic(err)
	}

	fmt.Printf("===>[VRF]New random string is %s\n", d)
	return MessageHashable{d}
}

// convert crypto.VrfProof([80]byte) to binary string
func BytesToBinaryString(bs crypto.VrfProof) string {
	buf := bytes.NewBuffer([]byte{})
	for _, v := range bs {
		buf.WriteString(fmt.Sprintf("%08b", v))
	}

	return buf.String()
}

// Get Port(20000 + id) for connection between entropy node and consensus node
func EntropyPortByID(id int) int {
	return 20000 + int(id)
}

// Hash any type message v, using SHA256
func Digest(v interface{}) []byte {
	h := sha256.New()
	h.Write([]byte(fmt.Sprintf("%v", v)))
	digest := h.Sum(nil)

	return digest
}

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
