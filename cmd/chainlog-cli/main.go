package main

import (
	"chainlog/core"
	"chainlog/crypto"  
	"chainlog/network"
	"fmt"
)

func main() {
	fmt.Println("Starting ChainLog Cryptography Test...")
	
	fmt.Println("\n1. Creating wallets...")
	wallet1, err := crypto.NewWallet()
	if err != nil {
		panic(err)
	}
	
	wallet2, err := crypto.NewWallet()
	if err != nil {
		panic(err)
	}
	
	wallet1.Display()
	wallet2.Display()
	
	fmt.Println("\n2. Creating blockchain...")
	bc := core.NewBlockchain()
	
	bc.AddBlock("First research finding: Gravity exists!")
	bc.AddBlock("Second finding: Water is wet")
	bc.AddBlock("System error: Database connection failed")
	
	fmt.Println("\n3. Creating signed transactions...")
	tx1, err := core.NewDataTransaction("User logged in", wallet1, 2)
	if err != nil {
		panic(err)
	}
	
	tx2, err := core.NewDataTransaction("Payment processed", wallet2, 3)
	if err != nil {
		panic(err)
	}
	
	bc.AddTransaction(tx1)
	bc.AddTransaction(tx2)
	
	fmt.Println("\n4. Starting Network Node...")
	node := network.NewNode("localhost:8080", wallet1, bc, true)
	err = node.Start()
	if err != nil {
		panic(err)
	}
	defer node.Stop() 
	
	node.AddPeer("localhost:8081")
	node.AddPeer("localhost:8082")
	node.AddPeer("localhost:8083")
	
	node.Display()
	
	fmt.Println("\n5. Broadcasting transactions to network...")
	node.BroadcastTransaction(tx1)
	node.BroadcastTransaction(tx2)
	
	fmt.Println("\n6. Blockchain state:")
	bc.Display()
	
	fmt.Println("\n7. Validating blockchain...")
	validator := core.NewValidator(bc)
	validator.ValidateBlockchain()
	
	fmt.Printf("\n8. Pending Transactions: %d\n", len(bc.GetPendingTransactions()))
	for i, tx := range bc.GetPendingTransactions() {
		fmt.Printf("\nTransaction %d:\n", i+1)
		tx.Display()
	}
	
	fmt.Println("\n9. Network Activity Demo...")
	
	fmt.Println("   Creating and broadcasting new network transaction...")
	tx3, err := core.NewDataTransaction("Network broadcast test!", wallet1, 1)
	if err != nil {
		panic(err)
	}
	bc.AddTransaction(tx3)
	node.BroadcastTransaction(tx3)
	
	fmt.Println("\n10. Final Network Status:")
	node.Display()
}
