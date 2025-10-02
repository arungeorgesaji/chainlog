package core

import (
	"fmt"
)

type Blockchain struct {
	Chain       []*Block        
	PendingTx   []*Transaction  
	Difficulty  int             
	BlockReward uint64          
}

func NewBlockchain() *Blockchain {
	genesisBlock := createGenesisBlock()
	
	return &Blockchain{
		Chain:       []*Block{genesisBlock},
		PendingTx:   []*Transaction{},
		Difficulty:  2,          
		BlockReward: 10,        
	}
}

func createGenesisBlock() *Block {
	genesisData := "Genesis Block - ChainLog Started!"
	genesisBlock := NewBlock(0, genesisData, "")
	genesisBlock.Miner = "system"
	return genesisBlock
}

func (bc *Blockchain) AddBlock(data string) {
	prevBlock := bc.Chain[len(bc.Chain)-1] 
	newBlock := NewBlock(prevBlock.Index+1, data, prevBlock.Hash)
	
	bc.Chain = append(bc.Chain, newBlock)
	fmt.Printf("Added Block %d to chain\n", newBlock.Index)
}

func (bc *Blockchain) AddTransaction(tx *Transaction) {
	bc.PendingTx = append(bc.PendingTx, tx)
	fmt.Printf("Added transaction to pending pool: %s\n", tx.Data)
}

func (bc *Blockchain) GetLastBlock() *Block {
	return bc.Chain[len(bc.Chain)-1]
}

func (bc *Blockchain) Display() {
	fmt.Printf("\n📚 CHAINLOG BLOCKCHAIN (%d blocks)\n", len(bc.Chain))
	fmt.Println("======================================")
	
	for _, block := range bc.Chain {
		block.Display()
	}
}

func (bc *Blockchain) IsValid() bool {
	for i := 1; i < len(bc.Chain); i++ {
		currentBlock := bc.Chain[i]
		previousBlock := bc.Chain[i-1]
		
		if currentBlock.Hash != currentBlock.CalculateHash() {
			fmt.Printf("Block %d hash is invalid!\n", currentBlock.Index)
			return false
		}
		
		if currentBlock.PrevHash != previousBlock.Hash {
			fmt.Printf("Block %d points to wrong previous hash!\n", currentBlock.Index)
			return false
		}
	}
	
	fmt.Println("Blockchain is valid - no tampering detected!")
	return true
}

func (bc *Blockchain) GetBlockCount() int {
	return len(bc.Chain)
}

func (bc *Blockchain) GetPendingTransactions() []*Transaction {
	return bc.PendingTx
}

func (bc *Blockchain) ClearPendingTransactions() {
	bc.PendingTx = []*Transaction{}
}
