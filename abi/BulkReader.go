// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package abi

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

// BulkReaderMetaData contains all meta data concerning the BulkReader contract.
var BulkReaderMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[],\"name\":\"getReserves\",\"outputs\":[{\"internalType\":\"uint256[2][]\",\"name\":\"\",\"type\":\"uint256[2][]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address[]\",\"name\":\"pairs\",\"type\":\"address[]\"}],\"name\":\"getReserves\",\"outputs\":[{\"internalType\":\"uint256[2][]\",\"name\":\"\",\"type\":\"uint256[2][]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address[]\",\"name\":\"_pairs\",\"type\":\"address[]\"}],\"name\":\"getReservesForBenchmark\",\"outputs\":[{\"internalType\":\"uint256[3][]\",\"name\":\"\",\"type\":\"uint256[3][]\"}],\"stateMutability\":\"view\",\"type\":\"function\"}]",
}

// BulkReaderABI is the input ABI used to generate the binding from.
// Deprecated: Use BulkReaderMetaData.ABI instead.
var BulkReaderABI = BulkReaderMetaData.ABI

// BulkReader is an auto generated Go binding around an Ethereum contract.
type BulkReader struct {
	BulkReaderCaller     // Read-only binding to the contract
	BulkReaderTransactor // Write-only binding to the contract
	BulkReaderFilterer   // Log filterer for contract events
}

// BulkReaderCaller is an auto generated read-only Go binding around an Ethereum contract.
type BulkReaderCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// BulkReaderTransactor is an auto generated write-only Go binding around an Ethereum contract.
type BulkReaderTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// BulkReaderFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type BulkReaderFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// BulkReaderSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type BulkReaderSession struct {
	Contract     *BulkReader       // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// BulkReaderCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type BulkReaderCallerSession struct {
	Contract *BulkReaderCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts     // Call options to use throughout this session
}

// BulkReaderTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type BulkReaderTransactorSession struct {
	Contract     *BulkReaderTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts     // Transaction auth options to use throughout this session
}

// BulkReaderRaw is an auto generated low-level Go binding around an Ethereum contract.
type BulkReaderRaw struct {
	Contract *BulkReader // Generic contract binding to access the raw methods on
}

// BulkReaderCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type BulkReaderCallerRaw struct {
	Contract *BulkReaderCaller // Generic read-only contract binding to access the raw methods on
}

// BulkReaderTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type BulkReaderTransactorRaw struct {
	Contract *BulkReaderTransactor // Generic write-only contract binding to access the raw methods on
}

// NewBulkReader creates a new instance of BulkReader, bound to a specific deployed contract.
func NewBulkReader(address common.Address, backend bind.ContractBackend) (*BulkReader, error) {
	contract, err := bindBulkReader(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &BulkReader{BulkReaderCaller: BulkReaderCaller{contract: contract}, BulkReaderTransactor: BulkReaderTransactor{contract: contract}, BulkReaderFilterer: BulkReaderFilterer{contract: contract}}, nil
}

// NewBulkReaderCaller creates a new read-only instance of BulkReader, bound to a specific deployed contract.
func NewBulkReaderCaller(address common.Address, caller bind.ContractCaller) (*BulkReaderCaller, error) {
	contract, err := bindBulkReader(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &BulkReaderCaller{contract: contract}, nil
}

// NewBulkReaderTransactor creates a new write-only instance of BulkReader, bound to a specific deployed contract.
func NewBulkReaderTransactor(address common.Address, transactor bind.ContractTransactor) (*BulkReaderTransactor, error) {
	contract, err := bindBulkReader(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &BulkReaderTransactor{contract: contract}, nil
}

// NewBulkReaderFilterer creates a new log filterer instance of BulkReader, bound to a specific deployed contract.
func NewBulkReaderFilterer(address common.Address, filterer bind.ContractFilterer) (*BulkReaderFilterer, error) {
	contract, err := bindBulkReader(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &BulkReaderFilterer{contract: contract}, nil
}

// bindBulkReader binds a generic wrapper to an already deployed contract.
func bindBulkReader(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(BulkReaderABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_BulkReader *BulkReaderRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _BulkReader.Contract.BulkReaderCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_BulkReader *BulkReaderRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _BulkReader.Contract.BulkReaderTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_BulkReader *BulkReaderRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _BulkReader.Contract.BulkReaderTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_BulkReader *BulkReaderCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _BulkReader.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_BulkReader *BulkReaderTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _BulkReader.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_BulkReader *BulkReaderTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _BulkReader.Contract.contract.Transact(opts, method, params...)
}

// GetReserves is a free data retrieval call binding the contract method 0x0902f1ac.
//
// Solidity: function getReserves() view returns(uint256[2][])
func (_BulkReader *BulkReaderCaller) GetReserves(opts *bind.CallOpts) ([][2]*big.Int, error) {
	var out []interface{}
	err := _BulkReader.contract.Call(opts, &out, "getReserves")

	if err != nil {
		return *new([][2]*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new([][2]*big.Int)).(*[][2]*big.Int)

	return out0, err

}

// GetReserves is a free data retrieval call binding the contract method 0x0902f1ac.
//
// Solidity: function getReserves() view returns(uint256[2][])
func (_BulkReader *BulkReaderSession) GetReserves() ([][2]*big.Int, error) {
	return _BulkReader.Contract.GetReserves(&_BulkReader.CallOpts)
}

// GetReserves is a free data retrieval call binding the contract method 0x0902f1ac.
//
// Solidity: function getReserves() view returns(uint256[2][])
func (_BulkReader *BulkReaderCallerSession) GetReserves() ([][2]*big.Int, error) {
	return _BulkReader.Contract.GetReserves(&_BulkReader.CallOpts)
}

// GetReserves0 is a free data retrieval call binding the contract method 0x407a4b08.
//
// Solidity: function getReserves(address[] pairs) view returns(uint256[2][])
func (_BulkReader *BulkReaderCaller) GetReserves0(opts *bind.CallOpts, pairs []common.Address) ([][2]*big.Int, error) {
	var out []interface{}
	err := _BulkReader.contract.Call(opts, &out, "getReserves0", pairs)

	if err != nil {
		return *new([][2]*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new([][2]*big.Int)).(*[][2]*big.Int)

	return out0, err

}

// GetReserves0 is a free data retrieval call binding the contract method 0x407a4b08.
//
// Solidity: function getReserves(address[] pairs) view returns(uint256[2][])
func (_BulkReader *BulkReaderSession) GetReserves0(pairs []common.Address) ([][2]*big.Int, error) {
	return _BulkReader.Contract.GetReserves0(&_BulkReader.CallOpts, pairs)
}

// GetReserves0 is a free data retrieval call binding the contract method 0x407a4b08.
//
// Solidity: function getReserves(address[] pairs) view returns(uint256[2][])
func (_BulkReader *BulkReaderCallerSession) GetReserves0(pairs []common.Address) ([][2]*big.Int, error) {
	return _BulkReader.Contract.GetReserves0(&_BulkReader.CallOpts, pairs)
}

// GetReservesForBenchmark is a free data retrieval call binding the contract method 0xef7b22d9.
//
// Solidity: function getReservesForBenchmark(address[] _pairs) view returns(uint256[3][])
func (_BulkReader *BulkReaderCaller) GetReservesForBenchmark(opts *bind.CallOpts, _pairs []common.Address) ([][3]*big.Int, error) {
	var out []interface{}
	err := _BulkReader.contract.Call(opts, &out, "getReservesForBenchmark", _pairs)

	if err != nil {
		return *new([][3]*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new([][3]*big.Int)).(*[][3]*big.Int)

	return out0, err

}

// GetReservesForBenchmark is a free data retrieval call binding the contract method 0xef7b22d9.
//
// Solidity: function getReservesForBenchmark(address[] _pairs) view returns(uint256[3][])
func (_BulkReader *BulkReaderSession) GetReservesForBenchmark(_pairs []common.Address) ([][3]*big.Int, error) {
	return _BulkReader.Contract.GetReservesForBenchmark(&_BulkReader.CallOpts, _pairs)
}

// GetReservesForBenchmark is a free data retrieval call binding the contract method 0xef7b22d9.
//
// Solidity: function getReservesForBenchmark(address[] _pairs) view returns(uint256[3][])
func (_BulkReader *BulkReaderCallerSession) GetReservesForBenchmark(_pairs []common.Address) ([][3]*big.Int, error) {
	return _BulkReader.Contract.GetReservesForBenchmark(&_BulkReader.CallOpts, _pairs)
}
