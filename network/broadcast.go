package network

import (
	"chainlog/core"
	"fmt"
)

func (n *Node) BroadcastTransaction(tx *core.Transaction) {
	fmt.Printf("Broadcasting transaction to %d peers...\n", n.GetPeerCount())
	
	for address, peer := range n.Peers {
		if peer.Connected {
			err := n.SendMessage(address, MsgNewTransaction, tx)
			if err != nil {
				fmt.Printf("Failed to broadcast to %s: %v\n", address, err)
				peer.Connected = false
			}
		}
	}
}

func (n *Node) BroadcastBlock(block *core.Block) {
	fmt.Printf("Broadcasting block %d to %d peers...\n", block.Index, n.GetPeerCount())
	
	for address, peer := range n.Peers {
		if peer.Connected {
			err := n.SendMessage(address, MsgNewBlock, block)
			if err != nil {
				fmt.Printf("Failed to broadcast to %s: %v\n", address, err)
				peer.Connected = false
			}
		}
	}
}

func (n *Node) RequestBlocks() {
	fmt.Printf("Requesting blocks from peers...\n")
	
	for address, peer := range n.Peers {
		if peer.Connected {
			err := n.SendMessage(address, MsgGetBlocks, map[string]interface{}{
				"from_height": n.Blockchain.GetBlockCount(),
			})
			if err != nil {
				fmt.Printf("Failed to request blocks from %s: %v\n", address, err)
			}
		}
	}
}

func (n *Node) RequestPeers() {
	fmt.Printf("Discovering new peers...\n")
	
	for address, peer := range n.Peers {
		if peer.Connected {
			err := n.SendMessage(address, MsgGetPeers, nil)
			if err != nil {
				fmt.Printf("Failed to request peers from %s: %v\n", address, err)
			}
		}
	}
}
