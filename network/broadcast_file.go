package network

import (
	"chainlog/core"
	"fmt"
	"os"
	"strings"
	"time"
)

func (n *Node) CheckBroadcastFile() {
    ticker := time.NewTicker(10 * time.Second) 
    defer ticker.Stop()
    
    for {
        select {
        case <-ticker.C:
            n.processBroadcastFile()
        case <-n.stopChan:
            return
        }
    }
}

func (n *Node) processBroadcastFile() {
    filePath := "pending_broadcasts.txt"
    data, err := os.ReadFile(filePath)
    if err != nil {
        return 
    }
    
    lines := strings.Split(string(data), "\n")
    for _, txID := range lines {
        txID = strings.TrimSpace(txID)
        if txID == "" {
            continue
        }
        
        if tx := n.findTransaction(txID); tx != nil {
            n.BroadcastTransaction(tx)
            fmt.Printf("Broadcast transaction: %s...\n", txID[:16])
        }
    }
    
    os.Remove(filePath)
}

func (n *Node) findTransaction(txID string) *core.Transaction {
    for _, tx := range n.Blockchain.GetPendingTransactions() {
        if strings.HasPrefix(tx.ID, txID) {
            return tx
        }
    }
    return nil
}
