// SPDX-License-Identifier: SEE LICENSE IN LICENSE
pragma solidity ^0.8.28;
/*
1.说明：
    Solidity 在出现错误时通过回退状态来处理，即当合约在运行时发生异常，合约的状态会回滚到调用前的状态，
    同时终止交易。这种错误处理机制保证了区块链上每个事务的原子性。
2.异常处理方式：
    Solidity在出现错误时通过回退状态来处理，即当合约在运行时发生异常，合约的状态会回滚到调用前的状态，
    同时终止交易。这种错误处理机制保证了区块链上每个事务的原子性。
    Solidity 提供了以下主要函数来进行错误处理：
        assert(bool condition)：检查内部错误或逻辑错误。如果断言失败，状态将回滚，并消耗所有剩余的 Gas。用于检查合约内部逻辑的错误或不应该发生的情况，通常在函数末尾或状态更改之后使用。
        require(bool condition)：用于检查外部输入或调用条件。如果条件不满足，状态回滚，并返还剩余的 Gas。推荐使用
        revert()：立即终止交易并回滚状态
        revert(string memory reason)：终止交易并回滚状态，同时返回错误信息。
    Solidity 0.6.0 引入的 try/catch 结构来处理外部调用中的异常。
    Solidity 0.8.0 引入的自定义错误机制及其优势，减少 Gas 消耗的错误处理方式。
*/
contract Errors {
    address public owner;
    //1.assert和require 
    /*
    assert与require的区别
        Gas 消耗assert失败时会消耗掉所有的剩余Gas，而require则会返还剩余的Gas给调用者。
        适用场景：
            assert：用于检查合约内部逻辑的错误或不应该发生的情况，通常在函数末尾或状态更改之后使用。
            require：用于检查输入参数、外部调用返回值等，通常在函数开头使用。
        操作符不同：assert 失败时执行无效操作（操作码 0xfe），require 失败时则执行回退操作（操作码 0xfd）。
     assert 与 require 的最佳实践：
        优先使用 require()：
            用于检查用户输入或外部合约调用的返回值。
            适合在函数开始时检查前置条件。
        优先使用 assert()：
            用于检测不应该发生的内部错误。
            适合在函数结尾或状态改变后使用。
    */
    constructor() public {
        owner = msg.sender;
    }
    function transferOwnership(address newOwner) public {
        require(msg.sender == owner, "Only the owner can transfer ownership."); // 检查调用者是否为合约所有者
        owner = newOwner;
    }
    function checkBalance(uint a, uint b) public pure returns (uint) {
        uint result = a + b;
        assert(result >= a); // 检查溢出错误
        return result;
    }
    uint256 public balance;
    uint256 public constant MAX_UINT = 2 ** 256 - 1;
    function deposit(uint256 _amount) public {
        uint256 oldBalance = balance;
        uint256 newBalance = balance + _amount;
        // balance + _amount does not overflow if balance + _amount >= balance
        require(newBalance >= oldBalance, "Overflow");
        balance = newBalance;
        assert(balance >= oldBalance);
    }
    /*2.revert 函数
        revert()和revert(string memory reason)函数可以用于立即停止执行并回滚状态。
        这通常用于在遇到某些无法满足的条件时终止函数。
    */
    function checkValue(uint value) public pure {
        if (value > 10) {
            revert("Value cannot exceed 10"); // 返回自定义错误信息
        }
    }
    /*3.自定义错误
        在 Solidity 0.8.0 之后，Solidity 引入了自定义错误机制（custom errors），
        提供了一种更加 Gas 高效的错误处理方式。自定义错误比require或revert的字符串消息消耗更少的 Gas，
        因为自定义错误只传递函数选择器和参数。
        自定义错误的优势：自定义错误不会在错误消息中传递冗长的字符串，因此相比传统的require和revert，节省了更多的Ga
    */
    error Unauthorized(address caller);  // 自定义错误
    function restrictedFunction() public {
        if (msg.sender != owner) {
            revert Unauthorized(msg.sender);  // 使用自定义错误
        }
    }
    /*4.try/catch
        Solidity 0.6.0 版本后引入了 try/catch 结构，用于捕获外部合约调用中的异常。
        此功能允许开发者捕获和处理外部调用中的错误，增强了智能合约编写的灵活性。
        try/catch 的使用场景：
            捕获外部合约调用失败时的错误，而不让整个交易失败。
            在同一个交易中可以对失败的调用进行处理或重试。
    */
    function getValue() public pure returns (uint) {//被外部合约调用
        return 42;
    }
    function willRevert() public pure {//被外部合约调用
        revert("This function always fails");
    }
}
contract TryCatchExample {//try示例
    Errors externalContract;//定义合约引用
    constructor() {
        externalContract = new Errors();//合约实例初始化
    }
    function tryCatchTest() public returns (uint, string memory) {
        try externalContract.getValue() returns (uint a) {//使用try处理
            return (a, "Success");
        } catch {
            return (0, "Failed");
        }
    }
    function tryCatchWithRevert() public returns (string memory) {
        try externalContract.willRevert() {
            return "This will not execute";
        } catch Error(string memory reason) {
            return reason;  // 捕获错误信息
        } catch {
            return "Unknown error";
        }
    }
}