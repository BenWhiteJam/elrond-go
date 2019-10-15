package mock

type FeeHandlerStub struct {
	MinGasPriceCalled func() uint64
	MinGasLimitCalled func() uint64
}

func (fhs *FeeHandlerStub) MinGasPrice() uint64 {
	return fhs.MinGasPriceCalled()
}

func (fhs *FeeHandlerStub) MinGasLimit() uint64 {
	return fhs.MinGasLimitCalled()
}

// IsInterfaceNil returns true if there is no value under the interface
func (fhs *FeeHandlerStub) IsInterfaceNil() bool {
	if fhs == nil {
		return true
	}
	return false
}
