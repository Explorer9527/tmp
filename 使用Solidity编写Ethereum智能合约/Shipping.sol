// SPDX-License-Identifier: MIT
pragma solidity >=0.4.25 <0.9.0; 
/*
将创建的智能合约⽤于跟踪从在线市场购买的商品的状态。 
    创建该合约时，装运状态设置为 Pending 。 
    装运商品后，则将装运状态设置为 Shipped 并会发出事件。 
    交货后，则将商品的装 运状态设置为 Delivered ，并发出另⼀个事件。

此文件编写完毕后编译：npx hardhat compile
测试环境启动：npx hardhat node
部署：npx hardhat ignition deploy ignition/modules/Shipping.js --network localhost
测试：npx hardhat test test/Shipping.js 

*/
contract Shipping { 

    enum ShippingStatus { Pending, //0
                          Shipped, //1
                          Delivered //2
                        } //定义枚举值

    ShippingStatus private status; //定义状态变量status

    event LogNewAlert(string description); //使用event记录日志

    constructor(){ //创建合约时
        status = ShippingStatus.Pending;//装载设置为Pending
    } 

    function Shipped() public {//装运商品后
        status = ShippingStatus.Shipped; //装运状态设置为Shipped
        emit LogNewAlert("Your package has been shipped"); //并发出事件
    } 

    function Delivered() public {//交货后
        status = ShippingStatus.Delivered; //将商品的装 运状态设置为 Delivered 
        emit LogNewAlert("Your package has arrived"); //发出⼀个事件。
    } 

    //获取当前商品状态，返回对应的描述
    function getStatus(ShippingStatus _status) internal pure returns (string memory statusText) { 
        if (ShippingStatus.Pending == _status) 
            return "Pending"; 
        if (ShippingStatus.Shipped == _status) 
            return "Shipped"; 
        if (ShippingStatus.Delivered == _status)
            return "Delivered"; 
    } 

    //获取商品状态的只读函数入口，不需要传参
    function Status() public view returns (string memory) { 
        ShippingStatus _status = status; 
        return getStatus(_status); 
    }
 }