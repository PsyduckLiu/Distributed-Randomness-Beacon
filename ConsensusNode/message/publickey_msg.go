package message

import (
	"consensusNode/signature"
	"crypto/ecdsa"
	"crypto/x509"
	"fmt"
)

// create new public key message of backups
func CreateIdentityMsg(t MType, id int64, sk *ecdsa.PrivateKey) *ConMessage {
	// sign message.Payload
	marshalledKey, err := x509.MarshalPKIXPublicKey(&sk.PublicKey)
	if err != nil {
		panic(fmt.Errorf("===>[ERROR from SetupConfig]setup conf curve(marshalled Key) failed, err:%s", err))
	}
	sig := signature.GenerateSig(marshalledKey, sk)
	identityMsg := &ConMessage{
		Typ:     t,
		Sig:     sig,
		From:    id,
		Payload: marshalledKey,
	}

	return identityMsg
}
