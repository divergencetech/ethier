// SPDX-License-Identifier: MIT
// Copyright (c) 2021 the ethier authors (github.com/divergencetech/ethier)
pragma solidity >=0.8.0 <0.9.0;

import "@openzeppelin/contracts/utils/math/Math.sol";
import "./Seller.sol";

/// @notice A seller module to limit the number of purchased tokens based on
/// per-tx and/or per-buyer limits.
abstract contract TxLimit is Seller {
    /// @notice Max number of purchases per transaction.
    /// @dev `_maxPerTx = 0` means this limit will be ignored.
    uint64 private _maxPerTx;

    /// @notice max number of purchases per address.
    /// @dev `_maxPerAddress = 0` means this limit will be ignored.
    uint64 private _maxPerAddress;

    /// @notice Tracks the number of items already bought by an address
    /// @dev This isn't public as it may be skewed due to differences in msg.sender
    /// and tx.origin, which it treats in the same way such that
    /// `sum(_bought) >= sum(num)`.
    mapping(address => uint64) private _bought;

    constructor(uint64 maxPerTx_, uint64 maxPerAddress_) {
        _setTxLimits(maxPerTx_, maxPerAddress_);
    }

    /// @notice Set's new limits.
    function _setTxLimits(uint64 maxPerTx_, uint64 maxPerAddress_) internal {
        _maxPerTx = maxPerTx_;
        _maxPerAddress = maxPerAddress_;
    }

    /// @notice The maximum number of purchases per transaction.
    function maxPerTx() public view returns (uint64) {
        return _maxPerTx;
    }

    /// @notice The maximum number of purchases per address.
    function maxPerAddress() public view returns (uint64) {
        return _maxPerAddress;
    }

    // -------------------------------------------------------------------------
    //
    //  Internals
    //
    // -------------------------------------------------------------------------

    /// @notice Checks if the number of requested purchases is below the limits.
    /// @dev Reverts otherwise.
    function _beforePurchase(
        address to,
        uint64 num,
        uint256 cost
    )
        internal
        virtual
        override(Seller)
        returns (
            address,
            uint64,
            uint256
        )
    {
        (to, num, cost) = Seller._beforePurchase(to, num, cost);
        require(num <= _capOnTxLimit(to, num), "TxLimit: To many requested");
        return (to, num, cost);
    }

    /// @notice Updating the number of tokens bought by the purchaser.
    function _afterPurchase(
        address to,
        uint64 num,
        uint256 cost
    ) internal virtual override(Seller) {
        Seller._afterPurchase(to, num, cost);

        if (_maxPerAddress > 0) {
            _bought[to] += num;

            bool alsoLimitSender = msg.sender != to;

            if (alsoLimitSender) {
                _bought[msg.sender] += num;
            }

            // solhint-disable avoid-tx-origin
            bool alsoLimitOrigin = tx.origin != msg.sender && tx.origin != to;
            if (alsoLimitOrigin) {
                _bought[tx.origin] += num;
            } // solhint-enable avoid-tx-origin
        }
    }

    /// @notice Computes the maximum number of purchases that can be performed in
    /// the current transaction based on previously bought tokens.
    /// @dev This function can be used to dynamically adapt the number of purchased
    /// tokens in case to many are requested.
    function _capOnTxLimit(address to, uint64 num)
        internal
        view
        returns (uint64)
    {
        num = _maxPerTx == 0 ? num : uint64(Math.min(num, _maxPerTx));

        if (_maxPerAddress > 0) {
            num = _capExtra(num, to, "Buyer limit");

            bool alsoLimitSender = msg.sender != to;
            if (alsoLimitSender) {
                num = _capExtra(num, msg.sender, "Sender limit");
            }

            // solhint-disable avoid-tx-origin
            // TODO why not only tx.origin != msg.sender?
            bool alsoLimitOrigin = tx.origin != msg.sender && tx.origin != to;
            if (alsoLimitOrigin) {
                // solhint-disable-next-line avoid-tx-origin
                num = _capExtra(num, tx.origin, "Origin limit");
            } // solhint-enable avoid-tx-origin
        }

        return num;
    }

    /**
    @notice Returns min(n, max(extra items addr can purchase)) and reverts if 0.
    @param zeroMsg The message with which to revert on 0 extra.
     */
    function _capExtra(
        uint64 requested,
        address addr,
        string memory zeroMsg
    ) private view returns (uint64) {
        uint64 remaining = _maxPerAddress - _bought[addr];
        require(
            remaining > 0,
            string(abi.encodePacked("TxLimiter: ", zeroMsg))
        );
        return uint64(Math.min(requested, remaining));
    }
}
