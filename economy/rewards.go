package economy

import (
	"chainlog/core"
	"fmt"
	"time"
)

type RewardManager struct {
	TotalRewardsDistributed uint64
	MinerAddress           string
}

func NewRewardManager(minerAddress string) *RewardManager {
	return &RewardManager{
		MinerAddress: minerAddress,
	}
}

func (rm *RewardManager) CreateBlockReward(blockHeight int64) *core.Transaction {
	blockReward := CalculateBlockReward(blockHeight)
	
	if blockReward == 0 {
		fmt.Println("Block reward has been reduced to 0 (all coins mined)")
		return nil
	}
	
	rewardTx := &core.Transaction{
		ID:        fmt.Sprintf("reward-%d-%d", blockHeight, time.Now().UnixNano()),
		Type:      core.RewardTx,
		Sender:    "NETWORK", 
		Receiver:  rm.MinerAddress,
		Amount:    blockReward,
		Fee:       0, 
		Timestamp: time.Now().Unix(),
		Nonce:     uint64(time.Now().UnixNano()),
	}
	
	rm.TotalRewardsDistributed += blockReward
	
	fmt.Printf("Block reward: %d LogCoins → %s\n", 
		blockReward, rm.MinerAddress[:8])
	
	return rewardTx
}

func (rm *RewardManager) CreateStakingReward(validatorAddress string, stakedAmount uint64) *core.Transaction {
	stakingReward := stakedAmount / 100 / 20 
	
	if stakingReward == 0 {
		return nil 
	}
	
	rewardTx := &core.Transaction{
		ID:        fmt.Sprintf("stake-reward-%d", time.Now().UnixNano()),
		Type:      core.RewardTx,
		Sender:    "NETWORK",
		Receiver:  validatorAddress,
		Amount:    stakingReward,
		Fee:       0,
		Timestamp: time.Now().Unix(),
		Nonce:     uint64(time.Now().UnixNano()),
	}
	
	fmt.Printf("Staking reward: %d LogCoins → %s\n", 
		stakingReward, validatorAddress[:8])
	
	return rewardTx
}

func (rm *RewardManager) GetRewardStatistics() uint64 {
	return rm.TotalRewardsDistributed
}

func (rm *RewardManager) DisplayRewardStats() {
	fmt.Printf("\nREWARD STATISTICS\n")
	fmt.Printf("├─ Total Rewards Distributed: %d LogCoins\n", rm.TotalRewardsDistributed)
	fmt.Printf("├─ Current Block Reward: %d LogCoins\n", CalculateBlockReward(0)) 
	fmt.Printf("└─ Miner: %s\n", rm.MinerAddress[:8])
}
