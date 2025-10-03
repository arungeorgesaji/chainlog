package storage

import (
	"chainlog/crypto"
	"fmt"
)

type AccountState struct {
	Address string `json:"address"`
	Balance uint64 `json:"balance"`
	Nonce   uint64 `json:"nonce"`
}

type StateManager struct {
	Accounts map[string]*AccountState `json:"accounts"`
}

func NewStateManager() *StateManager {
	return &StateManager{
		Accounts: make(map[string]*AccountState),
	}
}

func (sm *StateManager) SaveState() error {
	if err := EnsureDataDir(); err != nil {
		return err
	}
	
	return SaveToFile(sm, StateFile)
}

func (sm *StateManager) LoadState() error {
	if !FileExists(StateFile) {
		fmt.Println("No saved state found, starting fresh")
		return nil
	}
	
	return LoadFromFile(sm, StateFile)
}

func (sm *StateManager) UpdateAccount(address string, balance uint64, nonce uint64) {
	sm.Accounts[address] = &AccountState{
		Address: address,
		Balance: balance,
		Nonce:   nonce,
	}
}

func (sm *StateManager) GetAccount(address string) (*AccountState, bool) {
	account, exists := sm.Accounts[address]
	return account, exists
}

func (sm *StateManager) GetBalance(address string) uint64 {
	if account, exists := sm.Accounts[address]; exists {
		return account.Balance
	}
	return 0
}

func (sm *StateManager) InitializeGenesisState(genesisWallets []*crypto.Wallet) {
	fmt.Println("Initializing genesis state...")
	
	for _, wallet := range genesisWallets {
		sm.UpdateAccount(wallet.GetAddress(), 1000, 0)
		fmt.Printf("   %s: 1000 LogCoins (genesis)\n", wallet.GetAddressShort())
	}
}

func (sm *StateManager) DisplayState() {
	fmt.Printf("\nACCOUNT STATES (%d accounts)\n", len(sm.Accounts))
	fmt.Println("=============================")
	
	if len(sm.Accounts) == 0 {
		fmt.Println("No accounts found")
		return
	}
	
	totalBalance := uint64(0)
	for _, account := range sm.Accounts {
		fmt.Printf("%s: %d LogCoins (nonce: %d)\n",
			account.Address[:8], account.Balance, account.Nonce)
		totalBalance += account.Balance
	}
	
	fmt.Printf("Total circulating: %d LogCoins\n", totalBalance)
}
