// SPDX-License-Identifier: SEE LICENSE IN LICENSE
pragma solidity ^0.8.28;
/*
    理解 Solidity 中常用的全局变量及 API。
    学习如何通过全局变量获取区块和交易的相关信息。
    掌握 Solidity 中的 ABI 编码解码函数及其应用场景。
    了解 Solidity 中的数学和密码学 API，熟悉哈希算法与签名恢复。
    理解 Solidity 中的时间单位及其在智能合约中的应用。
*/
contract GlobalVariableAndAPI {
    //1.区块和交易属性API,Solidity提供了一些全局变量来获取当前区块链状态的信息，常用于获取区块信息和交易细节
    //1.1区块属性
    bytes32 hash = blockhash(100);//返回指定区块的哈希值，仅支持最近的 256 个区块，且不包括当前区块。
    address miner = block.coinbase;//返回挖出当前区块的矿工地址。
    uint difficulty = block.prevrandao; //返回当前区块的难度。difficulty被prevrandao替换
    uint gasLimit = block.gaslimit;//返回当前区块的 Gas 上限。
    uint blockNumber = block.number;//返回当前区块号。
    uint timestamp = block.timestamp;//返回当前区块的时间戳（单位：秒）。常用于时间条件判断。
    //1.2交易属性
    uint256 remainingGas = gasleft();//返回当前合约执行剩余的 Gas 数量。
    bytes data = msg.data;//返回当前调用的完整 calldata。
    address sender = msg.sender;//返回当前调用的发送者地址。
    bytes4 functionSelector = msg.sig;//返回当前调用的函数选择器。
    uint sentValue = msg.value;//返回此次调用发送的以太币数量（单位：wei）。
    uint gasPrice = tx.gasprice;//返回当前交易的 Gas 价格。
    address payable origin = tx.origin;//返回交易的最初发起者地址。如果只有一个调用，tx.origin 与 msg.sender 相同；否则，tx.origin 始终是最初的交易发起者。
    
    //2.ABI编码及解码函数API,ABI（应用二进制接口）函数用于编码和解码 Solidity 中的数据类型，特别适用于合约间交互时处理复杂数据结构
      //2.1编码函数
        //abi.encode(...) returns (bytes)：对输入的参数进行 ABI 编码，返回字节数组。
        bytes encodedData = abi.encode(uint(1), address(0x123));
        // abi.encodePacked(...) returns (bytes)：将多个参数进行紧密打包编码，不填充到 32 字节。适用于哈希计算。
        bytes packedData = abi.encodePacked(uint(1), address(0x123));
        // abi.encodeWithSelector(bytes4 selector, ...) returns (bytes)：将参数编码，并在前面加上函数选择器（用于外部调用）。
        bytes4 selector = bytes4(keccak256("transfer(address,uint256)")); 
        bytes encodedWithSelector = abi.encodeWithSelector(selector, address(0x123), 100);
        // abi.encodeWithSignature(string signature, ...) returns (bytes)：通过函数签名生成函数选择器，并将参数编码。
        bytes encodedWithSignature = abi.encodeWithSignature("transfer(address,uint256)", address(0x123), 100);
      //2.2 解码函数
        // abi.decode(bytes memory encodedData, (...)) returns (...)：对编码的数据进行解码，返回解码后的参数。
        // (uint a, address b) = abi.decode(encodedData, (uint, address));//函数中用
    //3.数学和密码学函数 API。Solidity提供了一些常用的数学与密码学函数，用于处理复杂运算和数据加密。
    //3.1数学函数
        // addmod(uint x, uint y, uint k) returns (uint)：计算 (x + y) % k，在任意精度下执行加法再取模，支持大数运算。
        uint result = addmod(10, 20, 7); // 结果为 2
        // mulmod(uint x, uint y, uint k) returns (uint)：计算 (x * y) % k，先进行乘法再取模。
        uint result = mulmod(10, 20, 7); // 结果为 6
    //3.2密码学哈希函数
        // keccak256(bytes memory) returns (bytes32)：使用 Keccak-256 算法计算哈希值（以太坊的主要哈希算法）。
        bytes32 hash = keccak256(abi.encodePacked("Hello, World!"));
        // sha256(bytes memory) returns (bytes32)：计算 SHA-256 哈希值。
        bytes32 hash = sha256(abi.encodePacked("Hello, World!"));
        // ripemd160(bytes memory) returns (bytes20)：计算 RIPEMD-160 哈希值，生成较短的 20 字节哈希值。
        bytes20 hash = ripemd160(abi.encodePacked("Hello, World!"));
    //3.3椭圆曲线签名恢复
        // ecrecover(bytes32 hash, uint8 v, bytes32 r, bytes32 s) returns (address)：通过椭圆曲线签名恢复公钥对应的地址，常用于验证签名。
        // address signer = ecrecover(hash, v, r, s);
    
    /*4.时间单位及其应用
      Solidity 提供了内置的时间单位，可以将秒转换为更常见的时间单位，如分钟、小时、天、周和年。
      这些单位在智能合约中可以用于表示和操作时间。时间单位及对应的秒数：
        1 seconds = 1 秒。
        1 minutes = 60 秒。
        1 hours = 3600 秒。
        1 days = 86400 秒。
        1 weeks = 604800 秒。
        1 years = 31536000 秒（非精确，忽略了闰年）。
    */
    uint public unlockTime; 
    address public owner;
    constructor(uint _lockTime) { 
        owner = msg.sender; 
        unlockTime = block.timestamp + _lockTime * 1 days; // 锁定指定天数 
    }
    function withdraw() public { 
        require(block.timestamp >= unlockTime, "Funds are locked."); 
        require(msg.sender == owner, "Only owner can withdraw."); // 执行提款操作
    }

}
