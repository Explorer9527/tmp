package main

import (
	"bytes"
	"context"
	"crypto/ecdsa"
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"log"
	"math/big"
	"os"
	"regexp"
	"runtime/debug"
	"strings"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/learn/init_order/store"
	"golang.org/x/crypto/sha3"
)

/*
如果go环境中没有import的包，执行go mod tidy会自动下载相关依赖包
下载完后编译：cd /Test1 && go build
*/
func main() {
	testFlg := 11 //测试的分支
	var client *ethclient.Client
	var err error
	/*1.创建一个以太坊客户端，用于与指定的以太坊节点建立连接，
		1.1连接公共节点(这里是连接到infura网关)
	        https://cloudflare-eth.com 为公共以太坊节点，开发者可以免费访问以太坊区块链数据，其局限性：
			主要面向以太坊主网，测试网支持较少（需确认具体端点）。
			仅支持 只读操作（如查询余额、交易历史）。
			不支持发送交易（需自己运行节点或使用钱包服务）。
			可能有速率限制（频繁请求可能被封禁）。
			HTTPS（请求-响应模式）、每次请求需建立新连接，延迟较高、适合低频查询（如获取单个区块信息）
			适用于简单查询、低吞吐量的主网交互（如获取区块数据、交易状态）场景。
		1.2通过 IPC（Inter-Process Communication）文件连接到本地运行的以太坊节点（如 Geth）
			client, err := ethclient.Dial("/home/user/.ethereum/geth.ipc")//Geth 节点默认的 IPC 文件路径。
			本地通信：仅适用于 同一台机器 上的进程间通信，无法跨网络。
			高效性：相比 HTTP/HTTPS，IPC 的延迟更低，适合高频操作。
			安全性：无需暴露节点端口到公网，减少被攻击风险。
			依赖节点：必须确保本地 Geth 节点正在运行，否则会连接失败。
			典型场景：开发需要高频交互或敏感操作（如管理私钥）的以太坊应用。通过 client 对象执行需要 本地节点权限 的操作（如发送交易、部署合约）。
		1.3通过Infura服务连接
			包括主网（如以太坊主网）和测试网（如 Ropsten、Goerli）。
			适用场景：高频交易、智能合约开发、需要稳定节点连接的 DApp。
			全双工实时通信、实时推送（如交易通知、区块更新）、适合高频操作（如订阅事件、交易监控）
	*/
	client, err = ethclient.Dial("https://cloudflare-eth.com") //HTTPS 连接主网
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(client, err)
	client0, err0 := ethclient.Dial("wss://ropsten.infura.io/ws") //WebSocket连接Ropsten测试网
	if err0 != nil {
		log.Fatal(err0)
	}
	//2.获取账户地址
	address := common.HexToAddress("0x71c7656ec7ab88b098defb751b7401b5f6d8976f")
	fmt.Println(address.Hex())   // 0x71C7656EC7ab88b098defB751B7401B5f6d8976F
	fmt.Println(address.Bytes()) // [113 199 101 110 199 171 136 176 152 222 251 117 27 116 1 181 246 216 151 111]
	if testFlg == 1 || testFlg == 2 {
		fmt.Println("1或2测试")
	} else if testFlg == 3 {
		//3.读取账户余额
		blance, err1 := client.BalanceAt(context.Background(), address, nil) //查询当前实时余额
		if err1 != nil {
			log.Println(string(debug.Stack()))
			fmt.Println(err1)
			// log.Fatal(err1) //有异常会跳出，暂时屏蔽
		}
		fmt.Println(blance)                                                          //25893180161173005034单位为wei
		balance2 := big.NewInt(5532993)                                              //通过指定区块号，可以获取账户在 过去某个时间点 的余额（而非当前最新值）。
		fmt.Println(balance2)                                                        // 25729324269165216042单位为wei
		eth := new(big.Float).Quo(new(big.Float).SetInt(blance), big.NewFloat(1e18)) //转化为eth，Quo浮点除法
		fmt.Println(eth)                                                             // 25.729324269165216042单位为eth
	} else if testFlg == 4 {
		//4.代币余额查询，需要等后面知识，暂时搁置
	} else if testFlg == 5 {
		//5.生成新钱包 - 依赖go-ethereum crypto 包
		privateKey, err1 := crypto.GenerateKey() //用于生成随机私钥
		if err1 != nil {
			log.Fatal(err1)
		}
		privateKeyBytes := crypto.FromECDSA(privateKey)    //使用 FromECDSA 方法将其转换为字节。
		fmt.Println(hexutil.Encode(privateKeyBytes)[2:])   //hexutil包将它转换为十六进制字符串,2:表示从字符串的第二个字符开始截取，即去掉 "0x" 前缀。
		publicKey := privateKey.Public()                   // 从私钥生成公钥
		publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey) // 将公钥转换为 ECDSA 公钥类型
		if !ok {
			log.Fatal("cannot assert type: publicKey is not of type *ecdsa.PublicKey")
		}
		publicKeyBytes := crypto.FromECDSAPub(publicKeyECDSA) // 使用crypto包的FromECDSAPub方法将公钥转换为字节
		fmt.Println(hexutil.Encode(publicKeyBytes)[4:])       //hexutil包将它转换为十六进制字符串,4:表示从字符串的第四个字符开始截取，即去掉 "0x" 前缀。
		address2 := crypto.PubkeyToAddress(*publicKeyECDSA)   // 从公钥生成地址
		fmt.Println(address2.Hex())                           // 0x96216849c49358B10257cb55b28eA603c874b05E

		hash := sha3.NewLegacyKeccak256()                     // 手动计算公钥地址，创建一个 Keccak-256 哈希对象
		hash.Write(publicKeyBytes[1:])                        // 跳过首字节（通常是 0x04）
		address3 := common.BytesToAddress(hash.Sum(nil)[12:]) // 从哈希结果中提取地址,取后 20 字节
		fmt.Println(address3.Hex())                           // 0x96216849c49358B10257cb55b28eA603c874b05E
	} else if testFlg == 6 {
		/*6.keystore，是一个包含经过加密了的私钥和元数据的 JSON 文件，用于安全地存储以太坊账户信息。
		    go-ethereum 中的 keystore，每个文件只能包含一个钱包密钥对
		    Keystore 的工作流程：
				1.生成 Keystore：
				  钱包生成私钥 → 用密码加密私钥 → 生成 Keystore JSON 文件。
				2.使用 Keystore：
				  用户输入密码 → 钱包用 Keystore 中的参数解密私钥 → 用私钥签名交易。
			JSON 格式，包含：
				地址：对应的以太坊地址。
				加密私钥：通过 scrypt 算法加密的私钥。
				盐值（salt）：随机生成的盐值，增强加密安全性。
				加密算法参数：如 n（迭代次数）、p（并行度）、r（块大小）。
		*/
		ks := keystore.NewKeyStore("./tmp", keystore.StandardScryptN, keystore.StandardScryptP) // 创建一个新的 keystore 实例
		password := "secret"                                                                    // 设置密码
		account, err1 := ks.NewAccount(password)                                                // 创建新账户
		if err1 != nil {
			log.Fatal(err1)
		}
		fmt.Println(account.Address.Hex()) // 0x527Fe43f1de42A7c5eBBd066E5D4847c31e92463

		file := "./wallets/UTC--2018-07-04T09-58-30.122808598Z--20f8d42fb0f667f2e53930fed426f225752453b3" // 指定 keystore 文件
		ks2 := keystore.NewKeyStore("./tmp", keystore.StandardScryptN, keystore.StandardScryptP)
		jsonBytes, err2 := ioutil.ReadFile(file) //读取keystore文件
		if err2 != nil {
			log.Fatal(err2) //有异常会跳出，暂时屏蔽
			fmt.Println(err2)
		}
		password = "secret"
		account2, err3 := ks2.Import(jsonBytes, password, password) // 导入账户
		if err3 != nil {
			log.Fatal(err3) //有异常会跳出，暂时屏蔽
			fmt.Println(err3)
		}
		fmt.Println(account2.Address.Hex()) // 0x3b703d0Ab69f0d938960a106641807201a298177
		err4 := os.Remove(file)             // 删除 keystore 文件
		if err4 != nil {
			log.Fatal(err4) //有异常会跳出，暂时屏蔽
			fmt.Println(err4)
		}
	} else if testFlg == 7 {
		//7.分层确定性(HD)钱包
		//https://github.com/miguelmota/go-ethereum-hdwallet
	} else if testFlg == 8 {
		//8.检查合约地址是否有效
		rege := regexp.MustCompile("^0x[0-9a-fA-F]{40}$")                           //检查地址是否符合以太坊地址的格式
		fmt.Println(rege.MatchString("0x39cb70F972E0EE920088AeF97Dbe5c6251a9c25D")) // true
		fmt.Println(rege.MatchString("0x39cb70F972E0EE920088AeF97Dbe5c6251a9c25X")) // false
		ec, err1 := ethclient.Dial("https://cloudflare-eth.com")                    //获取客户端对象
		if err1 != nil {
			log.Fatal(err1)
		}
		address1 := common.HexToAddress("0x39cb70F972E0EE920088AeF97Dbe5c6251a9c25D") //获取合约地址
		bytesCode, err2 := ec.CodeAt(context.Background(), address1, nil)             //获取合约代码
		if err2 != nil {
			log.Fatal(err2) //有异常会跳出，暂时屏蔽
			fmt.Println(err2)
		}
		fmt.Println(len(bytesCode) > 0) // true
	} else if testFlg == 9 {
		//9.查询区块信息
		blockNumber := big.NewInt(5671744)                                     //获取区块号
		block, err1 := client.BlockByNumber(context.Background(), blockNumber) //获取区块信息
		if err1 != nil {
			log.Fatal(err1) //有异常会跳出，暂时屏蔽
			fmt.Println(err1)
		}
		fmt.Println(block.Number().Uint64())                                       // 5671744 区块号
		fmt.Println(block.Time())                                                  // 1530515591 直接使用 block.Time()，它已经是 uint64 类型
		fmt.Println(block.Difficulty().Uint64())                                   // 17179869184 难度
		fmt.Println(block.Hash().Hex())                                            // 0x9e8751ebb5069389b855bba72d94902cc385042661498a415979b7b6ee9ba4b9
		fmt.Println(block.Transactions().Len())                                    // 1 交易数量
		count, err2 := client.TransactionCount(context.Background(), block.Hash()) // 获取区块的交易数量
		if err2 != nil {
			log.Fatal(err2)
		}
		fmt.Println(count) // 1
	} else if testFlg == 10 {
		//10.获取区块上的交易信息
		blockNumber := big.NewInt(5671744)                                     //获取区块号
		block, err1 := client.BlockByNumber(context.Background(), blockNumber) //获取区块信息
		if err1 != nil {
			log.Fatal(err1) //有异常会跳出，暂时屏蔽
			fmt.Println(err1)
		}
		for idx, tx := range block.Transactions() {
			fmt.Println(idx, "->", tx.Hash().Hex())                  // 0 -> 0x5d49fcaa394c97ec8a9c3e7bd9e8388d420fb050a52083ca52ff24b3b65bc9c2
			fmt.Println(tx.Value())                                  // 1000000000000000000 交易金额
			fmt.Println(tx.Gas())                                    // 21000 交易消耗的 gas
			fmt.Println(tx.GasPrice())                               // 1000000000000000000 gas 价格
			fmt.Println(tx.Nonce())                                  // 110644 交易 nonce
			fmt.Println(tx.Data())                                   // 0x 交易数据
			fmt.Println(tx.To().Hex())                               // 0x39cb70F972E0EE920088AeF97Dbe5c6251a9c25D 接收地址
			chaindId, err2 := client.NetworkID(context.Background()) // 查询当前连接的区块链网络的唯一标识符,不同网络的ChainID不同
			if err2 != nil {
				log.Fatal(err2)
			}
			signer := types.NewEIP155Signer(chaindId) // 创建签名器
			from, err3 := types.Sender(signer, tx)    // 获取交易的发送方地址
			if err3 != nil {
				log.Fatal(err3)
			}
			fmt.Println(from.Hex())                                                     // 发送地址
			receipt, err4 := client.TransactionReceipt(context.Background(), tx.Hash()) // 获取交易回执
			if err4 != nil {
				log.Fatal(err4)
			}
			fmt.Println(receipt.Status) // 交易状态
		}
		//获取区块上交易的第二种方式
		hash1 := common.HexToHash("0x9e8751ebb5069389b855bba72d94902cc385042661498a415979b7b6ee9ba4b9") //获取区块哈希
		count1, err5 := client.TransactionCount(context.Background(), hash1)
		if err5 != nil {
			log.Fatal(err5)
		}
		for i := uint(0); i < count1; i++ {
			tx, err6 := client.TransactionInBlock(context.Background(), hash1, i) // 获取区块中的交易
			if err6 != nil {
				log.Fatal(err6)
			}
			fmt.Println(tx.Hash().Hex()) // 0x5d49fcaa394c97ec8a9c3e7bd9e8388d420fb050a52083ca52ff24b3b65bc9c2
		}
		//根据交易hash获取交易信息
		txHash := common.HexToHash("0x5d49fcaa394c97ec8a9c3e7bd9e8388d420fb050a52083ca52ff24b3b65bc9c2") //交易hash
		tx, isPending, err7 := client.TransactionByHash(context.Background(), txHash)                    //获取交易信息
		if err7 != nil {
			log.Fatal(err7)
		}
		fmt.Println(tx.Hash().Hex()) // 0x5d49fcaa394c97ec8a9c3e7bd9e8388d420fb050a52083ca52ff24b3b65bc9c2
		fmt.Println(isPending)       // false
	} else if testFlg == 11 {
		//11.以太坊转账
		privateKey2, err1 := crypto.HexToECDSA("fad9c8855b740a0b7ed4c221dbad0f33a83a49cad6b3fe8d5817ac83d38b6a19") //私钥
		if err1 != nil {
			log.Println(string(debug.Stack()))
			log.Fatal(err1)
		}
		publicKey2 := privateKey2.Public()                   //私钥生成公钥
		publicKeyECDSA2, ok := publicKey2.(*ecdsa.PublicKey) //公钥转成公钥结构体PublicKey类型
		if !ok {
			log.Fatal("cannot assert type: publicKey is not of type *ecdsa.PublicKey")
		}
		address4 := crypto.PubkeyToAddress(*publicKeyECDSA2)                  //转换成公钥地址
		nonce, err2 := client0.PendingNonceAt(context.Background(), address4) //获取nonce
		if err2 != nil {
			log.Println(string(debug.Stack()))
			fmt.Println(err2)
			// log.Fatal(err2)//有异常会跳出，暂时屏蔽
		}
		value := big.NewInt(1e18)                                       // 1个以太币
		gasLimit := uint64(21000)                                       // 交易消耗的gas上限
		gasPrice, err3 := client0.SuggestGasPrice(context.Background()) //基于近几个块和交易池交易动态计算出较合理的gas价格
		if err3 != nil {
			log.Fatal(err3)
		}
		tx2 := types.NewTransaction(nonce, //交易序号
			common.HexToAddress("0x39cb70F972E0EE920088AeF97Dbe5c6251a9c25D"), //接收地址
			value,    //转账金额
			gasLimit, //gas上限
			gasPrice, //gas价格
			nil)      //创建交易
		chainID, err4 := client0.NetworkID(context.Background()) //获取网络ID
		if err4 != nil {
			log.Fatal(err4)
		}
		signedTx, err5 := types.SignTx(tx2, types.NewEIP155Signer(chainID), privateKey2) //对交易进行签名
		fmt.Println(signedTx.Hash().Hex())                                               // 0x8f2a55949038a9610f50fb23b5883af3b4ecb3c3bb792cbcefbd1542c692be63
		if err5 != nil {
			log.Fatal(err5)
		}
		err6 := client0.SendTransaction(context.Background(), signedTx) //广播交易到区块链网络
		if err6 != nil {
			log.Fatal(err6)
		}
	} else if testFlg == 12 {
		//12.代币转账
		privateKey2, err1 := crypto.HexToECDSA("fad9c8855b740a0b7ed4c221dbad0f33a83a49cad6b3fe8d5817ac83d38b6a19") //私钥
		if err1 != nil {
			log.Println(string(debug.Stack()))
			log.Fatal(err1)
		}
		publicKey2 := privateKey2.Public()                   //私钥生成公钥
		publicKeyECDSA2, ok := publicKey2.(*ecdsa.PublicKey) //公钥转成公钥结构体PublicKey类型
		if !ok {
			log.Fatal("cannot assert type: publicKey is not of type *ecdsa.PublicKey")
		}
		address4 := crypto.PubkeyToAddress(*publicKeyECDSA2)                  //转换成公钥地址
		nonce, err2 := client0.PendingNonceAt(context.Background(), address4) //获取nonce
		if err2 != nil {
			log.Println(string(debug.Stack()))
			fmt.Println(err2)
			// log.Fatal(err2)//有异常会跳出，暂时屏蔽
		}
		value := big.NewInt(0) //代币不需要转账金额
		toAddress := common.HexToAddress(
			"0x4592d8f8d7b001e72cb26a73e4fa1806a51ac79d") //代币的接收地址
		tokenAddress := common.HexToAddress(
			"0x28b149020d2152179873ec60bed6bf7cd705775d") //代币的转账合约地址
		transferFnSignature := []byte("transfer(address,uint256)") //代币转账函数
		hash := sha3.NewLegacyKeccak256()                          //获取函数签名对象
		hash.Write(transferFnSignature)                            //对函数方法进行哈希
		methodID := hash.Sum(nil)[:4]                              //获取函数名称id,只获取前四个字节
		fmt.Println(hexutil.Encode(methodID))                      // 0xa9059cbb
		paddedAddress := common.LeftPadBytes(methodID, 40)         //对函数名称id向左进行补0到40位
		fmt.Println(hexutil.Encode(paddedAddress))                 // 0x0000000000000000000000004592d8f8d7b001e72cb26a73e4fa1806a51ac79d
		amount := new(big.Int)                                     //代币转账金额
		amount.SetString("1000000000000000000", 10)                //转一个代币，单位是wei
		paddedAmount := common.LeftPadBytes(amount.Bytes(), 32)    //对转账金额进行补0到32位
		fmt.Println(hexutil.Encode(paddedAmount))                  //0x00000000000000000000000000000000000000000000003635c9adc5dea00000
		var data []byte
		data = append(data, methodID...)                                              //充填id
		data = append(data, paddedAddress...)                                         //充填地址
		data = append(data, paddedAmount...)                                          //充填转账金额
		gaslimit, err3 := client0.EstimateGas(context.Background(), ethereum.CallMsg{ //计算gas上限
			To:   &toAddress,
			Data: data,
		})
		if err3 != nil {
			log.Println(string(debug.Stack()))
			log.Fatal(err3)
		}
		fmt.Println(gaslimit)
		tx := types.NewTransaction(nonce, tokenAddress, value, gaslimit, nil, data) //构建交易
		chainID, err4 := client0.NetworkID(context.Background())                    //获取网络ID
		if err4 != nil {
			log.Fatal(err4)
		}
		signedTx, err5 := types.SignTx(tx, types.NewEIP155Signer(chainID), privateKey2) //对交易进行签名
		if err5 != nil {
			log.Println(string(debug.Stack()))
			log.Fatal(err5)
		}
		err6 := client0.SendTransaction(context.Background(), signedTx) //广播交易
		if err6 != nil {
			log.Println(string(debug.Stack()))
		}
	} else if testFlg == 13 {
		//13.订阅新区块
		headers := make(chan *types.Header)
		headersSub, err1 := client0.SubscribeNewHead(context.Background(), headers) //获取最新区块的订阅对象
		if err1 != nil {
			log.Println(string(debug.Stack()))
			log.Fatal(err1)
		}
		for { //循环获取信息
			select {
			case err2 := <-headersSub.Err(): //订阅出错
				log.Println(string(debug.Stack()))
				log.Fatal(err2)
			case header := <-headers: //订阅到新区块
				fmt.Println(header.Hash().Hex())    //打印新区块hash
				fmt.Println(header.Number.String()) //打印新区块高度
				fmt.Println(header.Time)            //打印新区块时间戳
				fmt.Println(header.Difficulty)      //打印新区块难度
				fmt.Println(header.Nonce)           //打印新区块nonce
				fmt.Println(header.TxHash)          //打印新区块交易hash
				fmt.Println(header.UncleHash)       //打印新区块叔块hash
				block, err3 := client0.BlockByHash(context.Background(),
					header.Hash()) //获取新区块对象
				if err3 != nil {
					log.Println(string(debug.Stack()))
					log.Fatal(err3)
				}
				fmt.Println(len(block.Transactions())) //打印新区块交易数量
				fmt.Println(block.ReceiptHash())       //打印新区块收据hash
			}
		}
	} else if testFlg == 14 {
		//14.发送多笔交易
		privateKey2, err1 := crypto.HexToECDSA("fad9c8855b740a0b7ed4c221dbad0f33a83a49cad6b3fe8d5817ac83d38b6a19") //私钥
		if err1 != nil {
			log.Println(string(debug.Stack()))
			log.Fatal(err1)
		}
		publicKey2 := privateKey2.Public()                   //私钥生成公钥
		publicKeyECDSA2, ok := publicKey2.(*ecdsa.PublicKey) //公钥转成公钥结构体PublicKey类型
		if !ok {
			log.Fatal("cannot assert type: publicKey is not of type *ecdsa.PublicKey")
		}
		address4 := crypto.PubkeyToAddress(*publicKeyECDSA2)                  //转换成公钥地址
		nonce, err2 := client0.PendingNonceAt(context.Background(), address4) //获取nonce
		if err2 != nil {
			log.Println(string(debug.Stack()))
			fmt.Println(err2)
			// log.Fatal(err2)//有异常会跳出，暂时屏蔽
		}
		value := big.NewInt(1e18)                                       // 1个以太币
		gasLimit := uint64(21000)                                       // 交易消耗的gas上限
		gasPrice, err3 := client0.SuggestGasPrice(context.Background()) //基于近几个块和交易池交易动态计算出较合理的gas价格
		if err3 != nil {
			log.Fatal(err3)
		}
		tx2 := types.NewTransaction(nonce, //交易序号
			common.HexToAddress("0x39cb70F972E0EE920088AeF97Dbe5c6251a9c25D"), //接收地址
			value,    //转账金额
			gasLimit, //gas上限
			gasPrice, //gas价格
			nil)      //创建交易
		chainID, err4 := client0.NetworkID(context.Background()) //获取网络ID
		if err4 != nil {
			log.Fatal(err4)
		}
		signedTx, err5 := types.SignTx(tx2, types.NewEIP155Signer(chainID), privateKey2) //对交易进行签名
		fmt.Println(signedTx.Hash().Hex())                                               // 0x8f2a55949038a9610f50fb23b5883af3b4ecb3c3bb792cbcefbd1542c692be63
		if err5 != nil {
			log.Fatal(err5)
		}
		//发送多笔交易
		ts := types.Transactions{signedTx}         //交易列表,可以放多笔交易
		var buf bytes.Buffer                       //缓冲区，用于存放编码后的交易列表
		ts.EncodeIndex(0, &buf)                    //编码交易列表
		rawHex := hex.EncodeToString(buf.Bytes())  //编码后的交易列表
		rawBytes, err6 := hex.DecodeString(rawHex) //解码
		if err6 != nil {
			log.Fatal(err6)
		}
		txs := new(types.Transaction)                              //交易列表
		rlp.DecodeBytes(rawBytes, txs)                             //Rlp解码
		err7 := client0.SendTransaction(context.Background(), txs) //广播交易到区块链网络
		if err7 != nil {
			log.Fatal(err7)
		}
	} else if testFlg == 15 {
		//15.部署合约-仅使用 ethclient 工具
		const (
			// store合约的字节码
			contractBytecode = "608060405234801561000f575f80fd5b5060405161087538038061087583398181016040528101906100319190610193565b805f908161003f91906103e7565b50506104b6565b5f604051905090565b5f80fd5b5f80fd5b5f80fd5b5f80fd5b5f601f19601f8301169050919050565b7f4e487b71000000000000000000000000000000000000000000000000000000005f52604160045260245ffd5b6100a58261005f565b810181811067ffffffffffffffff821117156100c4576100c361006f565b5b80604052505050565b5f6100d6610046565b90506100e2828261009c565b919050565b5f67ffffffffffffffff8211156101015761010061006f565b5b61010a8261005f565b9050602081019050919050565b8281835e5f83830152505050565b5f610137610132846100e7565b6100cd565b9050828152602081018484840111156101535761015261005b565b5b61015e848285610117565b509392505050565b5f82601f83011261017a57610179610057565b5b815161018a848260208601610125565b91505092915050565b5f602082840312156101a8576101a761004f565b5b5f82015167ffffffffffffffff8111156101c5576101c4610053565b5b6101d184828501610166565b91505092915050565b5f81519050919050565b7f4e487b71000000000000000000000000000000000000000000000000000000005f52602260045260245ffd5b5f600282049050600182168061022857607f821691505b60208210810361023b5761023a6101e4565b5b50919050565b5f819050815f5260205f209050919050565b5f6020601f8301049050919050565b5f82821b905092915050565b5f6008830261029d7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff82610262565b6102a78683610262565b95508019841693508086168417925050509392505050565b5f819050919050565b5f819050919050565b5f6102eb6102e66102e1846102bf565b6102c8565b6102bf565b9050919050565b5f819050919050565b610304836102d1565b610318610310826102f2565b84845461026e565b825550505050565b5f90565b61032c610320565b6103378184846102fb565b505050565b5b8181101561035a5761034f5f82610324565b60018101905061033d565b5050565b601f82111561039f5761037081610241565b61037984610253565b81016020851015610388578190505b61039c61039485610253565b83018261033c565b50505b505050565b5f82821c905092915050565b5f6103bf5f19846008026103a4565b1980831691505092915050565b5f6103d783836103b0565b9150826002028217905092915050565b6103f0826101da565b67ffffffffffffffff8111156104095761040861006f565b5b6104138254610211565b61041e82828561035e565b5f60209050601f83116001811461044f575f841561043d578287015190505b61044785826103cc565b8655506104ae565b601f19841661045d86610241565b5f5b828110156104845784890151825560018201915060208501945060208101905061045f565b868310156104a1578489015161049d601f8916826103b0565b8355505b6001600288020188555050505b505050505050565b6103b2806104c35f395ff3fe608060405234801561000f575f80fd5b506004361061003f575f3560e01c806348f343f31461004357806354fd4d5014610073578063f56256c714610091575b5f80fd5b61005d600480360381019061005891906101d7565b6100ad565b60405161006a9190610211565b60405180910390f35b61007b6100c2565b604051610088919061029a565b60405180910390f35b6100ab60048036038101906100a691906102ba565b61014d565b005b6001602052805f5260405f205f915090505481565b5f80546100ce90610325565b80601f01602080910402602001604051908101604052809291908181526020018280546100fa90610325565b80156101455780601f1061011c57610100808354040283529160200191610145565b820191905f5260205f20905b81548152906001019060200180831161012857829003601f168201915b505050505081565b8060015f8481526020019081526020015f20819055507fe79e73da417710ae99aa2088575580a60415d359acfad9cdd3382d59c80281d48282604051610194929190610355565b60405180910390a15050565b5f80fd5b5f819050919050565b6101b6816101a4565b81146101c0575f80fd5b50565b5f813590506101d1816101ad565b92915050565b5f602082840312156101ec576101eb6101a0565b5b5f6101f9848285016101c3565b91505092915050565b61020b816101a4565b82525050565b5f6020820190506102245f830184610202565b92915050565b5f81519050919050565b5f82825260208201905092915050565b8281835e5f83830152505050565b5f601f19601f8301169050919050565b5f61026c8261022a565b6102768185610234565b9350610286818560208601610244565b61028f81610252565b840191505092915050565b5f6020820190508181035f8301526102b28184610262565b905092915050565b5f80604083850312156102d0576102cf6101a0565b5b5f6102dd858286016101c3565b92505060206102ee858286016101c3565b9150509250929050565b7f4e487b71000000000000000000000000000000000000000000000000000000005f52602260045260245ffd5b5f600282049050600182168061033c57607f821691505b60208210810361034f5761034e6102f8565b5b50919050565b5f6040820190506103685f830185610202565b6103756020830184610202565b939250505056fea26469706673582212205aae308f77654b000c9d222eff2d9f2bd2ac18d990b10774842e4309d4e3e15664736f6c634300081a0033"
		)
		privateKey, err1 := crypto.HexToECDSA(
			"fad9c8855b740a0b7ed4c221dbad0f33a83a49cad6b3fe8d5817ac83d38b6a19") //私钥
		if err1 != nil {
			log.Println(string(debug.Stack()))
			log.Fatal(err1)
		}
		publicKey := privateKey.Public()                   //私钥生成公钥
		publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey) //公钥转成公钥结构体PublicKey类型
		if !ok {
			log.Fatal("cannot assert type: publicKey is not of type *ecdsa.PublicKey")
		}
		fromAddress := crypto.PubkeyToAddress(*publicKeyECDSA)                   //公钥转成地址
		nonce, err2 := client0.PendingNonceAt(context.Background(), fromAddress) //获取nonce
		if err2 != nil {
			log.Println(string(debug.Stack()))
			log.Fatal(err2)
		}
		gasPrice, err3 := client0.SuggestGasPrice(context.Background()) //基于近几个块和交易池交易动态计算出较合理的gas价格
		if err3 != nil {
			log.Fatal(err3)
		}
		data, err4 := hex.DecodeString(contractBytecode) // 解码合约字节码
		if err4 != nil {
			log.Fatal(err4)
		}
		tx := types.NewContractCreation(nonce, big.NewInt(0), 3000000, gasPrice, data) //创建交易
		chainid, err5 := client0.NetworkID(context.Background())                       //获取网络ID
		if err5 != nil {
			log.Fatal(err5)
		}
		signedTx, err6 := types.SignTx(tx, types.NewEIP155Signer(chainid), privateKey) //对交易进行签名
		if err6 != nil {
			log.Fatal(err6)
		}
		err7 := client0.SendTransaction(context.Background(), signedTx) //广播交易到区块链网络
		if err7 != nil {
			log.Fatal(err7)
		}
	} else if testFlg == 16 {
		//16.加载合约
		const ( //合约地址
			contractAddr = "0x8D4141ec2b522dE5Cf42705C3010541B4B3EC24e"
		)
		storeContract, err1 := store.NewStore(common.HexToAddress(contractAddr), client0)
		if err1 != nil { //store包引入异常，查看网址提示404 go list -m -versions github.com/learn/init_order
			log.Fatal(err1)
		}
		_ = storeContract // 调用合约方法
	} else if testFlg == 17 {
		//17.执行合约
		const ( //合约地址
			contractAddr = "<deployed contract address>"
		)
		storeContract, err := store.NewStore(common.HexToAddress(contractAddr), client0)
		if err != nil { //store包引入异常，查看网址提示404 go list -m -versions github.com/learn/init_order
			log.Fatal(err)
		}
		privateKey, err := crypto.HexToECDSA("<your private key>") //私钥
		if err != nil {
			log.Fatal(err)
		}
		var key [32]byte //
		var value [32]byte
		copy(key[:], []byte("demo_save_key"))
		copy(value[:], []byte("demo_save_value11111"))
		opt, err := bind.NewKeyedTransactorWithChainID(privateKey, big.NewInt(11155111)) //创建交易
		if err != nil {
			log.Fatal(err)
		}
		tx, err := storeContract.SetItem(opt, key, value) //调用合约方法
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println("tx hash:", tx.Hash().Hex())
		callOpt := &bind.CallOpts{Context: context.Background()}  //调用合约方法
		valueInContract, err := storeContract.Items(callOpt, key) //调用合约方法
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println("is value saving in contract equals to origin value:", valueInContract == value)
	} else if testFlg == 18 {
		//18.合约事件
		//18.1查询事件
		contractAddress := common.HexToAddress("0x2958d15bc5b64b11Ec65e623Ac50C198519f8742") //合约地址
		query := ethereum.FilterQuery{                                                       //查询条件
			// BlockHash
			FromBlock: big.NewInt(6920583), //区块号
			ToBlock:   big.NewInt(2394201), //区块号
			Addresses: []common.Address{ //合约地址
				contractAddress,
			},
		}
		logs, err1 := client0.FilterLogs(context.Background(), query) //查询事件
		if err1 != nil {
			log.Fatal(err1)
		}
		abiConstructor, err2 := abi.JSON(strings.NewReader(store.StoreABI)) //合约ABI
		if err2 != nil {
			log.Fatal(err2)
		}
		for _, vLog := range logs {
			fmt.Println(vLog.BlockHash)   //
			fmt.Println(vLog.BlockNumber) //

			event := struct {
				Key   [32]byte
				Value [32]byte
			}{}
			err3 := abiConstructor.UnpackIntoInterface(&event, "ItemSet", vLog.Data) //解析事件
			if err3 != nil {
				log.Fatal(err3)
			}
			fmt.Println("key:", string(event.Key[:]))     //
			fmt.Println("value:", string(event.Value[:])) //
			var topics []string
			for _, vLogTopic := range vLog.Topics { //解析主题
				topics = append(topics, vLogTopic.Hex())
			}
			fmt.Println("topics[0]=", topics[0]) //第一个主题总是事件的签名
			if len(topics) > 1 {                 //示例合约不包含可索引的事件，但如果它确实包含，这是如何读取事件主题。
				fmt.Println("indexed topics:", topics[1:])
			}
		}
		eventSignature := []byte("ItemSet(bytes32,bytes32)") //事件签名：首个主题只是被哈希过的事件签名
		hash := crypto.Keccak256Hash(eventSignature)         //事件签名hash
		fmt.Println(hash.Hex())
		//18.2订阅事件：订阅事件日志，和订阅区块一样，需要 websocket RPC URL
		contractAddress2 := common.HexToAddress("0x2958d15bc5b64b11Ec65e623Ac50C198519f8742")
		query2 := ethereum.FilterQuery{ //查询所有与 contractAddress 对应的合约地址所有的合约事件
			Addresses: []common.Address{contractAddress2},
		}
		logs2 := make(chan types.Log)                                                 //订阅事件日志-通过 channel 接收事件日志
		sub, err3 := client0.SubscribeFilterLogs(context.Background(), query2, logs2) //从客户端调用 SubscribeFilterLogs 方法来订阅
		if err3 != nil {
			log.Fatal(err3)
		}
		for { //轮询获取新的日志事件
			select {
			case err := <-sub.Err():
				log.Fatal(err)
			case vLog := <-logs2:
				fmt.Println(vLog.BlockHash)   //
				fmt.Println(vLog.BlockNumber) //
				event2 := struct {
					Key   [32]byte
					Value [32]byte
				}{}
				err4 := abiConstructor.UnpackIntoInterface(&event2, "ItemSet", vLog.Data) //解析事件
				if err4 != nil {
					log.Fatal(err4)
				}
				fmt.Println("key:", string(event2.Key[:]))
				fmt.Println("value:", string(event2.Value[:]))
				var topics []string
				for _, vLogTopic := range vLog.Topics { //解析主题
					topics = append(topics, vLogTopic.Hex())
				}
				fmt.Println("topics[0]=", topics[0]) //第一个主题总是事件的签名
				if len(topics) > 1 {                 //示例合约不包含可索引的事件，但如果它确实包含，这是如何读取事件主题。
					fmt.Println("indexed topics:", topics[1:])
				}
			}
		}
	} else {
		fmt.Println("请输入正确的测试分支")
	}
}
