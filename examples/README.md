# Examples

## 推荐示例

### Kit 示例 ⭐ 推荐

最简单、最强大的使用方式。

```bash
cd kit
go run main.go
```

这个示例展示了：
- 创建 Kit
- 获取地址和私钥
- 查询链信息
- 获取余额
- 发送交易
- 等待交易确认
- 单位转换

## 旧示例（已过时）

以下示例使用旧的 API，部分功能可能无法直接运行：

- `basic/` - 基础示例（需要更新）
- `advanced/` - 高级示例（需要更新）
- `erc20/` - ERC20 示例（需要更新）

**建议：** 使用 `Kit` 示例作为参考，它包含了所有常用功能。

## 快速开始

```go
package main

import (
    "context"
    "fmt"
    "log"
    etherkit "github.com/guanzhenxing/go-evm-kit"
)

func main() {
    // 创建 Kit
    kit, err := etherkit.NewKit("your_private_key", "https://eth-mainnet.g.alchemy.com/v2/your-api-key")
    if err != nil {
        log.Fatal(err)
    }
    defer kit.CloseWallet()
    
    ctx := context.Background()
    
    // 获取地址
    address := kit.GetAddress()
    fmt.Printf("地址: %s\n", address.Hex())
    
    // 获取链信息
    chainID, err := kit.GetChainID(ctx)
    if err != nil {
        log.Fatal(err)
    }
    fmt.Printf("Chain ID: %s\n", chainID)
    
    // 获取余额
    balance, err := kit.GetBalance(ctx)
    if err != nil {
        log.Fatal(err)
    }
    fmt.Printf("余额: %s Wei\n", balance)
    
    // 获取 ETH 单位的余额
    ethBalance, err := kit.GetBalanceInEther(ctx)
    if err != nil {
        log.Fatal(err)
    }
    fmt.Printf("余额: %.6f ETH\n", ethBalance)
}
```

## 架构说明

```
Provider (链交互)     →  只需要查询链数据
Signer (签名)         →  离线签名
Wallet (钱包)         →  完整钱包功能
Kit (推荐)            →  最便捷的使用方式 ⭐
```

## 更多示例

查看 `kit/main.go` 获取完整的使用示例。
