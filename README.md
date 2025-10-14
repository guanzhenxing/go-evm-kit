# go-evm-kit

[![Go Version](https://img.shields.io/badge/Go-1.24+-blue.svg)](https://golang.org)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](LICENSE)

**go-evm-kit** 是一个功能强大的以太坊及 EVM 兼容网络开发工具包，提供简洁易用的 API 来进行链上交互、钱包管理和智能合约操作。

## ✨ 特性

- 🔐 **钱包管理**：支持私钥、助记词、随机生成等多种方式创建账户
- 🌐 **网络连接**：轻松连接以太坊主网、测试网及其他 EVM 兼容网络  
- 💰 **交易操作**：完整的交易构建、签名、发送流程
- 📄 **智能合约**：合约调用、事件监听、ABI 处理
- 🪙 **代币支持**：内置 ERC20 代币操作支持
- 🔧 **实用工具**：单位转换、地址验证、签名验证等
- ⚡ **自动化**：自动计算 nonce、gas price 等参数
- 🔍 **链上查询**：区块、交易、余额等数据查询

## 📦 安装

```bash
go get github.com/guanzhenxing/go-evm-kit
```

## 🚀 快速开始

### 方式1：使用 Kit（推荐）⭐

**Kit** 是最便捷的使用方式，所有方法可以直接调用。

```go
package main

import (
    "context"
    "fmt"
    "log"
    "time"
    "github.com/guanzhenxing/go-evm-kit"
)

func main() {
    privateKey := "your_private_key_here"
    rpcURL := "https://eth-mainnet.g.alchemy.com/v2/your-api-key"
    
    // 创建 Kit
    kit, err := etherkit.NewKit(privateKey, rpcURL)
    if err != nil {
        log.Fatal(err)
    }
    defer kit.CloseWallet()
    
    ctx := context.Background()
    
    // 所有方法直接调用 - 简洁明了！
    address := kit.GetAddress()                    // 来自 Signer
    chainID, _ := kit.GetChainID(ctx)              // 来自 Provider
    balance, _ := kit.GetBalance(ctx)              // 来自 Wallet
    ethBalance, _ := kit.GetBalanceInEther(ctx)    // 增强功能
    
    fmt.Printf("地址: %s\n", address.Hex())
    fmt.Printf("Chain ID: %s\n", chainID)
    fmt.Printf("余额: %s Wei (%.6f ETH)\n", balance, ethBalance)
    
    // 发送交易并等待确认（增强功能）
    toAddress := common.HexToAddress("0x...")
    receipt, err := kit.SendTxAndWait(
        ctx, toAddress, 
        0, 0, nil,              // nonce, gasLimit, gasPrice 自动计算
        big.NewInt(1000000),    // value
        nil,                    // data
        30*time.Second,         // 超时时间
    )
    if err != nil {
        log.Fatal(err)
    }
    fmt.Printf("交易已确认，状态: %d\n", receipt.Status)
}
```

### 方式2：使用原始 Wallet

```go
package main

import (
    "context"
    "fmt"
    "log"
    "github.com/guanzhenxing/go-evm-kit"
)

func main() {
    // 使用私钥创建钱包
    privateKey := "your_private_key_here"
    rpcURL := "https://eth-mainnet.g.alchemy.com/v2/your-api-key"
    
    wallet, err := etherkit.NewWallet(privateKey, rpcURL)
    if err != nil {
        log.Fatal(err)
    }
    defer wallet.CloseWallet()
    
    ctx := context.Background()
    
    // 获取账户地址
    address := wallet.GetAddress()
    fmt.Printf("钱包地址: %s\n", address.Hex())
    
    // 获取余额
    balance, err := wallet.GetBalance(ctx)
    if err != nil {
        log.Fatal(err)
    }
    fmt.Printf("ETH 余额: %s\n", etherkit.ToDecimal(balance, etherkit.EthDecimals))
}
```

### 发送 ETH 转账

```go
func sendETH(wallet *etherkit.Wallet) {
    ctx := context.Background()
    toAddress := common.HexToAddress("0x742F35Cc6634C0532925a3b8D6dA2e")
    amount := etherkit.ToWei("0.1", etherkit.EthDecimals) // 0.1 ETH
    
    txHash, err := wallet.SendTx(
        ctx,           // context
        toAddress,     // 收款地址
        0,             // nonce (0 表示自动计算)
        0,             // gasLimit (0 表示自动估算)
        nil,           // gasPrice (nil 表示自动获取)
        amount,        // 转账金额
        nil,           // 交易数据
    )
    if err != nil {
        log.Fatal(err)
    }
    
    fmt.Printf("交易哈希: %s\n", txHash.Hex())
}
```

### ERC20 代币操作

```go
import (
    "context"
    "github.com/guanzhenxing/go-evm-kit/contracts/erc20"
)

func transferToken(wallet *etherkit.Wallet) {
    ctx := context.Background()
    tokenAddress := common.HexToAddress("0xA0b86a33E6411b6dE9C80e7F8DeD6c") // USDC 地址
    
    // 创建 ERC20 合约实例
    token, err := erc20.NewIERC20(tokenAddress, wallet.GetClient())
    if err != nil {
        log.Fatal(err)
    }
    
    // 构建交易选项
    opts, err := wallet.BuildTxOpts(
        ctx,              // context
        big.NewInt(0),    // value
        nil,              // nonce (自动计算)
        nil,              // gasPrice (自动获取)
    )
    if err != nil {
        log.Fatal(err)
    }
    
    // 转账代币
    toAddress := common.HexToAddress("0x742F35Cc6634C0532925a3b8D6dA2e")
    amount := etherkit.ToWei("100", etherkit.USDCDecimals) // 100 USDC
    
    tx, err := token.Transfer(opts, toAddress, amount)
    if err != nil {
        log.Fatal(err)
    }
    
    fmt.Printf("代币转账交易: %s\n", tx.Hash().Hex())
}
```

### 智能合约调用

```go
func callContract(wallet *etherkit.Wallet) {
    ctx := context.Background()
    contractAddress := common.HexToAddress("0x...")
    abiString := `[{"inputs":[],"name":"totalSupply","outputs":[{"type":"uint256"}],"type":"function"}]`
    
    // 获取合约 ABI
    contractAbi, err := etherkit.GetABI(abiString)
    if err != nil {
        log.Fatal(err)
    }
    
    // 调用合约方法 (只读)
    result, err := wallet.CallContract(ctx, contractAddress, contractAbi, "totalSupply")
    if err != nil {
        log.Fatal(err)
    }
    
    fmt.Printf("总供应量: %v\n", result[0])
}
```

## 📚 API 文档

### Provider (网络提供者)

```go
// 创建网络连接
provider, err := etherkit.NewProvider("https://eth-mainnet.g.alchemy.com/v2/your-api-key")
provider, err := etherkit.NewProviderWithChainId("https://polygon-rpc.com", 137)

// 基本查询（需要 context）
ctx := context.Background()
chainID, err := provider.GetChainID(ctx)
blockNumber, err := provider.GetBlockNumber(ctx) 
gasPrice, err := provider.GetSuggestGasPrice(ctx)
block, err := provider.GetBlockByNumber(ctx, big.NewInt(123456))
receipt, err := provider.GetTransactionReceipt(ctx, txHash)
```

### 私钥和地址工具函数

```go
// 多种方式生成/导入私钥
privateKey, err := etherkit.GeneratePrivateKey()                              // 随机生成
privateKey, err := etherkit.BuildPrivateKeyFromHex("0x...")                   // 从十六进制
privateKey, err := etherkit.BuildPrivateKeyFromMnemonic("word1 word2...")     // 从助记词

// 获取地址
address := etherkit.PrivateKeyToAddress(privateKey)

// 获取私钥的十六进制字符串
hexPk := etherkit.GetHexPrivateKey(privateKey)
```

### Wallet (钱包)

```go
// 创建钱包
wallet, err := etherkit.NewWallet(privateKey, rpcURL)

// 账户操作
address := wallet.GetAddress()
balance, err := wallet.GetBalance(ctx)
nonce, err := wallet.GetNonce(ctx)

// 直接调用 Provider 的方法（需要两步）
chainID, err := wallet.GetEthProvider().GetChainID(ctx)
blockNum, err := wallet.GetEthProvider().GetBlockNumber(ctx)

// 直接调用钱包方法
address := wallet.GetAddress()
privateKey := wallet.GetPrivateKey()

// 交易操作
tx, err := wallet.NewTx(ctx, toAddr, nonce, gasLimit, gasPrice, value, data)
txHash, err := wallet.SendTx(ctx, toAddr, nonce, gasLimit, gasPrice, value, data)
signedTx, err := wallet.SignTx(ctx, tx)
```

### Kit (开发工具包) ⭐

**推荐使用！** 提供最便捷的 API，所有方法可直接调用，无需通过 `GetEthProvider()`。

```go
// 创建 Kit
kit, err := etherkit.NewKit(privateKey, rpcURL)

// ========== 所有方法都可以直接调用 ==========

// 钱包方法（来自 Wallet）
address := kit.GetAddress()                // 获取地址
pk := kit.GetPrivateKey()                  // 获取私钥
balance, err := kit.GetBalance(ctx)        // 获取余额
nonce, err := kit.GetNonce(ctx)           // 获取 nonce

// Provider 方法（无需 GetEthProvider()）
chainID, err := kit.GetChainID(ctx)        // 获取链 ID
blockNum, err := kit.GetBlockNumber(ctx)   // 获取区块高度
receipt, err := kit.GetTransactionReceipt(ctx, txHash)  // 获取交易回执

// 发送交易
txHash, err := kit.SendTx(ctx, toAddr, 0, 0, nil, value, nil)

// ========== 增强功能 ==========

// 1. 以太币单位的余额
ethBalance, err := kit.GetBalanceInEther(ctx)
fmt.Printf("余额: %.6f ETH\n", ethBalance)

// 2. 便捷的以太币转账
txHash, err := kit.TransferEther(ctx, toAddr, 0.1)  // 转 0.1 ETH

// 3. 发送交易并等待确认
receipt, err := kit.SendTxAndWait(
    ctx, toAddr, 
    0, 0, nil,           // nonce, gasLimit, gasPrice 自动计算
    value, data,
    30*time.Second,      // 超时时间
)

// 4. 等待交易确认（带超时）
receipt, err := kit.WaitForReceipt(ctx, txHash, 30*time.Second)

// 5. 一次性获取链信息
chainID, networkID, blockNum, err := kit.GetChainInfo(ctx)

// 6. 简化的合约调用
result, err := kit.CallContractSimple(
    ctx, 
    contractAddr,
    abiJSON,        // JSON 字符串
    "balanceOf",
    myAddress,
)
```

**对比表格：**

| 功能 | Wallet | Kit |
|------|--------|-----|
| 调用 Provider 方法 | `wallet.GetEthProvider().GetChainID(ctx)` | `kit.GetChainID(ctx)` ✅ |
| 获取地址/私钥 | `wallet.GetAddress()` | `kit.GetAddress()` ✅ |
| 获取 ETH 余额 | 手动转换 Wei | `kit.GetBalanceInEther(ctx)` ✅ |
| 等待交易确认 | 手动轮询 | `kit.WaitForReceipt(ctx, hash, timeout)` ✅ |
| 转账以太币 | 手动转换单位 | `kit.TransferEther(ctx, to, 0.1)` ✅ |

### 工具函数

```go
// 单位转换
wei := etherkit.ToWei("1.5", etherkit.EthDecimals)     // 1.5 ETH 转 wei
eth := etherkit.ToDecimal(wei, etherkit.EthDecimals)   // wei 转 ETH

// 地址验证
isValid := etherkit.IsValidAddress("0x...")

// 签名验证  
isValid := etherkit.VerifySignature(address, data, signature)

// 合约工具
methodID := etherkit.GetContractMethodId("transfer(address,uint256)")
eventTopic := etherkit.GetEventTopic("Transfer(address,address,uint256)")

// 常量使用
chainID := etherkit.MainnetChainID  // 主网链ID
gasPrice := etherkit.DefaultGasPriceBig  // 默认Gas价格
```

## 📁 项目结构

```
go-evm-kit/
├── provider.go        # Provider 层 - 链交互
├── wallet.go          # Wallet 层 - 完整钱包功能
├── kit.go             # Kit 层 - 最便捷的工具包
├── crypto.go          # 加密工具（私钥、签名、地址）
├── convert.go         # 单位转换工具
├── contract.go        # 合约工具
├── transaction.go     # 交易工具
├── address.go         # 地址相关工具
├── constants.go       # 常量定义
├── errors.go          # 错误定义
├── contracts/         # 智能合约绑定
│   └── erc20/        # ERC20 合约
├── examples/          # 使用示例
│   └── kit/          # Kit 使用示例
├── *_test.go         # 单元测试文件
├── go.mod
├── go.sum
├── Makefile          # 构建和开发工具
├── LICENSE
└── README.md
```

## 🏗️ 架构设计

### 设计理念

go-evm-kit 采用**简洁实用**的设计理念，专注于提供便捷的以太坊开发工具，而不是构建复杂的抽象层。

### 核心架构

```
crypto.go (工具函数)
   ├─ GeneratePrivateKey()
   ├─ BuildPrivateKeyFromHex()
   └─ PrivateKeyToAddress()
      ↓
provider.go (链交互)
   ├─ GetChainID()
   ├─ GetBlockNumber()
   └─ GetTransactionReceipt()
      ↓
wallet.go (完整钱包)
   ├─ privateKey, address (私钥/地址)
   ├─ Provider (链交互)
   └─ SendTx(), SignTx()
      ↓
kit.go (推荐使用) ⭐
   ├─ *Wallet (嵌入)
   ├─ EtherProvider (嵌入接口)
   └─ 增强功能
```

### 三层架构

| 层级 | 文件 | 用途 | 何时使用 |
|------|------|------|---------|
| **Provider** | provider.go | 链交互查询 | 只需要查询链数据，不需要私钥 |
| **Wallet** | wallet.go | 完整钱包功能 | 需要发送交易、签名等完整功能 |
| **Kit** ⭐ | kit.go | 最便捷的工具包 | **推荐使用** - 所有方法直接调用 |

### 设计决策

**为什么将私钥直接放在 Wallet 中？**
1. **简化架构** - 移除不必要的 Signer 抽象层
2. **开发工具定位** - 专注于 99% 的开发场景，不需要支持硬件钱包等复杂场景
3. **更直观** - Wallet（钱包）本身就应该包含私钥和地址
4. **易于维护** - 更少的接口和类型

**为什么推荐使用 Kit？**
1. **最简单** - 所有方法直接调用，无需 `GetEthProvider()`
2. **最强大** - 包含所有功能 + 增强功能
3. **最实用** - 覆盖 99% 的使用场景

### 扩展性

虽然架构简化了，但仍保持良好的扩展性：

```go
// 自定义 Provider（添加缓存）
type CachedProvider struct {
    *Provider
    cache Cache
}

// 基于 Kit 扩展（DeFi 功能）
type DeFiKit struct {
    *Kit
}

func (dk *DeFiKit) SwapTokens(...) error {
    // 实现 DeFi 操作
}
```

## 🚀 最佳实践

### 1. 优先使用 Kit（推荐）

```go
// ✅ 推荐 - 最简单的使用方式
kit, _ := NewKit(privateKey, rpcURL)
balance, _ := kit.GetBalance(ctx)
chainID, _ := kit.GetChainID(ctx)  // 直接调用，无需 GetEthProvider()
```

### 2. 只读场景使用 Provider

```go
// ✅ 只需要查询链数据时
provider, _ := NewProvider(rpcURL)
blockNum, _ := provider.GetBlockNumber(ctx)
```

### 3. 复杂场景使用 Wallet

```go
// ✅ 需要自定义逻辑时
privateKey, _ := BuildPrivateKeyFromHex(hexPk)
provider, _ := NewProvider(rpcURL)
wallet, _ := NewWalletWithComponents(privateKey, provider)
```

## 🆕 最新改进

### v2.0 - 架构简化
- ✅ **简化架构** - 移除 Signer 抽象层，从 4 层减少到 3 层
- ✅ **Kit 工具包** - 提供最便捷的 API，所有方法直接调用
- ✅ **更少的文件** - 核心只有 3 个文件：provider.go、wallet.go、kit.go
- ✅ **更易理解** - Wallet 直接包含私钥和地址，概念更清晰

### 代码质量
- ✅ **完整单元测试** - 覆盖率 34%+
- ✅ **详细使用示例** - examples/kit 完整示例
- ✅ **标准化命名** - 符合 Go 语言习惯

## 🌐 支持的网络

| 网络名称 | Chain ID | 符号 | 区块时间 | 确认数 |
|---------|----------|------|----------|--------|
| Ethereum Mainnet | 1 | ETH | 12s | 12 |
| Goerli Testnet | 5 | ETH | 12s | 3 |
| Sepolia Testnet | 11155111 | ETH | 12s | 3 |
| Polygon | 137 | MATIC | 2s | 20 |
| BSC | 56 | BNB | 3s | 15 |
| Arbitrum One | 42161 | ETH | - | - |
| Optimism | 10 | ETH | - | - |

使用预定义常量：
```go
// 直接使用链ID常量
provider := etherkit.NewProviderWithChainId(rpcURL, etherkit.MainnetChainID)

// 获取网络配置
config := etherkit.NetworkConfigs[etherkit.PolygonChainID]
fmt.Printf("网络: %s, 符号: %s\n", config.Name, config.Symbol)
```

## 🔧 高级用法

### 批量操作

```go
// 批量查询余额
addresses := []common.Address{addr1, addr2, addr3}
for _, addr := range addresses {
    balance, _ := provider.GetEthClient().BalanceAt(context.Background(), addr, nil)
    fmt.Printf("地址 %s 余额: %s ETH\n", addr.Hex(), etherkit.ToDecimal(balance, 18))
}
```

### 事件监听

```go
// 监听 ERC20 Transfer 事件
query := ethereum.FilterQuery{
    Addresses: []common.Address{tokenAddress},
    Topics: [][]common.Hash{
        {common.HexToHash("0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef")},
    },
}

logs := make(chan types.Log)
sub, err := provider.GetEthClient().SubscribeFilterLogs(context.Background(), query, logs)
if err != nil {
    log.Fatal(err)
}

for {
    select {
    case err := <-sub.Err():
        log.Fatal(err)
    case vLog := <-logs:
        fmt.Printf("发现 Transfer 事件: %s\n", vLog.TxHash.Hex())
    }
}
```

## 🤝 贡献

欢迎提交 Issue 和 Pull Request！

1. Fork 项目
2. 创建特性分支 (`git checkout -b feature/AmazingFeature`)
3. 提交更改 (`git commit -m 'Add some AmazingFeature'`)
4. 推送到分支 (`git push origin feature/AmazingFeature`)
5. 打开 Pull Request

## 📄 许可证

本项目基于 MIT 许可证开源 - 查看 [LICENSE](LICENSE) 文件了解详情。

## 🔗 相关资源

- [以太坊官方文档](https://ethereum.org/developers/)
- [go-ethereum 文档](https://geth.ethereum.org/docs/)
- [Web3 开发指南](https://web3.guide/)

## 📞 支持

如有问题或建议，请：

- 提交 [Issue](https://github.com/guanzhenxing/go-evm-kit/issues)
- 发送邮件至 [your-email@example.com]
- 加入我们的讨论群组

---

⭐ 如果这个项目对你有帮助，请给个 Star！