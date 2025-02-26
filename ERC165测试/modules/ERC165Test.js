const { buildModule } = require("@nomicfoundation/hardhat-ignition/modules"); 
module.exports = buildModule("ERC165Test", (m) => { 
    const t = m.contract("ERC165Test", []); //获取ERC165Test实例
    m.call(t, "test", []);//调用ERC165Test.test函数
    return { t }; //返回当前商品状态
}); 