package crypto

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"math/big"
)

func GenerateKeyPair() (*ecdsa.PrivateKey, error) {
	privateKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		return nil, fmt.Errorf("failed to generate key pair: %v", err)
	}
	return privateKey, nil
}

func PublicKeyToAddress(publicKey *ecdsa.PublicKey) string {
	publicKeyBytes := append(
		publicKey.X.Bytes(),
		publicKey.Y.Bytes()...,
	)
	
	hash := sha256.Sum256(publicKeyBytes)
	
	addressBytes := hash[:20]
	
	return hex.EncodeToString(addressBytes)
}

func PrivateKeyToString(privateKey *ecdsa.PrivateKey) string {
	return hex.EncodeToString(privateKey.D.Bytes())
}

func StringToPrivateKey(privateKeyHex string) (*ecdsa.PrivateKey, error) {
	bytes, err := hex.DecodeString(privateKeyHex)
	if err != nil {
		return nil, err
	}
	
	privateKey := new(ecdsa.PrivateKey)
	privateKey.PublicKey.Curve = elliptic.P256()
	privateKey.D = new(big.Int).SetBytes(bytes)
	privateKey.PublicKey.X, privateKey.PublicKey.Y = privateKey.PublicKey.Curve.ScalarBaseMult(bytes)
	
	return privateKey, nil
}
