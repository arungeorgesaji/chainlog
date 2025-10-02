package core

import (
	"fmt"
	"time"
)

type Validator struct {
	blockchain *Blockchain
}

func NewValidator(bc *Blockchain) *Validator {
	return &Validator{
		blockchain: bc,
	}
}

func (v *Validator) ValidateBlock(block *Block) bool {
	fmt.Printf("Validating Block %d...\n", block.Index)
	
	if block.Index < 0 {
		fmt.Println("Invalid block index")
		return false
	}
	
	previousBlock := v.blockchain.Chain[block.Index-1]
	
	if block.Index == 0 {
		if block.PrevHash != "" {
			fmt.Println("Genesis block should not have previous hash")
			return false
		}
	} else {
		if block.Index != previousBlock.Index+1 {
			fmt.Printf("Block index out of order. Expected %d, got %d\n", 
				previousBlock.Index+1, block.Index)
			return false
		}
		
		if block.PrevHash != previousBlock.Hash {
			fmt.Println("Block points to wrong previous hash")
			return false
		}
	}
	
	calculatedHash := block.CalculateHash()
	if block.Hash != calculatedHash {
		fmt.Printf("Block hash is invalid. Expected %s, got %s\n", 
			calculatedHash, block.Hash)
		return false
	}
	
	currentTime := time.Now().Unix()
	if block.Timestamp > currentTime+3600 { 
		fmt.Println("Block timestamp is in future")
		return false
	}
	
	if block.Index > 0 && block.Timestamp < previousBlock.Timestamp {
		fmt.Println("Block timestamp is before previous block")
		return false
	}
	
	fmt.Printf("Block %d validation passed!\n", block.Index)
	return true
}

func (v *Validator) ValidateTransaction(tx *Transaction) bool {
	fmt.Printf("Validating Transaction %s...\n", tx.ID[:16])
	
	if tx.Data == "" && tx.Type == DataTx {
		fmt.Println("Data transaction has no data")
		return false
	}
	
	if tx.Fee < 1 || tx.Fee > 5 {
		fmt.Println("Transaction fee must be 1-5 LogCoins")
		return false
	}
	
	if tx.ID != tx.CalculateID() {
		fmt.Println("Transaction ID is invalid")
		return false
	}
	
	currentTime := GetCurrentTimestamp()
	if tx.Timestamp > currentTime+300 { 
		fmt.Println("Transaction timestamp is in future")
		return false
	}
	
	fmt.Printf("Transaction validation passed!\n")
	return true
}

func (v *Validator) ValidateBlockchain() bool {
	fmt.Println("Validating entire blockchain...")
	
	for i := 1; i < len(v.blockchain.Chain); i++ {
		currentBlock := v.blockchain.Chain[i]
		
		if !v.ValidateBlock(currentBlock) {
			fmt.Printf("Blockchain invalid at block %d\n", currentBlock.Index)
			return false
		}
	}
	
	fmt.Println("Entire blockchain validation passed!")
	return true
}

func GetCurrentTimestamp() int64 {
	return time.Now().Unix()
}
