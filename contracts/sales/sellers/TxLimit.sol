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
abstract contract TxLimit is Seller {
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

    function _beforePurchase(
        address to,
        uint256 num,
        uint256 cost
    )
        internal
        virtual
        override(Seller)
        returns (
            address,
            uint256,
            uint256
        )
    {
        (to, num, cost) = Seller._beforePurchase(to, num, cost);
        require(num <= _capOnTxLimit(to, num), "TxLimit: To many requested");
        return (to, num, cost);
    }

    function _afterPurchase(
        address to,
        uint256 num,
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

    function _capOnTxLimit(address to, uint256 num)
        internal
        view
        returns (uint256)
    {
        num = _maxPerTx == 0 ? num : Math.min(num, _maxPerTx);

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
        uint256 n,
        address addr,
        string memory zeroMsg
    ) private view returns (uint256) {
        uint256 remaining = _maxPerAddress - _bought[addr];
        require(
            remaining > 0,
            string(abi.encodePacked("TxLimiter: ", zeroMsg))
        );
        return Math.min(n, remaining);
    }
}
