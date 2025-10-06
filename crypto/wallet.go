package crypto

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"encoding/hex"
	"fmt"
	"math/big"
)

type Wallet struct {
	PrivateKey *ecdsa.PrivateKey
	PublicKey  *ecdsa.PublicKey
	Address    string
}

func NewWallet() (*Wallet, error) {
	privateKey, err := GenerateKeyPair()
	if err != nil {
		return nil, err
	}
	
	wallet := &Wallet{
		PrivateKey: privateKey,
		PublicKey:  &privateKey.PublicKey,
		Address:    PublicKeyToAddress(&privateKey.PublicKey),
	}
	
	return wallet, nil
}

func WalletFromPrivateKey(privateKeyHex string) (*Wallet, error) {
	if len(privateKeyHex) != 64 {
		return nil, fmt.Errorf("invalid private key length: expected 64 hex characters, got %d", len(privateKeyHex))
	}
	
	keyBytes, err := hex.DecodeString(privateKeyHex)
	if err != nil {
		return nil, fmt.Errorf("invalid private key format: not valid hexadecimal")
	}
	
	if isAllZeros(privateKeyHex) {
		return nil, fmt.Errorf("invalid private key: cannot be all zeros")
	}
	
	keyInt := new(big.Int).SetBytes(keyBytes)
	
	if !isValidPrivateKeyRange(keyInt) {
		return nil, fmt.Errorf("invalid private key: out of valid range for secp256k1 curve")
	}
	
	privateKey, err := StringToPrivateKey(privateKeyHex)
	if err != nil {
		return nil, err
	}
	
	wallet := &Wallet{
		PrivateKey: privateKey,
		PublicKey:  &privateKey.PublicKey,
		Address:    PublicKeyToAddress(&privateKey.PublicKey),
	}
	
	return wallet, nil
}

func isAllZeros(hexString string) bool {
	for _, char := range hexString {
		if char != '0' {
			return false
		}
	}
	return true
}

func isValidPrivateKeyRange(keyInt *big.Int) bool {
	curve := elliptic.P256() 
	
	nMinusOne := new(big.Int).Sub(curve.Params().N, big.NewInt(1))
	
	return keyInt.Cmp(big.NewInt(1)) >= 0 && keyInt.Cmp(nMinusOne) <= 0
}

func (w *Wallet) GetAddress() string {
	return w.Address
}

func (w *Wallet) GetAddressShort() string {
	if len(w.Address) >= 8 {
		return w.Address[:8] + "..."
	}
	return w.Address
}

func PublicKeyToString(pubKey *ecdsa.PublicKey) string {
	if pubKey == nil || pubKey.X == nil || pubKey.Y == nil {
		return ""
	}
	return fmt.Sprintf("%x", pubKey.X.Bytes()) + fmt.Sprintf("%x", pubKey.Y.Bytes())
}

func (w *Wallet) Display() {
	fmt.Printf("WALLET INFORMATION\n")
	fmt.Printf("├─ Address: %s\n", w.Address)
	fmt.Printf("├─ Short: %s\n", w.GetAddressShort())
	fmt.Printf("└─ Private Key: %s (SAVE THIS SECURELY!)\n", PrivateKeyToString(w.PrivateKey))
}
