// SPDX-License-Identifier: MIT
// Copyright (c) 2021 the ethier authors (github.com/divergencetech/ethier)
pragma solidity >=0.8.0 <0.9.0;

// import "../utils/Monotonic.sol";
// import "../utils/OwnerPausable.sol";
import "@openzeppelin/contracts/utils/Address.sol";
import "@openzeppelin/contracts/utils/Context.sol";
import "@openzeppelin/contracts/utils/math/Math.sol";
import "@openzeppelin/contracts/utils/Strings.sol";
import "@openzeppelin/contracts/access/Ownable.sol";
import "./Seller.sol";

/**
@notice An abstract contract providing the _purchase() function to:
 - Enforce per-wallet / per-transaction limits
 - Calculate required cost, forwarding to a beneficiary, and refunding extra
 */
abstract contract TxLimit is Context {
    uint64 private _maxPerTx;
    uint64 private _maxPerAddress;

    /*
    @notice Tracks the number of items already bought by an address, regardless
    of transferring out (in the case of ERC721).
    @dev This isn't public as it may be skewed due to differences in msg.sender
    and tx.origin, which it treats in the same way such that
    sum(_bought)>=totalSold().
     */
    mapping(address => uint256) private _bought;

    constructor(uint64 maxPerTx_, uint64 maxPerAddress_) {
        _setTxLimits(maxPerTx_, maxPerAddress_);
    }

    function _setTxLimits(uint64 maxPerTx_, uint64 maxPerAddress_) internal {
        _maxPerTx = maxPerTx_;
        _maxPerAddress = maxPerAddress_;
    }

    function maxPerTx() public view returns (uint64) {
        return _maxPerTx;
    }

    function maxPerAddress() public view returns (uint64) {
        return _maxPerAddress;
    }

    // -------------------------------------------------------------------------
    //
    //  Internals
    //
    // -------------------------------------------------------------------------

    function _capRequested(address to, uint256 num)
        internal
        virtual
        returns (uint256)
    {
        num = _maxPerTx == 0 ? num : Math.min(num, _maxPerTx);

        if (_maxPerAddress > 0) {
            bool alsoLimitSender = _msgSender() != to;
            // solhint-disable-next-line avoid-tx-origin
            bool alsoLimitOrigin = tx.origin != _msgSender() && tx.origin != to;

            num = _capExtra(num, to, "Buyer limit");
            if (alsoLimitSender) {
                num = _capExtra(num, _msgSender(), "Sender limit");
            }
            if (alsoLimitOrigin) {
                // solhint-disable-next-line avoid-tx-origin
                num = _capExtra(num, tx.origin, "Origin limit");
            }

            _bought[to] += num;
            if (alsoLimitSender) {
                _bought[_msgSender()] += num;
            }
            if (alsoLimitOrigin) {
                // solhint-disable-next-line avoid-tx-origin
                _bought[tx.origin] += num;
            }
        }

        return num;
    }

    /**
    @notice Returns min(n, max(extra items addr can purchase)) and reverts if 0.
    @param zeroMsg The message with which to revert on 0 extra.
     */
    function _capExtra(
        uint256 n,
        address addr,
        string memory zeroMsg
    ) private view returns (uint256) {
        uint256 extra = _maxPerAddress - _bought[addr];
        if (extra == 0) {
            revert(string(abi.encodePacked("TxLimiter: ", zeroMsg)));
        }
        return Math.min(n, extra);
    }
}
