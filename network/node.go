package network

import (
	"chainlog/core"
	"chainlog/crypto"
	"encoding/json"
	"fmt"
	"net"
	"sync"
	"time"
)

type Node struct {
	ID           string
	Address      string          
	Peers        map[string]*Peer
	Blockchain   *core.Blockchain
	Wallet       *crypto.Wallet
	IsMiner      bool
	peerFile   string
	Server       net.Listener
	mutex        sync.Mutex      
	stopChan     chan bool
}

type Peer struct {
	ID        string
	Address   string
	Connected bool
	Conn      net.Conn           
	Encoder   *json.Encoder        
	Decoder   *json.Decoder      
	LastSeen  time.Time
}

func NewNode(address string, wallet *crypto.Wallet, bc *core.Blockchain, isMiner bool) *Node {
	return &Node{
		ID:         wallet.GetAddressShort(), 
		Address:    address,
		Peers:      make(map[string]*Peer),
		Blockchain: bc,
		Wallet:     wallet,
		IsMiner:    isMiner,
		stopChan:   make(chan bool),
	}
}

func (n *Node) Start() error {
	server, err := net.Listen("tcp", n.Address)
	if err != nil {
		return fmt.Errorf("failed to start node: %v", err)
	}
	n.Server = server

	fmt.Printf(" Node %s started on %s\n", n.ID, n.Address)
	
	go n.acceptConnections()
	
	return nil
}

func (n *Node) Stop() {
	close(n.stopChan)
	if n.Server != nil {
		n.Server.Close()
	}
	fmt.Printf("Node %s stopped\n", n.ID)
}

func (n *Node) acceptConnections() {
	for {
		conn, err := n.Server.Accept()
		if err != nil {
			select {
			case <-n.stopChan:
				return 
			default:
				fmt.Printf("Connection error: %v\n", err)
				continue
			}
		}
		
		go n.handleConnection(conn)
	}
}

func (n *Node) handleConnection(conn net.Conn) {
	defer conn.Close()
	
	conn.Write([]byte("Hello from ChainLog node " + n.ID + "\n"))
	fmt.Printf("New connection from %s\n", conn.RemoteAddr().String())

	buffer := make([]byte, 1024*1024)

	for {
			nBytes, err := conn.Read(buffer)
			if err != nil {
					fmt.Printf("Connection closed: %v\n", err)
					return
			}
			
			var msg Message
			if err := json.Unmarshal(buffer[:nBytes], &msg); err != nil {
					fmt.Printf("Error parsing message: %v\n", err)
					continue
			}
			
			n.HandleMessage(msg, conn)
	}
}

func (n *Node) AddPeer(address string) {
	n.mutex.Lock()
	defer n.mutex.Unlock()
	
	if _, exists := n.Peers[address]; !exists {
		n.Peers[address] = &Peer{
			ID:        "peer-" + address,
			Address:   address,
			Connected: false,
			LastSeen:  time.Now(),
		}
		fmt.Printf("Added peer: %s\n", address)
	}
}

func (n *Node) ConnectToPeer(address string) error {
	conn, err := net.Dial("tcp", address)
	if err != nil {
			return err
	}
	
	peer := &Peer{
			ID:        "peer-" + address,
			Address:   address,
			Connected: true,
			Conn:      conn,
			Encoder:   json.NewEncoder(conn),
			Decoder:   json.NewDecoder(conn),
			LastSeen:  time.Now(),
	}
	
	n.Peers[address] = peer
	
	go n.listenToPeer(peer)
	
	return nil
}

func (n *Node) listenToPeer(peer *Peer) {
	for {
		var msg Message
		if err := peer.Decoder.Decode(&msg); err != nil {
				peer.Connected = false
				peer.Conn.Close()
				return
		}
		n.HandleMessage(msg, peer.Conn)
	}
}

func (n *Node) GetPeerCount() int {
	n.mutex.Lock()
	defer n.mutex.Unlock()
	return len(n.Peers)
}

func (n *Node) Display() {
	fmt.Printf("\nNODE INFORMATION\n")
	fmt.Printf("├─ ID: %s\n", n.ID)
	fmt.Printf("├─ Address: %s\n", n.Address)
	fmt.Printf("├─ Wallet: %s\n", n.Wallet.GetAddressShort())
	fmt.Printf("├─ Miner: %t\n", n.IsMiner)
	fmt.Printf("├─ Blockchain Height: %d\n", n.Blockchain.GetBlockCount())
	fmt.Printf("└─ Peers: %d\n", n.GetPeerCount())
	
	if n.GetPeerCount() > 0 {
		fmt.Printf("   Connected Peers:\n")
		n.mutex.Lock()
		for addr, peer := range n.Peers {
			status := "Offline"
			if peer.Connected {
				status = "Online"
			}
			fmt.Printf("   %s %s (last seen: %v ago)\n", 
				status, addr, time.Since(peer.LastSeen).Round(time.Second))
		}
		n.mutex.Unlock()
	}
}
