package network

import (
	"chainlog/core"
	"encoding/json"
	"fmt"
	"net"
)

type MessageType string

const (
	MsgNewBlock    MessageType = "NEW_BLOCK"
	MsgNewTransaction MessageType = "NEW_TRANSACTION"
	MsgGetBlocks   MessageType = "GET_BLOCKS"
	MsgBlocks      MessageType = "BLOCKS"
	MsgGetPeers    MessageType = "GET_PEERS"
	MsgPeers       MessageType = "PEERS"
)

type Message struct {
	Type    MessageType   `json:"type"`
	Data    interface{}   `json:"data"`
	From    string        `json:"from"`
	Version string        `json:"version"`
}

func (n *Node) SendMessage(peerAddress string, msgType MessageType, data interface{}) error {
	conn, err := net.Dial("tcp", peerAddress)
	if err != nil {
		return err
	}
	defer conn.Close()
	
	message := Message{
		Type:    msgType,
		Data:    data,
		From:    n.Address,
		Version: "1.0",
	}
	
	jsonData, err := json.Marshal(message)
	if err != nil {
		return err
	}
	
	_, err = conn.Write(jsonData)
	if err != nil {
		return err
	}
	
	fmt.Printf("Sent %s message to %s\n", msgType, peerAddress)
	return nil
}

func (n *Node) HandleMessage(msg Message, conn net.Conn) {
	fmt.Printf("Received %s message from %s\n", msg.Type, msg.From)
	
	switch msg.Type {
	case MsgNewBlock:
		n.handleNewBlock(msg)
	case MsgNewTransaction:
		n.handleNewTransaction(msg)
	case MsgGetBlocks:
		n.handleGetBlocks(msg, conn)
	case MsgGetPeers:
		n.handleGetPeers(msg, conn)
	default:
		fmt.Printf("Unknown message type: %s\n", msg.Type)
	}
}

func (n *Node) handleNewBlock(msg Message) {
    fmt.Printf("   ↳ New block received from network\n")
    
    blockData, err := json.Marshal(msg.Data)
    if err != nil {
        fmt.Printf("Error marshaling block data: %v\n", err)
        return
    }
    
    var block core.Block
    if err := json.Unmarshal(blockData, &block); err != nil {
        fmt.Printf("Error unmarshaling block: %v\n", err)
        return
    }
    
    validator := core.NewValidator(n.Blockchain)
    if !validator.ValidateBlock(&block) {
        fmt.Printf("Invalid block received: %d\n", block.Index)
        return
    }
    
    currentHeight := n.Blockchain.GetBlockCount()
    if block.Index <= int64(currentHeight-1) {
        fmt.Printf("Already have block %d\n", block.Index)
        return
    }
    
    n.Blockchain.Chain = append(n.Blockchain.Chain, &block)
    fmt.Printf("Added block %d from network\n", block.Index)
    
    n.Blockchain.ClearPendingTransactions()
}

func (n *Node) handleNewTransaction(msg Message) {
    fmt.Printf("   ↳ New transaction received from network\n")
    
    txData, err := json.Marshal(msg.Data)
    if err != nil {
        fmt.Printf("Error marshaling transaction data: %v\n", err)
        return
    }
    
    var tx core.Transaction
    if err := json.Unmarshal(txData, &tx); err != nil {
        fmt.Printf("Error unmarshaling transaction: %v\n", err)
        return
    }
    
    validator := core.NewValidator(n.Blockchain)
    if !validator.ValidateTransaction(&tx) {
        fmt.Printf("Invalid transaction received: %s\n", tx.ID[:16])
        return
    }
    
    for _, pendingTx := range n.Blockchain.PendingTx {
        if pendingTx.ID == tx.ID {
            fmt.Printf("Transaction already in pool: %s\n", tx.ID[:16])
            return
        }
    }
    
    n.Blockchain.PendingTx = append(n.Blockchain.PendingTx, &tx)
    fmt.Printf("Added transaction to pool: %s\n", tx.ID[:16])
}

func (n *Node) handleGetBlocks(msg Message, conn net.Conn) {
    fmt.Printf("   ↳ Sending blockchain to peer\n")
    
    fromHeight := 0
    if data, ok := msg.Data.(map[string]interface{}); ok {
        if fh, ok := data["from_height"].(float64); ok {
            fromHeight = int(fh)
        }
    }
    
    // Send blocks starting from requested height
    var blocksToSend []*core.Block
    for i := fromHeight; i < n.Blockchain.GetBlockCount(); i++ {
        if i < len(n.Blockchain.Chain) {
            blocksToSend = append(blocksToSend, n.Blockchain.Chain[i])
        }
    }
    
    response := Message{
        Type:    MsgBlocks,
        Data:    blocksToSend,
        From:    n.Address,
        Version: "1.0",
    }
    
    jsonData, err := json.Marshal(response)
    if err != nil {
        fmt.Printf("Error marshaling blocks response: %v\n", err)
        return
    }
    
    conn.Write(jsonData)
    fmt.Printf("Sent %d blocks to peer\n", len(blocksToSend))
}

func (n *Node) handleGetPeers(msg Message, conn net.Conn) {
    fmt.Printf("   ↳ Sending peer list to peer\n")
    
		peerAddresses := make([]string, 0, len(n.Peers))
    for addr := range n.Peers {
        peerAddresses = append(peerAddresses, addr)
    }
    
    response := Message{
        Type:    MsgPeers,
        Data:    peerAddresses,
        From:    n.Address,
        Version: "1.0",
    }
    
    jsonData, err := json.Marshal(response)
    if err != nil {
        fmt.Printf("Error marshaling peers response: %v\n", err)
        return
    }
    
    conn.Write(jsonData)
    fmt.Printf("Sent %d peer addresses to peer\n", len(peerAddresses))
}
