// SPDX-License-Identifier: SEE LICENSE IN LICENSE
pragma solidity ^0.8.28;
/*
Solidity中使用contract关键字定义合约，类似于其他编程语言中的类
*/
contract T {
    //1.this关键字
    function getAddress() public view returns (address) {
        return address(this);  // 返回当前合约的地址
    }
    //2.使用selfdestruct函数销毁合约，同时将合约中的以太币发送到指定地址
    function destroyContract(address payable recipient) public {
        selfdestruct(recipient);  // 销毁合约并发送以太币
    }
    //3.Solidity 0.6 版本开始，可以通过 type(C) 获取合约的类型信息
    function getContractInfo() public pure returns (string memory) {
        //存在循环引用，报错,猜测这里应该使用其他合约
        return (type(T).name);
        /*
        type(T).name: 获取合约的名字。
        type(T).creationCode: 获取创建合约的字节码。
        type(T).runtimeCode: 获取合约运行时的字节码。
        */
    }
    //4.通过 extcodesize 操作码判断一个地址是否为合约地址。
    function isContract(address addr) internal view returns (bool) {
        uint256 size;
        assembly { size := extcodesize(addr) }  // 获取地址的代码大小
        return size > 0;  // 大于 0 说明是合约地址
        /*
        在合约外部，可以使用 Web3.js 的 getCode 方法判断地址类型：
            web3.eth.getCode("0x1234567890123456789012345678901234567890").then(console.log);
            输出 "0x" 表示外部账号地址，其他字节码表示合约地址
        */
    }
    

}