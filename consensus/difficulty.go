package consensus

import (
	"chainlog/core"
	"fmt"
	"time"
)

type DifficultyManager struct {
	Blockchain        *core.Blockchain
	TargetBlockTime   time.Duration 
	AdjustmentBlocks  int           
}

func NewDifficultyManager(bc *core.Blockchain) *DifficultyManager {
	return &DifficultyManager{
		Blockchain:       bc,
		TargetBlockTime:  10 * time.Second, 
		AdjustmentBlocks: 10,               
	}
}

func (dm *DifficultyManager) CalculateNewDifficulty() int {
	if dm.Blockchain.GetBlockCount() < 2 {
		return initialDifficulty
	}
	
	if dm.Blockchain.GetBlockCount()%dm.AdjustmentBlocks != 0 {
		return dm.Blockchain.Difficulty
	}
	
	recentBlocks := 5
	if dm.Blockchain.GetBlockCount() < recentBlocks {
		recentBlocks = dm.Blockchain.GetBlockCount() - 1
	}
	
	var totalTime int64
	for i := 0; i < recentBlocks; i++ {
		block := dm.Blockchain.Chain[dm.Blockchain.GetBlockCount()-1-i]
		prevBlock := dm.Blockchain.Chain[dm.Blockchain.GetBlockCount()-2-i]
		blockTime := block.Timestamp - prevBlock.Timestamp
		totalTime += blockTime
	}
	
	averageBlockTime := totalTime / int64(recentBlocks)
	currentDifficulty := dm.Blockchain.Difficulty
	
	if averageBlockTime < int64(dm.TargetBlockTime.Seconds())/2 {
		return currentDifficulty + 1
	} else if averageBlockTime > int64(dm.TargetBlockTime.Seconds())*2 {
		if currentDifficulty > 1 {
			return currentDifficulty - 1
		}
	}
	
	return currentDifficulty 
}

func (dm *DifficultyManager) UpdateDifficulty() {
	newDifficulty := dm.CalculateNewDifficulty()
	if newDifficulty != dm.Blockchain.Difficulty {
		oldDiff := dm.Blockchain.Difficulty
		dm.Blockchain.Difficulty = newDifficulty
		fmt.Printf("Difficulty adjusted: %d → %d\n", oldDiff, newDifficulty)
	}
}

const initialDifficulty = 2
