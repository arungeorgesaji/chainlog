package economy

import (
	"chainlog/core"
	"chainlog/storage"
	"fmt"
)

type TransactionProcessor struct {
	FeeCollector string 
	BurnAddress  string 
	StateManager *storage.StateManager 
}

func NewTransactionProcessor(feeCollector string, stateManager *storage.StateManager) *TransactionProcessor {
	return &TransactionProcessor{
		FeeCollector: feeCollector,
		BurnAddress:  "BURN0000000000000000000000000000",
		StateManager: stateManager, 
	}
}

func (tp *TransactionProcessor) CalculateFeeDistribution(fee uint64) (uint64, uint64) {
	minerShare := (fee * 80 + 50) / 100  
	burnAmount := fee - minerShare  

	return minerShare, burnAmount
}

func (tp *TransactionProcessor) ProcessTransaction(tx *core.Transaction) error {
	if tx.ID == "" {
    return fmt.Errorf("transaction ID is required")
}

	if tx.Sender == "" {
			return fmt.Errorf("sender address is required")
	}

	if tx.Timestamp == 0 {
			return fmt.Errorf("transaction timestamp is required")
	}

	if tx.ID != tx.CalculateID() {
			return fmt.Errorf("transaction ID is invalid")
	}

	if !ValidateTransactionFee(tx.Fee) {
		return fmt.Errorf("invalid transaction fee: %d (must be %d-%d)", 
			tx.Fee, MinTransactionFee, MaxTransactionFee)
	}
	
	if !tp.CheckSufficientBalance(tx.Sender, tx.Amount + tx.Fee) {
		currentBalance := tp.StateManager.GetBalance(tx.Sender)
		return fmt.Errorf("insufficient balance: have %d, need %d (amount: %d + fee: %d)", 
			currentBalance, tx.Amount + tx.Fee, tx.Amount, tx.Fee)
	}
	
	if tx.Amount > 0 && tx.Receiver == "" {
		return fmt.Errorf("transfer transaction must have a receiver")
	}

	switch tx.Type {
		case core.DataTx:
				if tx.Data == "" {
						return fmt.Errorf("data transactions must contain data")
				}
		case core.TransferTx:
				if tx.Receiver == "" {
						return fmt.Errorf("transfer transactions require a receiver")
				}
				if tx.Amount == 0 {
						return fmt.Errorf("transfer amount must be greater than 0")
				}
				if tx.Sender == tx.Receiver {
						return fmt.Errorf("cannot transfer to self")
				}
		case core.StakeTx:
				if tx.Amount == 0 {
						return fmt.Errorf("staking amount must be greater than 0")
				}
				if tx.Amount < 100 { 
						return fmt.Errorf("minimum stake amount is 100 LogCoins")
				}
		}
	
	switch tx.Type {
		case core.DataTx:
			return tp.processDataTransaction(tx)
		case core.TransferTx: 
			return tp.processTransferTransaction(tx)
		case core.StakeTx:
			return tp.processStakeTransaction(tx)
		default:
			return fmt.Errorf("unsupported transaction type: %d", tx.Type)
	}
}


func (tp *TransactionProcessor) CheckSufficientBalance(address string, totalAmount uint64) bool {
	balance := tp.StateManager.GetBalance(address)
	return balance >= totalAmount
}

func (tp *TransactionProcessor) processDataTransaction(tx *core.Transaction) error {
	if err := tp.deductBalance(tx.Sender, tx.Fee); err != nil {
		return err
	}
	
	fmt.Printf("Data transaction processed: %s paid %d LogCoin fee\n", 
		tx.Sender[:8], tx.Fee)
	return nil
}

func (tp *TransactionProcessor) processTransferTransaction(tx *core.Transaction) error {
	totalDeduct := tx.Amount + tx.Fee
	if err := tp.deductBalance(tx.Sender, totalDeduct); err != nil {
		return err
	}
	
	if err := tp.addBalance(tx.Receiver, tx.Amount); err != nil {
		return err
	}
	
	fmt.Printf("Transfer processed: %s → %s: %d LogCoins (fee: %d)\n", 
		tx.Sender[:8], tx.Receiver[:8], tx.Amount, tx.Fee)
	return nil
}

func (tp *TransactionProcessor) processStakeTransaction(tx *core.Transaction) error {
	if err := tp.deductBalance(tx.Sender, tx.Amount); err != nil {
		return err
	}
	
	fmt.Printf("Stake processed: %s staked %d LogCoins\n", 
		tx.Sender[:8], tx.Amount)
	return nil
}

func (tp *TransactionProcessor) deductBalance(address string, amount uint64) error {
	account, exists := tp.StateManager.GetAccount(address)
	if !exists {
		return fmt.Errorf("account not found: %s", address)
	}
	
	if account.Balance < amount {
		return fmt.Errorf("insufficient balance: have %d, need %d", account.Balance, amount)
	}
	
	account.Balance -= amount
	account.Nonce++
	return nil
}

func (tp *TransactionProcessor) addBalance(address string, amount uint64) error {
	account, exists := tp.StateManager.GetAccount(address)
	if exists {
		if ^uint64(0) - account.Balance < amount {
			return fmt.Errorf("balance overflow: cannot add %d to current balance %d", amount, account.Balance)
    }
		account.Balance += amount
	} else {
		tp.StateManager.UpdateAccount(address, amount, 0)
	}
	return nil
}

func (tp *TransactionProcessor) CreateFeeDistributionTransactions(tx *core.Transaction) []*core.Transaction {
	minerShare, burnAmount := tp.CalculateFeeDistribution(tx.Fee)
	
	var feeTxs []*core.Transaction
	
	if minerShare > 0 {
		tp.addBalance(tp.FeeCollector, minerShare)
		
		feeTxs = append(feeTxs, &core.Transaction{
			ID:        fmt.Sprintf("fee-%s", tx.ID),
			Type:      core.FeeTx,
			Sender:    tx.Sender,
			Receiver:  tp.FeeCollector,
			Amount:    minerShare,
			Fee:       0, 
			Timestamp: tx.Timestamp,
			Nonce:     tx.Nonce + 1, 
		})
	}
	
	if burnAmount > 0 {
		feeTxs = append(feeTxs, &core.Transaction{
			ID:        fmt.Sprintf("burn-%s", tx.ID),
			Type:      core.FeeTx,
			Sender:    tx.Sender,
			Receiver:  tp.BurnAddress,
			Amount:    burnAmount,
			Fee:       0,
			Timestamp: tx.Timestamp,
			Nonce:     tx.Nonce + 2,
		})
	}
	
	return feeTxs
}
