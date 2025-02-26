// SPDX-License-Identifier: SEE LICENSE IN LICENSE
pragma solidity ^0.8.28;
import "@openzeppelin/contracts/token/ERC1155/ERC1155.sol"; 
import "./ERC165Test.sol";
/*

*/

contract Test is ERC1155("TT"){

    bytes4 ERC1155_Interface_ID = 0xd9b67a26;
    ERC165Test t ;

    constructor() {
        t = new ERC165Test();
    }

    function checkIsERC1155()  external view returns (bool) {//可本地部署测试
        return t.supportsInterface(ERC1155_Interface_ID);
    }

}
