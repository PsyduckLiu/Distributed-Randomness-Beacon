package p2pnetwork

import (
	"consensusNode/config"
	"consensusNode/message"
	"consensusNode/signature"
	"consensusNode/util"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net"
	_ "net/http/pprof"
	"runtime"
	"strconv"
	"sync"
	"time"
)

type P2pNetwork interface {
	GetPrimaryConn(primaryId int64) *net.TCPConn
	GetMySecretkey() *ecdsa.PrivateKey
	SendUniqueNode(conn *net.TCPConn, v interface{}) error
	BroadCast(v interface{}) error
}

// [SrvHub]: contains all TCP connections with other nodes
// [Peers]: map TCP connect to an int number
// [MsgChan]: a channel connects [p2p] with [state(consensus)], deliver consensus message, corresponding to [ch] in [state(consensus)]
type SimpleP2p struct {
	NodeId         int64
	SrvHub         *net.TCPListener
	Peers          map[string]*net.TCPConn
	Ip2Id          map[string]int64
	PrivateKey     *ecdsa.PrivateKey
	PeerPublicKeys map[int64]*ecdsa.PublicKey
	MsgChan        chan<- *message.ConMessage
	mutex          sync.Mutex
}

// new simple P2P liarary
func NewSimpleP2pLib(id int64, msgChan chan<- *message.ConMessage) P2pNetwork {
	// get specified curve
	marshalledCurve := config.GetCurve()
	pub, err := x509.ParsePKIXPublicKey([]byte(marshalledCurve))
	if err != nil {
		panic(fmt.Errorf("===>[ERROR from NewSimpleP2pLib]Parse elliptic curve error:%s", err))
	}
	normalPublicKey := pub.(*ecdsa.PublicKey)
	curve := normalPublicKey.Curve
	fmt.Printf("Curve is %v\n", curve.Params())

	// generate private key
	privateKey, err := ecdsa.GenerateKey(curve, rand.Reader)
	if err != nil {
		panic(fmt.Errorf("===>[ERROR from NewSimpleP2pLib]Generate private key err:%s", err))
	}
	fmt.Printf("===>[P2P]My own key is: %v\n", privateKey)

	// listen port 30000+id
	port := util.PortByID(id)
	s, err := net.ListenTCP("tcp4", &net.TCPAddr{
		Port: port,
	})
	if err != nil {
		panic(fmt.Errorf("===>[ERROR from NewSimpleP2pLib]Listen TCP err:%s", err))
	}
	fmt.Printf("===>[P2P]Node[%d] is waiting at:%s:%d\n", id, util.MyIPAddr, port)

	// write new node details into config
	if id == 0 {
		config.NewConsensusNode(id, util.MyIPAddr+":"+strconv.FormatInt(int64(port), 10), hex.EncodeToString(elliptic.Marshal(curve, privateKey.PublicKey.X, privateKey.PublicKey.Y)))
	}

	sp := &SimpleP2p{
		NodeId:         id,
		SrvHub:         s,
		Peers:          make(map[string]*net.TCPConn),
		Ip2Id:          make(map[string]int64),
		PrivateKey:     privateKey,
		PeerPublicKeys: make(map[int64]*ecdsa.PublicKey),
		MsgChan:        msgChan,
		mutex:          sync.Mutex{},
	}

	go sp.monitor(id)
	sp.dialTcp(id)

	return sp
}

// connect to all known nodes
func (sp *SimpleP2p) dialTcp(id int64) {
	nodeConfig := config.GetConsensusNode()

	for i := 0; i < len(nodeConfig); i++ {
		if int64(i) != id && nodeConfig[i].Ip != "0" {
			// resolve TCP address and dial TCP conn
			addr, err := net.ResolveTCPAddr("tcp4", nodeConfig[i].Ip)
			if err != nil {
				panic(fmt.Errorf("===>[ERROR from dialTcp]Resolve TCP Addr err:%s", err))
			}
			conn, err := net.DialTCP("tcp", nil, addr)
			if err != nil {
				// panic(fmt.Errorf("===>[ERROR from dialTcp]DialTCP err:%s", err))
				fmt.Println("===>[dialTcp]Failed to connect with", addr)
				continue
			}

			// new identity message
			// send identity message to origin nodes
			kMsg := message.CreateIdentityMsg(message.MTIdentity, sp.NodeId, sp.PrivateKey)
			if err := sp.SendUniqueNode(conn, kMsg); err != nil {
				panic(fmt.Errorf("===>[ERROR from dialTcp]Send Identity message error:%s", err))
			}

			go sp.waitData(conn)
		}
	}
}

// add new node OR remove old node
func (sp *SimpleP2p) monitor(id int64) {
	fmt.Printf("===>[P2P]Consensus node[%d] is waiting at:%s\n", id, sp.SrvHub.Addr().String())

	for {
		conn, err := sp.SrvHub.AcceptTCP()

		// remove a node
		if err != nil {
			fmt.Printf("===>[P2P]P2p network accept err:%s\n", err)
			if err == io.EOF {
				fmt.Printf("===>[P2P] Node[%d] Remove peer node[%d]%s\n", sp.NodeId, sp.Ip2Id[conn.RemoteAddr().String()], conn.RemoteAddr().String())
				// config.RemoveConsensusNode(sp.Ip2Id[conn.RemoteAddr().String()])
				delete(sp.Peers, conn.RemoteAddr().String())
				delete(sp.PeerPublicKeys, sp.Ip2Id[conn.RemoteAddr().String()])
				delete(sp.Ip2Id, conn.RemoteAddr().String())
			}
			continue
		}

		// add a new node
		sp.Peers[conn.RemoteAddr().String()] = conn

		// new identity message
		// send identity message to origin nodes
		kMsg := message.CreateIdentityMsg(message.MTIdentity, sp.NodeId, sp.PrivateKey)
		if err := sp.SendUniqueNode(conn, kMsg); err != nil {
			panic(fmt.Errorf("===>[ERROR from dialTcp]Send Identity message error:%s", err))
		}

		go sp.waitData(conn)
	}
}

// remove old node AND deliver consensus mseeage by [MsgChan]
func (sp *SimpleP2p) waitData(conn *net.TCPConn) {
	buf := make([]byte, 4096)
	for {
		n, err := conn.Read(buf)

		// remove a node
		if err != nil {
			fmt.Printf("===>[P2P]P2p network capture data err:%s\n", err)
			if err == io.EOF {
				fmt.Printf("===>[P2P]Node[%d] Remove peer node[%d]%s\n", sp.NodeId, sp.Ip2Id[conn.RemoteAddr().String()], conn.RemoteAddr().String())
				// config.RemoveConsensusNode(sp.Ip2Id[conn.RemoteAddr().String()])
				delete(sp.Peers, conn.RemoteAddr().String())
				delete(sp.PeerPublicKeys, sp.Ip2Id[conn.RemoteAddr().String()])
				delete(sp.Ip2Id, conn.RemoteAddr().String())
				return
			}
			continue
		}

		if n == 0 {
			fmt.Println("empty message!!!")
			continue
		}

		// handle a consensus message
		conMsg := &message.ConMessage{}
		fmt.Println("read from", conn.RemoteAddr().String(), time.Now())
		cMsgZip, err := util.Decode(buf[:n])
		if err != nil {
			panic(fmt.Errorf("===>[ERROR from waitData]Decode data err:%s", err))
		}
		if err := json.Unmarshal(cMsgZip, conMsg); err != nil {
			fmt.Println(string(cMsgZip))
			panic(fmt.Errorf("===>[ERROR from waitData]Unmarshal data err:%s", err))
		}

		switch conMsg.Typ {
		// handle new identity message from backups
		case message.MTIdentity:
			nodeConfig := config.GetConsensusNode()

			// unmarshal public key
			pub, err := x509.ParsePKIXPublicKey([]byte(conMsg.Payload))
			if err != nil {
				panic(fmt.Errorf("===>[ERROR from waitData]Parse public key error:%s", err))
			}
			newPublicKey := pub.(*ecdsa.PublicKey)

			// verify signature
			verify := signature.VerifySig(conMsg.Payload, conMsg.Sig, newPublicKey)
			if !verify {
				fmt.Printf("===>[ERROR from waitData]Verify new public key Signature failed, From Node[%d], IP[%s]\n", conMsg.From, conn.RemoteAddr().String())
				break
			}

			// store remote ip-conn-id-pk relation
			sp.mutex.Lock()
			sp.Peers[conn.RemoteAddr().String()] = conn
			sp.Ip2Id[conn.RemoteAddr().String()] = int64(conMsg.From)
			sp.PeerPublicKeys[int64(conMsg.From)] = newPublicKey
			sp.mutex.Unlock()

			// write new node details into config
			if sp.NodeId == 0 {
				config.NewConsensusNode(conMsg.From, nodeConfig[conMsg.From].Ip, hex.EncodeToString(elliptic.Marshal(newPublicKey.Curve, newPublicKey.X, newPublicKey.Y)))
			}

			fmt.Printf("===>[P2P]Get new public key from Node[%d], IP[%s]\n", conMsg.From, conn.RemoteAddr().String())
			fmt.Printf("===>[P2P]Node[%d<=>%d]Connected=[%s<=>%s]\n", sp.NodeId, conMsg.From, conn.LocalAddr().String(), conn.RemoteAddr().String())

		// handle consensus message from backups
		default:
			sp.MsgChan <- conMsg
		}
	}
}

// BroadCast message to all connected nodes
func (sp *SimpleP2p) BroadCast(v interface{}) error {
	if v == nil {
		return fmt.Errorf("===>[ERROR from BroadCast]empty msg body")
	}

	data, err := json.Marshal(v)
	dataZip, err := util.Encode(data)
	if err != nil {
		return fmt.Errorf("===>[ERROR from BroadCast]Marshal data err:%s", err)
	}

	for name, conn := range sp.Peers {
		time.Sleep(100 * time.Millisecond)
		go WriteTCP(conn, dataZip, name)
	}

	return nil
}

// BroadCast message to all connected nodes
func (sp *SimpleP2p) SendUniqueNode(conn *net.TCPConn, v interface{}) error {
	if v == nil {
		return fmt.Errorf("===>[ERROR from SendUniqueNode]empty msg body")
	}

	data, err := json.Marshal(v)
	dataZip, err := util.Encode(data)
	if err != nil {
		return fmt.Errorf("===>[ERROR from SendUniqueNode]Marshal data err:%s", err)
	}

	go WriteTCP(conn, dataZip, conn.RemoteAddr().String())

	return nil
}

// write message by TCP connection channel
func WriteTCP(conn *net.TCPConn, v []byte, name string) {
	_, err := conn.Write(v)
	if err != nil {
		fmt.Printf("===>[ERROR from WriteTCP]write to node[%s] err:%s\n", name, err)
		panic(err)
	}

	fmt.Println(time.Now())
	fmt.Printf("===>[Sending]Send request to Address[%s] success\n", conn.RemoteAddr().String())
	runtime.Goexit()
}

// Get Peer Publickey
func (sp *SimpleP2p) GetPrimaryConn(primaryId int64) *net.TCPConn {
	var conn *net.TCPConn
	for ip, id := range sp.Ip2Id {
		if id == primaryId {
			conn = sp.Peers[ip]
		}
	}

	return conn
}

// Get My Secret key
func (sp *SimpleP2p) GetMySecretkey() *ecdsa.PrivateKey {
	return sp.PrivateKey
}
