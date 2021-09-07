package systemcontract

import (
	"github.com/modernizingpark/go-modernizingpark/common"
	"github.com/modernizingpark/go-modernizingpark/core"
	"github.com/modernizingpark/go-modernizingpark/core/state"
	"github.com/modernizingpark/go-modernizingpark/core/types"
	"github.com/modernizingpark/go-modernizingpark/core/vm"
	"github.com/modernizingpark/go-modernizingpark/log"
	"github.com/modernizingpark/go-modernizingpark/params"
	"math"
	"math/big"
)

var (
	devAdmin        = common.HexToAddress("0x6a0e7ae8eded108a4ec0e88d355d9279fad93fce")
	devAdminTestnet = common.HexToAddress("0x0dA5ac74D30D5b3c5ca9167A8666Ca98Fd58d9fb")
)

const (
	addressListCode = "0x608060405234801561001057600080fd5b506004361061010b5760003560e01c80634fb9e9b7116100a25780639e23c209116100715780639e23c209146102e2578063c4d66de814610308578063db6619b01461032e578063f851a44014610336578063fb48270c1461033e5761010b565b80634fb9e9b71461025f5780635eca4a70146102855780636dfb5176146102ab57806370b03fc5146102da5761010b565b806326782247116100de57806326782247146101fc578063327564b614610220578063349cb7111461022857806343e0c73a146102575761010b565b8063143d79b614610110578063158ef93e1461016057806318c662121461017c57806322fbf1e8146101d4575b600080fd5b6101366004803603602081101561012657600080fd5b50356001600160a01b0316610346565b60405180831515815260200182600281111561014e57fe5b81526020019250505060405180910390f35b6101686103c8565b604080519115158252519081900360200190f35b6101846103d1565b60408051602080825283518183015283519192839290830191858101910280838360005b838110156101c05781810151838201526020016101a8565b505050509050019250505060405180910390f35b6101fa600480360360208110156101ea57600080fd5b50356001600160a01b0316610433565b005b61020461052f565b604080516001600160a01b039092168252519081900360200190f35b61016861053e565b6101fa6004803603604081101561023e57600080fd5b5080356001600160a01b0316906020013560ff1661054c565b6101fa61077b565b6101fa6004803603602081101561027557600080fd5b50356001600160a01b0316610851565b6101686004803603602081101561029b57600080fd5b50356001600160a01b03166108ed565b6101fa600480360360408110156102c157600080fd5b5080356001600160a01b0316906020013560ff1661090b565b610184610bfb565b6101fa600480360360208110156102f857600080fd5b50356001600160a01b0316610c5b565b6101fa6004803603602081101561031e57600080fd5b50356001600160a01b0316610d55565b6101fa610dd4565b610204610eaf565b6101fa610ec4565b6001600160a01b0381166000908152600560209081526040808320546006909252822054829115801591151590829061037c5750805b1561039057600160029350935050506103c3565b81156103a557600160009350935050506103c3565b80156103b9576001809350935050506103c3565b6000809350935050505b915091565b60005460ff1681565b6060600380548060200260200160405190810160405280929190818152602001828054801561042957602002820191906000526020600020905b81546001600160a01b0316815260019091019060200180831161040b575b5050505050905090565b6000546201000090046001600160a01b03163314610485576040805162461bcd60e51b815260206004820152600a60248201526941646d696e206f6e6c7960b01b604482015290519081900360640190fd5b6001600160a01b03811660009081526002602052604090205460ff16156104e3576040805162461bcd60e51b815260206004820152600d60248201526c185b1c9958591e481859191959609a1b604482015290519081900360640190fd5b6001600160a01b038116600081815260026020526040808220805460ff19166001179055517f058fdae480ed8e99b762bceb2d39835a68ee3a4789cd84e5c90cd59722ba02099190a250565b6001546001600160a01b031681565b600054610100900460ff1681565b6000546201000090046001600160a01b0316331461059e576040805162461bcd60e51b815260206004820152600a60248201526941646d696e206f6e6c7960b01b604482015290519081900360640190fd5b60028160028111156105ac57fe5b141561068d576001600160a01b03821660009081526005602052604090205461060f576040805162461bcd60e51b815260206004820152601060248201526f1b9bdd081a5b88199c9bdb481b1a5cdd60821b604482015290519081900360640190fd5b6001600160a01b03821660009081526006602052604090205461066a576040805162461bcd60e51b815260206004820152600e60248201526d1b9bdd081a5b881d1bc81b1a5cdd60921b604482015290519081900360640190fd5b61067960036005846000610f7e565b61068860046006846001610f7e565b610777565b600081600281111561069b57fe5b141561070d576001600160a01b0382166000908152600560205260409020546106fe576040805162461bcd60e51b815260206004820152601060248201526f1b9bdd081a5b88199c9bdb481b1a5cdd60821b604482015290519081900360640190fd5b61068860036005846000610f7e565b6001600160a01b038216600090815260066020526040902054610768576040805162461bcd60e51b815260206004820152600e60248201526d1b9bdd081a5b881d1bc81b1a5cdd60921b604482015290519081900360640190fd5b61077760046006846001610f7e565b5050565b6000546201000090046001600160a01b031633146107cd576040805162461bcd60e51b815260206004820152600a60248201526941646d696e206f6e6c7960b01b604482015290519081900360640190fd5b600054610100900460ff1661081c576040805162461bcd60e51b815260206004820152601060248201526f185b1c9958591e48191a5cd8589b195960821b604482015290519081900360640190fd5b6000805461ff00191681556040517f733a7f99819dc7466bff56e7c0b6753b43b750a692f2a5bb4fe373815a0c7845908290a2565b6000546201000090046001600160a01b031633146108a3576040805162461bcd60e51b815260206004820152600a60248201526941646d696e206f6e6c7960b01b604482015290519081900360640190fd5b600180546001600160a01b0319166001600160a01b0383169081179091556040517faefcaa6215f99fe8c2f605dd268ee4d23a5b596bbca026e25ce8446187f4f1ba90600090a250565b6001600160a01b031660009081526002602052604090205460ff1690565b6000546201000090046001600160a01b0316331461095d576040805162461bcd60e51b815260206004820152600a60248201526941646d696e206f6e6c7960b01b604482015290519081900360640190fd5b6000546001600160a01b03838116620100009092041614156109c6576040805162461bcd60e51b815260206004820152601d60248201527f63616e6e6f74206164642061646d696e20746f20626c61636b6c697374000000604482015290519081900360640190fd5b60028160028111156109d457fe5b1415610abb576001600160a01b03821660009081526005602052604090205415610a3c576040805162461bcd60e51b8152602060048201526014602482015273185b1c9958591e481a5b88199c9bdb481b1a5cdd60621b604482015290519081900360640190fd5b6001600160a01b03821660009081526006602052604090205415610a9c576040805162461bcd60e51b8152602060048201526012602482015271185b1c9958591e481a5b881d1bc81b1a5cdd60721b604482015290519081900360640190fd5b610aa960036005846110d0565b610ab660046006846110d0565b610bab565b6000816002811115610ac957fe5b1415610b3e576001600160a01b03821660009081526005602052604090205415610b31576040805162461bcd60e51b8152602060048201526014602482015273185b1c9958591e481a5b88199c9bdb481b1a5cdd60621b604482015290519081900360640190fd5b610ab660036005846110d0565b6001600160a01b03821660009081526006602052604090205415610b9e576040805162461bcd60e51b8152602060048201526012602482015271185b1c9958591e481a5b881d1bc81b1a5cdd60721b604482015290519081900360640190fd5b610bab60046006846110d0565b816001600160a01b03167f4bb8845da5ed7c2df200814ba7a0f3db11326cc817cf9a042fa54d4e5f6f29bb8260405180826002811115610be757fe5b815260200191505060405180910390a25050565b60606004805480602002602001604051908101604052809291908181526020018280548015610429576020028201919060005260206000209081546001600160a01b0316815260019091019060200180831161040b575050505050905090565b6000546201000090046001600160a01b03163314610cad576040805162461bcd60e51b815260206004820152600a60248201526941646d696e206f6e6c7960b01b604482015290519081900360640190fd5b6001600160a01b03811660009081526002602052604090205460ff16610d0c576040805162461bcd60e51b815260206004820152600f60248201526e3737ba1030903232bb32b637b832b960891b604482015290519081900360640190fd5b6001600160a01b038116600081815260026020526040808220805460ff19169055517f110a48e3e347ae018d4d40446e4e917b416f912dec489da19b4507bb9bb18cd49190a250565b60005460ff1615610da3576040805162461bcd60e51b8152602060048201526013602482015272105b1c9958591e481a5b9a5d1a585b1a5e9959606a1b604482015290519081900360640190fd5b6000805460ff196001600160a01b03909316620100000262010000600160b01b031990911617919091166001179055565b6000546201000090046001600160a01b03163314610e26576040805162461bcd60e51b815260206004820152600a60248201526941646d696e206f6e6c7960b01b604482015290519081900360640190fd5b600054610100900460ff1615610e75576040805162461bcd60e51b815260206004820152600f60248201526e185b1c9958591e48195b98589b1959608a1b604482015290519081900360640190fd5b6000805461ff0019166101001781556040516001917f733a7f99819dc7466bff56e7c0b6753b43b750a692f2a5bb4fe373815a0c784591a2565b6000546201000090046001600160a01b031681565b6001546001600160a01b03163314610f14576040805162461bcd60e51b815260206004820152600e60248201526d4e65772061646d696e206f6e6c7960901b604482015290519081900360640190fd5b600180546000805462010000600160b01b0319166001600160a01b0380841662010000908102929092178084556001600160a01b03199094169094556040519204909216917f7ce7ec0b50378fb6c0186ffb5f48325f6593fcb4ca4386f21861af3129188f5c91a2565b6001600160a01b0382166000908152602084905260408120805491905584546000199182019101811461104e57845485906000198101908110610fbd57fe5b9060005260206000200160009054906101000a90046001600160a01b0316858281548110610fe757fe5b9060005260206000200160006101000a8154816001600160a01b0302191690836001600160a01b031602179055508060010184600087848154811061102857fe5b60009182526020808320909101546001600160a01b031683528201929092526040019020555b8480548061105857fe5b600082815260209020810160001990810180546001600160a01b03191690550190556040516001600160a01b038416907f91b762fba034b39c8b14c1e6463a15b1f4c211dcd0023f7fa2f4ae2928dfc44d908490808260028111156110b957fe5b815260200191505060405180910390a25050505050565b82546001810184556000848152602080822090920180546001600160a01b039094166001600160a01b03199094168417905593549184529190915260409091205556fea2646970667358221220f12504ef643371812a137430fb60f43b1015de59a16a39edcc4fda651b53e11a64736f6c634300060c0033"
)

type hardForkAddressList struct {
}

func (s *hardForkAddressList) GetName() string {
	return AddressListContractName
}

func (s *hardForkAddressList) Update(config *params.ChainConfig, height *big.Int, state *state.StateDB) (err error) {
	contractCode := common.FromHex(addressListCode)

	//write addressListCode to sys contract
	state.SetCode(AddressListContractAddr, contractCode)
	log.Debug("Write code to system contract account", "addr", AddressListContractAddr.String(), "code", addressListCode)

	return
}

func (s *hardForkAddressList) getAdminByChainId(chainId *big.Int) common.Address {
	if chainId.Cmp(params.MainnetChainConfig.ChainID) == 0 {
		return devAdmin
	}

	return devAdminTestnet
}

func (s *hardForkAddressList) Execute(state *state.StateDB, header *types.Header, chainContext core.ChainContext, config *params.ChainConfig) (err error) {

	method := "initialize"
	data, err := GetInteractiveABI()[AddressListContractName].Pack(method, s.getAdminByChainId(config.ChainID))
	if err != nil {
		log.Error("Can't pack data for initialize", "error", err)
		return err
	}

	msg := types.NewMessage(header.Coinbase, &AddressListContractAddr, 0, new(big.Int), math.MaxUint64, new(big.Int), data, false)

	context := core.NewEVMContext(msg, header, chainContext, nil)
	evm := vm.NewEVM(context, state, config, vm.Config{})

	_, _, err = evm.Call(vm.AccountRef(msg.From()), *msg.To(), msg.Data(), msg.Gas(), msg.Value())

	return
}
