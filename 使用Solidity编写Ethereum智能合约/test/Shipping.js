const { expect } = require("chai"); 
const hre = require("hardhat"); 
describe("Shipping", function () { 
    let shippingContract; //定义实例变量
    before(async () => { 
        // ⽣成合约实例并且复⽤ 
        shippingContract = await hre.ethers.deployContract("Shipping", []); 
    }); 
    it("should return the status Pending", async function () {//检查初始状态是否为Pending
        expect(await shippingContract.Status()).to.equal("Pending");//调用Shipping的Status方法
    }); 
    it("should return the status Shipped", async () => {//商品状态变更为Shipped
        await shippingContract.Shipped(); ////调用Shipping的Shipped方法，修改上屏状态为1
        expect(await shippingContract.Status()).to.equal("Shipped"); //调用Shipping的Status方法
    }); 
    it("should return correct event description", async () => { //商品状态变更为Delivered
        await expect(shippingContract.Delivered()) // 验证事件是否被触发 
        .to.emit(shippingContract, "LogNewAlert") // 验证事件的参数是否符合预期 
        .withArgs("Your package has arrived"); 
    }); 
}); 