package network

import (
	"encoding/json"
	"os"
	"sync"
	"time"
)

type NodeInfo struct {
	Address   string `json:"address"`
	Port      string `json:"port"`
	Wallet    string `json:"wallet"`
	StartedAt int64  `json:"started_at"`
	LastSeen  int64  `json:"last_seen"`
}

type NodeRegistry struct {
	Nodes map[string]*NodeInfo `json:"nodes"` 
	mu    sync.RWMutex
}

var (
	registry     *NodeRegistry
	registryOnce sync.Once
)

func GetNodeRegistry() *NodeRegistry {
	registryOnce.Do(func() {
		registry = &NodeRegistry{
			Nodes: make(map[string]*NodeInfo),
		}
		registry.load()
	})
	return registry
}

func (nr *NodeRegistry) RegisterNode(address, port, wallet string) {
	nr.mu.Lock()
	defer nr.mu.Unlock()
	
	key := address + ":" + port
	nr.Nodes[key] = &NodeInfo{
		Address:   address,
		Port:      port,
		Wallet:    wallet,
		StartedAt: time.Now().Unix(),
		LastSeen:  time.Now().Unix(),
	}
	nr.save()
}

func (nr *NodeRegistry) UpdateLastSeen(address, port string) {
	nr.mu.Lock()
	defer nr.mu.Unlock()
	
	key := address + ":" + port
	if node, exists := nr.Nodes[key]; exists {
		node.LastSeen = time.Now().Unix()
		nr.save()
	}
}

func (nr *NodeRegistry) GetRunningNodes() []*NodeInfo {
	nr.mu.RLock()
	defer nr.mu.RUnlock()
	
	var running []*NodeInfo
	now := time.Now().Unix()
	
	for _, node := range nr.Nodes {
		if now-node.LastSeen < 120 {
			running = append(running, node)
		}
	}
	return running
}

func (nr *NodeRegistry) save() {
	data, _ := json.Marshal(nr.Nodes)
	os.WriteFile("chainlog-data/nodes.json", data, 0644)
}

func (nr *NodeRegistry) load() {
	data, err := os.ReadFile("chainlog-data/nodes.json")
	if err != nil {
		return
	}
	json.Unmarshal(data, &nr.Nodes)
}
