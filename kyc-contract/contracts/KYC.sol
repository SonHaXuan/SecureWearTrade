// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

import "@openzeppelin/contracts/token/ERC20/ERC20.sol";
import "@openzeppelin/contracts/token/ERC721/ERC721.sol";
import "@openzeppelin/contracts/token/ERC721/extensions/ERC721URIStorage.sol";
import "@openzeppelin/contracts/access/Ownable.sol";
import "hardhat/console.sol";

contract WasteManagement is Ownable, ERC721URIStorage {
    IERC20 _token;
    address public _owner;
    mapping(uint256 => bool) public processed;
    mapping(uint256 => uint256) public wasteBatchWeight;
    event WasteBatchProcessed(address processor, uint256 weight, uint256 id);

    constructor() ERC721("NFT Waste Batch", "WASTE") {
        _owner = msg.sender;
        _token = IERC20(_owner);
    }

    function mint(uint256 _tokenId, string memory metadataURI)
        external
        onlyOwner
    {
        // processed[_tokenId] = true;
        // wasteBatchWeight[_tokenId] = _weight;

        super._safeMint(_owner, _tokenId);
        super._setTokenURI(_tokenId, metadataURI);
    }

    function transfer(address _to, uint256 _tokenId) external payable {
        _validate(_tokenId);

        // emit WasteBatchProcessed(address(this), wasteBatchWeight[_tokenId], _tokenId);
        emit WasteBatchProcessed(_owner, 1, _tokenId);
        super._transfer(_owner, _to, _tokenId);

        // get owner of the token
        address payable sendTo = payable(ownerOf(_tokenId));
        sendTo.transfer(msg.value);
    }

    function _validate(uint256 _id) internal view {
        require(_exists(_id), "Error, wrong Token id");
        require(_ownerOf(_id) == _owner, "You are not the owner of the Nft");
        // require(processed[_id], "Error, Token is already processed");
        // require(msg.value >= wasteBatchWeight[_id], "Error, Insufficient processing fee");
    }
}
