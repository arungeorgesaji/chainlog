package storage

import (
	"chainlog/crypto"
	"fmt"
	"sync"
	"time"
)

type StoredWallet struct {
	Address    string `json:"address"`
	PrivateKey string `json:"privateKey"` 
	PublicKey  string `json:"publicKey"`
	Label      string `json:"label"`
	CreatedAt  int64  `json:"createdAt"`
}

type WalletManager struct {
	Wallets map[string]*StoredWallet `json:"wallets"` 
	mu      sync.RWMutex
}

var (
	walletManager *WalletManager
	managerOnce   sync.Once
)

func GetWalletManager() *WalletManager {
	managerOnce.Do(func() {
		walletManager = &WalletManager{
			Wallets: make(map[string]*StoredWallet),
		}
		walletManager.loadWallets()
	})
	return walletManager
}

func (wm *WalletManager) SaveWallet(wallet *crypto.Wallet, label string) error { 
	wm.mu.Lock()
	defer wm.mu.Unlock()

	stored := &StoredWallet{
		Address:    wallet.Address,
		PrivateKey: crypto.PrivateKeyToString(wallet.PrivateKey), 
		PublicKey:  crypto.PublicKeyToString(wallet.PublicKey),   
		Label:      label,
		CreatedAt:  time.Now().Unix(),
	}

	wm.Wallets[wallet.Address] = stored
	return wm.saveToFile()
}

func (wm *WalletManager) GetWallet(address string) (*StoredWallet, bool) {
	wm.mu.RLock()
	defer wm.mu.RUnlock()

	wallet, exists := wm.Wallets[address]
	return wallet, exists
}

func (wm *WalletManager) GetAllWallets() []*StoredWallet {
	wm.mu.RLock()
	defer wm.mu.RUnlock()

	wallets := make([]*StoredWallet, 0, len(wm.Wallets))
	for _, wallet := range wm.Wallets {
		wallets = append(wallets, wallet)
	}
	return wallets
}

func (wm *WalletManager) DeleteWallet(address string) bool {
	wm.mu.Lock()
	defer wm.mu.Unlock()

	if _, exists := wm.Wallets[address]; exists {
		delete(wm.Wallets, address)
		wm.saveToFile()
		return true
	}
	return false
}

func (wm *WalletManager) loadWallets() error {
	if err := EnsureDataDir(); err != nil { 
		return err
	}

	if !FileExists(WalletsFile) { 
		return nil 
	}

	var wallets map[string]*StoredWallet
	if err := LoadFromFile(&wallets, WalletsFile); err != nil { 
		return err
	}

	wm.Wallets = wallets
	return nil
}

func (wm *WalletManager) saveToFile() error {
	return SaveToFile(wm.Wallets, WalletsFile) 
}

func (wm *WalletManager) WalletCount() int {
	wm.mu.RLock()
	defer wm.mu.RUnlock()
	return len(wm.Wallets)
}

func LoadWalletFromStorage(address string) (*crypto.Wallet, error) { 
	wm := GetWalletManager()
	stored, exists := wm.GetWallet(address)
	if !exists {
		return nil, fmt.Errorf("wallet not found in storage: %s", address)
	}

	return crypto.WalletFromPrivateKey(stored.PrivateKey)
}

func GetDefaultWallet() (*crypto.Wallet, error) { 
	wm := GetWalletManager()
	wallets := wm.GetAllWallets()
	if len(wallets) == 0 {
		return nil, fmt.Errorf("no wallets found in storage")
	}

	return LoadWalletFromStorage(wallets[0].Address)
}
