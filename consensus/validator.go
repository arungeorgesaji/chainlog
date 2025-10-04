package consensus

import (
	"chainlog/core"
	"math/rand"
	"time"
)

func (sm *StakingManager) SelectValidator() *Validator {
    sm.mutex.Lock()
    defer sm.mutex.Unlock()
    
    if len(sm.Validators) == 0 {
        return nil
    }
    
    var validators []*Validator
    var weights []uint64
    
    for _, validator := range sm.Validators {
        if validator.Active {
            validators = append(validators, validator)
            weights = append(weights, validator.VotingPower)
        }
    }
    
    if len(validators) == 0 {
        return nil
    }
    
    rand.Seed(time.Now().UnixNano())
    totalWeight := uint64(0)
    for _, w := range weights {
        totalWeight += w
    }
    
    r := rand.Uint64() % totalWeight
    for i, validator := range validators {
        if r < weights[i] {
            return validator
        }
        r -= weights[i]
    }
    
    return validators[0] 
}

func (sm *StakingManager) ValidateBlock(block *core.Block) bool {
    if block.Miner == "unknown" || block.Miner == "system" {
        return true 
    }
    
    validator, exists := sm.Validators[block.Miner]
    if !exists || !validator.Active {
        return false
    }
    
    return true
}
