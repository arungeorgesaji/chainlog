package economy

import (
	"fmt"
)

type CoinEconomics struct {
	TotalSupply     uint64
	Circulating     uint64
	Burned          uint64
	BlockReward     uint64
	NextHalving     int64
}

const (
	GenesisSupply    uint64 = 1000000   
	MaxSupply        uint64 = 21000000  
	InitialBlockReward uint64 = 10      
	HalvingInterval  int64  = 100000    
	MinTransactionFee uint64 = 1        
	MaxTransactionFee uint64 = 5        
	MinStakeAmount   uint64 = 100       
)

var (
	Economics = &CoinEconomics{
		TotalSupply: GenesisSupply,
		Circulating: GenesisSupply,
		Burned:      0,
		BlockReward: InitialBlockReward,
		NextHalving: HalvingInterval,
	}
)

func CalculateBlockReward(currentHeight int64) uint64 {
	halvings := currentHeight / HalvingInterval
	reward := InitialBlockReward
	
	for i := int64(0); i < halvings; i++ {
		reward /= 2
		if reward == 0 {
			return 0
		}
	}
	
	return reward
}

func CalculateNextHalving(currentHeight int64) int64 {
	nextHalving := (currentHeight/HalvingInterval + 1) * HalvingInterval
	return nextHalving
}

func UpdateEconomics(currentHeight int64, feesBurned uint64) {
	Economics.BlockReward = CalculateBlockReward(currentHeight)
	Economics.NextHalving = CalculateNextHalving(currentHeight)
	
	Economics.Burned += feesBurned
	
	Economics.Circulating = Economics.TotalSupply - Economics.Burned
}

func ValidateTransactionFee(fee uint64) bool {
	return fee >= MinTransactionFee && fee <= MaxTransactionFee
}

func DisplayEconomics() {
	fmt.Printf("\nLOGCOIN ECONOMICS\n")
	fmt.Printf("├─ Total Supply: %d LogCoins\n", Economics.TotalSupply)
	fmt.Printf("├─ Circulating: %d LogCoins\n", Economics.Circulating)
	fmt.Printf("├─ Burned: %d LogCoins\n", Economics.Burned)
	fmt.Printf("├─ Current Block Reward: %d LogCoins\n", Economics.BlockReward)
	fmt.Printf("├─ Next Halving: Block %d\n", Economics.NextHalving)
	fmt.Printf("└─ Transaction Fees: %d-%d LogCoins\n", MinTransactionFee, MaxTransactionFee)
}
