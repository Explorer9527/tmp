// SPDX-License-Identifier: MIT
pragma solidity ^0.8.28;
/*
address是一个 20 字节（共160位）的值，代表以太坊区块链上的一个账户地址。
    一般使用十六进制表示，则为40位：十六进制一位相当于二进制的4位

address类型在智能合约中的实际应用，如白名单机制、支付合约等。
*/
contract AddressType{
    address public defaultAddr; // 0x0000000000000000000000000000000000000000
    address myAddress2 = msg.sender;  // 当前合约调用者的地址

    //1.常规用法
    function demo()  public {
        //1.使用 == 和 != 进行 address 类型的比较。
        if(defaultAddr == myAddress2){
        }

        //2.address payable是可以接收以太币的地址类型。address类型不能直接发送以太币，必须显式转换
        address payable payableAddress = payable(myAddress2);//address转换为address payable

        //3.获取地址的以太坊余额（单位为 wei）
        uint256 balance = myAddress2.balance;  // 获取地址余额（单位：wei）

        //4.使用 transfer() 方法将以太币转移到另一个地址，推荐使用这种方法。
        address payable recipient = payable(myAddress2);
        recipient.transfer(1 ether);  // 转移 1 ETH

        //5.使用 send() 方法转移以太币，返回布尔值表示转移是否成功。由于没有自动回退机制，不推荐使用。
        bool success = recipient.send(1 ether);  // 转移 1 ETH，返回成功与否
        require(success, "Transfer failed.");

        //6.使用 call() 进行低级别调用，讨论其安全性问题以及与 send() 和 transfer() 的区别。
        (bool successs, ) = recipient.call{value: 1 ether}("");  // 使用 call 转移 1 ETH
        require(successs, "Transfer failed.");
    }
    //2.白名单机制
    mapping (address => bool) whitelist;//定义白名单映射
    function addToWhitelist(address _address) public {//加入白名单
        whitelist[_address] = true;//将某个地址加入白名单
    }
    function pay(address payable recipient) public payable {//授权支付合约
        //支付前先验证是否在白名单内，不在就报错
        require(whitelist[recipient], "Recipient is not whitelisted.");
        //如果在白名单内就允许转移以太币
        recipient.transfer(msg.value);
    }

    //3.防重入攻击
    mapping (address => uint) balances;
    function withdraw(uint256 amount) public {
        //检查余额是否足够
        require(balances[msg.sender] >= amount, "Insufficient balance.");
        // 更新余额
        balances[msg.sender] -= amount;
        // 转账
        (bool success, ) = msg.sender.call{value: amount}("");
        require(success, "Transfer failed.");
    }
}