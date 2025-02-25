// SPDX-License-Identifier: SEE LICENSE IN LICENSE
pragma solidity ^0.8.28;
/*
在 Solidity 中，函数不仅是合约行为的基本单位，也是一种特殊的值类型。函数类型可用于变量声明、
作为参数传递和返回值。然而，与其他类型不同，函数类型在某些场景下的使用受到限制，特别是在 ABI 编码时。
重点：
    函数类型的声明与用法。
    函数选择器的概念与生成。
    函数类型的限制与解决方案。
*/
contract Functions{
    /*
    1.函数选择器的概念和生成
        函数选择器是通过对函数签名（函数名及其参数类型）进行 Keccak256 哈希计算，
        并截取前 4 个字节生成的唯一标识符。它用于识别和调用特定函数。
    */
    bytes4 selector1 = bytes4(keccak256("test(uint256)"));//获取函数aa的选择器
    bytes4 selector2 = this.test.selector;//Solidity内置获取选择器
    function test(uint256 i) external  {
    }
    //2.函数选择器的使用和限制
    function square(uint x) external returns (uint) {//定义一个平方函数
        return x**2;
    }
    function double(uint x) external returns (uint) {//定义一个乘法函数
        return x*2;
    }
    bytes4 ss = this.square.selector;//定义square函数选择器
    bytes4 ds = this.double.selector;//定义double函数选择器
    function select(uint8 flg,uint x) external returns (uint z) {//定义动态执行函数
        bytes4 s;
        if(flg == 1){//当输入1的时候使用square函数
            s = ss;
        }else if(flg == 2){//当输入2的时候使用double函数
            s = ds;
        }
        /*注意：
            函数类型变量的传递与返回可能会受到 ABI 编码器的限制，通常使用函数选择器代替。
            函数类型变量不能直接作为参数进行传递和返回，通常需要配合低级 call 来实现动态调用。
        */
        (bool success, bytes memory data) = address(this).call(abi.encodeWithSelector(s, x));
        require(success, "Function call failed");
        return abi.decode(data, (uint));
        
    }
    //3.函数选择器的限制和安全性
    // 限制：函数类型变量的传递与返回不能像其他基本类型那样自由，通常需要借助函数选择器来实现。
    // 安全性考虑：使用函数选择器与 call 时，需要确保调用的安全性，防止恶意代码执行。

 
 
}