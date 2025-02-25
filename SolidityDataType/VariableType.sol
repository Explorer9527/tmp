// SPDX-License-Identifier: MIT
pragma solidity ^0.8.28;
//Local variable
/*
 1.solidity属于静态语言，每个声明的变量都有一个基于类型的默认值，没有类似java中的null

 2.变量有三种类型：
    1.状态变量，合约内部
    2.局部变量，函数内部
    3.全局变量，全局命名空间内
 3.变量的作用域：变量的作用域指的是变量在代码中可见和有效的范围。Solidity 使用了 C99 作用域规则，变量从它们被声明后开始可见，直到包含它们的代码块 {} 结束。
                在 for 循环中初始化的变量，其作用域仅限于 for 循环的范围。
    局部变量作用域：在函数内部声明的变量为局部变量，其作用域仅限于函数内部或更小的代码块。
    状态变量作用域：状态变量是合约的一部分，定义在合约体内但不在任何函数内。状态变量的作用域分为三种类型：
        Public（公共）：公共状态变量可以在合约内部访问，也可以通过消息（如外部调用）访问。定义公共状态变量时，
                        Solidity 自动为其生成一个 getter 函数。
        Internal（内部）：内部状态变量只能在当前合约或其继承的子合约中访问，不能从外部直接访问。
        Private（私有）：私有状态变量只能在定义它们的合约内部访问，不能在子合约中访问。
 */
contract VariablesType {//合约名称
    uint16 public stateVariable = 1;//stateVariable状态参数
    address sender = msg.sender; //msg.sender全局参数

    function t() public view returns (uint) {
        uint104 localVariable = 1;//localVariable局部变量
        uint256 timestamp = block.timestamp;//block.timestamp全局参数
        return localVariable + timestamp;
    }
    /*
    为了规范变量名称的书写，在为变量命名时，需要记住以下规则：
        1.不应使用 Solidity 保留关键字作为变量名。例如，break 或 boolean，这类变量名是无效的。
        2.不应以数字(0-9)开头，必须以字母或下划线开头。例如，123test 是一个无效的变量名，但是_123test 是一个有效的变量名。变量名区分大小写。
        3.变量名区分大小写，例如，Name 和 name 是两个不同的变量。
    */
    //1.变量作用域举例
    uint public data = 30;         // 公共状态变量
    uint internal iData = 10;      // 内部状态变量
    function x() public returns (uint) {
        data = 3;                  // 内部访问公共变量
        return data;
    }
}
    contract Caller {//子合约
        VariablesType c = new VariablesType();
        function f() public view returns (uint) {
            return c.data();          // 外部访问公共变量
        }
    }
    contract D is VariablesType {//子合约
        function y() public returns (uint) {
            iData = 3;               // 派生合约内部访问内部变量
            return iData;
        }
    }