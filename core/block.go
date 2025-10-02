package core

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"time"
)

type Block struct {
	Index        int64         
	Timestamp    int64         
	Data         string        
	PrevHash     string        
	Hash         string        
	Nonce        int64         
	Difficulty   int           
	Miner        string        
}

func NewBlock(index int64, data string, prevHash string) *Block {
	block := &Block{
		Index:     index,
		Timestamp: time.Now().Unix(),  
		Data:      data,
		PrevHash:  prevHash,
		Nonce:     0,                  
		Difficulty: 2,                 
		Miner: "unknown",
	}
	
	block.Hash = block.CalculateHash()
	return block
}

func (b *Block) CalculateHash() string {
	blockData := fmt.Sprintf("%d%d%s%s%d", 
		b.Index, 
		b.Timestamp, 
		b.Data, 
		b.PrevHash, 
		b.Nonce)
	
	hash := sha256.Sum256([]byte(blockData))
	
	return hex.EncodeToString(hash[:])
}

func (b *Block) Display() {
	fmt.Printf("┌─── BLOCK %d ───\n", b.Index)
	fmt.Printf("│ Timestamp: %s\n", time.Unix(b.Timestamp, 0).Format("2006-01-02 15:04:05"))
	fmt.Printf("│ Data: %s\n", b.Data)
	
	if b.PrevHash == "" {
		fmt.Printf("│ Previous: [GENESIS]\n")
	} else {
		fmt.Printf("│ Previous: %s...\n", b.PrevHash[:16])
	}
	
	if b.Hash == "" {
		fmt.Printf("│ Hash: [NOT CALCULATED]\n")
	} else {
		fmt.Printf("│ Hash: %s...\n", b.Hash[:16])
	}
	
	fmt.Printf("│ Nonce: %d\n", b.Nonce)
	
	if b.Miner == "" {
		fmt.Printf("│ Miner: [UNKNOWN]\n")
	} else {
		fmt.Printf("│ Miner: %s\n", b.Miner)
	}
	
	fmt.Printf("└%s┘\n", "─────────────────")
}
