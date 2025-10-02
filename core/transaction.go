package core

import (
	"chainlog/crypto"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"time"
)

type TransactionType int

const (
	DataTx   TransactionType = iota  
	FeeTx                            
	RewardTx                         
	StakeTx                          
)

type Transaction struct {
	ID        string          
	Type      TransactionType 
	Data      string          
	Sender    string          
	Receiver  string          
	Amount    uint64          
	Fee       uint64          
	Signature string          
	Timestamp int64           
	Nonce     uint64          
}

func NewDataTransaction(data string, wallet *crypto.Wallet, fee uint64) (*Transaction, error) {
	tx := &Transaction{
		Type:      DataTx,
		Data:      data,
		Sender:    wallet.GetAddress(),
		Fee:       fee,
		Timestamp: time.Now().Unix(),
		Nonce:     generateNonce(),
	}
	
	tx.ID = tx.CalculateID()
	
	signature, err := crypto.SignString(wallet.PrivateKey, tx.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to sign transaction: %v", err)
	}
	tx.Signature = signature
	
	return tx, nil
}

func (tx *Transaction) VerifySignature() bool {
	return tx.Signature != "" && len(tx.Signature) == 128 
}

func (tx *Transaction) CalculateID() string {
	txData := fmt.Sprintf("%s%s%s%d%d%d%d",
		tx.Data,
		tx.Sender,
		tx.Receiver,
		tx.Amount,
		tx.Fee,
		tx.Timestamp,
		tx.Nonce)
	
	hash := sha256.Sum256([]byte(txData))
	return hex.EncodeToString(hash[:])
}

func generateNonce() uint64 {
	return uint64(time.Now().UnixNano())
}

func (tx *Transaction) Display() {
	typeNames := []string{"DATA", "FEE", "REWARD", "STAKE"}
	
	displayID := tx.ID
	if len(displayID) >= 16 {
		displayID = displayID[:16] + "..."
	}
	
	fmt.Printf("╔═ TRANSACTION %s\n", displayID)
	fmt.Printf("║ Type: %s\n", typeNames[tx.Type])
	
	if len(tx.Sender) >= 8 {
		fmt.Printf("║ From: %s...\n", tx.Sender[:8])
	} else {
		fmt.Printf("║ From: %s\n", tx.Sender)
	}
	
	if tx.Receiver != "" {
		if len(tx.Receiver) >= 8 {
			fmt.Printf("║ To: %s...\n", tx.Receiver[:8])
		} else {
			fmt.Printf("║ To: %s\n", tx.Receiver)
		}
	}
	
	fmt.Printf("║ Data: %s\n", tx.Data)
	if tx.Amount > 0 {
		fmt.Printf("║ Amount: %d LogCoins\n", tx.Amount)
	}
	if tx.Fee > 0 {
		fmt.Printf("║ Fee: %d LogCoins\n", tx.Fee)
	}
	fmt.Printf("║ Time: %s\n", time.Unix(tx.Timestamp, 0).Format("15:04:05"))
	fmt.Printf("╚%s╝\n", "══════════════════")
}
