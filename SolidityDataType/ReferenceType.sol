// SPDX-License-Identifier: SEE LICENSE IN LICENSE
pragma solidity ^0.8.28;
/*
1.说明：
    在 Solidity 中，引用类型包括结构体（struct）、数组（array）和映射（mapping）。
    与值类型不同，引用类型在赋值时不会直接复制值，而是创建一个指向原数据的引用。这样可以避免对大型数据的多次拷贝，节省 Gas。
2.关键点：
    了解引用类型与值类型的区别。理解引用类型在合约中的应用场景及其优点。
3.数据位置：
    引用类型的一个重要特点是其数据存储位置。Solidity 中存在三种数据位置，每一种位置都有其独特的特点和使用场景。
    memory（内存）：
        数据仅在函数调用期间存在，函数调用结束后自动释放。用于局部变量，不能用于外部调用。
    storage（存储）：
        数据存储在合约的持久化存储中，状态变量默认存储在这里。只要合约存在，数据就一直保存在区块链上。
    calldata（调用数据）：
        用于存储函数参数的特殊数据位置。是一个不可修改的、非持久的存储区域，通常用于外部函数调用时的参数传递。
4.引用类型的赋值规则
    在 Solidity 中，不同的数据位置之间的赋值行为有所不同，这直接影响合约的执行效率和正确性。主要有以下几种：
        从 storage 到 memory：创建一份独立的拷贝。 示例：solidity uint[] storageArray = x; uint[] memory memArray = storageArray;
        从 memory 到 storage：创建一份独立的拷贝。 示例：solidity uint[] memory memArray = new uint[](10); x = memArray;
        从 memory 到 memory：只创建一个引用，更改其中一个变量会影响所有指向该数据的其他变量。 示例：solidity uint[] memory memArray1 = new uint[](10); uint[] memory memArray2 = memArray1;
        从 storage 到 storage：只创建一个引用，指向相同的存储位置。 示例：solidity uint[] storage y = x;
    赋值过程中的数据拷贝与引用：
        同一数据位置的赋值：通常只增加一个引用，多个变量指向同一个数据。
        跨数据位置的赋值：例如从 memory 到 storage，则会创建独立的拷贝。
5.数据位置与 Gas 消耗
    storage：永久保存合约状态变量，开销最大。使用 storage 来存储需要长期保存的数据，但要注意其较高的 Gas 消耗。
    memory：仅保存临时变量，函数调用结束后释放，开销较小。使用 memory 来存储临时数据，以减少合约的持久化存储开销。
    calldata：保存函数参数，几乎免费使用。如果可能，尽量使用 calldata 来存储数据，因为它既节省 Gas，又确保数据不可修改。
6.优化 Gas 消耗的建议
    使用 calldata：在外部函数调用中尽量使用 calldata 传递参数，节省 Gas。
    谨慎使用 storage：在必要时才使用 storage，并尽量减少不必要的状态变量修改。
    
*/
contract ReferenceType {
    uint[] x; //x是状态变量，其数据存储位置是 storage
    function f(uint[] memory memoryArray) public {
        x = memoryArray; // 将整个数组拷贝到 storage 中
        uint[] storage y = x; // 分配一个指针，y 的数据存储位置是 storage
        y[0] = 100; // 修改 y 的值，同时 x 的值也会改变
        delete x; // 清除 x，同时影响 y
    }
    //通过 g 函数和 h 函数展示如何在不同数据位置之间传递数据
    function g(uint[] storage) internal pure {}
    function h(uint[] memory) public pure {}

/**********************************************/

    uint[] public data; // 存储在storage中的动态数组
    function updateData(uint[] memory newData) public {// 将memory中的数组内容复制到storage中的data数组
        data = newData; // 从memory到storage的赋值，创建了独立的拷贝
    }
    function getData() public view returns (uint[] memory) {// 返回storage中的data数组
        return data; // 返回storage中的数据作为memory中的数组
    }
    function modifyStorageData(uint index, uint value) public {//修改storage中的data数组中指定索引位置的值
        require(index < data.length, "Index out of bounds");
        data[index] = value; // 修改storage中的值，开销较大
    }
    function modifyMemoryData(uint[] memory memData) public pure returns (uint[] memory) {//尝试修改传入的memory数组，并返回修改后的数组
        if (memData.length > 0) {
            memData[0] = 999; // 修改memory中的值，开销较小
        }
        return memData; // 返回修改后的memory数组
    }
}