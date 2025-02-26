// SPDX-License-Identifier: SEE LICENSE IN LICENSE
pragma solidity ^0.8.28;
import "@openzeppelin/contracts/utils/introspection/IERC165.sol";
/*
1.EIP与ERC的概念
    EIP:
        EIP（Ethereum Improvement Proposal，以太坊改进提案）是提出对以太坊网络进行改进的建议。
        EIP 包含协议更新、应用标准（如合约接口标准）等类别。提案流程：社区成员在 GitHub EIP 库中提交问题，
        经过讨论和共识后，提案可成为以太坊协议的一部分或社区推荐标准。
    ERC:
        ERC（Ethereum Request for Comment）是一种 EIP 提案，主要用于标准化合约接口，
        ERC 旨在为开发者提供通用接口和标准，但不具强制性。常见的 ERC 标准有 ERC20（同质化代币标准）
        和 ERC721（非同质化代币标准）。
2.ERC165的背景与用途
    ERC165 的定义：
        ERC165 是一个用于声明合约接口的标准。它允许合约查询自己或其他合约实现了哪些接口，
        使用 supportsInterface 函数返回是否支持指定的接口。主要用途：通过合约接口的查询功能，
        其他合约或应用可以通过 supportsInterface 判断该合约是否支持特定功能或接口。
        接口标识符计算：
            接口标识符 interfaceID 是对接口内所有函数的函数选择器进行异或（XOR）操作计算得出的。
            函数选择器是函数签名的 Keccak-256 哈希值的前 4 个字节。例如，对于一个简单的接口定义 
            interface IMyInterface { function myFunction() external returns (bool); }，
            先计算 myFunction() 的函数选择器，再得出该接口的标识符。
    ERC165 主要解决的问题：
        提供一种通用方式，使得智能合约能够以标准化的方式声明其实现了哪些接口，
        避免开发者和 DApp 需要手动确认合约是否兼容某些接口。

3.ERC165 的实际应用场景
    与其他 ERC 标准结合使用：
        ERC721、ERC1155 等标准依赖 ERC165 来声明自己实现了某些接口。比如，ERC721 的合约可以通过 ERC165 声明其支持
        ERC721 的接口 ID，从而方便钱包和其他应用查询。
    如何在实际项目中使用 ERC165：
        场景 1：当您开发一个合约，需要与多个不同类型的合约进行交互时，可以使用 supportsInterface 
               来动态检测目标合约是否支持某些功能。
        场景 2：开发 NFT 市场或钱包时，通过 ERC165 可以判断某个合约是否支持 ERC721 或 ERC1155 标准，
               以决定如何与该合约进行交互。

    此文件编写完毕后编译：npx hardhat compile
    测试环境启动：npx hardhat node
    部署：npx hardhat ignition deploy ignition/modules/ERC165Test.js --network localhost
    测试：npx hardhat test test/ERC165Test.js 

*/
contract ERC165Test is IERC165{

    mapping(bytes4 => bool) private _supportedInterfaces;

    constructor(){//注册多个接口ID
        _registerInterface(0x36372b07);//ERC20
        _registerInterface(0x80ac58cd );//ERC721
        _registerInterface(0xd9b67a26);//ERC1155
        _registerInterface(0x89afcb44);//IERC777
        _registerInterface(0x01ffc9a7);//IERC165
    }
    function _registerInterface(bytes4 interfaceId) internal {//注册逻辑
        require(interfaceId != 0xffffffff, "ERC165: invalid interface id");
        _supportedInterfaces[interfaceId] = true;
    }
    function supportsInterface(bytes4 interfaceId) external view returns (bool) {//对外提供验证接口
        return _supportedInterfaces[interfaceId];
    }
    function test() external view returns (bool) {
        return _supportedInterfaces[0x36372b07];
    }
}
