const { expect } = require("chai"); 
const hre = require("hardhat"); 
describe("ERC165Test", function () { 
    let erc165; //定义实例变量
    before(async () => { 
        // ⽣成合约实例并且复⽤ 
        erc165 = await hre.ethers.deployContract("ERC165Test", []); 
    }); 
    it("0x36372b07存在注册接口中", async () => {//erc20
        expect(await erc165.supportsInterface("0x36372b07")).to.be.true; //true
    }); 
    it("0xffffffff不在注册接口中", async () => {//
        expect(await erc165.supportsInterface("0xffffffff")).to.be.false; //false
    }); 
    
}); 