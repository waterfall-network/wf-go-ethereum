package operation

type withdrawalFromValStateOperation struct{}

func NewWithdrawalFromValStateOperation() WithdrawalFromValState {
	return &withdrawalFromValStateOperation{}
}

func (op *withdrawalFromValStateOperation) MarshalBinary() ([]byte, error) {
	data := make([]byte, 0)

	return data, nil
}

func (op *withdrawalFromValStateOperation) UnmarshalBinary(data []byte) error {
	return nil
}

func (op *withdrawalFromValStateOperation) OpCode() Code {
	return WithdrawalFromValStateCode
}
