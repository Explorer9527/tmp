// SPDX-License-Identifier: SEE LICENSE IN LICENSE
pragma solidity ^0.8.28;
/*
1.数组是一种用于存储相同类型元素的集合。在 Solidity 中，数组类型可以通过在数据类型后添加 [] 来定义。
  Solidity 支持两种数组类型：静态数组（Fixed-size Arrays）和动态数组（Dynamic Arrays）



*/
contract ArrayType{
    //1.存储（Storage）数组:存储在区块链上，生命周期与合约生命周期相同
    int[5] public a = [-1,2,3,4,5];//静态数组
    uint[] public b = [1,2,3];//动态数组

    function array(bytes calldata bb)  public {
        //2.内存（Memory）数组:临时存在于函数调用中，生命周期与函数相同，函数执行完毕后销毁
        uint[] memory c = new uint[](2);
        c[1] = 2;//赋值
        /*
        3.特殊数组类型： bytes 和 string，但注意：
              bytes 和 string 不支持使用下标索引进行访问。
              使用长度受限的字节数组时，建议使用 bytes1 到 bytes32 类型，以减少 gas 费用。
        */
        bytes memory bs = "abc\x22\x22";  // 通过十六进制字符串初始化
        bytes memory _data = new bytes(10);  // 创建一个长度为 10 的字节数组
        string memory str0;
        string memory str1 = "TinyXiong\u718A";  // 使用Unicode编码值
        /*
        4.数组成员属性和函数：
            length 属性：返回数组当前长度（只读），动态数组的长度可以动态改变。
            push()：用于动态数组，在数组末尾添加新元素并返回元素引用。
            pop()：用于动态数组，从数组末尾删除元素，并减少数组长度。
        */
        uint l = b.length;  // 获取数组长度
        b.push(1);  // 向数组添加元素
        b.pop();  // 删除数组最后一个元素
        delete a[1];//删除某个元素的值，不改变数组的长度
        /*删除数组中某个元素的逻辑：
            1.从这个元素位置开始，每个位置的值使用后一个元素替代，直到最后
            2.删除最后一个元素
        */
        for (uint256 i = 2; i < b.length - 1; i++) {
            b[i] = b[i + 1];//[1,3,4,5]删除位置2的元素4，首先将后续元素变为[1,3,5,5],
        }
        b.pop();//然后删除最后的元素[1,3,5]
        /*
        5.多维数组与数组切片
            支持多维数组，可以使用多个方括号表示，例如 uint[][5] 表示长度为 5 的变长数组的数组。
            数组切片是数组的一段连续部分，通过 [start:end] 的方式定义。
        */
        uint[][5] memory multiArray;  // 一个元素为变长数组的静态数组
        uint element = multiArray[2][1];  // 访问第三个动态数组的第二个元素
        bytes memory slice = bb[1:2];  // 创建数组切片
        bytes4 sig = abi.decode(bb[:4], (bytes4));  // 解码函数选择器
        address owner = abi.decode(bb[4:], (address));  // 解码地址
    }

}