package crypto

import (
	"crypto/sha256"
	"crypto/ecdsa"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"math/big"
)

func SignData(privateKey *ecdsa.PrivateKey, data []byte) (string, error) {
	hash := sha256.Sum256(data)
	
	r, s, err := ecdsa.Sign(rand.Reader, privateKey, hash[:])
	if err != nil {
		return "", fmt.Errorf("failed to sign data: %v", err)
	}
	
	signature := append(r.Bytes(), s.Bytes()...)
	return hex.EncodeToString(signature), nil
}

func VerifySignature(publicKey *ecdsa.PublicKey, data []byte, signatureHex string) bool {
	signature, err := hex.DecodeString(signatureHex)
	if err != nil {
		return false
	}
	
	if len(signature) != 64 {
		return false
	}
	
	hash := sha256.Sum256(data)
	
	r := new(big.Int).SetBytes(signature[:32])
	s := new(big.Int).SetBytes(signature[32:])
	
	return ecdsa.Verify(publicKey, hash[:], r, s)
}

func SignString(privateKey *ecdsa.PrivateKey, data string) (string, error) {
	return SignData(privateKey, []byte(data))
}

func VerifyStringSignature(publicKey *ecdsa.PublicKey, data string, signatureHex string) bool {
	return VerifySignature(publicKey, []byte(data), signatureHex)
}
