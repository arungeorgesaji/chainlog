package storage

import (
	"encoding/json"
	"os"
	"path/filepath"
)

type NodeState struct {
	IsRunning bool   `json:"is_running"`
	Port      string `json:"port"`
	DataDir   string `json:"data_dir"`
}

func IsNodeRunning() bool {
	path := filepath.Join(DataDir, "node_state.json")
	_, err := os.Stat(path)
	return err == nil
}

func SaveNodeState(port string) error {
	state := &NodeState{
		IsRunning: true,
		Port:      port,
		DataDir:   DataDir,
	}
	
	path := filepath.Join(DataDir, "node_state.json")
	data, err := json.Marshal(state)
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0644)
}

func DeleteNodeState() error {
	path := filepath.Join(DataDir, "node_state.json")
	return os.Remove(path)
}
