package storage

import (
	"chainlog/core"
	"fmt"
)

type LedgerManager struct {
	Blockchain *core.Blockchain
}

func NewLedgerManager(bc *core.Blockchain) *LedgerManager {
	return &LedgerManager{
		Blockchain: bc,
	}
}

func (lm *LedgerManager) SaveBlockchain() error {
	if err := EnsureDataDir(); err != nil {
		return err
	}
	
	blockchainData := struct {
		Chain       []*core.Block  `json:"chain"`
		PendingTx   []*core.Transaction `json:"pending_transactions"`
		Difficulty  int            `json:"difficulty"`
		BlockReward uint64         `json:"block_reward"`
	}{
		Chain:       lm.Blockchain.Chain,
		PendingTx:   lm.Blockchain.PendingTx,
		Difficulty:  lm.Blockchain.Difficulty,
		BlockReward: lm.Blockchain.BlockReward,
	}
	
	if err := SaveToFile(blockchainData, BlocksFile); err != nil {
		return err
	}
	
	fmt.Printf("Saved blockchain: %d blocks, %d pending transactions\n",
		len(lm.Blockchain.Chain), len(lm.Blockchain.PendingTx))
	
	return nil
}

func (lm *LedgerManager) LoadBlockchain() error {
	if !FileExists(BlocksFile) {
		return fmt.Errorf("blockchain file does not exist, starting fresh")
	}
	
	var blockchainData struct {
		Chain       []*core.Block  `json:"chain"`
		PendingTx   []*core.Transaction `json:"pending_transactions"`
		Difficulty  int            `json:"difficulty"`
		BlockReward uint64         `json:"block_reward"`
	}
	
	if err := LoadFromFile(&blockchainData, BlocksFile); err != nil {
		return err
	}
	
	lm.Blockchain.Chain = blockchainData.Chain
	lm.Blockchain.PendingTx = blockchainData.PendingTx
	lm.Blockchain.Difficulty = blockchainData.Difficulty
	lm.Blockchain.BlockReward = blockchainData.BlockReward
	
	fmt.Printf("Loaded blockchain: %d blocks, %d pending transactions\n",
		len(lm.Blockchain.Chain), len(lm.Blockchain.PendingTx))
	
	return nil
}

func (lm *LedgerManager) SaveBlock(block *core.Block) error {
	return lm.SaveBlockchain()
}

func (lm *LedgerManager) GetBlockchainInfo() (int, int, error) {
	_, err := GetDataSize()
	if err != nil {
		return 0, 0, err
	}
	
	return len(lm.Blockchain.Chain), len(lm.Blockchain.PendingTx), nil
}

func (lm *LedgerManager) DisplayStorageInfo() {
	blockCount, pendingTx, err := lm.GetBlockchainInfo()
	if err != nil {
		fmt.Printf("Storage error: %v\n", err)
		return
	}
	
	dataSize, _ := GetDataSize()
	
	fmt.Printf("\nSTORAGE INFORMATION\n")
	fmt.Printf("├─ Data Directory: %s\n", DataDir)
	fmt.Printf("├─ Blocks: %d\n", blockCount)
	fmt.Printf("├─ Pending Transactions: %d\n", pendingTx)
	fmt.Printf("├─ Data Size: %.2f KB\n", float64(dataSize)/1024)
	fmt.Printf("└─ Persisted: %t\n", FileExists(BlocksFile))
}
