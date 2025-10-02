package network

import (
	"fmt"
	"time"
)

func (n *Node) Bootstrap(bootstrapNodes []string) {
	fmt.Printf("Bootstrapping with %d nodes...\n", len(bootstrapNodes))
	
	for _, nodeAddr := range bootstrapNodes {
		n.AddPeer(nodeAddr)
		
		err := n.ConnectToPeer(nodeAddr)
		if err != nil {
			fmt.Printf("Failed to bootstrap with %s: %v\n", nodeAddr, err)
		} else {
			fmt.Printf("Successfully connected to %s\n", nodeAddr)
		}
		
		time.Sleep(100 * time.Millisecond)
	}
}

func (n *Node) DiscoverPeers() {
	ticker := time.NewTicker(30 * time.Second) 
	defer ticker.Stop()
	
	for {
		select {
		case <-ticker.C:
			if n.GetPeerCount() < 10 { 
				n.RequestPeers()
			}
		case <-n.stopChan:
			return
		}
	}
}

func (n *Node) MaintainConnections() {
	ticker := time.NewTicker(60 * time.Second) 
	defer ticker.Stop()
	
	for {
		select {
		case <-ticker.C:
			n.reconnectDisconnectedPeers()
		case <-n.stopChan:
			return
		}
	}
}

func (n *Node) reconnectDisconnectedPeers() {
	n.mutex.Lock()
	defer n.mutex.Unlock()
	
	for addr, peer := range n.Peers {
		if !peer.Connected && time.Since(peer.LastSeen) > 2*time.Minute {
			fmt.Printf("Attempting to reconnect to %s...\n", addr)
			go n.ConnectToPeer(addr)
		}
	}
}
