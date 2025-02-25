// SPDX-License-Identifier: SEE LICENSE IN LICENSE
pragma solidity ^0.8.28;
/*
1.solidity中两个关键字：constant（常量）和immutable（不变量）。
    状态变量声明这个两个关键字之后，不能在合约后更改数值；并且还可以节省gas。
    只有数值变量可以声明constant和immutable；
    string和bytes可以声明为constant，但不能为immutable。
    constant​
2.区别
    constant变量必须在声明的时候初始化，之后再也不能改变。尝试改变的话，编译不通过。
    immutable变量可以在声明时或构造函数中初始化，因此更加灵活。

*/
contract Others {
    uint256 constant CONSTANT_NUM = 10;//不可变
    string constant CONSTANT_STRING = "0xAA";
    bytes constant CONSTANT_BYTES = "WTF";
    address constant CONSTANT_ADDRESS = 0x0000000000000000000000000000000000000000;

    uint256 public immutable IMMUTABLE_NUM = 9999999999;
    address public immutable IMMUTABLE_ADDRESS;
    uint256 public immutable IMMUTABLE_BLOCK;

    constructor(){//在构造方法中初始化immutable变量
        IMMUTABLE_ADDRESS = address(this);
        IMMUTABLE_BLOCK = block.number;
    }
}