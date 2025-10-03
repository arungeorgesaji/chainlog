package network

import (
	"chainlog/core"
	"chainlog/network"
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
}

func (n *Node) handleNewTransaction(msg Message) {
	fmt.Printf("   ↳ New transaction received from network\n")
}

func (n *Node) handleGetBlocks(msg Message, conn net.Conn) {
	fmt.Printf("   ↳ Sending blockchain to peer\n")
}

func (n *Node) handleGetPeers(msg Message, conn net.Conn) {
	fmt.Printf("   ↳ Sending peer list to peer\n")
}
