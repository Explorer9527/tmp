// SPDX-License-Identifier: SEE LICENSE IN LICENSE
pragma solidity ^0.8.28;
/*
1.函数
    定义：
        function 函数名(< 参数类型 > < 参数名 >) < 可见性 > < 状态可变性 > [returns(< 返回类型 >)] {
        // 函数体
        }
    参数：函数可以包含输入参数、输出参数、可见性修饰符、状态可变性修饰符和返回类型。
    自由函数：函数不仅可以在合约内部定义，还可以作为自由函数在合约外部定义。
             自由函数的使用可以帮助分离业务逻辑，使代码更具模块化。
*/

contract FunctionUse{

   //1.函数返回值详解
    function returnMany() public pure returns (uint256, bool, uint256) {//返回多个参数
        return (1, true, 2);
    }
    function named() public pure returns (uint256 x, bool b, uint256 y) {//返回多个参数，return中有变量名称
        return (1, true, 2);
    }
    function assigned() public pure returns (uint256 x, bool b, uint256 y) {//赋值返回
        x = 1;
        b = true;
        y = 2;
    }
    function destructuringAssignments() public pure returns (uint256, bool, uint256, uint256, uint256){
        (uint256 i, bool b, uint256 j) = returnMany();//结构化赋值
        (uint256 x,, uint256 y) = (4, 5, 6);
        return (i, b, j, x, y);
    }
    /*
    2.可见性修饰符：Solidity中的函数可见性修饰符有四种，决定了函数在何处可以被访问：
        private （私有）：只能在定义该函数的合约内部调用。
        internal （内部）：可在定义该函数的合约内部调用，也可从继承该合约的子合约中调用。
        external （外部）：只能从合约外部调用。如果需要从合约内部调用，必须使用this关键字。
        public （公开）：可以从任何地方调用，包括合约内部、继承合约和合约外部。
    */
    function privateFunction() private pure returns (string memory) {
        return "Private";
    }
    function internalFunction() internal pure returns (string memory) {
        return "Internal";
    }
    function externalFunction() external pure returns (string memory) {
        return "External";
    }
    function publicFunction() public pure returns (string memory) {
        return "Public";
    }

    //3.函数可变性修饰符
    int8 public t0 = 10;//定义一个变量
    /*
    3.1 view关键字用于声明函数是视图函数，即函数不修改合约的状态变量，但可以读取合约的状态。
                视图函数用于查询合约状态或计算结果，而不会改变合约的状态。
        视图函数不会改变合约状态，也不会发送交易或调用其他合约。
        视图函数可以读取合约的状态变量和其他视图函数的返回值。
        视图函数内部不能修改状态变量的值。
    使用view关键字可以提供以下好处：
        在编译时进行静态检查，确保函数不会修改状态。
        允许在函数中访问合约的状态，并进行相应的计算和查询操作。
        允许 Solidity 编译器进行更多的优化。
    */
    function name() public view returns(int8 s){//定义一个只读函数
        // t0 = 11;//不能修改
        return t0;
    }
    /*
    3.2 pure关键字用于声明函数是纯函数，即函数不读取或修改合约的状态变量，并且不与外部合约进行交互。只根据输入参数计算结果，并返回一个值。
        纯函数不会改变合约状态，也不会发送交易或调用其他合约。
        纯函数内部不能访问 msg、block 和 tx 这些全局变量。
        纯函数在不同的块上执行时，给定相同的输入参数，总是返回相同的结果。
    使用pure关键字可以提供以下好处：
        在编译时进行静态检查，确保函数不会修改状态或与外部合约交互。
        提供更好的可读性和可理解性，明确函数的行为和约束。
        允许 Solidity 编译器进行更多的优化，提高代码执行效率。
    */
    function add(uint256 i, uint256 j) public pure returns (uint256) {
        // int8 b = a;//报错
        return i + j;
    }
    //3.3 payable声明函数可以接受以太币，如果没有该修饰符，函数将拒绝任何发送到它的以太币。
    function deposit() external payable {
    // 接收以太币
    }

}
