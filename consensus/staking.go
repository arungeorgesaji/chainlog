package consensus

import (
	"fmt"
	"sync"
)

type Validator struct {
	Address    string
	Staked     uint64  
	VotingPower uint64 
	Active     bool
}

type StakingManager struct {
	Validators map[string]*Validator
	MinStake   uint64
	mutex      sync.Mutex
}

func NewStakingManager() *StakingManager {
	return &StakingManager{
		Validators: make(map[string]*Validator),
		MinStake:   100, 
	}
}

func (sm *StakingManager) AddStake(address string, amount uint64) error {
	sm.mutex.Lock()
	defer sm.mutex.Unlock()
	
	if amount < sm.MinStake {
		return fmt.Errorf("need at least %d LogCoins to stake", sm.MinStake)
	}
	
	if existing, exists := sm.Validators[address]; exists {
		existing.Staked += amount
		existing.VotingPower = existing.Staked
		fmt.Printf("Validator %s increased stake to %d\n", address[:8], existing.Staked)
	} else {
		sm.Validators[address] = &Validator{
			Address:    address,
			Staked:     amount,
			VotingPower: amount,
			Active:     true,
		}
		fmt.Printf("New validator: %s with %d LogCoins\n", address[:8], amount)
	}
	
	return nil
}

func (sm *StakingManager) GetTotalStaked() uint64 {
	total := uint64(0)
	for _, validator := range sm.Validators {
		total += validator.Staked
	}
	return total
}

func (sm *StakingManager) DisplayValidators() {
	fmt.Printf("\nVALIDATORS (%d total)\n", len(sm.Validators))
	fmt.Println("========================")
	
	if len(sm.Validators) == 0 {
		fmt.Println("No active validators")
		return
	}
	
	for _, validator := range sm.Validators {
		status := Active
		if !validator.Active {
			status = "Inactive"
		}
		fmt.Printf("%s %s: %d LogCoins (%d votes)\n",
			status, validator.Address[:8], validator.Staked, validator.VotingPower)
	}
	fmt.Printf("Total staked: %d LogCoins\n", sm.GetTotalStaked())
}
