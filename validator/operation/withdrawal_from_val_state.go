package operation

import (
	"math/big"

	"gitlab.waterfall.network/waterfall/protocol/gwat/common"
)

type withdrawalFromValStateOperation struct {
	withdrawalAddress common.Address
	amount            *big.Int
}

func (op *withdrawalFromValStateOperation) init(
	withdrawalAddress common.Address,
	amount *big.Int,
) error {
	if withdrawalAddress == (common.Address{}) {
		return ErrNoWithdrawalAddress
	}

	if amount == nil {
		return ErrNoAmount
	}

	op.withdrawalAddress = withdrawalAddress
	op.amount = amount

	return nil
}

func NewWithdrawalFromValStateOperation(
	withdrawalAddress common.Address,
	amount *big.Int,
) (WithdrawalFromValState, error) {
	op := &withdrawalFromValStateOperation{}
	if err := op.init(withdrawalAddress, amount); err != nil {
		return nil, err
	}

	return op, nil
}

func (op *withdrawalFromValStateOperation) MarshalBinary() ([]byte, error) {
	data := make([]byte, 0)

	data = append(data, op.withdrawalAddress.Bytes()...)

	data = append(data, op.amount.Bytes()...)

	return data, nil
}

func (op *withdrawalFromValStateOperation) UnmarshalBinary(data []byte) error {
	withdrawalAddress := common.BytesToAddress(data[:common.AddressLength])

	amount := new(big.Int).SetBytes(data[common.AddressLength:])

	return op.init(withdrawalAddress, amount)
}

func (op *withdrawalFromValStateOperation) OpCode() Code {
	return WithdrawalFromValStateCode
}

func (op *withdrawalFromValStateOperation) WithdrawalAddress() common.Address {
	return op.withdrawalAddress
}

func (op *withdrawalFromValStateOperation) Amount() *big.Int {
	return op.amount
}
