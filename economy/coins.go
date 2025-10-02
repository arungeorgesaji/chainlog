package economy

const (
	GenesisSupply    = 1000000
	MaxSupply        = 21000000
	BlockReward      = 10
	HalvingInterval  = 100000
	MinTransactionFee = 1
	MaxTransactionFee = 5
	MinStakeAmount   = 100
)

func CalculateBlockReward(currentHeight int64) uint64 {
	halvings := currentHeight / HalvingInterval
	reward := uint64(BlockReward)
	
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
