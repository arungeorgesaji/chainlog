package consensus

import (
	"chainlog/core"
	"chainlog/crypto"
	"chainlog/economy"
	"fmt"
	"time"
)

type Miner struct {
	Wallet     *crypto.Wallet
	Blockchain *core.Blockchain
	Address    string
	IsMining   bool
	stopChan   chan bool
}

func NewMiner(wallet *crypto.Wallet, bc *core.Blockchain) *Miner {
	return &Miner{
		Wallet:     wallet,
		Blockchain: bc,
		Address:    wallet.GetAddress(),
		IsMining:   false,
		stopChan:   make(chan bool),
	}
}

func (m *Miner) MineBlock() (*core.Block, error) {
	if len(m.Blockchain.PendingTx) == 0 {
		return nil, fmt.Errorf("no transactions to mine")
	}
	
	lastBlock := m.Blockchain.GetLastBlock()
	
	newBlock := &core.Block{
		Index:        lastBlock.Index + 1,
		Timestamp:    time.Now().Unix(),
		Transactions: m.Blockchain.PendingTx,
		PrevHash:     lastBlock.Hash,
		Difficulty:   m.Blockchain.Difficulty,
		Miner:        m.Address,
	}
	
	pow := NewProofOfWork(newBlock, m.Blockchain.Difficulty)
	nonce, hash := pow.Run()
	
	if nonce == -1 {
		return nil, fmt.Errorf("failed to mine block")
	}
	
	newBlock.Nonce = nonce
	newBlock.Hash = hash
	
	rewardTx := m.createRewardTransaction()
	newBlock.Transactions = append([]*core.Transaction{rewardTx}, newBlock.Transactions...)
	
	return newBlock, nil
}

func (m *Miner) createRewardTransaction() *core.Transaction {
	currentHeight := int64(m.Blockchain.GetBlockCount())
  blockReward := economy.CalculateBlockReward(currentHeight)
	
	return &core.Transaction{
		ID:        fmt.Sprintf("reward-%d", time.Now().UnixNano()),
		Type:      core.RewardTx,
		Sender:    "network", 
		Receiver:  m.Address,
		Amount:    blockReward,
		Fee:       0,
		Timestamp: time.Now().Unix(),
		Nonce:     uint64(time.Now().UnixNano()),
	}
}

func (m *Miner) StartMining() {
	m.IsMining = true
	fmt.Printf("Miner %s started mining...\n", m.Wallet.GetAddressShort())
	
	go func() {
		for m.IsMining {
			select {
			case <-m.stopChan:
				return
			default:
				if len(m.Blockchain.PendingTx) > 0 {
					block, err := m.MineBlock()
					if err == nil {
						m.Blockchain.Chain = append(m.Blockchain.Chain, block)
						m.Blockchain.ClearPendingTransactions()
						currentHeight := int64(m.Blockchain.GetBlockCount() - 1)
						fmt.Printf("Mined block %d! Reward: %d LogCoins\n", 
							block.Index, economy.CalculateBlockReward(currentHeight))
					}
				}
				time.Sleep(1 * time.Second) 
			}
		}
	}()
}

func (m *Miner) StopMining() {
	m.IsMining = false
	close(m.stopChan)
	fmt.Printf("Miner %s stopped mining\n", m.Wallet.GetAddressShort())
}

func (m *Miner) GetMiningStatus() string {
	if m.IsMining {
		return "MINING"
	}
	return "STOPPED"
}

func (m *Miner) Display() {
	fmt.Printf("\nMINER INFORMATION\n")
	fmt.Printf("├─ Address: %s\n", m.Wallet.GetAddressShort())
	fmt.Printf("├─ Status: %s\n", m.GetMiningStatus())
	fmt.Printf("├─ Blockchain Height: %d\n", m.Blockchain.GetBlockCount())
	fmt.Printf("├─ Pending Transactions: %d\n", len(m.Blockchain.PendingTx))
	fmt.Printf("└─ Current Difficulty: %d\n", m.Blockchain.Difficulty)
}
