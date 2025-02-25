// SPDX-License-Identifier: SEE LICENSE IN LICENSE
pragma solidity ^0.8.28;
/*
1.说明
    EVM有一个日志功能，用于将数据“写”到智能合约之外的数据结构中。其中一个重要的数据是Solidity事件。
    事件允许我们“打印”在区块链上的信息，这种方式比在智能合约中保存到公共存储变量更容易搜索，且更省gas费。
    当被发送事件（调用）时，会触发参数存储到交易的日志中（一种区块链上的特殊数据结构）。
    这些日志与合约的地址关联，并记录到区块链中.
2.事件有 3 中主要的使用场景：
    智能合约给用户的返回值
    异常触发
    更便宜的数据存储
*/
contract Events {
    //定义一个event
    event TestLog(address indexed sender, string message);
    //定义一个默认event
    event AnotherLog();

    //触发事件
    function test() public {
        emit TestLog(msg.sender, "Hello World!");
        emit TestLog(msg.sender, "Hello EVM!");
        emit AnotherLog();
    }
}