// SPDX-License-Identifier: SEE LICENSE IN LICENSE
pragma solidity ^0.8.28;
/*
1.定义：
    Solidity 中的映射类型在功能上类似于 Java 的 Map 和 Python 的 Dict。映射是一种键值对存储结构，用于根据键快速访问值。
    定义形式：mapping(KT => KV) mappingName;
2.区别：
    在Java中，Map提供了丰富的功能，如键集合、值集合、长度检查等。
    在Python中，Dict也具有类似功能，并且可以动态添加或删除键值对。
    Solidity的映射则相对简单，没有键集合或值集合的概念，也无法直接获取映射的长度。
3.应用场景：
    代币合约：使用映射存储账户余额。
    游戏合约：使用映射存储玩家的等级或其他属性。
*/
contract MappingType {
    
    mapping(uint => string) public idName;
    mapping(address => uint) public balances;
    mapping(address => uint) public userLevel;
    uint balance = balances[msg.sender];//获取余额
    mapping(int => mapping(string => int)) public nested;//多重嵌套映射
    /*
    1.映射的局限性
        键类型的限制：键类型不能是映射、变长数组、合约、枚举或结构体。
            例如，以下定义是非法的：
            mapping(mapping(uint => string) => uint) illegalMapping;
        值类型的无限制：值类型可以是任何类型，包括映射类型。没有长度和键集合/值集合的概念：
            Solidity中的映射没有内建的键集合或值集合，也无法获取映射的长度。这与Java和Python中的映射结构不同。
        删除操作的特殊性：
            从映射中删除一个键的值，只需使用delete关键字，但键本身不会被移除，只是值被重置为默认值。
            delete balances[userAddr];
        访问操作：映射中的值访问方式与数组类似，如balances[userAddr]。但映射没有提供长度或集合访问的方法。
    */
   function test() public {
        //映射基本操作
        string memory str = idName[1];//获取映射中的值
        idName[1] = "123";//赋值
        delete idName[1];//只删除里面的值，恢复为默认值
        //多重嵌套映射操作
        int a  = nested[1]["sd"];//获取多重嵌套中的值
        nested[1]["a"] = 33;//多重嵌套赋值
        delete nested[1]["a"];//删除元素值
   }

}