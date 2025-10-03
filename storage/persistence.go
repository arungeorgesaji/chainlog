package storage

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
)

const (
	DataDir        = "./chainlog-data"
	BlocksFile     = "blocks.json"
	StateFile      = "state.json"
	WalletsFile    = "wallets.json"
)

func EnsureDataDir() error {
	if err := os.MkdirAll(DataDir, 0755); err != nil {
		return fmt.Errorf("failed to create data directory: %v", err)
	}
	fmt.Printf("Data directory: %s\n", DataDir)
	return nil
}

func SaveToFile(data interface{}, filename string) error {
	filePath := filepath.Join(DataDir, filename)
	
	jsonData, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal data: %v", err)
	}
	
	if err := ioutil.WriteFile(filePath, jsonData, 0644); err != nil {
		return fmt.Errorf("failed to write file: %v", err)
	}
	
	fmt.Printf("Saved data to: %s\n", filePath)
	return nil
}

func LoadFromFile(data interface{}, filename string) error {
	filePath := filepath.Join(DataDir, filename)
	
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return fmt.Errorf("file does not exist: %s", filePath)
	}
	
	jsonData, err := ioutil.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("failed to read file: %v", err)
	}
	
	if err := json.Unmarshal(jsonData, data); err != nil {
		return fmt.Errorf("failed to unmarshal data: %v", err)
	}
	
	fmt.Printf("Loaded data from: %s\n", filePath)
	return nil
}

func FileExists(filename string) bool {
	filePath := filepath.Join(DataDir, filename)
	_, err := os.Stat(filePath)
	return !os.IsNotExist(err)
}

func GetDataSize() (int64, error) {
	var totalSize int64
	
	err := filepath.Walk(DataDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			totalSize += info.Size()
		}
		return nil
	})
	
	if err != nil {
		return 0, err
	}
	
	return totalSize, nil
}
