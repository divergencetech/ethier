// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package chainlinktestabi

import (
	"errors"
	"math/big"
	"strings"

	ethereum "github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/event"
)

// Reference imports to suppress errors if they are not otherwise used.
var (
	_ = errors.New
	_ = big.NewInt
	_ = strings.NewReader
	_ = ethereum.NotFound
	_ = bind.Bind
	_ = common.Big1
	_ = types.BloomLookup
	_ = event.NewSubscription
)

// ChainlinkMetaData contains all meta data concerning the Chainlink contract.
var ChainlinkMetaData = &bind.MetaData{
	ABI: "[]",
	Bin: "0x60566050600b82828239805160001a6073146043577f4e487b7100000000000000000000000000000000000000000000000000000000600052600060045260246000fd5b30600052607381538281f3fe73000000000000000000000000000000000000000030146080604052600080fdfea2646970667358221220d96aff142f5b971a201ac252989e7d7df30b9c991e3e192b253d7bf2ee2eda9e64736f6c634300080b0033",
}

// ChainlinkABI is the input ABI used to generate the binding from.
// Deprecated: Use ChainlinkMetaData.ABI instead.
var ChainlinkABI = ChainlinkMetaData.ABI

// ChainlinkBin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use ChainlinkMetaData.Bin instead.
var ChainlinkBin = ChainlinkMetaData.Bin

// DeployChainlink deploys a new Ethereum contract, binding an instance of Chainlink to it.
func DeployChainlink(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *Chainlink, error) {
	parsed, err := ChainlinkMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(ChainlinkBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &Chainlink{ChainlinkCaller: ChainlinkCaller{contract: contract}, ChainlinkTransactor: ChainlinkTransactor{contract: contract}, ChainlinkFilterer: ChainlinkFilterer{contract: contract}}, nil
}

// Chainlink is an auto generated Go binding around an Ethereum contract.
type Chainlink struct {
	ChainlinkCaller     // Read-only binding to the contract
	ChainlinkTransactor // Write-only binding to the contract
	ChainlinkFilterer   // Log filterer for contract events
}

// ChainlinkCaller is an auto generated read-only Go binding around an Ethereum contract.
type ChainlinkCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ChainlinkTransactor is an auto generated write-only Go binding around an Ethereum contract.
type ChainlinkTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ChainlinkFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type ChainlinkFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ChainlinkSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type ChainlinkSession struct {
	Contract     *Chainlink        // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// ChainlinkCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type ChainlinkCallerSession struct {
	Contract *ChainlinkCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts    // Call options to use throughout this session
}

// ChainlinkTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type ChainlinkTransactorSession struct {
	Contract     *ChainlinkTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts    // Transaction auth options to use throughout this session
}

// ChainlinkRaw is an auto generated low-level Go binding around an Ethereum contract.
type ChainlinkRaw struct {
	Contract *Chainlink // Generic contract binding to access the raw methods on
}

// ChainlinkCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type ChainlinkCallerRaw struct {
	Contract *ChainlinkCaller // Generic read-only contract binding to access the raw methods on
}

// ChainlinkTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type ChainlinkTransactorRaw struct {
	Contract *ChainlinkTransactor // Generic write-only contract binding to access the raw methods on
}

// NewChainlink creates a new instance of Chainlink, bound to a specific deployed contract.
func NewChainlink(address common.Address, backend bind.ContractBackend) (*Chainlink, error) {
	contract, err := bindChainlink(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &Chainlink{ChainlinkCaller: ChainlinkCaller{contract: contract}, ChainlinkTransactor: ChainlinkTransactor{contract: contract}, ChainlinkFilterer: ChainlinkFilterer{contract: contract}}, nil
}

// NewChainlinkCaller creates a new read-only instance of Chainlink, bound to a specific deployed contract.
func NewChainlinkCaller(address common.Address, caller bind.ContractCaller) (*ChainlinkCaller, error) {
	contract, err := bindChainlink(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &ChainlinkCaller{contract: contract}, nil
}

// NewChainlinkTransactor creates a new write-only instance of Chainlink, bound to a specific deployed contract.
func NewChainlinkTransactor(address common.Address, transactor bind.ContractTransactor) (*ChainlinkTransactor, error) {
	contract, err := bindChainlink(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &ChainlinkTransactor{contract: contract}, nil
}

// NewChainlinkFilterer creates a new log filterer instance of Chainlink, bound to a specific deployed contract.
func NewChainlinkFilterer(address common.Address, filterer bind.ContractFilterer) (*ChainlinkFilterer, error) {
	contract, err := bindChainlink(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &ChainlinkFilterer{contract: contract}, nil
}

// bindChainlink binds a generic wrapper to an already deployed contract.
func bindChainlink(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(ChainlinkABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Chainlink *ChainlinkRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Chainlink.Contract.ChainlinkCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Chainlink *ChainlinkRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Chainlink.Contract.ChainlinkTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Chainlink *ChainlinkRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Chainlink.Contract.ChainlinkTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Chainlink *ChainlinkCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Chainlink.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Chainlink *ChainlinkTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Chainlink.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Chainlink *ChainlinkTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Chainlink.Contract.contract.Transact(opts, method, params...)
}

// ContextMetaData contains all meta data concerning the Context contract.
var ContextMetaData = &bind.MetaData{
	ABI: "[]",
}

// ContextABI is the input ABI used to generate the binding from.
// Deprecated: Use ContextMetaData.ABI instead.
var ContextABI = ContextMetaData.ABI

// Context is an auto generated Go binding around an Ethereum contract.
type Context struct {
	ContextCaller     // Read-only binding to the contract
	ContextTransactor // Write-only binding to the contract
	ContextFilterer   // Log filterer for contract events
}

// ContextCaller is an auto generated read-only Go binding around an Ethereum contract.
type ContextCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ContextTransactor is an auto generated write-only Go binding around an Ethereum contract.
type ContextTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ContextFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type ContextFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ContextSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type ContextSession struct {
	Contract     *Context          // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// ContextCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type ContextCallerSession struct {
	Contract *ContextCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts  // Call options to use throughout this session
}

// ContextTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type ContextTransactorSession struct {
	Contract     *ContextTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts  // Transaction auth options to use throughout this session
}

// ContextRaw is an auto generated low-level Go binding around an Ethereum contract.
type ContextRaw struct {
	Contract *Context // Generic contract binding to access the raw methods on
}

// ContextCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type ContextCallerRaw struct {
	Contract *ContextCaller // Generic read-only contract binding to access the raw methods on
}

// ContextTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type ContextTransactorRaw struct {
	Contract *ContextTransactor // Generic write-only contract binding to access the raw methods on
}

// NewContext creates a new instance of Context, bound to a specific deployed contract.
func NewContext(address common.Address, backend bind.ContractBackend) (*Context, error) {
	contract, err := bindContext(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &Context{ContextCaller: ContextCaller{contract: contract}, ContextTransactor: ContextTransactor{contract: contract}, ContextFilterer: ContextFilterer{contract: contract}}, nil
}

// NewContextCaller creates a new read-only instance of Context, bound to a specific deployed contract.
func NewContextCaller(address common.Address, caller bind.ContractCaller) (*ContextCaller, error) {
	contract, err := bindContext(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &ContextCaller{contract: contract}, nil
}

// NewContextTransactor creates a new write-only instance of Context, bound to a specific deployed contract.
func NewContextTransactor(address common.Address, transactor bind.ContractTransactor) (*ContextTransactor, error) {
	contract, err := bindContext(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &ContextTransactor{contract: contract}, nil
}

// NewContextFilterer creates a new log filterer instance of Context, bound to a specific deployed contract.
func NewContextFilterer(address common.Address, filterer bind.ContractFilterer) (*ContextFilterer, error) {
	contract, err := bindContext(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &ContextFilterer{contract: contract}, nil
}

// bindContext binds a generic wrapper to an already deployed contract.
func bindContext(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(ContextABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Context *ContextRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Context.Contract.ContextCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Context *ContextRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Context.Contract.ContextTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Context *ContextRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Context.Contract.ContextTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Context *ContextCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Context.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Context *ContextTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Context.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Context *ContextTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Context.Contract.contract.Transact(opts, method, params...)
}

// LinkTokenInterfaceMetaData contains all meta data concerning the LinkTokenInterface contract.
var LinkTokenInterfaceMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"spender\",\"type\":\"address\"}],\"name\":\"allowance\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"remaining\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"spender\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"approve\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"success\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"}],\"name\":\"balanceOf\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"balance\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"decimals\",\"outputs\":[{\"internalType\":\"uint8\",\"name\":\"decimalPlaces\",\"type\":\"uint8\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"spender\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"addedValue\",\"type\":\"uint256\"}],\"name\":\"decreaseApproval\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"success\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"spender\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"subtractedValue\",\"type\":\"uint256\"}],\"name\":\"increaseApproval\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"name\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"tokenName\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"symbol\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"tokenSymbol\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"totalSupply\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"totalTokensIssued\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"transfer\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"success\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"data\",\"type\":\"bytes\"}],\"name\":\"transferAndCall\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"success\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"transferFrom\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"success\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
}

// LinkTokenInterfaceABI is the input ABI used to generate the binding from.
// Deprecated: Use LinkTokenInterfaceMetaData.ABI instead.
var LinkTokenInterfaceABI = LinkTokenInterfaceMetaData.ABI

// LinkTokenInterface is an auto generated Go binding around an Ethereum contract.
type LinkTokenInterface struct {
	LinkTokenInterfaceCaller     // Read-only binding to the contract
	LinkTokenInterfaceTransactor // Write-only binding to the contract
	LinkTokenInterfaceFilterer   // Log filterer for contract events
}

// LinkTokenInterfaceCaller is an auto generated read-only Go binding around an Ethereum contract.
type LinkTokenInterfaceCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// LinkTokenInterfaceTransactor is an auto generated write-only Go binding around an Ethereum contract.
type LinkTokenInterfaceTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// LinkTokenInterfaceFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type LinkTokenInterfaceFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// LinkTokenInterfaceSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type LinkTokenInterfaceSession struct {
	Contract     *LinkTokenInterface // Generic contract binding to set the session for
	CallOpts     bind.CallOpts       // Call options to use throughout this session
	TransactOpts bind.TransactOpts   // Transaction auth options to use throughout this session
}

// LinkTokenInterfaceCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type LinkTokenInterfaceCallerSession struct {
	Contract *LinkTokenInterfaceCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts             // Call options to use throughout this session
}

// LinkTokenInterfaceTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type LinkTokenInterfaceTransactorSession struct {
	Contract     *LinkTokenInterfaceTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts             // Transaction auth options to use throughout this session
}

// LinkTokenInterfaceRaw is an auto generated low-level Go binding around an Ethereum contract.
type LinkTokenInterfaceRaw struct {
	Contract *LinkTokenInterface // Generic contract binding to access the raw methods on
}

// LinkTokenInterfaceCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type LinkTokenInterfaceCallerRaw struct {
	Contract *LinkTokenInterfaceCaller // Generic read-only contract binding to access the raw methods on
}

// LinkTokenInterfaceTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type LinkTokenInterfaceTransactorRaw struct {
	Contract *LinkTokenInterfaceTransactor // Generic write-only contract binding to access the raw methods on
}

// NewLinkTokenInterface creates a new instance of LinkTokenInterface, bound to a specific deployed contract.
func NewLinkTokenInterface(address common.Address, backend bind.ContractBackend) (*LinkTokenInterface, error) {
	contract, err := bindLinkTokenInterface(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &LinkTokenInterface{LinkTokenInterfaceCaller: LinkTokenInterfaceCaller{contract: contract}, LinkTokenInterfaceTransactor: LinkTokenInterfaceTransactor{contract: contract}, LinkTokenInterfaceFilterer: LinkTokenInterfaceFilterer{contract: contract}}, nil
}

// NewLinkTokenInterfaceCaller creates a new read-only instance of LinkTokenInterface, bound to a specific deployed contract.
func NewLinkTokenInterfaceCaller(address common.Address, caller bind.ContractCaller) (*LinkTokenInterfaceCaller, error) {
	contract, err := bindLinkTokenInterface(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &LinkTokenInterfaceCaller{contract: contract}, nil
}

// NewLinkTokenInterfaceTransactor creates a new write-only instance of LinkTokenInterface, bound to a specific deployed contract.
func NewLinkTokenInterfaceTransactor(address common.Address, transactor bind.ContractTransactor) (*LinkTokenInterfaceTransactor, error) {
	contract, err := bindLinkTokenInterface(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &LinkTokenInterfaceTransactor{contract: contract}, nil
}

// NewLinkTokenInterfaceFilterer creates a new log filterer instance of LinkTokenInterface, bound to a specific deployed contract.
func NewLinkTokenInterfaceFilterer(address common.Address, filterer bind.ContractFilterer) (*LinkTokenInterfaceFilterer, error) {
	contract, err := bindLinkTokenInterface(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &LinkTokenInterfaceFilterer{contract: contract}, nil
}

// bindLinkTokenInterface binds a generic wrapper to an already deployed contract.
func bindLinkTokenInterface(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(LinkTokenInterfaceABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_LinkTokenInterface *LinkTokenInterfaceRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _LinkTokenInterface.Contract.LinkTokenInterfaceCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_LinkTokenInterface *LinkTokenInterfaceRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _LinkTokenInterface.Contract.LinkTokenInterfaceTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_LinkTokenInterface *LinkTokenInterfaceRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _LinkTokenInterface.Contract.LinkTokenInterfaceTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_LinkTokenInterface *LinkTokenInterfaceCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _LinkTokenInterface.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_LinkTokenInterface *LinkTokenInterfaceTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _LinkTokenInterface.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_LinkTokenInterface *LinkTokenInterfaceTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _LinkTokenInterface.Contract.contract.Transact(opts, method, params...)
}

// Allowance is a free data retrieval call binding the contract method 0xdd62ed3e.
//
// Solidity: function allowance(address owner, address spender) view returns(uint256 remaining)
func (_LinkTokenInterface *LinkTokenInterfaceCaller) Allowance(opts *bind.CallOpts, owner common.Address, spender common.Address) (*big.Int, error) {
	var out []interface{}
	err := _LinkTokenInterface.contract.Call(opts, &out, "allowance", owner, spender)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// Allowance is a free data retrieval call binding the contract method 0xdd62ed3e.
//
// Solidity: function allowance(address owner, address spender) view returns(uint256 remaining)
func (_LinkTokenInterface *LinkTokenInterfaceSession) Allowance(owner common.Address, spender common.Address) (*big.Int, error) {
	return _LinkTokenInterface.Contract.Allowance(&_LinkTokenInterface.CallOpts, owner, spender)
}

// Allowance is a free data retrieval call binding the contract method 0xdd62ed3e.
//
// Solidity: function allowance(address owner, address spender) view returns(uint256 remaining)
func (_LinkTokenInterface *LinkTokenInterfaceCallerSession) Allowance(owner common.Address, spender common.Address) (*big.Int, error) {
	return _LinkTokenInterface.Contract.Allowance(&_LinkTokenInterface.CallOpts, owner, spender)
}

// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
//
// Solidity: function balanceOf(address owner) view returns(uint256 balance)
func (_LinkTokenInterface *LinkTokenInterfaceCaller) BalanceOf(opts *bind.CallOpts, owner common.Address) (*big.Int, error) {
	var out []interface{}
	err := _LinkTokenInterface.contract.Call(opts, &out, "balanceOf", owner)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
//
// Solidity: function balanceOf(address owner) view returns(uint256 balance)
func (_LinkTokenInterface *LinkTokenInterfaceSession) BalanceOf(owner common.Address) (*big.Int, error) {
	return _LinkTokenInterface.Contract.BalanceOf(&_LinkTokenInterface.CallOpts, owner)
}

// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
//
// Solidity: function balanceOf(address owner) view returns(uint256 balance)
func (_LinkTokenInterface *LinkTokenInterfaceCallerSession) BalanceOf(owner common.Address) (*big.Int, error) {
	return _LinkTokenInterface.Contract.BalanceOf(&_LinkTokenInterface.CallOpts, owner)
}

// Decimals is a free data retrieval call binding the contract method 0x313ce567.
//
// Solidity: function decimals() view returns(uint8 decimalPlaces)
func (_LinkTokenInterface *LinkTokenInterfaceCaller) Decimals(opts *bind.CallOpts) (uint8, error) {
	var out []interface{}
	err := _LinkTokenInterface.contract.Call(opts, &out, "decimals")

	if err != nil {
		return *new(uint8), err
	}

	out0 := *abi.ConvertType(out[0], new(uint8)).(*uint8)

	return out0, err

}

// Decimals is a free data retrieval call binding the contract method 0x313ce567.
//
// Solidity: function decimals() view returns(uint8 decimalPlaces)
func (_LinkTokenInterface *LinkTokenInterfaceSession) Decimals() (uint8, error) {
	return _LinkTokenInterface.Contract.Decimals(&_LinkTokenInterface.CallOpts)
}

// Decimals is a free data retrieval call binding the contract method 0x313ce567.
//
// Solidity: function decimals() view returns(uint8 decimalPlaces)
func (_LinkTokenInterface *LinkTokenInterfaceCallerSession) Decimals() (uint8, error) {
	return _LinkTokenInterface.Contract.Decimals(&_LinkTokenInterface.CallOpts)
}

// Name is a free data retrieval call binding the contract method 0x06fdde03.
//
// Solidity: function name() view returns(string tokenName)
func (_LinkTokenInterface *LinkTokenInterfaceCaller) Name(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _LinkTokenInterface.contract.Call(opts, &out, "name")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// Name is a free data retrieval call binding the contract method 0x06fdde03.
//
// Solidity: function name() view returns(string tokenName)
func (_LinkTokenInterface *LinkTokenInterfaceSession) Name() (string, error) {
	return _LinkTokenInterface.Contract.Name(&_LinkTokenInterface.CallOpts)
}

// Name is a free data retrieval call binding the contract method 0x06fdde03.
//
// Solidity: function name() view returns(string tokenName)
func (_LinkTokenInterface *LinkTokenInterfaceCallerSession) Name() (string, error) {
	return _LinkTokenInterface.Contract.Name(&_LinkTokenInterface.CallOpts)
}

// Symbol is a free data retrieval call binding the contract method 0x95d89b41.
//
// Solidity: function symbol() view returns(string tokenSymbol)
func (_LinkTokenInterface *LinkTokenInterfaceCaller) Symbol(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _LinkTokenInterface.contract.Call(opts, &out, "symbol")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// Symbol is a free data retrieval call binding the contract method 0x95d89b41.
//
// Solidity: function symbol() view returns(string tokenSymbol)
func (_LinkTokenInterface *LinkTokenInterfaceSession) Symbol() (string, error) {
	return _LinkTokenInterface.Contract.Symbol(&_LinkTokenInterface.CallOpts)
}

// Symbol is a free data retrieval call binding the contract method 0x95d89b41.
//
// Solidity: function symbol() view returns(string tokenSymbol)
func (_LinkTokenInterface *LinkTokenInterfaceCallerSession) Symbol() (string, error) {
	return _LinkTokenInterface.Contract.Symbol(&_LinkTokenInterface.CallOpts)
}

// TotalSupply is a free data retrieval call binding the contract method 0x18160ddd.
//
// Solidity: function totalSupply() view returns(uint256 totalTokensIssued)
func (_LinkTokenInterface *LinkTokenInterfaceCaller) TotalSupply(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _LinkTokenInterface.contract.Call(opts, &out, "totalSupply")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// TotalSupply is a free data retrieval call binding the contract method 0x18160ddd.
//
// Solidity: function totalSupply() view returns(uint256 totalTokensIssued)
func (_LinkTokenInterface *LinkTokenInterfaceSession) TotalSupply() (*big.Int, error) {
	return _LinkTokenInterface.Contract.TotalSupply(&_LinkTokenInterface.CallOpts)
}

// TotalSupply is a free data retrieval call binding the contract method 0x18160ddd.
//
// Solidity: function totalSupply() view returns(uint256 totalTokensIssued)
func (_LinkTokenInterface *LinkTokenInterfaceCallerSession) TotalSupply() (*big.Int, error) {
	return _LinkTokenInterface.Contract.TotalSupply(&_LinkTokenInterface.CallOpts)
}

// Approve is a paid mutator transaction binding the contract method 0x095ea7b3.
//
// Solidity: function approve(address spender, uint256 value) returns(bool success)
func (_LinkTokenInterface *LinkTokenInterfaceTransactor) Approve(opts *bind.TransactOpts, spender common.Address, value *big.Int) (*types.Transaction, error) {
	return _LinkTokenInterface.contract.Transact(opts, "approve", spender, value)
}

// Approve is a paid mutator transaction binding the contract method 0x095ea7b3.
//
// Solidity: function approve(address spender, uint256 value) returns(bool success)
func (_LinkTokenInterface *LinkTokenInterfaceSession) Approve(spender common.Address, value *big.Int) (*types.Transaction, error) {
	return _LinkTokenInterface.Contract.Approve(&_LinkTokenInterface.TransactOpts, spender, value)
}

// Approve is a paid mutator transaction binding the contract method 0x095ea7b3.
//
// Solidity: function approve(address spender, uint256 value) returns(bool success)
func (_LinkTokenInterface *LinkTokenInterfaceTransactorSession) Approve(spender common.Address, value *big.Int) (*types.Transaction, error) {
	return _LinkTokenInterface.Contract.Approve(&_LinkTokenInterface.TransactOpts, spender, value)
}

// DecreaseApproval is a paid mutator transaction binding the contract method 0x66188463.
//
// Solidity: function decreaseApproval(address spender, uint256 addedValue) returns(bool success)
func (_LinkTokenInterface *LinkTokenInterfaceTransactor) DecreaseApproval(opts *bind.TransactOpts, spender common.Address, addedValue *big.Int) (*types.Transaction, error) {
	return _LinkTokenInterface.contract.Transact(opts, "decreaseApproval", spender, addedValue)
}

// DecreaseApproval is a paid mutator transaction binding the contract method 0x66188463.
//
// Solidity: function decreaseApproval(address spender, uint256 addedValue) returns(bool success)
func (_LinkTokenInterface *LinkTokenInterfaceSession) DecreaseApproval(spender common.Address, addedValue *big.Int) (*types.Transaction, error) {
	return _LinkTokenInterface.Contract.DecreaseApproval(&_LinkTokenInterface.TransactOpts, spender, addedValue)
}

// DecreaseApproval is a paid mutator transaction binding the contract method 0x66188463.
//
// Solidity: function decreaseApproval(address spender, uint256 addedValue) returns(bool success)
func (_LinkTokenInterface *LinkTokenInterfaceTransactorSession) DecreaseApproval(spender common.Address, addedValue *big.Int) (*types.Transaction, error) {
	return _LinkTokenInterface.Contract.DecreaseApproval(&_LinkTokenInterface.TransactOpts, spender, addedValue)
}

// IncreaseApproval is a paid mutator transaction binding the contract method 0xd73dd623.
//
// Solidity: function increaseApproval(address spender, uint256 subtractedValue) returns()
func (_LinkTokenInterface *LinkTokenInterfaceTransactor) IncreaseApproval(opts *bind.TransactOpts, spender common.Address, subtractedValue *big.Int) (*types.Transaction, error) {
	return _LinkTokenInterface.contract.Transact(opts, "increaseApproval", spender, subtractedValue)
}

// IncreaseApproval is a paid mutator transaction binding the contract method 0xd73dd623.
//
// Solidity: function increaseApproval(address spender, uint256 subtractedValue) returns()
func (_LinkTokenInterface *LinkTokenInterfaceSession) IncreaseApproval(spender common.Address, subtractedValue *big.Int) (*types.Transaction, error) {
	return _LinkTokenInterface.Contract.IncreaseApproval(&_LinkTokenInterface.TransactOpts, spender, subtractedValue)
}

// IncreaseApproval is a paid mutator transaction binding the contract method 0xd73dd623.
//
// Solidity: function increaseApproval(address spender, uint256 subtractedValue) returns()
func (_LinkTokenInterface *LinkTokenInterfaceTransactorSession) IncreaseApproval(spender common.Address, subtractedValue *big.Int) (*types.Transaction, error) {
	return _LinkTokenInterface.Contract.IncreaseApproval(&_LinkTokenInterface.TransactOpts, spender, subtractedValue)
}

// Transfer is a paid mutator transaction binding the contract method 0xa9059cbb.
//
// Solidity: function transfer(address to, uint256 value) returns(bool success)
func (_LinkTokenInterface *LinkTokenInterfaceTransactor) Transfer(opts *bind.TransactOpts, to common.Address, value *big.Int) (*types.Transaction, error) {
	return _LinkTokenInterface.contract.Transact(opts, "transfer", to, value)
}

// Transfer is a paid mutator transaction binding the contract method 0xa9059cbb.
//
// Solidity: function transfer(address to, uint256 value) returns(bool success)
func (_LinkTokenInterface *LinkTokenInterfaceSession) Transfer(to common.Address, value *big.Int) (*types.Transaction, error) {
	return _LinkTokenInterface.Contract.Transfer(&_LinkTokenInterface.TransactOpts, to, value)
}

// Transfer is a paid mutator transaction binding the contract method 0xa9059cbb.
//
// Solidity: function transfer(address to, uint256 value) returns(bool success)
func (_LinkTokenInterface *LinkTokenInterfaceTransactorSession) Transfer(to common.Address, value *big.Int) (*types.Transaction, error) {
	return _LinkTokenInterface.Contract.Transfer(&_LinkTokenInterface.TransactOpts, to, value)
}

// TransferAndCall is a paid mutator transaction binding the contract method 0x4000aea0.
//
// Solidity: function transferAndCall(address to, uint256 value, bytes data) returns(bool success)
func (_LinkTokenInterface *LinkTokenInterfaceTransactor) TransferAndCall(opts *bind.TransactOpts, to common.Address, value *big.Int, data []byte) (*types.Transaction, error) {
	return _LinkTokenInterface.contract.Transact(opts, "transferAndCall", to, value, data)
}

// TransferAndCall is a paid mutator transaction binding the contract method 0x4000aea0.
//
// Solidity: function transferAndCall(address to, uint256 value, bytes data) returns(bool success)
func (_LinkTokenInterface *LinkTokenInterfaceSession) TransferAndCall(to common.Address, value *big.Int, data []byte) (*types.Transaction, error) {
	return _LinkTokenInterface.Contract.TransferAndCall(&_LinkTokenInterface.TransactOpts, to, value, data)
}

// TransferAndCall is a paid mutator transaction binding the contract method 0x4000aea0.
//
// Solidity: function transferAndCall(address to, uint256 value, bytes data) returns(bool success)
func (_LinkTokenInterface *LinkTokenInterfaceTransactorSession) TransferAndCall(to common.Address, value *big.Int, data []byte) (*types.Transaction, error) {
	return _LinkTokenInterface.Contract.TransferAndCall(&_LinkTokenInterface.TransactOpts, to, value, data)
}

// TransferFrom is a paid mutator transaction binding the contract method 0x23b872dd.
//
// Solidity: function transferFrom(address from, address to, uint256 value) returns(bool success)
func (_LinkTokenInterface *LinkTokenInterfaceTransactor) TransferFrom(opts *bind.TransactOpts, from common.Address, to common.Address, value *big.Int) (*types.Transaction, error) {
	return _LinkTokenInterface.contract.Transact(opts, "transferFrom", from, to, value)
}

// TransferFrom is a paid mutator transaction binding the contract method 0x23b872dd.
//
// Solidity: function transferFrom(address from, address to, uint256 value) returns(bool success)
func (_LinkTokenInterface *LinkTokenInterfaceSession) TransferFrom(from common.Address, to common.Address, value *big.Int) (*types.Transaction, error) {
	return _LinkTokenInterface.Contract.TransferFrom(&_LinkTokenInterface.TransactOpts, from, to, value)
}

// TransferFrom is a paid mutator transaction binding the contract method 0x23b872dd.
//
// Solidity: function transferFrom(address from, address to, uint256 value) returns(bool success)
func (_LinkTokenInterface *LinkTokenInterfaceTransactorSession) TransferFrom(from common.Address, to common.Address, value *big.Int) (*types.Transaction, error) {
	return _LinkTokenInterface.Contract.TransferFrom(&_LinkTokenInterface.TransactOpts, from, to, value)
}

// OwnableMetaData contains all meta data concerning the Ownable contract.
var OwnableMetaData = &bind.MetaData{
	ABI: "[{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"previousOwner\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"renounceOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
}

// OwnableABI is the input ABI used to generate the binding from.
// Deprecated: Use OwnableMetaData.ABI instead.
var OwnableABI = OwnableMetaData.ABI

// Ownable is an auto generated Go binding around an Ethereum contract.
type Ownable struct {
	OwnableCaller     // Read-only binding to the contract
	OwnableTransactor // Write-only binding to the contract
	OwnableFilterer   // Log filterer for contract events
}

// OwnableCaller is an auto generated read-only Go binding around an Ethereum contract.
type OwnableCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// OwnableTransactor is an auto generated write-only Go binding around an Ethereum contract.
type OwnableTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// OwnableFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type OwnableFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// OwnableSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type OwnableSession struct {
	Contract     *Ownable          // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// OwnableCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type OwnableCallerSession struct {
	Contract *OwnableCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts  // Call options to use throughout this session
}

// OwnableTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type OwnableTransactorSession struct {
	Contract     *OwnableTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts  // Transaction auth options to use throughout this session
}

// OwnableRaw is an auto generated low-level Go binding around an Ethereum contract.
type OwnableRaw struct {
	Contract *Ownable // Generic contract binding to access the raw methods on
}

// OwnableCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type OwnableCallerRaw struct {
	Contract *OwnableCaller // Generic read-only contract binding to access the raw methods on
}

// OwnableTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type OwnableTransactorRaw struct {
	Contract *OwnableTransactor // Generic write-only contract binding to access the raw methods on
}

// NewOwnable creates a new instance of Ownable, bound to a specific deployed contract.
func NewOwnable(address common.Address, backend bind.ContractBackend) (*Ownable, error) {
	contract, err := bindOwnable(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &Ownable{OwnableCaller: OwnableCaller{contract: contract}, OwnableTransactor: OwnableTransactor{contract: contract}, OwnableFilterer: OwnableFilterer{contract: contract}}, nil
}

// NewOwnableCaller creates a new read-only instance of Ownable, bound to a specific deployed contract.
func NewOwnableCaller(address common.Address, caller bind.ContractCaller) (*OwnableCaller, error) {
	contract, err := bindOwnable(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &OwnableCaller{contract: contract}, nil
}

// NewOwnableTransactor creates a new write-only instance of Ownable, bound to a specific deployed contract.
func NewOwnableTransactor(address common.Address, transactor bind.ContractTransactor) (*OwnableTransactor, error) {
	contract, err := bindOwnable(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &OwnableTransactor{contract: contract}, nil
}

// NewOwnableFilterer creates a new log filterer instance of Ownable, bound to a specific deployed contract.
func NewOwnableFilterer(address common.Address, filterer bind.ContractFilterer) (*OwnableFilterer, error) {
	contract, err := bindOwnable(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &OwnableFilterer{contract: contract}, nil
}

// bindOwnable binds a generic wrapper to an already deployed contract.
func bindOwnable(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(OwnableABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Ownable *OwnableRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Ownable.Contract.OwnableCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Ownable *OwnableRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Ownable.Contract.OwnableTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Ownable *OwnableRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Ownable.Contract.OwnableTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Ownable *OwnableCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Ownable.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Ownable *OwnableTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Ownable.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Ownable *OwnableTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Ownable.Contract.contract.Transact(opts, method, params...)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_Ownable *OwnableCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _Ownable.contract.Call(opts, &out, "owner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_Ownable *OwnableSession) Owner() (common.Address, error) {
	return _Ownable.Contract.Owner(&_Ownable.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_Ownable *OwnableCallerSession) Owner() (common.Address, error) {
	return _Ownable.Contract.Owner(&_Ownable.CallOpts)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_Ownable *OwnableTransactor) RenounceOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Ownable.contract.Transact(opts, "renounceOwnership")
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_Ownable *OwnableSession) RenounceOwnership() (*types.Transaction, error) {
	return _Ownable.Contract.RenounceOwnership(&_Ownable.TransactOpts)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_Ownable *OwnableTransactorSession) RenounceOwnership() (*types.Transaction, error) {
	return _Ownable.Contract.RenounceOwnership(&_Ownable.TransactOpts)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_Ownable *OwnableTransactor) TransferOwnership(opts *bind.TransactOpts, newOwner common.Address) (*types.Transaction, error) {
	return _Ownable.contract.Transact(opts, "transferOwnership", newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_Ownable *OwnableSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _Ownable.Contract.TransferOwnership(&_Ownable.TransactOpts, newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_Ownable *OwnableTransactorSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _Ownable.Contract.TransferOwnership(&_Ownable.TransactOpts, newOwner)
}

// OwnableOwnershipTransferredIterator is returned from FilterOwnershipTransferred and is used to iterate over the raw logs and unpacked data for OwnershipTransferred events raised by the Ownable contract.
type OwnableOwnershipTransferredIterator struct {
	Event *OwnableOwnershipTransferred // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *OwnableOwnershipTransferredIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(OwnableOwnershipTransferred)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(OwnableOwnershipTransferred)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *OwnableOwnershipTransferredIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *OwnableOwnershipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// OwnableOwnershipTransferred represents a OwnershipTransferred event raised by the Ownable contract.
type OwnableOwnershipTransferred struct {
	PreviousOwner common.Address
	NewOwner      common.Address
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterOwnershipTransferred is a free log retrieval operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_Ownable *OwnableFilterer) FilterOwnershipTransferred(opts *bind.FilterOpts, previousOwner []common.Address, newOwner []common.Address) (*OwnableOwnershipTransferredIterator, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _Ownable.contract.FilterLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return &OwnableOwnershipTransferredIterator{contract: _Ownable.contract, event: "OwnershipTransferred", logs: logs, sub: sub}, nil
}

// WatchOwnershipTransferred is a free log subscription operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_Ownable *OwnableFilterer) WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *OwnableOwnershipTransferred, previousOwner []common.Address, newOwner []common.Address) (event.Subscription, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _Ownable.contract.WatchLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(OwnableOwnershipTransferred)
				if err := _Ownable.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseOwnershipTransferred is a log parse operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_Ownable *OwnableFilterer) ParseOwnershipTransferred(log types.Log) (*OwnableOwnershipTransferred, error) {
	event := new(OwnableOwnershipTransferred)
	if err := _Ownable.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// SimulatedChainlinkMetaData contains all meta data concerning the SimulatedChainlink contract.
var SimulatedChainlinkMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[],\"name\":\"linkToken\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"vrfCoordinator\",\"outputs\":[{\"internalType\":\"contractSimulatedVRFCoordinator\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"}]",
	Bin: "0x60c060405234801561001057600080fd5b5060405161001d906103bb565b604051809103906000f080158015610039573d6000803e3d6000fd5b5073ffffffffffffffffffffffffffffffffffffffff1660808173ffffffffffffffffffffffffffffffffffffffff1681525050604051610079906103c8565b604051809103906000f080158015610095573d6000803e3d6000fd5b5073ffffffffffffffffffffffffffffffffffffffff1660a08173ffffffffffffffffffffffffffffffffffffffff168152505060a05173ffffffffffffffffffffffffffffffffffffffff1663f2fde38b336040518263ffffffff1660e01b81526004016101049190610416565b600060405180830381600087803b15801561011e57600080fd5b505af1158015610132573d6000803e3d6000fd5b5050505061014861023d60201b6100bf1760201c565b73ffffffffffffffffffffffffffffffffffffffff1660805173ffffffffffffffffffffffffffffffffffffffff16146101b7576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004016101ae906104b4565b60405180910390fd5b6101c96102fc60201b61017e1760201c565b73ffffffffffffffffffffffffffffffffffffffff1660a05173ffffffffffffffffffffffffffffffffffffffff1614610238576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040161022f90610546565b60405180910390fd5b610566565b60004660018114610270576004811461028c57608981146102a8576201388181146102c45761053981146102e0576102f8565b73514910771af9ca656af840dff83e8264ecf986ca91506102f8565b7301be23585060835e02b77ef475b0cc51aa1e070991506102f8565b73b0897686c545045afc77cf20ec7a532e3120e0f191506102f8565b73326c977e6efc84e512bb9c30f76e30c160ed06fb91506102f8565b7355b04d60213bcfdc383a6411ceff3f759ab366d691505b5090565b6000466001811461032f576004811461034b57608981146103675762013881811461038357610539811461039f576103b7565b73f0d54349addcf704f77ae15b96510dea15cb795291506103b7565b73b3dccb4cf7a26f6cf6b120cf5a73875b7bbc655b91506103b7565b733d2341adb2d31f1c5530cdc622016af293177ae091506103b7565b738c7382f9d8f56b33781fe506e897a4f1e2d1725591506103b7565b735ffd760b2b48575f3869722cd816d8b3f94ddb4891505b5090565b61158d806108d283390190565b610b6b80611e5f83390190565b600073ffffffffffffffffffffffffffffffffffffffff82169050919050565b6000610400826103d5565b9050919050565b610410816103f5565b82525050565b600060208201905061042b6000830184610407565b92915050565b600082825260208201905092915050565b7f53696d756c61746564436861696e4c696e6b3a20756e6578706563746564204c60008201527f494e4b20746f6b656e2061646472657373000000000000000000000000000000602082015250565b600061049e603183610431565b91506104a982610442565b604082019050919050565b600060208201905081810360008301526104cd81610491565b9050919050565b7f53696d756c61746564436861696e4c696e6b3a20756e6578706563746564205660008201527f5246436f6f7264696e61746f7220616464726573730000000000000000000000602082015250565b6000610530603583610431565b915061053b826104d4565b604082019050919050565b6000602082019050818103600083015261055f81610523565b9050919050565b60805160a0516103496105896000396000609d01526000607901526103496000f3fe608060405234801561001057600080fd5b50600436106100365760003560e01c806357970e931461003b578063a3e56fa814610059575b600080fd5b610043610077565b604051610050919061027e565b60405180910390f35b61006161009b565b60405161006e91906102f8565b60405180910390f35b7f000000000000000000000000000000000000000000000000000000000000000081565b7f000000000000000000000000000000000000000000000000000000000000000081565b600046600181146100f2576004811461010e576089811461012a576201388181146101465761053981146101625761017a565b73514910771af9ca656af840dff83e8264ecf986ca915061017a565b7301be23585060835e02b77ef475b0cc51aa1e0709915061017a565b73b0897686c545045afc77cf20ec7a532e3120e0f1915061017a565b73326c977e6efc84e512bb9c30f76e30c160ed06fb915061017a565b7355b04d60213bcfdc383a6411ceff3f759ab366d691505b5090565b600046600181146101b157600481146101cd57608981146101e95762013881811461020557610539811461022157610239565b73f0d54349addcf704f77ae15b96510dea15cb79529150610239565b73b3dccb4cf7a26f6cf6b120cf5a73875b7bbc655b9150610239565b733d2341adb2d31f1c5530cdc622016af293177ae09150610239565b738c7382f9d8f56b33781fe506e897a4f1e2d172559150610239565b735ffd760b2b48575f3869722cd816d8b3f94ddb4891505b5090565b600073ffffffffffffffffffffffffffffffffffffffff82169050919050565b60006102688261023d565b9050919050565b6102788161025d565b82525050565b6000602082019050610293600083018461026f565b92915050565b6000819050919050565b60006102be6102b96102b48461023d565b610299565b61023d565b9050919050565b60006102d0826102a3565b9050919050565b60006102e2826102c5565b9050919050565b6102f2816102d7565b82525050565b600060208201905061030d60008301846102e9565b9291505056fea26469706673582212202b60bcdfe4fc616963f180dc0e7f3762452e382cc4cee217702b2afbf3ef6de764736f6c634300080b003360a060405234801561001057600080fd5b503373ffffffffffffffffffffffffffffffffffffffff1660808173ffffffffffffffffffffffffffffffffffffffff168152505060805161152d6100606000396000610101015261152d6000f3fe608060405234801561001057600080fd5b506004361061004c5760003560e01c80634000aea01461005157806370a08231146100815780637b56c2b2146100b1578063a9059cbb146100cd575b600080fd5b61006b60048036038101906100669190610baa565b6100fd565b6040516100789190610c39565b60405180910390f35b61009b60048036038101906100969190610c54565b610327565b6040516100a89190610c90565b60405180910390f35b6100cb60048036038101906100c69190610cab565b61033f565b005b6100e760048036038101906100e29190610cab565b610398565b6040516100f49190610c39565b60405180910390f35b60007f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff1663a3e56fa86040518163ffffffff1660e01b8152600401602060405180830381865afa15801561016a573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019061018e9190610d29565b73ffffffffffffffffffffffffffffffffffffffff168573ffffffffffffffffffffffffffffffffffffffff16146101fb576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004016101f290610dd9565b60405180910390fd5b6102036105ba565b8414610244576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040161023b90610e6b565b60405180910390fd5b61024c61062a565b600060405160200161025f929190610ed7565b604051602081830303815290604052805190602001208383604051610285929190610f42565b6040518091039020146102cd576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004016102c490610fcd565b60405180910390fd5b6102d78585610398565b503373ffffffffffffffffffffffffffffffffffffffff167f6871e329198b319adcce5196458d56b81939284c334a427d776ebdc356dd5acb60405160405180910390a260019050949350505050565b60006020528060005260406000206000915090505481565b806000808473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff168152602001908152602001600020600082825461038d919061101c565b925050819055505050565b6000816000803373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002054101561040660143373ffffffffffffffffffffffffffffffffffffffff1661070a90919063ffffffff16565b61045e655af3107a40006000803373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020016000205461045991906110a1565b610946565b610478655af3107a40008661047391906110a1565b610946565b6104a260148873ffffffffffffffffffffffffffffffffffffffff1661070a90919063ffffffff16565b6040516020016104b59493929190611230565b60405160208183030381529060405290610505576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004016104fc91906112d9565b60405180910390fd5b50816000803373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff168152602001908152602001600020600082825461055491906112fb565b92505081905550816000808573ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060008282546105a9919061101c565b925050819055506001905092915050565b60008046905060018114806105d0575061053981145b156105e657671bc16d674ec80000915050610627565b60898114806105f757506201388181145b1561060b57655af3107a4000915050610627565b60048114156106255767016345785d8a0000915050610627565b505b90565b6000466001811461065d576004811461068557608981146106ad576201388181146106d55761053981146106fd57610706565b7faa77729d3466ca35ae8d28b3bbac7cc36a5031efdc430821c02bc31a238af4459150610706565b7f2ed0feb3e7fd2022120aa84fab1945545a9f2ffc9076fd6156fa96eaff4c13119150610706565b7ff86195cf7690c55907b2b611ebb7343a6f649bff128701cc542f0569e2c549da9150610706565b7f6e75b569a01ef56d18cab6a8e71e6600d6ce853834d4a5748b720d06f878b3a49150610706565b60026113372091505b5090565b60606000600283600261071d919061132f565b610727919061101c565b67ffffffffffffffff8111156107405761073f611389565b5b6040519080825280601f01601f1916602001820160405280156107725781602001600182028036833780820191505090505b5090507f3000000000000000000000000000000000000000000000000000000000000000816000815181106107aa576107a96113b8565b5b60200101907effffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff1916908160001a9053507f78000000000000000000000000000000000000000000000000000000000000008160018151811061080e5761080d6113b8565b5b60200101907effffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff1916908160001a9053506000600184600261084e919061132f565b610858919061101c565b90505b60018111156108f8577f3031323334353637383961626364656600000000000000000000000000000000600f86166010811061089a576108996113b8565b5b1a60f81b8282815181106108b1576108b06113b8565b5b60200101907effffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff1916908160001a905350600485901c9450806108f1906113e7565b905061085b565b506000841461093c576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004016109339061145d565b60405180910390fd5b8091505092915050565b6060600082141561098e576040518060400160405280600181526020017f30000000000000000000000000000000000000000000000000000000000000008152509050610aa2565b600082905060005b600082146109c05780806109a99061147d565b915050600a826109b991906110a1565b9150610996565b60008167ffffffffffffffff8111156109dc576109db611389565b5b6040519080825280601f01601f191660200182016040528015610a0e5781602001600182028036833780820191505090505b5090505b60008514610a9b57600182610a2791906112fb565b9150600a85610a3691906114c6565b6030610a42919061101c565b60f81b818381518110610a5857610a576113b8565b5b60200101907effffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff1916908160001a905350600a85610a9491906110a1565b9450610a12565b8093505050505b919050565b600080fd5b600080fd5b600073ffffffffffffffffffffffffffffffffffffffff82169050919050565b6000610adc82610ab1565b9050919050565b610aec81610ad1565b8114610af757600080fd5b50565b600081359050610b0981610ae3565b92915050565b6000819050919050565b610b2281610b0f565b8114610b2d57600080fd5b50565b600081359050610b3f81610b19565b92915050565b600080fd5b600080fd5b600080fd5b60008083601f840112610b6a57610b69610b45565b5b8235905067ffffffffffffffff811115610b8757610b86610b4a565b5b602083019150836001820283011115610ba357610ba2610b4f565b5b9250929050565b60008060008060608587031215610bc457610bc3610aa7565b5b6000610bd287828801610afa565b9450506020610be387828801610b30565b935050604085013567ffffffffffffffff811115610c0457610c03610aac565b5b610c1087828801610b54565b925092505092959194509250565b60008115159050919050565b610c3381610c1e565b82525050565b6000602082019050610c4e6000830184610c2a565b92915050565b600060208284031215610c6a57610c69610aa7565b5b6000610c7884828501610afa565b91505092915050565b610c8a81610b0f565b82525050565b6000602082019050610ca56000830184610c81565b92915050565b60008060408385031215610cc257610cc1610aa7565b5b6000610cd085828601610afa565b9250506020610ce185828601610b30565b9150509250929050565b6000610cf682610ad1565b9050919050565b610d0681610ceb565b8114610d1157600080fd5b50565b600081519050610d2381610cfd565b92915050565b600060208284031215610d3f57610d3e610aa7565b5b6000610d4d84828501610d14565b91505092915050565b600082825260208201905092915050565b7f53696d756c61746564204c494e4b20746f6b656e3a20696e636f72726563742060008201527f56524620436f6f7264696e61746f720000000000000000000000000000000000602082015250565b6000610dc3602f83610d56565b9150610dce82610d67565b604082019050919050565b60006020820190508181036000830152610df281610db6565b9050919050565b7f53696d756c61746564204c494e4b20746f6b656e3a20696e636f72726563742060008201527f66656520666f7220565246000000000000000000000000000000000000000000602082015250565b6000610e55602b83610d56565b9150610e6082610df9565b604082019050919050565b60006020820190508181036000830152610e8481610e48565b9050919050565b6000819050919050565b6000819050919050565b610eb0610eab82610e8b565b610e95565b82525050565b6000819050919050565b610ed1610ecc82610b0f565b610eb6565b82525050565b6000610ee38285610e9f565b602082019150610ef38284610ec0565b6020820191508190509392505050565b600081905092915050565b82818337600083830152505050565b6000610f298385610f03565b9350610f36838584610f0e565b82840190509392505050565b6000610f4f828486610f1d565b91508190509392505050565b7f53696d756c61746564204c494e4b20746f6b656e3a20696e76616c696420646160008201527f7461000000000000000000000000000000000000000000000000000000000000602082015250565b6000610fb7602283610d56565b9150610fc282610f5b565b604082019050919050565b60006020820190508181036000830152610fe681610faa565b9050919050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052601160045260246000fd5b600061102782610b0f565b915061103283610b0f565b9250827fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff0382111561106757611066610fed565b5b828201905092915050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052601260045260246000fd5b60006110ac82610b0f565b91506110b783610b0f565b9250826110c7576110c6611072565b5b828204905092915050565b600081519050919050565b600081905092915050565b60005b838110156111065780820151818401526020810190506110eb565b83811115611115576000848401525b50505050565b6000611126826110d2565b61113081856110dd565b93506111408185602086016110e8565b80840191505092915050565b7f2068617320696e73756666696369656e742062616c616e636520000000000000600082015250565b6000611182601a836110dd565b915061118d8261114c565b601a82019050919050565b7f65313420746f207472616e736665722000000000000000000000000000000000600082015250565b60006111ce6010836110dd565b91506111d982611198565b601082019050919050565b7f65313420746f2000000000000000000000000000000000000000000000000000600082015250565b600061121a6007836110dd565b9150611225826111e4565b600782019050919050565b600061123c828761111b565b915061124782611175565b9150611253828661111b565b915061125e826111c1565b915061126a828561111b565b91506112758261120d565b9150611281828461111b565b915081905095945050505050565b6000601f19601f8301169050919050565b60006112ab826110d2565b6112b58185610d56565b93506112c58185602086016110e8565b6112ce8161128f565b840191505092915050565b600060208201905081810360008301526112f381846112a0565b905092915050565b600061130682610b0f565b915061131183610b0f565b92508282101561132457611323610fed565b5b828203905092915050565b600061133a82610b0f565b915061134583610b0f565b9250817fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff048311821515161561137e5761137d610fed565b5b828202905092915050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052604160045260246000fd5b7f4e487b7100000000000000000000000000000000000000000000000000000000600052603260045260246000fd5b60006113f282610b0f565b9150600082141561140657611405610fed565b5b600182039050919050565b7f537472696e67733a20686578206c656e67746820696e73756666696369656e74600082015250565b6000611447602083610d56565b915061145282611411565b602082019050919050565b600060208201905081810360008301526114768161143a565b9050919050565b600061148882610b0f565b91507fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff8214156114bb576114ba610fed565b5b600182019050919050565b60006114d182610b0f565b91506114dc83610b0f565b9250826114ec576114eb611072565b5b82820690509291505056fea26469706673582212200d4a10c570d0a4850ec1aacd324bd9eef115d804c400bb9092216198cacc3e8764736f6c634300080b0033608060405234801561001057600080fd5b5061002d61002261003260201b60201c565b61003a60201b60201c565b6100fe565b600033905090565b60008060009054906101000a900473ffffffffffffffffffffffffffffffffffffffff169050816000806101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff1602179055508173ffffffffffffffffffffffffffffffffffffffff168173ffffffffffffffffffffffffffffffffffffffff167f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e060405160405180910390a35050565b610a5e8061010d6000396000f3fe608060405234801561001057600080fd5b506004361061004c5760003560e01c806344e03a1f14610051578063715018a61461006d5780638da5cb5b14610077578063f2fde38b14610095575b600080fd5b61006b600480360381019061006691906106c8565b6100b1565b005b61007561028f565b005b61007f610317565b60405161008c9190610704565b60405180910390f35b6100af60048036038101906100aa919061074b565b610340565b005b6100b9610438565b73ffffffffffffffffffffffffffffffffffffffff166100d7610317565b73ffffffffffffffffffffffffffffffffffffffff161461012d576040517f08c379a0000000000000000000000000000000000000000000000000000000008152600401610124906107d5565b60405180910390fd5b6000819050600061018761013f610440565b600084600160008773ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002054610520565b9050600160008373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060008154809291906101d99061082e565b919050555060006101f16101eb610440565b8361055c565b90508373ffffffffffffffffffffffffffffffffffffffff166394985ddd828360405160200161022191906108a2565b6040516020818303038152906040528051906020012060001c6040518363ffffffff1660e01b81526004016102579291906108db565b600060405180830381600087803b15801561027157600080fd5b505af1158015610285573d6000803e3d6000fd5b5050505050505050565b610297610438565b73ffffffffffffffffffffffffffffffffffffffff166102b5610317565b73ffffffffffffffffffffffffffffffffffffffff161461030b576040517f08c379a0000000000000000000000000000000000000000000000000000000008152600401610302906107d5565b60405180910390fd5b610315600061058f565b565b60008060009054906101000a900473ffffffffffffffffffffffffffffffffffffffff16905090565b610348610438565b73ffffffffffffffffffffffffffffffffffffffff16610366610317565b73ffffffffffffffffffffffffffffffffffffffff16146103bc576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004016103b3906107d5565b60405180910390fd5b600073ffffffffffffffffffffffffffffffffffffffff168173ffffffffffffffffffffffffffffffffffffffff16141561042c576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040161042390610976565b60405180910390fd5b6104358161058f565b50565b600033905090565b60004660018114610473576004811461049b57608981146104c3576201388181146104eb5761053981146105135761051c565b7faa77729d3466ca35ae8d28b3bbac7cc36a5031efdc430821c02bc31a238af445915061051c565b7f2ed0feb3e7fd2022120aa84fab1945545a9f2ffc9076fd6156fa96eaff4c1311915061051c565b7ff86195cf7690c55907b2b611ebb7343a6f649bff128701cc542f0569e2c549da915061051c565b7f6e75b569a01ef56d18cab6a8e71e6600d6ce853834d4a5748b720d06f878b3a4915061051c565b60026113372091505b5090565b6000848484846040516020016105399493929190610996565b6040516020818303038152906040528051906020012060001c9050949350505050565b600082826040516020016105719291906109fc565b60405160208183030381529060405280519060200120905092915050565b60008060009054906101000a900473ffffffffffffffffffffffffffffffffffffffff169050816000806101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff1602179055508173ffffffffffffffffffffffffffffffffffffffff168173ffffffffffffffffffffffffffffffffffffffff167f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e060405160405180910390a35050565b600080fd5b600073ffffffffffffffffffffffffffffffffffffffff82169050919050565b600061068382610658565b9050919050565b600061069582610678565b9050919050565b6106a58161068a565b81146106b057600080fd5b50565b6000813590506106c28161069c565b92915050565b6000602082840312156106de576106dd610653565b5b60006106ec848285016106b3565b91505092915050565b6106fe81610678565b82525050565b600060208201905061071960008301846106f5565b92915050565b61072881610678565b811461073357600080fd5b50565b6000813590506107458161071f565b92915050565b60006020828403121561076157610760610653565b5b600061076f84828501610736565b91505092915050565b600082825260208201905092915050565b7f4f776e61626c653a2063616c6c6572206973206e6f7420746865206f776e6572600082015250565b60006107bf602083610778565b91506107ca82610789565b602082019050919050565b600060208201905081810360008301526107ee816107b2565b9050919050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052601160045260246000fd5b6000819050919050565b600061083982610824565b91507fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff82141561086c5761086b6107f5565b5b600182019050919050565b6000819050919050565b6000819050919050565b61089c61089782610877565b610881565b82525050565b60006108ae828461088b565b60208201915081905092915050565b6108c681610877565b82525050565b6108d581610824565b82525050565b60006040820190506108f060008301856108bd565b6108fd60208301846108cc565b9392505050565b7f4f776e61626c653a206e6577206f776e657220697320746865207a65726f206160008201527f6464726573730000000000000000000000000000000000000000000000000000602082015250565b6000610960602683610778565b915061096b82610904565b604082019050919050565b6000602082019050818103600083015261098f81610953565b9050919050565b60006080820190506109ab60008301876108bd565b6109b860208301866108cc565b6109c560408301856106f5565b6109d260608301846108cc565b95945050505050565b6000819050919050565b6109f66109f182610824565b6109db565b82525050565b6000610a08828561088b565b602082019150610a1882846109e5565b602082019150819050939250505056fea26469706673582212206529d122f2f8d2b0379431c5cdbeb3deb5df96030d0cefbb709af7a4d457279464736f6c634300080b0033",
}

// SimulatedChainlinkABI is the input ABI used to generate the binding from.
// Deprecated: Use SimulatedChainlinkMetaData.ABI instead.
var SimulatedChainlinkABI = SimulatedChainlinkMetaData.ABI

// SimulatedChainlinkBin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use SimulatedChainlinkMetaData.Bin instead.
var SimulatedChainlinkBin = SimulatedChainlinkMetaData.Bin

// DeploySimulatedChainlink deploys a new Ethereum contract, binding an instance of SimulatedChainlink to it.
func DeploySimulatedChainlink(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *SimulatedChainlink, error) {
	parsed, err := SimulatedChainlinkMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(SimulatedChainlinkBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &SimulatedChainlink{SimulatedChainlinkCaller: SimulatedChainlinkCaller{contract: contract}, SimulatedChainlinkTransactor: SimulatedChainlinkTransactor{contract: contract}, SimulatedChainlinkFilterer: SimulatedChainlinkFilterer{contract: contract}}, nil
}

// SimulatedChainlink is an auto generated Go binding around an Ethereum contract.
type SimulatedChainlink struct {
	SimulatedChainlinkCaller     // Read-only binding to the contract
	SimulatedChainlinkTransactor // Write-only binding to the contract
	SimulatedChainlinkFilterer   // Log filterer for contract events
}

// SimulatedChainlinkCaller is an auto generated read-only Go binding around an Ethereum contract.
type SimulatedChainlinkCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// SimulatedChainlinkTransactor is an auto generated write-only Go binding around an Ethereum contract.
type SimulatedChainlinkTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// SimulatedChainlinkFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type SimulatedChainlinkFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// SimulatedChainlinkSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type SimulatedChainlinkSession struct {
	Contract     *SimulatedChainlink // Generic contract binding to set the session for
	CallOpts     bind.CallOpts       // Call options to use throughout this session
	TransactOpts bind.TransactOpts   // Transaction auth options to use throughout this session
}

// SimulatedChainlinkCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type SimulatedChainlinkCallerSession struct {
	Contract *SimulatedChainlinkCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts             // Call options to use throughout this session
}

// SimulatedChainlinkTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type SimulatedChainlinkTransactorSession struct {
	Contract     *SimulatedChainlinkTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts             // Transaction auth options to use throughout this session
}

// SimulatedChainlinkRaw is an auto generated low-level Go binding around an Ethereum contract.
type SimulatedChainlinkRaw struct {
	Contract *SimulatedChainlink // Generic contract binding to access the raw methods on
}

// SimulatedChainlinkCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type SimulatedChainlinkCallerRaw struct {
	Contract *SimulatedChainlinkCaller // Generic read-only contract binding to access the raw methods on
}

// SimulatedChainlinkTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type SimulatedChainlinkTransactorRaw struct {
	Contract *SimulatedChainlinkTransactor // Generic write-only contract binding to access the raw methods on
}

// NewSimulatedChainlink creates a new instance of SimulatedChainlink, bound to a specific deployed contract.
func NewSimulatedChainlink(address common.Address, backend bind.ContractBackend) (*SimulatedChainlink, error) {
	contract, err := bindSimulatedChainlink(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &SimulatedChainlink{SimulatedChainlinkCaller: SimulatedChainlinkCaller{contract: contract}, SimulatedChainlinkTransactor: SimulatedChainlinkTransactor{contract: contract}, SimulatedChainlinkFilterer: SimulatedChainlinkFilterer{contract: contract}}, nil
}

// NewSimulatedChainlinkCaller creates a new read-only instance of SimulatedChainlink, bound to a specific deployed contract.
func NewSimulatedChainlinkCaller(address common.Address, caller bind.ContractCaller) (*SimulatedChainlinkCaller, error) {
	contract, err := bindSimulatedChainlink(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &SimulatedChainlinkCaller{contract: contract}, nil
}

// NewSimulatedChainlinkTransactor creates a new write-only instance of SimulatedChainlink, bound to a specific deployed contract.
func NewSimulatedChainlinkTransactor(address common.Address, transactor bind.ContractTransactor) (*SimulatedChainlinkTransactor, error) {
	contract, err := bindSimulatedChainlink(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &SimulatedChainlinkTransactor{contract: contract}, nil
}

// NewSimulatedChainlinkFilterer creates a new log filterer instance of SimulatedChainlink, bound to a specific deployed contract.
func NewSimulatedChainlinkFilterer(address common.Address, filterer bind.ContractFilterer) (*SimulatedChainlinkFilterer, error) {
	contract, err := bindSimulatedChainlink(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &SimulatedChainlinkFilterer{contract: contract}, nil
}

// bindSimulatedChainlink binds a generic wrapper to an already deployed contract.
func bindSimulatedChainlink(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(SimulatedChainlinkABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_SimulatedChainlink *SimulatedChainlinkRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _SimulatedChainlink.Contract.SimulatedChainlinkCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_SimulatedChainlink *SimulatedChainlinkRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _SimulatedChainlink.Contract.SimulatedChainlinkTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_SimulatedChainlink *SimulatedChainlinkRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _SimulatedChainlink.Contract.SimulatedChainlinkTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_SimulatedChainlink *SimulatedChainlinkCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _SimulatedChainlink.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_SimulatedChainlink *SimulatedChainlinkTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _SimulatedChainlink.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_SimulatedChainlink *SimulatedChainlinkTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _SimulatedChainlink.Contract.contract.Transact(opts, method, params...)
}

// LinkToken is a free data retrieval call binding the contract method 0x57970e93.
//
// Solidity: function linkToken() view returns(address)
func (_SimulatedChainlink *SimulatedChainlinkCaller) LinkToken(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _SimulatedChainlink.contract.Call(opts, &out, "linkToken")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// LinkToken is a free data retrieval call binding the contract method 0x57970e93.
//
// Solidity: function linkToken() view returns(address)
func (_SimulatedChainlink *SimulatedChainlinkSession) LinkToken() (common.Address, error) {
	return _SimulatedChainlink.Contract.LinkToken(&_SimulatedChainlink.CallOpts)
}

// LinkToken is a free data retrieval call binding the contract method 0x57970e93.
//
// Solidity: function linkToken() view returns(address)
func (_SimulatedChainlink *SimulatedChainlinkCallerSession) LinkToken() (common.Address, error) {
	return _SimulatedChainlink.Contract.LinkToken(&_SimulatedChainlink.CallOpts)
}

// VrfCoordinator is a free data retrieval call binding the contract method 0xa3e56fa8.
//
// Solidity: function vrfCoordinator() view returns(address)
func (_SimulatedChainlink *SimulatedChainlinkCaller) VrfCoordinator(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _SimulatedChainlink.contract.Call(opts, &out, "vrfCoordinator")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// VrfCoordinator is a free data retrieval call binding the contract method 0xa3e56fa8.
//
// Solidity: function vrfCoordinator() view returns(address)
func (_SimulatedChainlink *SimulatedChainlinkSession) VrfCoordinator() (common.Address, error) {
	return _SimulatedChainlink.Contract.VrfCoordinator(&_SimulatedChainlink.CallOpts)
}

// VrfCoordinator is a free data retrieval call binding the contract method 0xa3e56fa8.
//
// Solidity: function vrfCoordinator() view returns(address)
func (_SimulatedChainlink *SimulatedChainlinkCallerSession) VrfCoordinator() (common.Address, error) {
	return _SimulatedChainlink.Contract.VrfCoordinator(&_SimulatedChainlink.CallOpts)
}

// SimulatedLinkTokenMetaData contains all meta data concerning the SimulatedLinkToken contract.
var SimulatedLinkTokenMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"by\",\"type\":\"address\"}],\"name\":\"RandomnessRequested\",\"type\":\"event\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"name\":\"balanceOf\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"recipient\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"faucet\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"recipient\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"transfer\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"data\",\"type\":\"bytes\"}],\"name\":\"transferAndCall\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"success\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
	Bin: "0x60a060405234801561001057600080fd5b503373ffffffffffffffffffffffffffffffffffffffff1660808173ffffffffffffffffffffffffffffffffffffffff168152505060805161152d6100606000396000610101015261152d6000f3fe608060405234801561001057600080fd5b506004361061004c5760003560e01c80634000aea01461005157806370a08231146100815780637b56c2b2146100b1578063a9059cbb146100cd575b600080fd5b61006b60048036038101906100669190610baa565b6100fd565b6040516100789190610c39565b60405180910390f35b61009b60048036038101906100969190610c54565b610327565b6040516100a89190610c90565b60405180910390f35b6100cb60048036038101906100c69190610cab565b61033f565b005b6100e760048036038101906100e29190610cab565b610398565b6040516100f49190610c39565b60405180910390f35b60007f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff1663a3e56fa86040518163ffffffff1660e01b8152600401602060405180830381865afa15801561016a573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019061018e9190610d29565b73ffffffffffffffffffffffffffffffffffffffff168573ffffffffffffffffffffffffffffffffffffffff16146101fb576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004016101f290610dd9565b60405180910390fd5b6102036105ba565b8414610244576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040161023b90610e6b565b60405180910390fd5b61024c61062a565b600060405160200161025f929190610ed7565b604051602081830303815290604052805190602001208383604051610285929190610f42565b6040518091039020146102cd576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004016102c490610fcd565b60405180910390fd5b6102d78585610398565b503373ffffffffffffffffffffffffffffffffffffffff167f6871e329198b319adcce5196458d56b81939284c334a427d776ebdc356dd5acb60405160405180910390a260019050949350505050565b60006020528060005260406000206000915090505481565b806000808473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff168152602001908152602001600020600082825461038d919061101c565b925050819055505050565b6000816000803373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002054101561040660143373ffffffffffffffffffffffffffffffffffffffff1661070a90919063ffffffff16565b61045e655af3107a40006000803373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020016000205461045991906110a1565b610946565b610478655af3107a40008661047391906110a1565b610946565b6104a260148873ffffffffffffffffffffffffffffffffffffffff1661070a90919063ffffffff16565b6040516020016104b59493929190611230565b60405160208183030381529060405290610505576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004016104fc91906112d9565b60405180910390fd5b50816000803373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff168152602001908152602001600020600082825461055491906112fb565b92505081905550816000808573ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060008282546105a9919061101c565b925050819055506001905092915050565b60008046905060018114806105d0575061053981145b156105e657671bc16d674ec80000915050610627565b60898114806105f757506201388181145b1561060b57655af3107a4000915050610627565b60048114156106255767016345785d8a0000915050610627565b505b90565b6000466001811461065d576004811461068557608981146106ad576201388181146106d55761053981146106fd57610706565b7faa77729d3466ca35ae8d28b3bbac7cc36a5031efdc430821c02bc31a238af4459150610706565b7f2ed0feb3e7fd2022120aa84fab1945545a9f2ffc9076fd6156fa96eaff4c13119150610706565b7ff86195cf7690c55907b2b611ebb7343a6f649bff128701cc542f0569e2c549da9150610706565b7f6e75b569a01ef56d18cab6a8e71e6600d6ce853834d4a5748b720d06f878b3a49150610706565b60026113372091505b5090565b60606000600283600261071d919061132f565b610727919061101c565b67ffffffffffffffff8111156107405761073f611389565b5b6040519080825280601f01601f1916602001820160405280156107725781602001600182028036833780820191505090505b5090507f3000000000000000000000000000000000000000000000000000000000000000816000815181106107aa576107a96113b8565b5b60200101907effffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff1916908160001a9053507f78000000000000000000000000000000000000000000000000000000000000008160018151811061080e5761080d6113b8565b5b60200101907effffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff1916908160001a9053506000600184600261084e919061132f565b610858919061101c565b90505b60018111156108f8577f3031323334353637383961626364656600000000000000000000000000000000600f86166010811061089a576108996113b8565b5b1a60f81b8282815181106108b1576108b06113b8565b5b60200101907effffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff1916908160001a905350600485901c9450806108f1906113e7565b905061085b565b506000841461093c576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004016109339061145d565b60405180910390fd5b8091505092915050565b6060600082141561098e576040518060400160405280600181526020017f30000000000000000000000000000000000000000000000000000000000000008152509050610aa2565b600082905060005b600082146109c05780806109a99061147d565b915050600a826109b991906110a1565b9150610996565b60008167ffffffffffffffff8111156109dc576109db611389565b5b6040519080825280601f01601f191660200182016040528015610a0e5781602001600182028036833780820191505090505b5090505b60008514610a9b57600182610a2791906112fb565b9150600a85610a3691906114c6565b6030610a42919061101c565b60f81b818381518110610a5857610a576113b8565b5b60200101907effffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff1916908160001a905350600a85610a9491906110a1565b9450610a12565b8093505050505b919050565b600080fd5b600080fd5b600073ffffffffffffffffffffffffffffffffffffffff82169050919050565b6000610adc82610ab1565b9050919050565b610aec81610ad1565b8114610af757600080fd5b50565b600081359050610b0981610ae3565b92915050565b6000819050919050565b610b2281610b0f565b8114610b2d57600080fd5b50565b600081359050610b3f81610b19565b92915050565b600080fd5b600080fd5b600080fd5b60008083601f840112610b6a57610b69610b45565b5b8235905067ffffffffffffffff811115610b8757610b86610b4a565b5b602083019150836001820283011115610ba357610ba2610b4f565b5b9250929050565b60008060008060608587031215610bc457610bc3610aa7565b5b6000610bd287828801610afa565b9450506020610be387828801610b30565b935050604085013567ffffffffffffffff811115610c0457610c03610aac565b5b610c1087828801610b54565b925092505092959194509250565b60008115159050919050565b610c3381610c1e565b82525050565b6000602082019050610c4e6000830184610c2a565b92915050565b600060208284031215610c6a57610c69610aa7565b5b6000610c7884828501610afa565b91505092915050565b610c8a81610b0f565b82525050565b6000602082019050610ca56000830184610c81565b92915050565b60008060408385031215610cc257610cc1610aa7565b5b6000610cd085828601610afa565b9250506020610ce185828601610b30565b9150509250929050565b6000610cf682610ad1565b9050919050565b610d0681610ceb565b8114610d1157600080fd5b50565b600081519050610d2381610cfd565b92915050565b600060208284031215610d3f57610d3e610aa7565b5b6000610d4d84828501610d14565b91505092915050565b600082825260208201905092915050565b7f53696d756c61746564204c494e4b20746f6b656e3a20696e636f72726563742060008201527f56524620436f6f7264696e61746f720000000000000000000000000000000000602082015250565b6000610dc3602f83610d56565b9150610dce82610d67565b604082019050919050565b60006020820190508181036000830152610df281610db6565b9050919050565b7f53696d756c61746564204c494e4b20746f6b656e3a20696e636f72726563742060008201527f66656520666f7220565246000000000000000000000000000000000000000000602082015250565b6000610e55602b83610d56565b9150610e6082610df9565b604082019050919050565b60006020820190508181036000830152610e8481610e48565b9050919050565b6000819050919050565b6000819050919050565b610eb0610eab82610e8b565b610e95565b82525050565b6000819050919050565b610ed1610ecc82610b0f565b610eb6565b82525050565b6000610ee38285610e9f565b602082019150610ef38284610ec0565b6020820191508190509392505050565b600081905092915050565b82818337600083830152505050565b6000610f298385610f03565b9350610f36838584610f0e565b82840190509392505050565b6000610f4f828486610f1d565b91508190509392505050565b7f53696d756c61746564204c494e4b20746f6b656e3a20696e76616c696420646160008201527f7461000000000000000000000000000000000000000000000000000000000000602082015250565b6000610fb7602283610d56565b9150610fc282610f5b565b604082019050919050565b60006020820190508181036000830152610fe681610faa565b9050919050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052601160045260246000fd5b600061102782610b0f565b915061103283610b0f565b9250827fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff0382111561106757611066610fed565b5b828201905092915050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052601260045260246000fd5b60006110ac82610b0f565b91506110b783610b0f565b9250826110c7576110c6611072565b5b828204905092915050565b600081519050919050565b600081905092915050565b60005b838110156111065780820151818401526020810190506110eb565b83811115611115576000848401525b50505050565b6000611126826110d2565b61113081856110dd565b93506111408185602086016110e8565b80840191505092915050565b7f2068617320696e73756666696369656e742062616c616e636520000000000000600082015250565b6000611182601a836110dd565b915061118d8261114c565b601a82019050919050565b7f65313420746f207472616e736665722000000000000000000000000000000000600082015250565b60006111ce6010836110dd565b91506111d982611198565b601082019050919050565b7f65313420746f2000000000000000000000000000000000000000000000000000600082015250565b600061121a6007836110dd565b9150611225826111e4565b600782019050919050565b600061123c828761111b565b915061124782611175565b9150611253828661111b565b915061125e826111c1565b915061126a828561111b565b91506112758261120d565b9150611281828461111b565b915081905095945050505050565b6000601f19601f8301169050919050565b60006112ab826110d2565b6112b58185610d56565b93506112c58185602086016110e8565b6112ce8161128f565b840191505092915050565b600060208201905081810360008301526112f381846112a0565b905092915050565b600061130682610b0f565b915061131183610b0f565b92508282101561132457611323610fed565b5b828203905092915050565b600061133a82610b0f565b915061134583610b0f565b9250817fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff048311821515161561137e5761137d610fed565b5b828202905092915050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052604160045260246000fd5b7f4e487b7100000000000000000000000000000000000000000000000000000000600052603260045260246000fd5b60006113f282610b0f565b9150600082141561140657611405610fed565b5b600182039050919050565b7f537472696e67733a20686578206c656e67746820696e73756666696369656e74600082015250565b6000611447602083610d56565b915061145282611411565b602082019050919050565b600060208201905081810360008301526114768161143a565b9050919050565b600061148882610b0f565b91507fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff8214156114bb576114ba610fed565b5b600182019050919050565b60006114d182610b0f565b91506114dc83610b0f565b9250826114ec576114eb611072565b5b82820690509291505056fea26469706673582212200d4a10c570d0a4850ec1aacd324bd9eef115d804c400bb9092216198cacc3e8764736f6c634300080b0033",
}

// SimulatedLinkTokenABI is the input ABI used to generate the binding from.
// Deprecated: Use SimulatedLinkTokenMetaData.ABI instead.
var SimulatedLinkTokenABI = SimulatedLinkTokenMetaData.ABI

// SimulatedLinkTokenBin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use SimulatedLinkTokenMetaData.Bin instead.
var SimulatedLinkTokenBin = SimulatedLinkTokenMetaData.Bin

// DeploySimulatedLinkToken deploys a new Ethereum contract, binding an instance of SimulatedLinkToken to it.
func DeploySimulatedLinkToken(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *SimulatedLinkToken, error) {
	parsed, err := SimulatedLinkTokenMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(SimulatedLinkTokenBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &SimulatedLinkToken{SimulatedLinkTokenCaller: SimulatedLinkTokenCaller{contract: contract}, SimulatedLinkTokenTransactor: SimulatedLinkTokenTransactor{contract: contract}, SimulatedLinkTokenFilterer: SimulatedLinkTokenFilterer{contract: contract}}, nil
}

// SimulatedLinkToken is an auto generated Go binding around an Ethereum contract.
type SimulatedLinkToken struct {
	SimulatedLinkTokenCaller     // Read-only binding to the contract
	SimulatedLinkTokenTransactor // Write-only binding to the contract
	SimulatedLinkTokenFilterer   // Log filterer for contract events
}

// SimulatedLinkTokenCaller is an auto generated read-only Go binding around an Ethereum contract.
type SimulatedLinkTokenCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// SimulatedLinkTokenTransactor is an auto generated write-only Go binding around an Ethereum contract.
type SimulatedLinkTokenTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// SimulatedLinkTokenFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type SimulatedLinkTokenFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// SimulatedLinkTokenSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type SimulatedLinkTokenSession struct {
	Contract     *SimulatedLinkToken // Generic contract binding to set the session for
	CallOpts     bind.CallOpts       // Call options to use throughout this session
	TransactOpts bind.TransactOpts   // Transaction auth options to use throughout this session
}

// SimulatedLinkTokenCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type SimulatedLinkTokenCallerSession struct {
	Contract *SimulatedLinkTokenCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts             // Call options to use throughout this session
}

// SimulatedLinkTokenTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type SimulatedLinkTokenTransactorSession struct {
	Contract     *SimulatedLinkTokenTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts             // Transaction auth options to use throughout this session
}

// SimulatedLinkTokenRaw is an auto generated low-level Go binding around an Ethereum contract.
type SimulatedLinkTokenRaw struct {
	Contract *SimulatedLinkToken // Generic contract binding to access the raw methods on
}

// SimulatedLinkTokenCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type SimulatedLinkTokenCallerRaw struct {
	Contract *SimulatedLinkTokenCaller // Generic read-only contract binding to access the raw methods on
}

// SimulatedLinkTokenTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type SimulatedLinkTokenTransactorRaw struct {
	Contract *SimulatedLinkTokenTransactor // Generic write-only contract binding to access the raw methods on
}

// NewSimulatedLinkToken creates a new instance of SimulatedLinkToken, bound to a specific deployed contract.
func NewSimulatedLinkToken(address common.Address, backend bind.ContractBackend) (*SimulatedLinkToken, error) {
	contract, err := bindSimulatedLinkToken(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &SimulatedLinkToken{SimulatedLinkTokenCaller: SimulatedLinkTokenCaller{contract: contract}, SimulatedLinkTokenTransactor: SimulatedLinkTokenTransactor{contract: contract}, SimulatedLinkTokenFilterer: SimulatedLinkTokenFilterer{contract: contract}}, nil
}

// NewSimulatedLinkTokenCaller creates a new read-only instance of SimulatedLinkToken, bound to a specific deployed contract.
func NewSimulatedLinkTokenCaller(address common.Address, caller bind.ContractCaller) (*SimulatedLinkTokenCaller, error) {
	contract, err := bindSimulatedLinkToken(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &SimulatedLinkTokenCaller{contract: contract}, nil
}

// NewSimulatedLinkTokenTransactor creates a new write-only instance of SimulatedLinkToken, bound to a specific deployed contract.
func NewSimulatedLinkTokenTransactor(address common.Address, transactor bind.ContractTransactor) (*SimulatedLinkTokenTransactor, error) {
	contract, err := bindSimulatedLinkToken(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &SimulatedLinkTokenTransactor{contract: contract}, nil
}

// NewSimulatedLinkTokenFilterer creates a new log filterer instance of SimulatedLinkToken, bound to a specific deployed contract.
func NewSimulatedLinkTokenFilterer(address common.Address, filterer bind.ContractFilterer) (*SimulatedLinkTokenFilterer, error) {
	contract, err := bindSimulatedLinkToken(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &SimulatedLinkTokenFilterer{contract: contract}, nil
}

// bindSimulatedLinkToken binds a generic wrapper to an already deployed contract.
func bindSimulatedLinkToken(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(SimulatedLinkTokenABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_SimulatedLinkToken *SimulatedLinkTokenRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _SimulatedLinkToken.Contract.SimulatedLinkTokenCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_SimulatedLinkToken *SimulatedLinkTokenRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _SimulatedLinkToken.Contract.SimulatedLinkTokenTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_SimulatedLinkToken *SimulatedLinkTokenRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _SimulatedLinkToken.Contract.SimulatedLinkTokenTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_SimulatedLinkToken *SimulatedLinkTokenCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _SimulatedLinkToken.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_SimulatedLinkToken *SimulatedLinkTokenTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _SimulatedLinkToken.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_SimulatedLinkToken *SimulatedLinkTokenTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _SimulatedLinkToken.Contract.contract.Transact(opts, method, params...)
}

// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
//
// Solidity: function balanceOf(address ) view returns(uint256)
func (_SimulatedLinkToken *SimulatedLinkTokenCaller) BalanceOf(opts *bind.CallOpts, arg0 common.Address) (*big.Int, error) {
	var out []interface{}
	err := _SimulatedLinkToken.contract.Call(opts, &out, "balanceOf", arg0)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
//
// Solidity: function balanceOf(address ) view returns(uint256)
func (_SimulatedLinkToken *SimulatedLinkTokenSession) BalanceOf(arg0 common.Address) (*big.Int, error) {
	return _SimulatedLinkToken.Contract.BalanceOf(&_SimulatedLinkToken.CallOpts, arg0)
}

// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
//
// Solidity: function balanceOf(address ) view returns(uint256)
func (_SimulatedLinkToken *SimulatedLinkTokenCallerSession) BalanceOf(arg0 common.Address) (*big.Int, error) {
	return _SimulatedLinkToken.Contract.BalanceOf(&_SimulatedLinkToken.CallOpts, arg0)
}

// Faucet is a paid mutator transaction binding the contract method 0x7b56c2b2.
//
// Solidity: function faucet(address recipient, uint256 amount) returns()
func (_SimulatedLinkToken *SimulatedLinkTokenTransactor) Faucet(opts *bind.TransactOpts, recipient common.Address, amount *big.Int) (*types.Transaction, error) {
	return _SimulatedLinkToken.contract.Transact(opts, "faucet", recipient, amount)
}

// Faucet is a paid mutator transaction binding the contract method 0x7b56c2b2.
//
// Solidity: function faucet(address recipient, uint256 amount) returns()
func (_SimulatedLinkToken *SimulatedLinkTokenSession) Faucet(recipient common.Address, amount *big.Int) (*types.Transaction, error) {
	return _SimulatedLinkToken.Contract.Faucet(&_SimulatedLinkToken.TransactOpts, recipient, amount)
}

// Faucet is a paid mutator transaction binding the contract method 0x7b56c2b2.
//
// Solidity: function faucet(address recipient, uint256 amount) returns()
func (_SimulatedLinkToken *SimulatedLinkTokenTransactorSession) Faucet(recipient common.Address, amount *big.Int) (*types.Transaction, error) {
	return _SimulatedLinkToken.Contract.Faucet(&_SimulatedLinkToken.TransactOpts, recipient, amount)
}

// Transfer is a paid mutator transaction binding the contract method 0xa9059cbb.
//
// Solidity: function transfer(address recipient, uint256 amount) returns(bool)
func (_SimulatedLinkToken *SimulatedLinkTokenTransactor) Transfer(opts *bind.TransactOpts, recipient common.Address, amount *big.Int) (*types.Transaction, error) {
	return _SimulatedLinkToken.contract.Transact(opts, "transfer", recipient, amount)
}

// Transfer is a paid mutator transaction binding the contract method 0xa9059cbb.
//
// Solidity: function transfer(address recipient, uint256 amount) returns(bool)
func (_SimulatedLinkToken *SimulatedLinkTokenSession) Transfer(recipient common.Address, amount *big.Int) (*types.Transaction, error) {
	return _SimulatedLinkToken.Contract.Transfer(&_SimulatedLinkToken.TransactOpts, recipient, amount)
}

// Transfer is a paid mutator transaction binding the contract method 0xa9059cbb.
//
// Solidity: function transfer(address recipient, uint256 amount) returns(bool)
func (_SimulatedLinkToken *SimulatedLinkTokenTransactorSession) Transfer(recipient common.Address, amount *big.Int) (*types.Transaction, error) {
	return _SimulatedLinkToken.Contract.Transfer(&_SimulatedLinkToken.TransactOpts, recipient, amount)
}

// TransferAndCall is a paid mutator transaction binding the contract method 0x4000aea0.
//
// Solidity: function transferAndCall(address to, uint256 value, bytes data) returns(bool success)
func (_SimulatedLinkToken *SimulatedLinkTokenTransactor) TransferAndCall(opts *bind.TransactOpts, to common.Address, value *big.Int, data []byte) (*types.Transaction, error) {
	return _SimulatedLinkToken.contract.Transact(opts, "transferAndCall", to, value, data)
}

// TransferAndCall is a paid mutator transaction binding the contract method 0x4000aea0.
//
// Solidity: function transferAndCall(address to, uint256 value, bytes data) returns(bool success)
func (_SimulatedLinkToken *SimulatedLinkTokenSession) TransferAndCall(to common.Address, value *big.Int, data []byte) (*types.Transaction, error) {
	return _SimulatedLinkToken.Contract.TransferAndCall(&_SimulatedLinkToken.TransactOpts, to, value, data)
}

// TransferAndCall is a paid mutator transaction binding the contract method 0x4000aea0.
//
// Solidity: function transferAndCall(address to, uint256 value, bytes data) returns(bool success)
func (_SimulatedLinkToken *SimulatedLinkTokenTransactorSession) TransferAndCall(to common.Address, value *big.Int, data []byte) (*types.Transaction, error) {
	return _SimulatedLinkToken.Contract.TransferAndCall(&_SimulatedLinkToken.TransactOpts, to, value, data)
}

// SimulatedLinkTokenRandomnessRequestedIterator is returned from FilterRandomnessRequested and is used to iterate over the raw logs and unpacked data for RandomnessRequested events raised by the SimulatedLinkToken contract.
type SimulatedLinkTokenRandomnessRequestedIterator struct {
	Event *SimulatedLinkTokenRandomnessRequested // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *SimulatedLinkTokenRandomnessRequestedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(SimulatedLinkTokenRandomnessRequested)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(SimulatedLinkTokenRandomnessRequested)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *SimulatedLinkTokenRandomnessRequestedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *SimulatedLinkTokenRandomnessRequestedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// SimulatedLinkTokenRandomnessRequested represents a RandomnessRequested event raised by the SimulatedLinkToken contract.
type SimulatedLinkTokenRandomnessRequested struct {
	By  common.Address
	Raw types.Log // Blockchain specific contextual infos
}

// FilterRandomnessRequested is a free log retrieval operation binding the contract event 0x6871e329198b319adcce5196458d56b81939284c334a427d776ebdc356dd5acb.
//
// Solidity: event RandomnessRequested(address indexed by)
func (_SimulatedLinkToken *SimulatedLinkTokenFilterer) FilterRandomnessRequested(opts *bind.FilterOpts, by []common.Address) (*SimulatedLinkTokenRandomnessRequestedIterator, error) {

	var byRule []interface{}
	for _, byItem := range by {
		byRule = append(byRule, byItem)
	}

	logs, sub, err := _SimulatedLinkToken.contract.FilterLogs(opts, "RandomnessRequested", byRule)
	if err != nil {
		return nil, err
	}
	return &SimulatedLinkTokenRandomnessRequestedIterator{contract: _SimulatedLinkToken.contract, event: "RandomnessRequested", logs: logs, sub: sub}, nil
}

// WatchRandomnessRequested is a free log subscription operation binding the contract event 0x6871e329198b319adcce5196458d56b81939284c334a427d776ebdc356dd5acb.
//
// Solidity: event RandomnessRequested(address indexed by)
func (_SimulatedLinkToken *SimulatedLinkTokenFilterer) WatchRandomnessRequested(opts *bind.WatchOpts, sink chan<- *SimulatedLinkTokenRandomnessRequested, by []common.Address) (event.Subscription, error) {

	var byRule []interface{}
	for _, byItem := range by {
		byRule = append(byRule, byItem)
	}

	logs, sub, err := _SimulatedLinkToken.contract.WatchLogs(opts, "RandomnessRequested", byRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(SimulatedLinkTokenRandomnessRequested)
				if err := _SimulatedLinkToken.contract.UnpackLog(event, "RandomnessRequested", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseRandomnessRequested is a log parse operation binding the contract event 0x6871e329198b319adcce5196458d56b81939284c334a427d776ebdc356dd5acb.
//
// Solidity: event RandomnessRequested(address indexed by)
func (_SimulatedLinkToken *SimulatedLinkTokenFilterer) ParseRandomnessRequested(log types.Log) (*SimulatedLinkTokenRandomnessRequested, error) {
	event := new(SimulatedLinkTokenRandomnessRequested)
	if err := _SimulatedLinkToken.contract.UnpackLog(event, "RandomnessRequested", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// SimulatedVRFCoordinatorMetaData contains all meta data concerning the SimulatedVRFCoordinator contract.
var SimulatedVRFCoordinatorMetaData = &bind.MetaData{
	ABI: "[{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"previousOwner\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"inputs\":[{\"internalType\":\"contractVRFConsumerBase\",\"name\":\"consumer\",\"type\":\"address\"}],\"name\":\"fulfill\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"renounceOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
	Bin: "0x608060405234801561001057600080fd5b5061002d61002261003260201b60201c565b61003a60201b60201c565b6100fe565b600033905090565b60008060009054906101000a900473ffffffffffffffffffffffffffffffffffffffff169050816000806101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff1602179055508173ffffffffffffffffffffffffffffffffffffffff168173ffffffffffffffffffffffffffffffffffffffff167f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e060405160405180910390a35050565b610a5e8061010d6000396000f3fe608060405234801561001057600080fd5b506004361061004c5760003560e01c806344e03a1f14610051578063715018a61461006d5780638da5cb5b14610077578063f2fde38b14610095575b600080fd5b61006b600480360381019061006691906106c8565b6100b1565b005b61007561028f565b005b61007f610317565b60405161008c9190610704565b60405180910390f35b6100af60048036038101906100aa919061074b565b610340565b005b6100b9610438565b73ffffffffffffffffffffffffffffffffffffffff166100d7610317565b73ffffffffffffffffffffffffffffffffffffffff161461012d576040517f08c379a0000000000000000000000000000000000000000000000000000000008152600401610124906107d5565b60405180910390fd5b6000819050600061018761013f610440565b600084600160008773ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002054610520565b9050600160008373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060008154809291906101d99061082e565b919050555060006101f16101eb610440565b8361055c565b90508373ffffffffffffffffffffffffffffffffffffffff166394985ddd828360405160200161022191906108a2565b6040516020818303038152906040528051906020012060001c6040518363ffffffff1660e01b81526004016102579291906108db565b600060405180830381600087803b15801561027157600080fd5b505af1158015610285573d6000803e3d6000fd5b5050505050505050565b610297610438565b73ffffffffffffffffffffffffffffffffffffffff166102b5610317565b73ffffffffffffffffffffffffffffffffffffffff161461030b576040517f08c379a0000000000000000000000000000000000000000000000000000000008152600401610302906107d5565b60405180910390fd5b610315600061058f565b565b60008060009054906101000a900473ffffffffffffffffffffffffffffffffffffffff16905090565b610348610438565b73ffffffffffffffffffffffffffffffffffffffff16610366610317565b73ffffffffffffffffffffffffffffffffffffffff16146103bc576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004016103b3906107d5565b60405180910390fd5b600073ffffffffffffffffffffffffffffffffffffffff168173ffffffffffffffffffffffffffffffffffffffff16141561042c576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040161042390610976565b60405180910390fd5b6104358161058f565b50565b600033905090565b60004660018114610473576004811461049b57608981146104c3576201388181146104eb5761053981146105135761051c565b7faa77729d3466ca35ae8d28b3bbac7cc36a5031efdc430821c02bc31a238af445915061051c565b7f2ed0feb3e7fd2022120aa84fab1945545a9f2ffc9076fd6156fa96eaff4c1311915061051c565b7ff86195cf7690c55907b2b611ebb7343a6f649bff128701cc542f0569e2c549da915061051c565b7f6e75b569a01ef56d18cab6a8e71e6600d6ce853834d4a5748b720d06f878b3a4915061051c565b60026113372091505b5090565b6000848484846040516020016105399493929190610996565b6040516020818303038152906040528051906020012060001c9050949350505050565b600082826040516020016105719291906109fc565b60405160208183030381529060405280519060200120905092915050565b60008060009054906101000a900473ffffffffffffffffffffffffffffffffffffffff169050816000806101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff1602179055508173ffffffffffffffffffffffffffffffffffffffff168173ffffffffffffffffffffffffffffffffffffffff167f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e060405160405180910390a35050565b600080fd5b600073ffffffffffffffffffffffffffffffffffffffff82169050919050565b600061068382610658565b9050919050565b600061069582610678565b9050919050565b6106a58161068a565b81146106b057600080fd5b50565b6000813590506106c28161069c565b92915050565b6000602082840312156106de576106dd610653565b5b60006106ec848285016106b3565b91505092915050565b6106fe81610678565b82525050565b600060208201905061071960008301846106f5565b92915050565b61072881610678565b811461073357600080fd5b50565b6000813590506107458161071f565b92915050565b60006020828403121561076157610760610653565b5b600061076f84828501610736565b91505092915050565b600082825260208201905092915050565b7f4f776e61626c653a2063616c6c6572206973206e6f7420746865206f776e6572600082015250565b60006107bf602083610778565b91506107ca82610789565b602082019050919050565b600060208201905081810360008301526107ee816107b2565b9050919050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052601160045260246000fd5b6000819050919050565b600061083982610824565b91507fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff82141561086c5761086b6107f5565b5b600182019050919050565b6000819050919050565b6000819050919050565b61089c61089782610877565b610881565b82525050565b60006108ae828461088b565b60208201915081905092915050565b6108c681610877565b82525050565b6108d581610824565b82525050565b60006040820190506108f060008301856108bd565b6108fd60208301846108cc565b9392505050565b7f4f776e61626c653a206e6577206f776e657220697320746865207a65726f206160008201527f6464726573730000000000000000000000000000000000000000000000000000602082015250565b6000610960602683610778565b915061096b82610904565b604082019050919050565b6000602082019050818103600083015261098f81610953565b9050919050565b60006080820190506109ab60008301876108bd565b6109b860208301866108cc565b6109c560408301856106f5565b6109d260608301846108cc565b95945050505050565b6000819050919050565b6109f66109f182610824565b6109db565b82525050565b6000610a08828561088b565b602082019150610a1882846109e5565b602082019150819050939250505056fea26469706673582212206529d122f2f8d2b0379431c5cdbeb3deb5df96030d0cefbb709af7a4d457279464736f6c634300080b0033",
}

// SimulatedVRFCoordinatorABI is the input ABI used to generate the binding from.
// Deprecated: Use SimulatedVRFCoordinatorMetaData.ABI instead.
var SimulatedVRFCoordinatorABI = SimulatedVRFCoordinatorMetaData.ABI

// SimulatedVRFCoordinatorBin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use SimulatedVRFCoordinatorMetaData.Bin instead.
var SimulatedVRFCoordinatorBin = SimulatedVRFCoordinatorMetaData.Bin

// DeploySimulatedVRFCoordinator deploys a new Ethereum contract, binding an instance of SimulatedVRFCoordinator to it.
func DeploySimulatedVRFCoordinator(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *SimulatedVRFCoordinator, error) {
	parsed, err := SimulatedVRFCoordinatorMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(SimulatedVRFCoordinatorBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &SimulatedVRFCoordinator{SimulatedVRFCoordinatorCaller: SimulatedVRFCoordinatorCaller{contract: contract}, SimulatedVRFCoordinatorTransactor: SimulatedVRFCoordinatorTransactor{contract: contract}, SimulatedVRFCoordinatorFilterer: SimulatedVRFCoordinatorFilterer{contract: contract}}, nil
}

// SimulatedVRFCoordinator is an auto generated Go binding around an Ethereum contract.
type SimulatedVRFCoordinator struct {
	SimulatedVRFCoordinatorCaller     // Read-only binding to the contract
	SimulatedVRFCoordinatorTransactor // Write-only binding to the contract
	SimulatedVRFCoordinatorFilterer   // Log filterer for contract events
}

// SimulatedVRFCoordinatorCaller is an auto generated read-only Go binding around an Ethereum contract.
type SimulatedVRFCoordinatorCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// SimulatedVRFCoordinatorTransactor is an auto generated write-only Go binding around an Ethereum contract.
type SimulatedVRFCoordinatorTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// SimulatedVRFCoordinatorFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type SimulatedVRFCoordinatorFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// SimulatedVRFCoordinatorSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type SimulatedVRFCoordinatorSession struct {
	Contract     *SimulatedVRFCoordinator // Generic contract binding to set the session for
	CallOpts     bind.CallOpts            // Call options to use throughout this session
	TransactOpts bind.TransactOpts        // Transaction auth options to use throughout this session
}

// SimulatedVRFCoordinatorCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type SimulatedVRFCoordinatorCallerSession struct {
	Contract *SimulatedVRFCoordinatorCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts                  // Call options to use throughout this session
}

// SimulatedVRFCoordinatorTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type SimulatedVRFCoordinatorTransactorSession struct {
	Contract     *SimulatedVRFCoordinatorTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts                  // Transaction auth options to use throughout this session
}

// SimulatedVRFCoordinatorRaw is an auto generated low-level Go binding around an Ethereum contract.
type SimulatedVRFCoordinatorRaw struct {
	Contract *SimulatedVRFCoordinator // Generic contract binding to access the raw methods on
}

// SimulatedVRFCoordinatorCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type SimulatedVRFCoordinatorCallerRaw struct {
	Contract *SimulatedVRFCoordinatorCaller // Generic read-only contract binding to access the raw methods on
}

// SimulatedVRFCoordinatorTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type SimulatedVRFCoordinatorTransactorRaw struct {
	Contract *SimulatedVRFCoordinatorTransactor // Generic write-only contract binding to access the raw methods on
}

// NewSimulatedVRFCoordinator creates a new instance of SimulatedVRFCoordinator, bound to a specific deployed contract.
func NewSimulatedVRFCoordinator(address common.Address, backend bind.ContractBackend) (*SimulatedVRFCoordinator, error) {
	contract, err := bindSimulatedVRFCoordinator(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &SimulatedVRFCoordinator{SimulatedVRFCoordinatorCaller: SimulatedVRFCoordinatorCaller{contract: contract}, SimulatedVRFCoordinatorTransactor: SimulatedVRFCoordinatorTransactor{contract: contract}, SimulatedVRFCoordinatorFilterer: SimulatedVRFCoordinatorFilterer{contract: contract}}, nil
}

// NewSimulatedVRFCoordinatorCaller creates a new read-only instance of SimulatedVRFCoordinator, bound to a specific deployed contract.
func NewSimulatedVRFCoordinatorCaller(address common.Address, caller bind.ContractCaller) (*SimulatedVRFCoordinatorCaller, error) {
	contract, err := bindSimulatedVRFCoordinator(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &SimulatedVRFCoordinatorCaller{contract: contract}, nil
}

// NewSimulatedVRFCoordinatorTransactor creates a new write-only instance of SimulatedVRFCoordinator, bound to a specific deployed contract.
func NewSimulatedVRFCoordinatorTransactor(address common.Address, transactor bind.ContractTransactor) (*SimulatedVRFCoordinatorTransactor, error) {
	contract, err := bindSimulatedVRFCoordinator(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &SimulatedVRFCoordinatorTransactor{contract: contract}, nil
}

// NewSimulatedVRFCoordinatorFilterer creates a new log filterer instance of SimulatedVRFCoordinator, bound to a specific deployed contract.
func NewSimulatedVRFCoordinatorFilterer(address common.Address, filterer bind.ContractFilterer) (*SimulatedVRFCoordinatorFilterer, error) {
	contract, err := bindSimulatedVRFCoordinator(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &SimulatedVRFCoordinatorFilterer{contract: contract}, nil
}

// bindSimulatedVRFCoordinator binds a generic wrapper to an already deployed contract.
func bindSimulatedVRFCoordinator(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(SimulatedVRFCoordinatorABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_SimulatedVRFCoordinator *SimulatedVRFCoordinatorRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _SimulatedVRFCoordinator.Contract.SimulatedVRFCoordinatorCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_SimulatedVRFCoordinator *SimulatedVRFCoordinatorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _SimulatedVRFCoordinator.Contract.SimulatedVRFCoordinatorTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_SimulatedVRFCoordinator *SimulatedVRFCoordinatorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _SimulatedVRFCoordinator.Contract.SimulatedVRFCoordinatorTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_SimulatedVRFCoordinator *SimulatedVRFCoordinatorCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _SimulatedVRFCoordinator.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_SimulatedVRFCoordinator *SimulatedVRFCoordinatorTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _SimulatedVRFCoordinator.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_SimulatedVRFCoordinator *SimulatedVRFCoordinatorTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _SimulatedVRFCoordinator.Contract.contract.Transact(opts, method, params...)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_SimulatedVRFCoordinator *SimulatedVRFCoordinatorCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _SimulatedVRFCoordinator.contract.Call(opts, &out, "owner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_SimulatedVRFCoordinator *SimulatedVRFCoordinatorSession) Owner() (common.Address, error) {
	return _SimulatedVRFCoordinator.Contract.Owner(&_SimulatedVRFCoordinator.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_SimulatedVRFCoordinator *SimulatedVRFCoordinatorCallerSession) Owner() (common.Address, error) {
	return _SimulatedVRFCoordinator.Contract.Owner(&_SimulatedVRFCoordinator.CallOpts)
}

// Fulfill is a paid mutator transaction binding the contract method 0x44e03a1f.
//
// Solidity: function fulfill(address consumer) returns()
func (_SimulatedVRFCoordinator *SimulatedVRFCoordinatorTransactor) Fulfill(opts *bind.TransactOpts, consumer common.Address) (*types.Transaction, error) {
	return _SimulatedVRFCoordinator.contract.Transact(opts, "fulfill", consumer)
}

// Fulfill is a paid mutator transaction binding the contract method 0x44e03a1f.
//
// Solidity: function fulfill(address consumer) returns()
func (_SimulatedVRFCoordinator *SimulatedVRFCoordinatorSession) Fulfill(consumer common.Address) (*types.Transaction, error) {
	return _SimulatedVRFCoordinator.Contract.Fulfill(&_SimulatedVRFCoordinator.TransactOpts, consumer)
}

// Fulfill is a paid mutator transaction binding the contract method 0x44e03a1f.
//
// Solidity: function fulfill(address consumer) returns()
func (_SimulatedVRFCoordinator *SimulatedVRFCoordinatorTransactorSession) Fulfill(consumer common.Address) (*types.Transaction, error) {
	return _SimulatedVRFCoordinator.Contract.Fulfill(&_SimulatedVRFCoordinator.TransactOpts, consumer)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_SimulatedVRFCoordinator *SimulatedVRFCoordinatorTransactor) RenounceOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _SimulatedVRFCoordinator.contract.Transact(opts, "renounceOwnership")
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_SimulatedVRFCoordinator *SimulatedVRFCoordinatorSession) RenounceOwnership() (*types.Transaction, error) {
	return _SimulatedVRFCoordinator.Contract.RenounceOwnership(&_SimulatedVRFCoordinator.TransactOpts)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_SimulatedVRFCoordinator *SimulatedVRFCoordinatorTransactorSession) RenounceOwnership() (*types.Transaction, error) {
	return _SimulatedVRFCoordinator.Contract.RenounceOwnership(&_SimulatedVRFCoordinator.TransactOpts)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_SimulatedVRFCoordinator *SimulatedVRFCoordinatorTransactor) TransferOwnership(opts *bind.TransactOpts, newOwner common.Address) (*types.Transaction, error) {
	return _SimulatedVRFCoordinator.contract.Transact(opts, "transferOwnership", newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_SimulatedVRFCoordinator *SimulatedVRFCoordinatorSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _SimulatedVRFCoordinator.Contract.TransferOwnership(&_SimulatedVRFCoordinator.TransactOpts, newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_SimulatedVRFCoordinator *SimulatedVRFCoordinatorTransactorSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _SimulatedVRFCoordinator.Contract.TransferOwnership(&_SimulatedVRFCoordinator.TransactOpts, newOwner)
}

// SimulatedVRFCoordinatorOwnershipTransferredIterator is returned from FilterOwnershipTransferred and is used to iterate over the raw logs and unpacked data for OwnershipTransferred events raised by the SimulatedVRFCoordinator contract.
type SimulatedVRFCoordinatorOwnershipTransferredIterator struct {
	Event *SimulatedVRFCoordinatorOwnershipTransferred // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *SimulatedVRFCoordinatorOwnershipTransferredIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(SimulatedVRFCoordinatorOwnershipTransferred)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(SimulatedVRFCoordinatorOwnershipTransferred)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *SimulatedVRFCoordinatorOwnershipTransferredIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *SimulatedVRFCoordinatorOwnershipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// SimulatedVRFCoordinatorOwnershipTransferred represents a OwnershipTransferred event raised by the SimulatedVRFCoordinator contract.
type SimulatedVRFCoordinatorOwnershipTransferred struct {
	PreviousOwner common.Address
	NewOwner      common.Address
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterOwnershipTransferred is a free log retrieval operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_SimulatedVRFCoordinator *SimulatedVRFCoordinatorFilterer) FilterOwnershipTransferred(opts *bind.FilterOpts, previousOwner []common.Address, newOwner []common.Address) (*SimulatedVRFCoordinatorOwnershipTransferredIterator, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _SimulatedVRFCoordinator.contract.FilterLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return &SimulatedVRFCoordinatorOwnershipTransferredIterator{contract: _SimulatedVRFCoordinator.contract, event: "OwnershipTransferred", logs: logs, sub: sub}, nil
}

// WatchOwnershipTransferred is a free log subscription operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_SimulatedVRFCoordinator *SimulatedVRFCoordinatorFilterer) WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *SimulatedVRFCoordinatorOwnershipTransferred, previousOwner []common.Address, newOwner []common.Address) (event.Subscription, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _SimulatedVRFCoordinator.contract.WatchLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(SimulatedVRFCoordinatorOwnershipTransferred)
				if err := _SimulatedVRFCoordinator.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseOwnershipTransferred is a log parse operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_SimulatedVRFCoordinator *SimulatedVRFCoordinatorFilterer) ParseOwnershipTransferred(log types.Log) (*SimulatedVRFCoordinatorOwnershipTransferred, error) {
	event := new(SimulatedVRFCoordinatorOwnershipTransferred)
	if err := _SimulatedVRFCoordinator.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// StringsMetaData contains all meta data concerning the Strings contract.
var StringsMetaData = &bind.MetaData{
	ABI: "[]",
	Bin: "0x60566050600b82828239805160001a6073146043577f4e487b7100000000000000000000000000000000000000000000000000000000600052600060045260246000fd5b30600052607381538281f3fe73000000000000000000000000000000000000000030146080604052600080fdfea264697066735822122018578ba17adbc840fc6f53d6c438cefd829f27a49e1588cfa275f616756bbf0f64736f6c634300080b0033",
}

// StringsABI is the input ABI used to generate the binding from.
// Deprecated: Use StringsMetaData.ABI instead.
var StringsABI = StringsMetaData.ABI

// StringsBin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use StringsMetaData.Bin instead.
var StringsBin = StringsMetaData.Bin

// DeployStrings deploys a new Ethereum contract, binding an instance of Strings to it.
func DeployStrings(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *Strings, error) {
	parsed, err := StringsMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(StringsBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &Strings{StringsCaller: StringsCaller{contract: contract}, StringsTransactor: StringsTransactor{contract: contract}, StringsFilterer: StringsFilterer{contract: contract}}, nil
}

// Strings is an auto generated Go binding around an Ethereum contract.
type Strings struct {
	StringsCaller     // Read-only binding to the contract
	StringsTransactor // Write-only binding to the contract
	StringsFilterer   // Log filterer for contract events
}

// StringsCaller is an auto generated read-only Go binding around an Ethereum contract.
type StringsCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// StringsTransactor is an auto generated write-only Go binding around an Ethereum contract.
type StringsTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// StringsFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type StringsFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// StringsSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type StringsSession struct {
	Contract     *Strings          // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// StringsCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type StringsCallerSession struct {
	Contract *StringsCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts  // Call options to use throughout this session
}

// StringsTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type StringsTransactorSession struct {
	Contract     *StringsTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts  // Transaction auth options to use throughout this session
}

// StringsRaw is an auto generated low-level Go binding around an Ethereum contract.
type StringsRaw struct {
	Contract *Strings // Generic contract binding to access the raw methods on
}

// StringsCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type StringsCallerRaw struct {
	Contract *StringsCaller // Generic read-only contract binding to access the raw methods on
}

// StringsTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type StringsTransactorRaw struct {
	Contract *StringsTransactor // Generic write-only contract binding to access the raw methods on
}

// NewStrings creates a new instance of Strings, bound to a specific deployed contract.
func NewStrings(address common.Address, backend bind.ContractBackend) (*Strings, error) {
	contract, err := bindStrings(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &Strings{StringsCaller: StringsCaller{contract: contract}, StringsTransactor: StringsTransactor{contract: contract}, StringsFilterer: StringsFilterer{contract: contract}}, nil
}

// NewStringsCaller creates a new read-only instance of Strings, bound to a specific deployed contract.
func NewStringsCaller(address common.Address, caller bind.ContractCaller) (*StringsCaller, error) {
	contract, err := bindStrings(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &StringsCaller{contract: contract}, nil
}

// NewStringsTransactor creates a new write-only instance of Strings, bound to a specific deployed contract.
func NewStringsTransactor(address common.Address, transactor bind.ContractTransactor) (*StringsTransactor, error) {
	contract, err := bindStrings(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &StringsTransactor{contract: contract}, nil
}

// NewStringsFilterer creates a new log filterer instance of Strings, bound to a specific deployed contract.
func NewStringsFilterer(address common.Address, filterer bind.ContractFilterer) (*StringsFilterer, error) {
	contract, err := bindStrings(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &StringsFilterer{contract: contract}, nil
}

// bindStrings binds a generic wrapper to an already deployed contract.
func bindStrings(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(StringsABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Strings *StringsRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Strings.Contract.StringsCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Strings *StringsRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Strings.Contract.StringsTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Strings *StringsRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Strings.Contract.StringsTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Strings *StringsCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Strings.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Strings *StringsTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Strings.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Strings *StringsTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Strings.Contract.contract.Transact(opts, method, params...)
}

// VRFConsumerBaseMetaData contains all meta data concerning the VRFConsumerBase contract.
var VRFConsumerBaseMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"requestId\",\"type\":\"bytes32\"},{\"internalType\":\"uint256\",\"name\":\"randomness\",\"type\":\"uint256\"}],\"name\":\"rawFulfillRandomness\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
}

// VRFConsumerBaseABI is the input ABI used to generate the binding from.
// Deprecated: Use VRFConsumerBaseMetaData.ABI instead.
var VRFConsumerBaseABI = VRFConsumerBaseMetaData.ABI

// VRFConsumerBase is an auto generated Go binding around an Ethereum contract.
type VRFConsumerBase struct {
	VRFConsumerBaseCaller     // Read-only binding to the contract
	VRFConsumerBaseTransactor // Write-only binding to the contract
	VRFConsumerBaseFilterer   // Log filterer for contract events
}

// VRFConsumerBaseCaller is an auto generated read-only Go binding around an Ethereum contract.
type VRFConsumerBaseCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// VRFConsumerBaseTransactor is an auto generated write-only Go binding around an Ethereum contract.
type VRFConsumerBaseTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// VRFConsumerBaseFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type VRFConsumerBaseFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// VRFConsumerBaseSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type VRFConsumerBaseSession struct {
	Contract     *VRFConsumerBase  // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// VRFConsumerBaseCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type VRFConsumerBaseCallerSession struct {
	Contract *VRFConsumerBaseCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts          // Call options to use throughout this session
}

// VRFConsumerBaseTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type VRFConsumerBaseTransactorSession struct {
	Contract     *VRFConsumerBaseTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts          // Transaction auth options to use throughout this session
}

// VRFConsumerBaseRaw is an auto generated low-level Go binding around an Ethereum contract.
type VRFConsumerBaseRaw struct {
	Contract *VRFConsumerBase // Generic contract binding to access the raw methods on
}

// VRFConsumerBaseCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type VRFConsumerBaseCallerRaw struct {
	Contract *VRFConsumerBaseCaller // Generic read-only contract binding to access the raw methods on
}

// VRFConsumerBaseTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type VRFConsumerBaseTransactorRaw struct {
	Contract *VRFConsumerBaseTransactor // Generic write-only contract binding to access the raw methods on
}

// NewVRFConsumerBase creates a new instance of VRFConsumerBase, bound to a specific deployed contract.
func NewVRFConsumerBase(address common.Address, backend bind.ContractBackend) (*VRFConsumerBase, error) {
	contract, err := bindVRFConsumerBase(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &VRFConsumerBase{VRFConsumerBaseCaller: VRFConsumerBaseCaller{contract: contract}, VRFConsumerBaseTransactor: VRFConsumerBaseTransactor{contract: contract}, VRFConsumerBaseFilterer: VRFConsumerBaseFilterer{contract: contract}}, nil
}

// NewVRFConsumerBaseCaller creates a new read-only instance of VRFConsumerBase, bound to a specific deployed contract.
func NewVRFConsumerBaseCaller(address common.Address, caller bind.ContractCaller) (*VRFConsumerBaseCaller, error) {
	contract, err := bindVRFConsumerBase(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &VRFConsumerBaseCaller{contract: contract}, nil
}

// NewVRFConsumerBaseTransactor creates a new write-only instance of VRFConsumerBase, bound to a specific deployed contract.
func NewVRFConsumerBaseTransactor(address common.Address, transactor bind.ContractTransactor) (*VRFConsumerBaseTransactor, error) {
	contract, err := bindVRFConsumerBase(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &VRFConsumerBaseTransactor{contract: contract}, nil
}

// NewVRFConsumerBaseFilterer creates a new log filterer instance of VRFConsumerBase, bound to a specific deployed contract.
func NewVRFConsumerBaseFilterer(address common.Address, filterer bind.ContractFilterer) (*VRFConsumerBaseFilterer, error) {
	contract, err := bindVRFConsumerBase(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &VRFConsumerBaseFilterer{contract: contract}, nil
}

// bindVRFConsumerBase binds a generic wrapper to an already deployed contract.
func bindVRFConsumerBase(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(VRFConsumerBaseABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_VRFConsumerBase *VRFConsumerBaseRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _VRFConsumerBase.Contract.VRFConsumerBaseCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_VRFConsumerBase *VRFConsumerBaseRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _VRFConsumerBase.Contract.VRFConsumerBaseTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_VRFConsumerBase *VRFConsumerBaseRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _VRFConsumerBase.Contract.VRFConsumerBaseTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_VRFConsumerBase *VRFConsumerBaseCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _VRFConsumerBase.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_VRFConsumerBase *VRFConsumerBaseTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _VRFConsumerBase.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_VRFConsumerBase *VRFConsumerBaseTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _VRFConsumerBase.Contract.contract.Transact(opts, method, params...)
}

// RawFulfillRandomness is a paid mutator transaction binding the contract method 0x94985ddd.
//
// Solidity: function rawFulfillRandomness(bytes32 requestId, uint256 randomness) returns()
func (_VRFConsumerBase *VRFConsumerBaseTransactor) RawFulfillRandomness(opts *bind.TransactOpts, requestId [32]byte, randomness *big.Int) (*types.Transaction, error) {
	return _VRFConsumerBase.contract.Transact(opts, "rawFulfillRandomness", requestId, randomness)
}

// RawFulfillRandomness is a paid mutator transaction binding the contract method 0x94985ddd.
//
// Solidity: function rawFulfillRandomness(bytes32 requestId, uint256 randomness) returns()
func (_VRFConsumerBase *VRFConsumerBaseSession) RawFulfillRandomness(requestId [32]byte, randomness *big.Int) (*types.Transaction, error) {
	return _VRFConsumerBase.Contract.RawFulfillRandomness(&_VRFConsumerBase.TransactOpts, requestId, randomness)
}

// RawFulfillRandomness is a paid mutator transaction binding the contract method 0x94985ddd.
//
// Solidity: function rawFulfillRandomness(bytes32 requestId, uint256 randomness) returns()
func (_VRFConsumerBase *VRFConsumerBaseTransactorSession) RawFulfillRandomness(requestId [32]byte, randomness *big.Int) (*types.Transaction, error) {
	return _VRFConsumerBase.Contract.RawFulfillRandomness(&_VRFConsumerBase.TransactOpts, requestId, randomness)
}

// VRFRequestIDBaseMetaData contains all meta data concerning the VRFRequestIDBase contract.
var VRFRequestIDBaseMetaData = &bind.MetaData{
	ABI: "[]",
	Bin: "0x6080604052348015600f57600080fd5b50603f80601d6000396000f3fe6080604052600080fdfea264697066735822122064f4d2fc289f631a9adb98822668df4f1fc3caf4764b20b1f61e134da434dbec64736f6c634300080b0033",
}

// VRFRequestIDBaseABI is the input ABI used to generate the binding from.
// Deprecated: Use VRFRequestIDBaseMetaData.ABI instead.
var VRFRequestIDBaseABI = VRFRequestIDBaseMetaData.ABI

// VRFRequestIDBaseBin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use VRFRequestIDBaseMetaData.Bin instead.
var VRFRequestIDBaseBin = VRFRequestIDBaseMetaData.Bin

// DeployVRFRequestIDBase deploys a new Ethereum contract, binding an instance of VRFRequestIDBase to it.
func DeployVRFRequestIDBase(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *VRFRequestIDBase, error) {
	parsed, err := VRFRequestIDBaseMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(VRFRequestIDBaseBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &VRFRequestIDBase{VRFRequestIDBaseCaller: VRFRequestIDBaseCaller{contract: contract}, VRFRequestIDBaseTransactor: VRFRequestIDBaseTransactor{contract: contract}, VRFRequestIDBaseFilterer: VRFRequestIDBaseFilterer{contract: contract}}, nil
}

// VRFRequestIDBase is an auto generated Go binding around an Ethereum contract.
type VRFRequestIDBase struct {
	VRFRequestIDBaseCaller     // Read-only binding to the contract
	VRFRequestIDBaseTransactor // Write-only binding to the contract
	VRFRequestIDBaseFilterer   // Log filterer for contract events
}

// VRFRequestIDBaseCaller is an auto generated read-only Go binding around an Ethereum contract.
type VRFRequestIDBaseCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// VRFRequestIDBaseTransactor is an auto generated write-only Go binding around an Ethereum contract.
type VRFRequestIDBaseTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// VRFRequestIDBaseFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type VRFRequestIDBaseFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// VRFRequestIDBaseSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type VRFRequestIDBaseSession struct {
	Contract     *VRFRequestIDBase // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// VRFRequestIDBaseCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type VRFRequestIDBaseCallerSession struct {
	Contract *VRFRequestIDBaseCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts           // Call options to use throughout this session
}

// VRFRequestIDBaseTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type VRFRequestIDBaseTransactorSession struct {
	Contract     *VRFRequestIDBaseTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts           // Transaction auth options to use throughout this session
}

// VRFRequestIDBaseRaw is an auto generated low-level Go binding around an Ethereum contract.
type VRFRequestIDBaseRaw struct {
	Contract *VRFRequestIDBase // Generic contract binding to access the raw methods on
}

// VRFRequestIDBaseCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type VRFRequestIDBaseCallerRaw struct {
	Contract *VRFRequestIDBaseCaller // Generic read-only contract binding to access the raw methods on
}

// VRFRequestIDBaseTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type VRFRequestIDBaseTransactorRaw struct {
	Contract *VRFRequestIDBaseTransactor // Generic write-only contract binding to access the raw methods on
}

// NewVRFRequestIDBase creates a new instance of VRFRequestIDBase, bound to a specific deployed contract.
func NewVRFRequestIDBase(address common.Address, backend bind.ContractBackend) (*VRFRequestIDBase, error) {
	contract, err := bindVRFRequestIDBase(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &VRFRequestIDBase{VRFRequestIDBaseCaller: VRFRequestIDBaseCaller{contract: contract}, VRFRequestIDBaseTransactor: VRFRequestIDBaseTransactor{contract: contract}, VRFRequestIDBaseFilterer: VRFRequestIDBaseFilterer{contract: contract}}, nil
}

// NewVRFRequestIDBaseCaller creates a new read-only instance of VRFRequestIDBase, bound to a specific deployed contract.
func NewVRFRequestIDBaseCaller(address common.Address, caller bind.ContractCaller) (*VRFRequestIDBaseCaller, error) {
	contract, err := bindVRFRequestIDBase(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &VRFRequestIDBaseCaller{contract: contract}, nil
}

// NewVRFRequestIDBaseTransactor creates a new write-only instance of VRFRequestIDBase, bound to a specific deployed contract.
func NewVRFRequestIDBaseTransactor(address common.Address, transactor bind.ContractTransactor) (*VRFRequestIDBaseTransactor, error) {
	contract, err := bindVRFRequestIDBase(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &VRFRequestIDBaseTransactor{contract: contract}, nil
}

// NewVRFRequestIDBaseFilterer creates a new log filterer instance of VRFRequestIDBase, bound to a specific deployed contract.
func NewVRFRequestIDBaseFilterer(address common.Address, filterer bind.ContractFilterer) (*VRFRequestIDBaseFilterer, error) {
	contract, err := bindVRFRequestIDBase(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &VRFRequestIDBaseFilterer{contract: contract}, nil
}

// bindVRFRequestIDBase binds a generic wrapper to an already deployed contract.
func bindVRFRequestIDBase(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(VRFRequestIDBaseABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_VRFRequestIDBase *VRFRequestIDBaseRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _VRFRequestIDBase.Contract.VRFRequestIDBaseCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_VRFRequestIDBase *VRFRequestIDBaseRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _VRFRequestIDBase.Contract.VRFRequestIDBaseTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_VRFRequestIDBase *VRFRequestIDBaseRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _VRFRequestIDBase.Contract.VRFRequestIDBaseTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_VRFRequestIDBase *VRFRequestIDBaseCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _VRFRequestIDBase.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_VRFRequestIDBase *VRFRequestIDBaseTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _VRFRequestIDBase.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_VRFRequestIDBase *VRFRequestIDBaseTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _VRFRequestIDBase.Contract.contract.Transact(opts, method, params...)
}
