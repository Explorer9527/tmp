const { buildModule } = require("@nomicfoundation/hardhat-ignition/modules"); 
module.exports = buildModule("ShippingModule", (m) => { 
    const shipping = m.contract("Shipping", []); //获取Shipping实例
    m.call(shipping, "Status", []);//调用Shipping.Status函数
    return { shipping }; //返回当前商品状态
}); 