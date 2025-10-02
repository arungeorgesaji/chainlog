package consensus

import (
	"chainlog/core"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"math/big"
	"strings"
)

type ProofOfWork struct {
	Block      *core.Block
	Target     *big.Int
	Difficulty int
}

func NewProofOfWork(block *core.Block, difficulty int) *ProofOfWork {
	target := big.NewInt(1)
	target.Lsh(target, uint(256-difficulty)) 
	
	return &ProofOfWork{
		Block:      block,
		Target:     target,
		Difficulty: difficulty,
	}
}

func (pow *ProofOfWork) Run() (int64, string) {
	var nonce int64 = 0
	var hash string
	
	fmt.Printf("Mining block %d (difficulty: %d)...\n", pow.Block.Index, pow.Difficulty)
	
	for nonce < maxNonce {
		pow.Block.Nonce = nonce
		hash = pow.Block.CalculateHash()
		
		if pow.ValidateHash(hash) {
			fmt.Printf("Block mined! Nonce: %d, Hash: %s\n", nonce, hash[:16])
			return nonce, hash
		}
		
		nonce++
		
		if nonce%100000 == 0 {
			fmt.Printf("  ...tried %d nonces...\n", nonce)
		}
	}
	
	return -1, "" 
}

func (pow *ProofOfWork) ValidateHash(hash string) bool {
	hashBytes, _ := hex.DecodeString(hash)
	hashInt := new(big.Int).SetBytes(hashBytes)
	
	return hashInt.Cmp(pow.Target) == -1 
}

func (pow *ProofOfWork) Validate() bool {
	hash := pow.Block.CalculateHash()
	return pow.ValidateHash(hash)
}

const maxNonce = 10000000
