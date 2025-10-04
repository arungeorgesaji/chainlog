package crypto

import (
	"crypto/ecdsa"
	"fmt"
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
	fmt.Printf("└─ Private Key: %s... (keep secret!)\n", PrivateKeyToString(w.PrivateKey)[:16])
}
