// SPDX-License-Identifier: SEE LICENSE IN LICENSE
pragma solidity ^0.8.28;
/*
枚举（enum）是solidity中用户定义的数据类型。它主要用于为uint分配名称，
使程序易于阅读和维护。它与C语言中的enum类似，使用名称来代替从0开始的uint。
枚举类型在Solidity合约中常用于定义状态机的状态，以及其他具有固定取值的常量。
enum的一个比较冷门的变量，很少使用。
*/
contract EnumsType {
    enum Status {
        Pending,//0
        Shipped//1
    }
    Status public status;//声明一个枚举变量，默认值为第一个值0
    function get() public view returns (Status) {//获取值
        return status;
    }
    function set(Status _status) public {//设置值
        status = _status;
    }
    function reset() public {
        delete status;//删除复位枚举值为0
    }
}