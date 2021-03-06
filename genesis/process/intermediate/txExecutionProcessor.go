package intermediate

import (
	"math"
	"math/big"

	"github.com/ElrondNetwork/elrond-go/core/check"
	"github.com/ElrondNetwork/elrond-go/data/state"
	transactionData "github.com/ElrondNetwork/elrond-go/data/transaction"
	"github.com/ElrondNetwork/elrond-go/genesis"
	"github.com/ElrondNetwork/elrond-go/process"
)

type txExecutionProcessor struct {
	txProcessor process.TransactionProcessor
	accounts    state.AccountsAdapter
}

// NewTxExecutionProcessor is able to execute a transaction
func NewTxExecutionProcessor(
	txProcessor process.TransactionProcessor,
	accounts state.AccountsAdapter,
) (*txExecutionProcessor, error) {
	if check.IfNil(txProcessor) {
		return nil, process.ErrNilTxProcessor
	}
	if check.IfNil(accounts) {
		return nil, process.ErrNilAccountsAdapter
	}

	return &txExecutionProcessor{
		txProcessor: txProcessor,
		accounts:    accounts,
	}, nil
}

// ExecuteTransaction will try to assemble a transaction and execute it against the accounts db
func (tep *txExecutionProcessor) ExecuteTransaction(
	nonce uint64,
	sndAddr []byte,
	rcvAddress []byte,
	value *big.Int,
	data []byte,
) error {
	tx := &transactionData.Transaction{
		Nonce:     nonce,
		SndAddr:   sndAddr,
		Value:     value,
		RcvAddr:   rcvAddress,
		GasPrice:  0,
		GasLimit:  math.MaxUint64,
		Data:      data,
		Signature: nil,
	}

	return tep.txProcessor.ProcessTransaction(tx)
}

// GetNonce returns the current nonce of the provided sender account
func (tep *txExecutionProcessor) GetNonce(senderBytes []byte) (uint64, error) {
	accnt, err := tep.accounts.LoadAccount(senderBytes)
	if err != nil {
		return 0, err
	}

	return accnt.GetNonce(), nil
}

// AccountExists returns if an account exists in the accounts DB
func (tep *txExecutionProcessor) AccountExists(address []byte) bool {
	_, err := tep.accounts.GetExistingAccount(address)
	return err == nil
}

// AddBalance adds the provided value on the balance field
func (tep *txExecutionProcessor) AddBalance(senderBytes []byte, value *big.Int) error {
	accnt, err := tep.accounts.LoadAccount(senderBytes)
	if err != nil {
		return err
	}

	userAccnt, ok := accnt.(state.UserAccountHandler)
	if !ok {
		return genesis.ErrWrongTypeAssertion
	}

	err = userAccnt.AddToBalance(value)
	if err != nil {
		return err
	}

	return tep.accounts.SaveAccount(userAccnt)
}

// AddNonce adds the provided value on the nonce field
func (tep *txExecutionProcessor) AddNonce(senderBytes []byte, nonce uint64) error {
	accnt, err := tep.accounts.LoadAccount(senderBytes)
	if err != nil {
		return err
	}

	accnt.IncreaseNonce(nonce)

	return tep.accounts.SaveAccount(accnt)
}

// IsInterfaceNil returns if underlying object is true
func (tep *txExecutionProcessor) IsInterfaceNil() bool {
	return tep == nil
}
