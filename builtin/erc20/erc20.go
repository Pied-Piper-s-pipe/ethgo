// Code generated by ethgo/abigen. DO NOT EDIT.
// Hash: a1a873d70d345feef023ee086fd6135b24d775444b950ee9d5ea411e72b0f373
// Version: 0.1.3
package erc20

import (
	"fmt"
	"math/big"

	"github.com/umbracle/ethgo"
	"github.com/umbracle/ethgo/contract"
	"github.com/umbracle/ethgo/jsonrpc"
)

var (
	_ = big.NewInt
	_ = jsonrpc.NewClient
)

// ERC20 is a solidity contract
type ERC20 struct {
	c *contract.Contract
}

// NewERC20 creates a new instance of the contract at a specific address
func NewERC20(addr ethgo.Address, opts ...contract.ContractOption) *ERC20 {
	return &ERC20{c: contract.NewContract(addr, abiERC20, opts...)}
}

// calls

// Allowance calls the allowance method in the solidity contract
func (e *ERC20) Allowance(owner ethgo.Address, spender ethgo.Address, block ...ethgo.BlockNumber) (retval0 *big.Int, err error) {
	var out map[string]interface{}
	var ok bool

	out, err = e.c.Call("allowance", ethgo.EncodeBlock(block...), owner, spender)
	if err != nil {
		return
	}

	// decode outputs
	retval0, ok = out["0"].(*big.Int)
	if !ok {
		err = fmt.Errorf("failed to encode output at index 0")
		return
	}
	
	return
}

// BalanceOf calls the balanceOf method in the solidity contract
func (e *ERC20) BalanceOf(owner ethgo.Address, block ...ethgo.BlockNumber) (retval0 *big.Int, err error) {
	var out map[string]interface{}
	var ok bool

	out, err = e.c.Call("balanceOf", ethgo.EncodeBlock(block...), owner)
	if err != nil {
		return
	}

	// decode outputs
	retval0, ok = out["balance"].(*big.Int)
	if !ok {
		err = fmt.Errorf("failed to encode output at index 0")
		return
	}
	
	return
}

// Decimals calls the decimals method in the solidity contract
func (e *ERC20) Decimals(block ...ethgo.BlockNumber) (retval0 uint8, err error) {
	var out map[string]interface{}
	var ok bool

	out, err = e.c.Call("decimals", ethgo.EncodeBlock(block...))
	if err != nil {
		return
	}

	// decode outputs
	retval0, ok = out["0"].(uint8)
	if !ok {
		err = fmt.Errorf("failed to encode output at index 0")
		return
	}
	
	return
}

// Name calls the name method in the solidity contract
func (e *ERC20) Name(block ...ethgo.BlockNumber) (retval0 string, err error) {
	var out map[string]interface{}
	var ok bool

	out, err = e.c.Call("name", ethgo.EncodeBlock(block...))
	if err != nil {
		return
	}

	// decode outputs
	retval0, ok = out["0"].(string)
	if !ok {
		err = fmt.Errorf("failed to encode output at index 0")
		return
	}
	
	return
}

// Symbol calls the symbol method in the solidity contract
func (e *ERC20) Symbol(block ...ethgo.BlockNumber) (retval0 string, err error) {
	var out map[string]interface{}
	var ok bool

	out, err = e.c.Call("symbol", ethgo.EncodeBlock(block...))
	if err != nil {
		return
	}

	// decode outputs
	retval0, ok = out["0"].(string)
	if !ok {
		err = fmt.Errorf("failed to encode output at index 0")
		return
	}
	
	return
}

// TotalSupply calls the totalSupply method in the solidity contract
func (e *ERC20) TotalSupply(block ...ethgo.BlockNumber) (retval0 *big.Int, err error) {
	var out map[string]interface{}
	var ok bool

	out, err = e.c.Call("totalSupply", ethgo.EncodeBlock(block...))
	if err != nil {
		return
	}

	// decode outputs
	retval0, ok = out["0"].(*big.Int)
	if !ok {
		err = fmt.Errorf("failed to encode output at index 0")
		return
	}
	
	return
}

// txns

// Approve sends a approve transaction in the solidity contract
func (e *ERC20) Approve(spender ethgo.Address, value *big.Int) (contract.Txn, error) {
	return e.c.Txn("approve", spender, value)
}

// Transfer sends a transfer transaction in the solidity contract
func (e *ERC20) Transfer(to ethgo.Address, value *big.Int) (contract.Txn, error) {
	return e.c.Txn("transfer", to, value)
}

// TransferFrom sends a transferFrom transaction in the solidity contract
func (e *ERC20) TransferFrom(from ethgo.Address, to ethgo.Address, value *big.Int) (contract.Txn, error) {
	return e.c.Txn("transferFrom", from, to, value)
}

// events

func (e *ERC20) ApprovalEventSig() ethgo.Hash {
	return e.c.GetABI().Events["Approval"].ID()
}

func (e *ERC20) TransferEventSig() ethgo.Hash {
	return e.c.GetABI().Events["Transfer"].ID()
}
