// Package revert provides means for testing Ethereum reverted execution.
package revert

import (
	"fmt"
	"strings"

	"github.com/ethereum/go-ethereum/crypto"

	"github.com/ethereum/go-ethereum/common"
	"github.com/h-fam/errdiff"
)

// A Checker checks that a transaction reverts with the specified string. The
// empty string is valid and only checks that a transaction reverted, but is
// agnostic to the message. Therefore the empty string IS NOT equivalent to a
// nil error.
type Checker string

const Any = Checker("execution reverted")

// Checkers for OpenZeppelin modifiers.
const (
	AlreadyInitialized = Checker("Initializable: contract is already initialized")
	OnlyOwner          = Checker("Ownable: caller is not the owner")
	Paused             = Checker("Pausable: paused")
	Reentrant          = Checker("ReentrancyGuard: reentrant call")
)

// MissingRole returns a Checker that checks that a transaction reverts because the account requires but does not have the specified role.
func MissingRole(account common.Address, role [32]byte) Checker {
	return Checker(fmt.Sprintf("AccessControl: account %s is missing role 0x%x", strings.ToLower(account.Hex()), role))
}

// MissingRoleByName returns a Checker that checks that a transaction reverts because the account requires but does not have the specified role name.
func MissingRoleByName(account common.Address, name string) Checker {
	h := crypto.Keccak256([]byte(name))
	var role [32]byte
	copy(role[:], h)
	return MissingRole(account, role)
}

// Checkers for ethier libraries and contracts.
const (
	ERC721ApproveOrOwner = Checker("ERC721ACommon: Not approved nor owner")
	InvalidSignature     = Checker("SignatureChecker: Invalid signature")
	NotStarted           = Checker("LinearDutchAuction: Not started")
	SoldOut              = Checker("Seller: Sold out")
)

// Checkers for wETH test double. Use the wethtest package to deploy a modified
// wETH contract that includes these revert messages.
const (
	WETHWithdraw  = Checker("WETH9Test: insufficient balance to withdraw")
	WETHTransfer  = Checker("WETH9Test: insufficient balance to transfer")
	WETHAllowance = Checker("WETH9Test: insufficient allowance")
)

// Diff returns a message describing the difference between err and the Checker
// string, using substring matching. The empty-string Checker is treated as
// DefaultChecker.
//
// The first argument to Diff is ignored but is present to allow transaction
// functions to be used directly as input, without assigning to intermediate
// variables. For example:
//
//  if diff := ethtest.OnlyOwner.Diff(contract.Foo(…)); diff != "" {
//	  t.Errorf("contract.Foo() %s", diff)
//  }
func (c Checker) Diff(_ interface{}, err error) string {
	if c == "" {
		c = Any
	}
	return errdiff.Substring(err, string(c))
}
