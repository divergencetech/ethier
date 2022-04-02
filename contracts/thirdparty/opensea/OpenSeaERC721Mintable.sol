// SPDX-License-Identifier: MIT
// Copyright (c) 2021 the ethier authors (github.com/divergencetech/ethier)
pragma solidity >=0.8.0 <0.9.0;

import "./OpenSeaGasFreeListing.sol";
import "../../utils/OwnerPausable.sol";
import "@openzeppelin/contracts/utils/Strings.sol";

/**
@notice An ERC721 extension that allows for minting directly from OpenSea using
"option IDs". See https://docs.opensea.io/docs/2-custom-item-sale-contract.
@dev All factory logic is abstracted away such that it's possible to deploy
without any interaction with the factory contract. Apart from requiring the
factory address for OpenSea listings, users of this contract can generally
ignore the existence of the extra contract.
 */
abstract contract OpenSeaERC721Mintable {
    /// @notice Factory contract deployed by this one's constructor.
    OpenSeaERC721Factory public factory;

    constructor(
        string memory factoryName,
        string memory factorySymbol,
        uint256 numFactoryOptions,
        string memory baseOptionURI
    ) {
        factory = new OpenSeaERC721Factory(
            factoryName,
            factorySymbol,
            baseOptionURI,
            msg.sender,
            numFactoryOptions
        );
    }

    /**
    @notice Returns whether the factory can currently mint the specified option.
     */
    function factoryCanMint(uint256 optionId)
        public
        view
        virtual
        returns (bool);

    /**
    @notice Mints the specified option for the recipient. Tip: although there is
    no field for the number of tokens to mint, this can be encoded in the option
    number.
    @dev Note that this has internal visibility and its access is subject to
    caller requirements so no further checks are necessary.
     */
    function _factoryMint(uint256 optionId, address to) internal virtual;

    /**
    @notice Mints the sp[ecified option for the recipient.
    @dev Only callable by the factory; instead use factory.mint().
     */
    function factoryMint(uint256 optionId, address to) external {
        require(
            msg.sender == address(factory),
            "OpenSeaERC721Mintable: only factory"
        );

        _factoryMint(optionId, to);
    }
}

/**
@notice Factory contract to mint OpenSeaERC721Mintable tokens.
@dev There is likely no need to use this contract directly; intead, inherit from
OpenSeaERC721Mintable and implement the necessary virtual functions.
 */
contract OpenSeaERC721Factory is OwnerPausable {
    using Strings for uint256;

    /// @notice Contract that deployed this factory.
    OpenSeaERC721Mintable public token;

    /// @notice Factory name and symbol.
    string private name_;
    string private symbol_;

    /// @notice Base URI for constructing tokenURI values for options.
    string private baseOptionURI;

    /**
    @notice Standard ERC721 Transfer event, used to trigger OpenSea into
    recognising the existence of the factory.
     */
    event Transfer(
        address indexed from,
        address indexed to,
        uint256 indexed tokenId
    );

    uint256 public immutable numOptions;

    /**
    @param owner Initial contract owner as it will be deployed by another
    contract but ownership should be transferred to an EOA.
     */
    constructor(
        string memory _name,
        string memory _symbol,
        string memory _baseOptionURI,
        address owner,
        uint256 _numOptions
    ) {
        _name = name_;
        _symbol = symbol_;
        token = OpenSeaERC721Mintable(msg.sender);
        setBaseOptionURI(_baseOptionURI);

        numOptions = _numOptions;

        super.transferOwnership(owner);
        emitTransfers(address(0), owner);
    }

    /// @notice Sets the base URI for constructing tokenURI values for options.
    function setBaseOptionURI(string memory _baseOptionURI) public onlyOwner {
        baseOptionURI = _baseOptionURI;
    }

    /// @notice Returns the factory name.
    function name() external view returns (string memory) {
        return name_;
    }

    /// @notice Returns the factory symbol.
    function symbol() external view returns (string memory) {
        return symbol_;
    }

    /**
    @notice Emits standard ERC721.Transfer events for all of the option "tokens"
    to induce correct OpenSea behaviour. These are first emitted upon contract
    deployment to signal "creation" of the option tokens, and on any ownership
    transfer of the contract.

     */
    function emitTransfers(address from, address to) internal {
        for (uint256 i = 0; i < numOptions; i++) {
            emit Transfer(from, to, i);
        }
    }

    /**
    @notice Transfers contract ownership just as with OpenZeppelin's Ownable,
    but also triggers Transfer events as OpenSea expects the option "tokens" to
    be owned by the contract owner.
     */
    function transferOwnership(address to) public override onlyOwner {
        emitTransfers(super.owner(), to);
        super.transferOwnership(to);
    }

    /**
    @notice Returns whether the option ID can be minted, deferring the logic to
    the factoryCanMint() method of the contract that deployed this factory.
     */
    function canMint(uint256 optionId) external view returns (bool) {
        return
            !paused() &&
            optionId < numOptions &&
            token.factoryCanMint(optionId);
    }

    /**
    @notice Returns a URL specifying option metadata, conforming to standard
    ERC721 metadata format.
     */
    function tokenURI(uint256 optionId) external view returns (string memory) {
        return string(abi.encodePacked(baseOptionURI, optionId.toString()));
    }

    /**
    @dev The OpenSea FactoryERC721 interface requires this instead of using
    EIP165 supportsInterface().
    @return true.
     */
    function supportsFactoryInterface() external pure returns (bool) {
        return true;
    }

    /**
    @notice Requires that the caller is either the owner or the owner's OpenSea
    Wyvern proxy, then proxies the call to the factoryMint() method of the
    contract that deployed this factory.
     */
    function mint(uint256 optionId, address to) public whenNotPaused {
        require(
            msg.sender == owner() ||
                msg.sender == OpenSeaGasFreeListing.proxyFor(owner()),
            "OpenSeaERC721Factory: only owner or proxy"
        );
        token.factoryMint(optionId, to);
    }

    /**
    @dev Calls mint(tokenId, to) to comply with OpenSea's overriding of the use
    of the ERC721 interface. Presumably they have hijacked their standard sale
    workflow to "transfer" a factory NFT to the buyer, which must have the
    effect of minting them a real NFT (this also explains the isApprovedForAll
    behaviour).
     */
    function transferFrom(
        address,
        address to,
        uint256 tokenId
    ) public {
        mint(tokenId, to);
    }

    /**
    @dev Returns true if owner is the contract owner, and either (a) operator is
    the OpenSea Wyver proxy for the owner; or (b) operator == owner. This is
    required to comply with OpenSea's overriding of the use of the ERC721
    interface. See comment on transferFrom().
     */
    function isApprovedForAll(address owner, address operator)
        public
        view
        returns (bool)
    {
        return
            owner == super.owner() &&
            (owner == operator ||
                OpenSeaGasFreeListing.isApprovedForAll(owner, operator));
    }

    /**
    @dev Always returns the contract owner. This is required to comply with
    OpenSea's overriding of the use of the ERC721 interface.
     */
    function ownerOf(uint256) public view returns (address) {
        return owner();
    }
}
