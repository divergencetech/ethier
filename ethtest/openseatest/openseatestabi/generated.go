// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package openseatestabi

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

// OwnableDelegateProxyMetaData contains all meta data concerning the OwnableDelegateProxy contract.
var OwnableDelegateProxyMetaData = &bind.MetaData{
	ABI: "[]",
	Bin: "0x6080604052348015600f57600080fd5b50603f80601d6000396000f3fe6080604052600080fdfea26469706673582212209091d38f303f5fa796dac40177149812a66a3d0f5613781575fa0a04d931ace464736f6c634300080f0033",
}

// OwnableDelegateProxyABI is the input ABI used to generate the binding from.
// Deprecated: Use OwnableDelegateProxyMetaData.ABI instead.
var OwnableDelegateProxyABI = OwnableDelegateProxyMetaData.ABI

// OwnableDelegateProxyBin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use OwnableDelegateProxyMetaData.Bin instead.
var OwnableDelegateProxyBin = OwnableDelegateProxyMetaData.Bin

// DeployOwnableDelegateProxy deploys a new Ethereum contract, binding an instance of OwnableDelegateProxy to it.
func DeployOwnableDelegateProxy(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *OwnableDelegateProxy, error) {
	parsed, err := OwnableDelegateProxyMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(OwnableDelegateProxyBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &OwnableDelegateProxy{OwnableDelegateProxyCaller: OwnableDelegateProxyCaller{contract: contract}, OwnableDelegateProxyTransactor: OwnableDelegateProxyTransactor{contract: contract}, OwnableDelegateProxyFilterer: OwnableDelegateProxyFilterer{contract: contract}}, nil
}

// OwnableDelegateProxy is an auto generated Go binding around an Ethereum contract.
type OwnableDelegateProxy struct {
	OwnableDelegateProxyCaller     // Read-only binding to the contract
	OwnableDelegateProxyTransactor // Write-only binding to the contract
	OwnableDelegateProxyFilterer   // Log filterer for contract events
}

// OwnableDelegateProxyCaller is an auto generated read-only Go binding around an Ethereum contract.
type OwnableDelegateProxyCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// OwnableDelegateProxyTransactor is an auto generated write-only Go binding around an Ethereum contract.
type OwnableDelegateProxyTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// OwnableDelegateProxyFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type OwnableDelegateProxyFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// OwnableDelegateProxySession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type OwnableDelegateProxySession struct {
	Contract     *OwnableDelegateProxy // Generic contract binding to set the session for
	CallOpts     bind.CallOpts         // Call options to use throughout this session
	TransactOpts bind.TransactOpts     // Transaction auth options to use throughout this session
}

// OwnableDelegateProxyCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type OwnableDelegateProxyCallerSession struct {
	Contract *OwnableDelegateProxyCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts               // Call options to use throughout this session
}

// OwnableDelegateProxyTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type OwnableDelegateProxyTransactorSession struct {
	Contract     *OwnableDelegateProxyTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts               // Transaction auth options to use throughout this session
}

// OwnableDelegateProxyRaw is an auto generated low-level Go binding around an Ethereum contract.
type OwnableDelegateProxyRaw struct {
	Contract *OwnableDelegateProxy // Generic contract binding to access the raw methods on
}

// OwnableDelegateProxyCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type OwnableDelegateProxyCallerRaw struct {
	Contract *OwnableDelegateProxyCaller // Generic read-only contract binding to access the raw methods on
}

// OwnableDelegateProxyTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type OwnableDelegateProxyTransactorRaw struct {
	Contract *OwnableDelegateProxyTransactor // Generic write-only contract binding to access the raw methods on
}

// NewOwnableDelegateProxy creates a new instance of OwnableDelegateProxy, bound to a specific deployed contract.
func NewOwnableDelegateProxy(address common.Address, backend bind.ContractBackend) (*OwnableDelegateProxy, error) {
	contract, err := bindOwnableDelegateProxy(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &OwnableDelegateProxy{OwnableDelegateProxyCaller: OwnableDelegateProxyCaller{contract: contract}, OwnableDelegateProxyTransactor: OwnableDelegateProxyTransactor{contract: contract}, OwnableDelegateProxyFilterer: OwnableDelegateProxyFilterer{contract: contract}}, nil
}

// NewOwnableDelegateProxyCaller creates a new read-only instance of OwnableDelegateProxy, bound to a specific deployed contract.
func NewOwnableDelegateProxyCaller(address common.Address, caller bind.ContractCaller) (*OwnableDelegateProxyCaller, error) {
	contract, err := bindOwnableDelegateProxy(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &OwnableDelegateProxyCaller{contract: contract}, nil
}

// NewOwnableDelegateProxyTransactor creates a new write-only instance of OwnableDelegateProxy, bound to a specific deployed contract.
func NewOwnableDelegateProxyTransactor(address common.Address, transactor bind.ContractTransactor) (*OwnableDelegateProxyTransactor, error) {
	contract, err := bindOwnableDelegateProxy(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &OwnableDelegateProxyTransactor{contract: contract}, nil
}

// NewOwnableDelegateProxyFilterer creates a new log filterer instance of OwnableDelegateProxy, bound to a specific deployed contract.
func NewOwnableDelegateProxyFilterer(address common.Address, filterer bind.ContractFilterer) (*OwnableDelegateProxyFilterer, error) {
	contract, err := bindOwnableDelegateProxy(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &OwnableDelegateProxyFilterer{contract: contract}, nil
}

// bindOwnableDelegateProxy binds a generic wrapper to an already deployed contract.
func bindOwnableDelegateProxy(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(OwnableDelegateProxyABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_OwnableDelegateProxy *OwnableDelegateProxyRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _OwnableDelegateProxy.Contract.OwnableDelegateProxyCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_OwnableDelegateProxy *OwnableDelegateProxyRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _OwnableDelegateProxy.Contract.OwnableDelegateProxyTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_OwnableDelegateProxy *OwnableDelegateProxyRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _OwnableDelegateProxy.Contract.OwnableDelegateProxyTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_OwnableDelegateProxy *OwnableDelegateProxyCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _OwnableDelegateProxy.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_OwnableDelegateProxy *OwnableDelegateProxyTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _OwnableDelegateProxy.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_OwnableDelegateProxy *OwnableDelegateProxyTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _OwnableDelegateProxy.Contract.contract.Transact(opts, method, params...)
}

// ProxyRegistryMetaData contains all meta data concerning the ProxyRegistry contract.
var ProxyRegistryMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"name\":\"proxies\",\"outputs\":[{\"internalType\":\"contractOwnableDelegateProxy\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"}]",
	Sigs: map[string]string{
		"c4552791": "proxies(address)",
	},
	Bin: "0x608060405234801561001057600080fd5b5060d38061001f6000396000f3fe6080604052348015600f57600080fd5b506004361060285760003560e01c8063c455279114602d575b600080fd5b60536038366004606f565b6000602081905290815260409020546001600160a01b031681565b6040516001600160a01b03909116815260200160405180910390f35b600060208284031215608057600080fd5b81356001600160a01b0381168114609657600080fd5b939250505056fea2646970667358221220a179d21fa1cbbd9f3e7a02e1ebcf36db5a83b36df4ed0e20443e06b3830b4b7c64736f6c634300080f0033",
}

// ProxyRegistryABI is the input ABI used to generate the binding from.
// Deprecated: Use ProxyRegistryMetaData.ABI instead.
var ProxyRegistryABI = ProxyRegistryMetaData.ABI

// Deprecated: Use ProxyRegistryMetaData.Sigs instead.
// ProxyRegistryFuncSigs maps the 4-byte function signature to its string representation.
var ProxyRegistryFuncSigs = ProxyRegistryMetaData.Sigs

// ProxyRegistryBin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use ProxyRegistryMetaData.Bin instead.
var ProxyRegistryBin = ProxyRegistryMetaData.Bin

// DeployProxyRegistry deploys a new Ethereum contract, binding an instance of ProxyRegistry to it.
func DeployProxyRegistry(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *ProxyRegistry, error) {
	parsed, err := ProxyRegistryMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(ProxyRegistryBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &ProxyRegistry{ProxyRegistryCaller: ProxyRegistryCaller{contract: contract}, ProxyRegistryTransactor: ProxyRegistryTransactor{contract: contract}, ProxyRegistryFilterer: ProxyRegistryFilterer{contract: contract}}, nil
}

// ProxyRegistry is an auto generated Go binding around an Ethereum contract.
type ProxyRegistry struct {
	ProxyRegistryCaller     // Read-only binding to the contract
	ProxyRegistryTransactor // Write-only binding to the contract
	ProxyRegistryFilterer   // Log filterer for contract events
}

// ProxyRegistryCaller is an auto generated read-only Go binding around an Ethereum contract.
type ProxyRegistryCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ProxyRegistryTransactor is an auto generated write-only Go binding around an Ethereum contract.
type ProxyRegistryTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ProxyRegistryFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type ProxyRegistryFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ProxyRegistrySession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type ProxyRegistrySession struct {
	Contract     *ProxyRegistry    // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// ProxyRegistryCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type ProxyRegistryCallerSession struct {
	Contract *ProxyRegistryCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts        // Call options to use throughout this session
}

// ProxyRegistryTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type ProxyRegistryTransactorSession struct {
	Contract     *ProxyRegistryTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts        // Transaction auth options to use throughout this session
}

// ProxyRegistryRaw is an auto generated low-level Go binding around an Ethereum contract.
type ProxyRegistryRaw struct {
	Contract *ProxyRegistry // Generic contract binding to access the raw methods on
}

// ProxyRegistryCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type ProxyRegistryCallerRaw struct {
	Contract *ProxyRegistryCaller // Generic read-only contract binding to access the raw methods on
}

// ProxyRegistryTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type ProxyRegistryTransactorRaw struct {
	Contract *ProxyRegistryTransactor // Generic write-only contract binding to access the raw methods on
}

// NewProxyRegistry creates a new instance of ProxyRegistry, bound to a specific deployed contract.
func NewProxyRegistry(address common.Address, backend bind.ContractBackend) (*ProxyRegistry, error) {
	contract, err := bindProxyRegistry(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &ProxyRegistry{ProxyRegistryCaller: ProxyRegistryCaller{contract: contract}, ProxyRegistryTransactor: ProxyRegistryTransactor{contract: contract}, ProxyRegistryFilterer: ProxyRegistryFilterer{contract: contract}}, nil
}

// NewProxyRegistryCaller creates a new read-only instance of ProxyRegistry, bound to a specific deployed contract.
func NewProxyRegistryCaller(address common.Address, caller bind.ContractCaller) (*ProxyRegistryCaller, error) {
	contract, err := bindProxyRegistry(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &ProxyRegistryCaller{contract: contract}, nil
}

// NewProxyRegistryTransactor creates a new write-only instance of ProxyRegistry, bound to a specific deployed contract.
func NewProxyRegistryTransactor(address common.Address, transactor bind.ContractTransactor) (*ProxyRegistryTransactor, error) {
	contract, err := bindProxyRegistry(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &ProxyRegistryTransactor{contract: contract}, nil
}

// NewProxyRegistryFilterer creates a new log filterer instance of ProxyRegistry, bound to a specific deployed contract.
func NewProxyRegistryFilterer(address common.Address, filterer bind.ContractFilterer) (*ProxyRegistryFilterer, error) {
	contract, err := bindProxyRegistry(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &ProxyRegistryFilterer{contract: contract}, nil
}

// bindProxyRegistry binds a generic wrapper to an already deployed contract.
func bindProxyRegistry(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(ProxyRegistryABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_ProxyRegistry *ProxyRegistryRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _ProxyRegistry.Contract.ProxyRegistryCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_ProxyRegistry *ProxyRegistryRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ProxyRegistry.Contract.ProxyRegistryTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_ProxyRegistry *ProxyRegistryRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _ProxyRegistry.Contract.ProxyRegistryTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_ProxyRegistry *ProxyRegistryCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _ProxyRegistry.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_ProxyRegistry *ProxyRegistryTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ProxyRegistry.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_ProxyRegistry *ProxyRegistryTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _ProxyRegistry.Contract.contract.Transact(opts, method, params...)
}

// Proxies is a free data retrieval call binding the contract method 0xc4552791.
//
// Solidity: function proxies(address ) view returns(address)
func (_ProxyRegistry *ProxyRegistryCaller) Proxies(opts *bind.CallOpts, arg0 common.Address) (common.Address, error) {
	var out []interface{}
	err := _ProxyRegistry.contract.Call(opts, &out, "proxies", arg0)

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Proxies is a free data retrieval call binding the contract method 0xc4552791.
//
// Solidity: function proxies(address ) view returns(address)
func (_ProxyRegistry *ProxyRegistrySession) Proxies(arg0 common.Address) (common.Address, error) {
	return _ProxyRegistry.Contract.Proxies(&_ProxyRegistry.CallOpts, arg0)
}

// Proxies is a free data retrieval call binding the contract method 0xc4552791.
//
// Solidity: function proxies(address ) view returns(address)
func (_ProxyRegistry *ProxyRegistryCallerSession) Proxies(arg0 common.Address) (common.Address, error) {
	return _ProxyRegistry.Contract.Proxies(&_ProxyRegistry.CallOpts, arg0)
}

// SimulatedProxyRegistryMetaData contains all meta data concerning the SimulatedProxyRegistry contract.
var SimulatedProxyRegistryMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"name\":\"proxies\",\"outputs\":[{\"internalType\":\"contractOwnableDelegateProxy\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"proxy\",\"type\":\"address\"}],\"name\":\"setProxyFor\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
	Sigs: map[string]string{
		"c4552791": "proxies(address)",
		"ab72167d": "setProxyFor(address,address)",
	},
	Bin: "0x608060405234801561001057600080fd5b50610165806100206000396000f3fe608060405234801561001057600080fd5b50600436106100365760003560e01c8063ab72167d1461003b578063c455279114610079575b600080fd5b6100776100493660046100da565b6001600160a01b03918216600090815260208190526040902080546001600160a01b03191691909216179055565b005b6100a261008736600461010d565b6000602081905290815260409020546001600160a01b031681565b6040516001600160a01b03909116815260200160405180910390f35b80356001600160a01b03811681146100d557600080fd5b919050565b600080604083850312156100ed57600080fd5b6100f6836100be565b9150610104602084016100be565b90509250929050565b60006020828403121561011f57600080fd5b610128826100be565b939250505056fea2646970667358221220cbcf690fc0c886597a91564f8c8a1b66eab9cf915080949c180d1a70bdeba15364736f6c634300080f0033",
}

// SimulatedProxyRegistryABI is the input ABI used to generate the binding from.
// Deprecated: Use SimulatedProxyRegistryMetaData.ABI instead.
var SimulatedProxyRegistryABI = SimulatedProxyRegistryMetaData.ABI

// Deprecated: Use SimulatedProxyRegistryMetaData.Sigs instead.
// SimulatedProxyRegistryFuncSigs maps the 4-byte function signature to its string representation.
var SimulatedProxyRegistryFuncSigs = SimulatedProxyRegistryMetaData.Sigs

// SimulatedProxyRegistryBin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use SimulatedProxyRegistryMetaData.Bin instead.
var SimulatedProxyRegistryBin = SimulatedProxyRegistryMetaData.Bin

// DeploySimulatedProxyRegistry deploys a new Ethereum contract, binding an instance of SimulatedProxyRegistry to it.
func DeploySimulatedProxyRegistry(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *SimulatedProxyRegistry, error) {
	parsed, err := SimulatedProxyRegistryMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(SimulatedProxyRegistryBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &SimulatedProxyRegistry{SimulatedProxyRegistryCaller: SimulatedProxyRegistryCaller{contract: contract}, SimulatedProxyRegistryTransactor: SimulatedProxyRegistryTransactor{contract: contract}, SimulatedProxyRegistryFilterer: SimulatedProxyRegistryFilterer{contract: contract}}, nil
}

// SimulatedProxyRegistry is an auto generated Go binding around an Ethereum contract.
type SimulatedProxyRegistry struct {
	SimulatedProxyRegistryCaller     // Read-only binding to the contract
	SimulatedProxyRegistryTransactor // Write-only binding to the contract
	SimulatedProxyRegistryFilterer   // Log filterer for contract events
}

// SimulatedProxyRegistryCaller is an auto generated read-only Go binding around an Ethereum contract.
type SimulatedProxyRegistryCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// SimulatedProxyRegistryTransactor is an auto generated write-only Go binding around an Ethereum contract.
type SimulatedProxyRegistryTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// SimulatedProxyRegistryFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type SimulatedProxyRegistryFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// SimulatedProxyRegistrySession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type SimulatedProxyRegistrySession struct {
	Contract     *SimulatedProxyRegistry // Generic contract binding to set the session for
	CallOpts     bind.CallOpts           // Call options to use throughout this session
	TransactOpts bind.TransactOpts       // Transaction auth options to use throughout this session
}

// SimulatedProxyRegistryCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type SimulatedProxyRegistryCallerSession struct {
	Contract *SimulatedProxyRegistryCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts                 // Call options to use throughout this session
}

// SimulatedProxyRegistryTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type SimulatedProxyRegistryTransactorSession struct {
	Contract     *SimulatedProxyRegistryTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts                 // Transaction auth options to use throughout this session
}

// SimulatedProxyRegistryRaw is an auto generated low-level Go binding around an Ethereum contract.
type SimulatedProxyRegistryRaw struct {
	Contract *SimulatedProxyRegistry // Generic contract binding to access the raw methods on
}

// SimulatedProxyRegistryCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type SimulatedProxyRegistryCallerRaw struct {
	Contract *SimulatedProxyRegistryCaller // Generic read-only contract binding to access the raw methods on
}

// SimulatedProxyRegistryTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type SimulatedProxyRegistryTransactorRaw struct {
	Contract *SimulatedProxyRegistryTransactor // Generic write-only contract binding to access the raw methods on
}

// NewSimulatedProxyRegistry creates a new instance of SimulatedProxyRegistry, bound to a specific deployed contract.
func NewSimulatedProxyRegistry(address common.Address, backend bind.ContractBackend) (*SimulatedProxyRegistry, error) {
	contract, err := bindSimulatedProxyRegistry(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &SimulatedProxyRegistry{SimulatedProxyRegistryCaller: SimulatedProxyRegistryCaller{contract: contract}, SimulatedProxyRegistryTransactor: SimulatedProxyRegistryTransactor{contract: contract}, SimulatedProxyRegistryFilterer: SimulatedProxyRegistryFilterer{contract: contract}}, nil
}

// NewSimulatedProxyRegistryCaller creates a new read-only instance of SimulatedProxyRegistry, bound to a specific deployed contract.
func NewSimulatedProxyRegistryCaller(address common.Address, caller bind.ContractCaller) (*SimulatedProxyRegistryCaller, error) {
	contract, err := bindSimulatedProxyRegistry(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &SimulatedProxyRegistryCaller{contract: contract}, nil
}

// NewSimulatedProxyRegistryTransactor creates a new write-only instance of SimulatedProxyRegistry, bound to a specific deployed contract.
func NewSimulatedProxyRegistryTransactor(address common.Address, transactor bind.ContractTransactor) (*SimulatedProxyRegistryTransactor, error) {
	contract, err := bindSimulatedProxyRegistry(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &SimulatedProxyRegistryTransactor{contract: contract}, nil
}

// NewSimulatedProxyRegistryFilterer creates a new log filterer instance of SimulatedProxyRegistry, bound to a specific deployed contract.
func NewSimulatedProxyRegistryFilterer(address common.Address, filterer bind.ContractFilterer) (*SimulatedProxyRegistryFilterer, error) {
	contract, err := bindSimulatedProxyRegistry(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &SimulatedProxyRegistryFilterer{contract: contract}, nil
}

// bindSimulatedProxyRegistry binds a generic wrapper to an already deployed contract.
func bindSimulatedProxyRegistry(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(SimulatedProxyRegistryABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_SimulatedProxyRegistry *SimulatedProxyRegistryRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _SimulatedProxyRegistry.Contract.SimulatedProxyRegistryCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_SimulatedProxyRegistry *SimulatedProxyRegistryRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _SimulatedProxyRegistry.Contract.SimulatedProxyRegistryTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_SimulatedProxyRegistry *SimulatedProxyRegistryRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _SimulatedProxyRegistry.Contract.SimulatedProxyRegistryTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_SimulatedProxyRegistry *SimulatedProxyRegistryCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _SimulatedProxyRegistry.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_SimulatedProxyRegistry *SimulatedProxyRegistryTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _SimulatedProxyRegistry.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_SimulatedProxyRegistry *SimulatedProxyRegistryTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _SimulatedProxyRegistry.Contract.contract.Transact(opts, method, params...)
}

// Proxies is a free data retrieval call binding the contract method 0xc4552791.
//
// Solidity: function proxies(address ) view returns(address)
func (_SimulatedProxyRegistry *SimulatedProxyRegistryCaller) Proxies(opts *bind.CallOpts, arg0 common.Address) (common.Address, error) {
	var out []interface{}
	err := _SimulatedProxyRegistry.contract.Call(opts, &out, "proxies", arg0)

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Proxies is a free data retrieval call binding the contract method 0xc4552791.
//
// Solidity: function proxies(address ) view returns(address)
func (_SimulatedProxyRegistry *SimulatedProxyRegistrySession) Proxies(arg0 common.Address) (common.Address, error) {
	return _SimulatedProxyRegistry.Contract.Proxies(&_SimulatedProxyRegistry.CallOpts, arg0)
}

// Proxies is a free data retrieval call binding the contract method 0xc4552791.
//
// Solidity: function proxies(address ) view returns(address)
func (_SimulatedProxyRegistry *SimulatedProxyRegistryCallerSession) Proxies(arg0 common.Address) (common.Address, error) {
	return _SimulatedProxyRegistry.Contract.Proxies(&_SimulatedProxyRegistry.CallOpts, arg0)
}

// SetProxyFor is a paid mutator transaction binding the contract method 0xab72167d.
//
// Solidity: function setProxyFor(address owner, address proxy) returns()
func (_SimulatedProxyRegistry *SimulatedProxyRegistryTransactor) SetProxyFor(opts *bind.TransactOpts, owner common.Address, proxy common.Address) (*types.Transaction, error) {
	return _SimulatedProxyRegistry.contract.Transact(opts, "setProxyFor", owner, proxy)
}

// SetProxyFor is a paid mutator transaction binding the contract method 0xab72167d.
//
// Solidity: function setProxyFor(address owner, address proxy) returns()
func (_SimulatedProxyRegistry *SimulatedProxyRegistrySession) SetProxyFor(owner common.Address, proxy common.Address) (*types.Transaction, error) {
	return _SimulatedProxyRegistry.Contract.SetProxyFor(&_SimulatedProxyRegistry.TransactOpts, owner, proxy)
}

// SetProxyFor is a paid mutator transaction binding the contract method 0xab72167d.
//
// Solidity: function setProxyFor(address owner, address proxy) returns()
func (_SimulatedProxyRegistry *SimulatedProxyRegistryTransactorSession) SetProxyFor(owner common.Address, proxy common.Address) (*types.Transaction, error) {
	return _SimulatedProxyRegistry.Contract.SetProxyFor(&_SimulatedProxyRegistry.TransactOpts, owner, proxy)
}

