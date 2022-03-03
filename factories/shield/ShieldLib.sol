// SPDX-License-Identifier: MIT
// Copyright (c) 2022 the ethier authors (github.com/divergencetech/ethier)
pragma solidity >=0.8.0 <0.9.0;

library ShieldLib {
    /**
    @dev The value to which the release block is set when an asset is frozen;
    i.e. max(uint256) = (2^256 - 1).
     */
    uint256 private constant FROZEN = ~uint256(0);

    /**
    @dev Criteria dictating when an asset MAY be reclaimed. When an asset is in
    a frozen state, earliestBlock is set to max(uint256), completely locking it
    within this contract. When an asset is unfrozen it begins a thawing period
    and earliestBlock is set to (block.number + thawPeriod). Reclaiming an asset
    MUST NOT be allowed if block.number < earliestBlock.
     */
    struct ReleaseCriteria {
        uint256 earliestBlock;
        uint256 thawPeriod;
    }

    /**
    @dev Sets the assets's earliest release block to max(uint256) and stores the
    thawing period for when unfreeze() is called. freeze() can be called at any
    time, even while already frozen or during the thaw period, but it MUST NOT
    reduce the shortest time to release and will otherwise revert.
    @param thawPeriod When unfreeze() is called, the asset's release block will
    be set to (block.number + thawPeriod).
     */
    function freeze(ReleaseCriteria storage release, uint256 thawPeriod)
        internal
    {
        uint256 soonest = release.earliestBlock == FROZEN
            ? release.thawPeriod
            : release.earliestBlock - block.number;
        require(thawPeriod >= soonest, "ShieldLib: thaw reduction");

        release.earliestBlock = FROZEN;
        release.thawPeriod = thawPeriod;
    }

    /**
    @dev Sets the assets's release block to current plus `thawPeriod` as
    passed to freeze(). After that many blocks, requireThawed() will no longer
    revert.
    @return The earliest block at which requireThawed() will no longer revert.
     */
    function unfreeze(ReleaseCriteria storage release)
        internal
        returns (uint256)
    {
        require(release.earliestBlock == FROZEN, "ShieldLib: not frozen");

        uint256 earliest = block.number + release.thawPeriod;
        release.earliestBlock = earliest;
        return earliest;
    }

    /**
    @dev Requires that release.earliestBlock has already been reached.
    @dev MUST be called by asset-reclaiming functions.
     */
    function requireThawed(ReleaseCriteria storage release) internal view {
        require(release.earliestBlock <= block.number, "ShieldLib: frozen");
    }
}
