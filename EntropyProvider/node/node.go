package node

import (
	"crypto/ecdsa"
	"crypto/rand"
	"crypto/x509"
	"encoding/json"
	"entropyNode/commitment"
	"entropyNode/config"
	"entropyNode/message"
	"entropyNode/util"
	"fmt"
	"net"
	"runtime"
	"strconv"
	"time"

	"github.com/algorand/go-algorand/crypto"
)

// initialize an entropy node
func StartEntropyNode(id int) {
	var signal chan interface{}
	fmt.Printf("===>[Initialization]Node[%d] starts running\n", id)

	// get specified curve
	// config.ReadConfig()
	marshalledCurve := config.GetCurve()
	pub, err := x509.ParsePKIXPublicKey([]byte(marshalledCurve))
	if err != nil {
		panic(fmt.Errorf("===>[ERROR from StartEntropyNode]Parse elliptic curve error:%s", err))
	}
	normalPublicKey := pub.(*ecdsa.PublicKey)
	curve := normalPublicKey.Curve
	// fmt.Printf("Curve is %v\n", curve.Params())

	// generate ECDSA private-public key pair based on public parameters
	ecdsaSK, err := ecdsa.GenerateKey(curve, rand.Reader)
	if err != nil {
		panic(fmt.Errorf("===>[ERROR from StartEntropyNode]Generate private key err:%s", err))
	}
	fmt.Println("===>[Initialization]ECDSA key pair is:", ecdsaSK)

	// generate VRF private-public key pair
	vrfPK, vrfSK := crypto.VrfKeygen()
	fmt.Println("===>[Initialization]VRF public key is:", vrfPK)
	fmt.Println("===>[Initialization]VRF secret key is:", vrfSK)

	// after initializing private-public key pairs, start watching output.yml
	go WatchConfig(ecdsaSK, vrfSK, id, signal)
	s := <-signal
	fmt.Printf("===>[EXIT from StartEntropyNode]Node[%d] exit because of:%s\n", id, s)
}

// entropy node continues watching on output.yml
// when [previousOutput] in output.yml changes, entropy node starts calculating VRF and sending TC
func WatchConfig(ecdsaSK *ecdsa.PrivateKey, vrfSK crypto.VrfPrivkey, id int, sig chan interface{}) {
	// get the earliest output written in output.yml
	previousOutput := config.GetPreviousOutput()
	fmt.Println("\n===>[Watching]The earliest output is", previousOutput)

	for {
		time.Sleep(500 * time.Millisecond)

		// when new output comes, entropy node starts calculating VRF and sending TC
		newOutput := config.GetPreviousOutput()
		if previousOutput != newOutput && newOutput != "" {
			fmt.Println("\n===>[Watching]Output changed,new output is", newOutput)
			previousOutput = newOutput

			// calculate VRF result
			msg := util.RandString()
			vrfResult, ok := vrfSK.Prove(msg)
			if !ok {
				panic(fmt.Errorf("===>[ERROR from WatchConfig]Failed to construct VRF proof"))
			}
			vrfResultBinary := util.BytesToBinaryString(vrfResult)
			fmt.Println("===>[Watching]VRF result is:", util.BytesToBinaryString(vrfResult))
			fmt.Println("===>[Watching]VRF result's last bit is:", vrfResultBinary[len(vrfResultBinary)-1:])

			// compare VRF result with difficulty
			difficulty := config.GetDifficulty()
			vrfResultTail, err := strconv.Atoi(vrfResultBinary[len(vrfResultBinary)-1:])
			if err != nil {
				panic(fmt.Errorf("===>[ERROR from WatchConfig]Failed to get VRF result's last bit:%s", err))
			}

			// VRF result passes difficulty requirement
			if vrfResultTail == difficulty {
				fmt.Println("===>[Watching]Pass but no use!!!")
			}
			if true {
				fmt.Println("===>[Watching]VRF result passes difficulty requirement!!!")

				// generate timed commitment
				groupLength := config.GetL()
				c, h, rKSubOne, rK, a1, a2, a3, z := commitment.GenerateTimedCommitment(groupLength)
				fmt.Println("\n===>[Watching]In TC, c is", c)
				fmt.Println("===>[Watching]In TC, h is", h)
				fmt.Println("===>[Watching]In TC, rKSubOne is", rKSubOne)
				fmt.Println("===>[Watching]In TC, rK is", rK)
				fmt.Println("===>[Watching]In TC, a1 is", a1)
				fmt.Println("===>[Watching]In TC, a2 is", a2)
				fmt.Println("===>[Watching]In TC, a3 is", a3)
				fmt.Println("===>[Watching]In TC, z is", z)

				// marshal timed commitment parameters
				cMarshal, _ := c.MarshalJSON()
				hMarshal, _ := h.MarshalJSON()
				rKSubOneMarshal, _ := rKSubOne.MarshalJSON()
				rKMarshal, _ := rK.MarshalJSON()
				a1Marshal, _ := a1.MarshalJSON()
				a2Marshal, _ := a2.MarshalJSON()
				a3Marshal, _ := a3.MarshalJSON()
				zMarshal, _ := z.MarshalJSON()
				fmt.Printf("===>[Watching]Time commitment is:%v,%v,%v,%v\n", cMarshal, hMarshal, rKSubOneMarshal, rKMarshal)

				// send VRF messages and TC messages
				fmt.Println()
				time.Sleep(5 * time.Second)
				sendVRFMsg(ecdsaSK, vrfSK, vrfResult, msg.Data, int64(id), sig)
				time.Sleep(1 * time.Second)
				fmt.Println()
				sendTCMsg(ecdsaSK, int64(id), cMarshal, hMarshal, rKSubOneMarshal, rKMarshal, a1Marshal, a2Marshal, a3Marshal, zMarshal, sig)
			} else {
				fmt.Println("===>[Watching]VRF result doesn't pass difficulty requirement***")
			}
		}
	}

}

// Entropy node sends VRF message to all Consensus nodes for verification
func sendVRFMsg(ecdsaSK *ecdsa.PrivateKey, vrfSK crypto.VrfPrivkey, vrfResult crypto.VRFProof, msg []byte, id int64, sig chan interface{}) {
	// new VRF message
	vrfMsg := &message.EntropyVRFMessage{
		PublicKey: vrfSK.Pubkey(),
		VRFResult: vrfResult,
		ClientID:  id,
		Msg:       msg,
	}

	// get Consensus nodes' information
	nodeConfig := config.GetConsensusNode()
	for i := 0; i < len(nodeConfig); i++ {
		time.Sleep(200 * time.Millisecond)

		// dial remote TCP port
		addr, err := net.ResolveTCPAddr("tcp4", nodeConfig[i].Ip)
		if err != nil {
			panic(fmt.Errorf("===>[ERROR from dialTcp]Resolve TCP Addr err:%s", err))
		}
		addr.Port = util.EntropyPortByID(i)

		conn, err := net.DialTCP("tcp", nil, addr)
		if err != nil {
			fmt.Println(time.Now())
			fmt.Println("===>[Fail from sendVRFMsg]Dial tcp err:", err)
			continue
		}

		// send VRF message to Consensus nodes
		cMsg := message.CreateConMsg(message.MTVRFVerify, vrfMsg, ecdsaSK, id)
		bs, err := json.Marshal(cMsg)
		// bsZip, err := util.Encode(bs)
		if err != nil {
			panic(fmt.Errorf("===>[ERROR from sendVRFMsg]Marshal consensus message failed:%s", err))
		}
		fmt.Println("===>[Sending]Length of marshalled consensus VRF message is:", len(bs))

		go WriteTCP(i, conn, bs)
	}
}

// Entropy node sends TC message to all Consensus nodes
func sendTCMsg(ecdsaSK *ecdsa.PrivateKey, id int64, cMarshal []byte, hMarshal []byte, rKSubOneMarshal []byte, rKMarshal []byte, a1Marshal []byte, a2Marshal []byte, a3Marshal []byte, zMarshal []byte, sig chan interface{}) {
	// new time commitment message
	tcMsg := &message.EntropyTCMessage{
		ClientID:               id,
		TimeCommitmentC:        string(cMarshal),
		TimeCommitmentH:        string(hMarshal),
		TimeCommitmentrKSubOne: string(rKSubOneMarshal),
		TimeCommitmentrK:       string(rKMarshal),
		TimeCommitmentA1:       string(a1Marshal),
		TimeCommitmentA2:       string(a2Marshal),
		TimeCommitmentA3:       string(a3Marshal),
		TimeCommitmentZ:        string(zMarshal),
	}

	// get consensus nodes' information
	nodeConfig := config.GetConsensusNode()
	for i := 0; i < len(nodeConfig); i++ {
		time.Sleep(500 * time.Millisecond)

		// dial remote TCP port
		addr, err := net.ResolveTCPAddr("tcp4", nodeConfig[i].Ip)
		if err != nil {
			panic(fmt.Errorf("===>[ERROR from dialTcp]Resolve TCP Addr err:%s", err))
		}
		addr.Port = util.EntropyPortByID(i)

		conn, err := net.DialTCP("tcp", nil, addr)
		if err != nil {
			fmt.Println(time.Now())
			fmt.Println("===>[Fail from sendTCMsg]Dial tcp err:", err)
			continue
		}

		// send time commitment message to Consensus nodes
		cMsg := message.CreateConMsg(message.MTCommitFromEntropy, tcMsg, ecdsaSK, id)
		bs, err := json.Marshal(cMsg)
		// bsZip, err := util.Encode(bs)
		if err != nil {
			panic(fmt.Errorf("===>[ERROR from sendTCMsg]Marshal consensus message failed:%s", err))
		}
		fmt.Println("===>[Sending]Length of marshalled consensus TC message is:", len(bs))
		fmt.Println(string(bs))

		go WriteTCP(i, conn, bs)
	}
}

// write message by TCP connection channel
func WriteTCP(id int, conn *net.TCPConn, v []byte) {
	_, err := conn.Write(v)
	if err != nil {
		panic(fmt.Errorf("===>[ERROR from WriteTCP]write to node failed:%s", err))
	}

	fmt.Println(time.Now())
	fmt.Printf("===>[Sending]Send request to Node[%d], Address[%s] success\n", id, conn.RemoteAddr().String())
	fmt.Println(len(v))
	runtime.Goexit()
}
