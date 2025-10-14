package main

import (
	"context"
	"fmt"
	"log"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/common"
	etherkit "github.com/guanzhenxing/go-evm-kit"
)

func main() {
	// 生成随机私钥用于演示
	pk, err := etherkit.GeneratePrivateKey()
	if err != nil {
		log.Fatalf("生成私钥失败: %v", err)
	}
	privateKey := etherkit.GetHexPrivateKey(pk)
	rpcURL := "https://eth-mainnet.g.alchemy.com/v2/demo"

	ctx := context.Background()

	fmt.Println("=== Kit 使用示例 ===")

	// ============ 示例1：基本使用 ============
	fmt.Println("示例1：创建 Kit 并获取基本信息")
	basicUsage(ctx, privateKey, rpcURL)

	// ============ 示例2：对比 Wallet 和 Kit ============
	fmt.Println("\n示例2：调用方式对比")
	comparisonExample(privateKey, rpcURL)

	// ============ 示例3：单位转换 ============
	fmt.Println("\n示例3：单位转换")
	conversionExample()

	fmt.Println("\n=== 所有示例完成 ===")
}

// basicUsage Kit 基本使用
func basicUsage(ctx context.Context, privateKey, rpcURL string) {
	// 创建 Kit
	kit, err := etherkit.NewKit(privateKey, rpcURL)
	if err != nil {
		log.Fatalf("创建 Kit 失败: %v", err)
	}
	defer kit.CloseWallet()

	// 获取地址 - 直接调用（来自 Wallet）
	address := kit.GetAddress()
	fmt.Printf("  地址: %s\n", address.Hex())

	// 获取私钥 - 直接调用
	pk := kit.GetPrivateKey()
	fmt.Printf("  私钥存在: %v\n", pk != nil)

	// 注意：以下需要实际的 RPC 连接，这里仅作演示
	fmt.Println("  （以下功能需要有效的 RPC 连接）")

	// 获取链信息 - 直接调用（来自 Provider）
	// chainID, err := kit.GetChainID(ctx)
	// blockNum, err := kit.GetBlockNumber(ctx)
	// balance, err := kit.GetBalance(ctx)
	// ethBalance, err := kit.GetBalanceInEther(ctx)
}

// comparisonExample 对比 Wallet 和 Kit
func comparisonExample(privateKey, rpcURL string) {
	fmt.Println("  原始 Wallet 的调用方式：")
	fmt.Println("    wallet.GetAddress()                     // 需要 Wallet 方法")
	fmt.Println("    wallet.GetEthProvider().GetChainID(ctx) // 需要两步")

	fmt.Println("\n  Kit 的调用方式：")
	fmt.Println("    kit.GetAddress()        // 直接调用 ✅")
	fmt.Println("    kit.GetChainID(ctx)     // 直接调用 ✅")
	fmt.Println("    kit.GetBlockNumber(ctx) // 直接调用 ✅")
}

// conversionExample 单位转换示例
func conversionExample() {
	// 1 ETH = 1000000000000000000 Wei
	oneEth := etherkit.ToWei(1.0, 18)
	fmt.Printf("  1 ETH = %s Wei\n", oneEth.String())

	// Wei 转 ETH
	ethValue := etherkit.ToDecimal(oneEth, 18)
	ethFloat, _ := ethValue.Float64()
	fmt.Printf("  %s Wei = %.1f ETH\n", oneEth.String(), ethFloat)

	// 0.5 ETH
	halfEth := etherkit.ToWei(0.5, 18)
	fmt.Printf("  0.5 ETH = %s Wei\n", halfEth.String())
}

// 以下是更多示例，需要有效的 RPC 连接才能运行

func chainInfoExample(ctx context.Context, kit *etherkit.Kit) {
	fmt.Println("\n=== 链信息查询 ===")

	// 方式1：单独查询
	chainID, err := kit.GetChainID(ctx)
	if err != nil {
		log.Printf("获取 ChainID 失败: %v", err)
		return
	}
	fmt.Printf("Chain ID: %s\n", chainID.String())

	blockNum, err := kit.GetBlockNumber(ctx)
	if err != nil {
		log.Printf("获取区块高度失败: %v", err)
		return
	}
	fmt.Printf("当前区块: %d\n", blockNum)

	// 方式2：一次性获取（便捷方法）
	chainID2, networkID, blockNum2, err := kit.GetChainInfo(ctx)
	if err != nil {
		log.Printf("获取链信息失败: %v", err)
		return
	}
	fmt.Printf("ChainID: %s, NetworkID: %s, 区块: %s\n", chainID2, networkID, blockNum2)
}

func balanceExample(ctx context.Context, kit *etherkit.Kit) {
	fmt.Println("\n=== 余额查询 ===")

	// 方式1：获取 Wei 余额
	balance, err := kit.GetBalance(ctx)
	if err != nil {
		log.Printf("获取余额失败: %v", err)
		return
	}
	fmt.Printf("余额 (Wei): %s\n", balance.String())

	// 方式2：获取 ETH 余额（便捷方法）
	ethBalance, err := kit.GetBalanceInEther(ctx)
	if err != nil {
		log.Printf("获取以太币余额失败: %v", err)
		return
	}
	fmt.Printf("余额 (ETH): %.6f\n", ethBalance)
}

func transferExample(ctx context.Context, kit *etherkit.Kit) {
	fmt.Println("\n=== 转账示例 ===")

	toAddress := common.HexToAddress("0x70997970C51812dc3A010C7d01b50e0d17dc79C8")

	// 方式1：使用 SendTx（完整参数）
	value := big.NewInt(1000000000000000) // 0.001 ETH
	txHash, err := kit.SendTx(
		ctx,
		toAddress,
		0, 0, nil, // nonce, gasLimit, gasPrice 自动计算
		value,
		nil, // data
	)
	if err != nil {
		log.Printf("发送交易失败: %v", err)
		return
	}
	fmt.Printf("交易已发送: %s\n", txHash.Hex())

	// 方式2：使用 TransferEther（便捷方法）
	txHash2, err := kit.TransferEther(ctx, toAddress, 0.001) // 转 0.001 ETH
	if err != nil {
		log.Printf("转账失败: %v", err)
		return
	}
	fmt.Printf("转账交易: %s\n", txHash2.Hex())

	// 方式3：发送并等待确认
	receipt, err := kit.SendTxAndWait(
		ctx,
		toAddress,
		0, 0, nil,
		value,
		nil,
		30*time.Second, // 超时时间
	)
	if err != nil {
		log.Printf("发送交易并等待确认失败: %v", err)
		return
	}
	fmt.Printf("交易已确认，状态: %d, Gas Used: %d\n", receipt.Status, receipt.GasUsed)
}

func waitReceiptExample(ctx context.Context, kit *etherkit.Kit, txHash common.Hash) {
	fmt.Println("\n=== 等待交易确认 ===")

	// 等待交易确认，最多等待 30 秒
	receipt, err := kit.WaitForReceipt(ctx, txHash, 30*time.Second)
	if err != nil {
		log.Printf("等待交易确认失败: %v", err)
		return
	}

	fmt.Printf("交易状态: %d\n", receipt.Status)
	fmt.Printf("Gas Used: %d\n", receipt.GasUsed)
	fmt.Printf("区块号: %d\n", receipt.BlockNumber.Uint64())
}

// printArchitecture 架构说明
func printArchitecture() {
	fmt.Print(`
Kit 架构：

Provider (链交互)
   ├─ GetChainID()
   ├─ GetBlockNumber()
   └─ GetTransactionReceipt()
   
Wallet (钱包) = 私钥/地址 + Provider
   ├─ GetAddress()
   ├─ GetPrivateKey()
   ├─ GetBalance()
   ├─ SendTx()
   └─ SignTx()
   
Kit (推荐) = Wallet + Provider 嵌入 + 增强功能
   ├─ 直接调用所有 Wallet 方法
   ├─ 直接调用所有 Provider 方法
   ├─ 增加便捷功能
   │   ├─ GetBalanceInEther()
   │   ├─ TransferEther()
   │   ├─ WaitForReceipt()
   │   ├─ SendTxAndWait()
   │   └─ GetChainInfo()
   └─ 最简单的使用方式 ⭐

选择建议：
- 只查询链数据 → Provider
- 需要钱包功能 → Wallet 或 Kit
- 追求便捷使用 → Kit ⭐⭐⭐
`)
}
