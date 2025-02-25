// SPDX-License-Identifier: SEE LICENSE IN LICENSE
pragma solidity ^0.8.28;
/*
1.Ether：
    目前以太币（Ether）主要分为这三个：wei、gwei以及ether。
    或许您之前还见过 finney 和 szabo，但这两个早在solidity 0.7.0就被删除了。
    而gwei却是solidity 0.6.11新添加的！
    使用方式：在数字后边跟上这些以太币单位，数字与以太币单位之间需要空格隔开的
2.gas：
    Gas 的工作原理：
        每当你执行一个操作，比如发送交易或执行智能合约，网络需要消耗计算资源来处理它。Gas就是这些计算工作的定价单位。
        想象你在玩一款视频游戏，每个动作都消耗一定量的能量点。
    Gas 限制和交易失败：
        如果为交易设置的Gas限制太低，交易可能因为“燃料耗尽”而失败。但你仍然需要为使用的Gas支付费用。
        这就像是你的赛车在终点线前没油了，虽然没完成比赛，但油费还是得付。
    Gas 价格波动：
        由于网络拥堵和其他因素，Gas价格会波动。在网络拥堵时，你可能需要支付更高的Gas价格来加快交易确认。
        想象在高峰期开车，你可能需要更多的油来应对交通拥堵。
    智能合约中的Gas优化：
        在智能合约编写过程中，优化代码以减少Gas消耗是至关重要的。这不仅节省费用，还提高了合约执行的效率。
*/
contract GasAndEther {
    //1.ether
    uint256 public oneWei = 1 wei;
    bool public isOneWei = (oneWei == 1);//最小单位
    uint256 public oneGwei = 1 gwei;
    bool public isOneGwei = (oneGwei == 1e9);//1 gwei = 1^9 wei
    uint256 public oneEther = 1 ether;
    bool public isOneEther = (oneEther == 1e18);//1 ether = 1^18 wei

    //2.gas
    function forever() public {//交易失败案例
        int i ;
        while (true) {//这里会耗尽gas
            i += 1;
        }
    }
}