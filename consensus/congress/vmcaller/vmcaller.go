package vmcaller

import (
	"github.com/modernizingpark/go-modernizingpark/core"
	"github.com/modernizingpark/go-modernizingpark/core/state"
	"github.com/modernizingpark/go-modernizingpark/core/types"
	"github.com/modernizingpark/go-modernizingpark/core/vm"
	"github.com/modernizingpark/go-modernizingpark/params"
)

// ExecuteMsg executes transaction sent to system contracts.
func ExecuteMsg(msg core.Message, state *state.StateDB, header *types.Header, chainContext core.ChainContext, chainConfig *params.ChainConfig) (ret []byte, err error) {
	// Set gas price to zero
	context := core.NewEVMContext(msg, header, chainContext, &(header.Coinbase))
	vmenv := vm.NewEVM(context, state, chainConfig, vm.Config{})

	ret, _, err = vmenv.Call(vm.AccountRef(msg.From()), *msg.To(), msg.Data(), msg.Gas(), msg.Value())

	return ret, err
}
