package etherkit

import (
	"context"
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"
)

// TestKitCreation 测试 Kit 的创建
func TestKitCreation(t *testing.T) {
	// 生成随机私钥用于测试
	pk, err := GeneratePrivateKey()
	if err != nil {
		t.Fatalf("生成私钥失败: %v", err)
	}
	testPrivateKey := GetHexPrivateKey(pk)
	testRPCURL := "https://eth-mainnet.g.alchemy.com/v2/demo"

	kit, err := NewKit(testPrivateKey, testRPCURL)
	if err != nil {
		t.Fatalf("创建 Kit 失败: %v", err)
	}
	defer kit.CloseWallet()

	if kit.Wallet == nil {
		t.Error("Wallet 不应该为 nil")
	}

	address := kit.GetAddress()
	if address == (common.Address{}) {
		t.Error("地址不应该为空")
	}

	t.Logf("创建成功，地址: %s", address.Hex())
}

// TestKitDirectAccess 测试 Kit 的直接访问
func TestKitDirectAccess(t *testing.T) {
	pk, err := GeneratePrivateKey()
	if err != nil {
		t.Fatalf("生成私钥失败: %v", err)
	}
	testPrivateKey := GetHexPrivateKey(pk)
	testRPCURL := "https://eth-mainnet.g.alchemy.com/v2/demo"

	kit, err := NewKit(testPrivateKey, testRPCURL)
	if err != nil {
		t.Fatalf("创建 Kit 失败: %v", err)
	}
	defer kit.CloseWallet()

	// 测试直接调用 Wallet 的方法
	address := kit.GetAddress()
	if address == (common.Address{}) {
		t.Error("地址不应该为空")
	}

	retrievedPk := kit.GetPrivateKey()
	if retrievedPk == nil {
		t.Error("私钥不应该为 nil")
	}

	t.Logf("地址: %s", address.Hex())
}

// TestKitWithComponents 测试使用组件创建 Kit
func TestKitWithComponents(t *testing.T) {
	privateKey, err := GeneratePrivateKey()
	if err != nil {
		t.Fatalf("生成私钥失败: %v", err)
	}
	testRPCURL := "https://eth-mainnet.g.alchemy.com/v2/demo"

	provider, err := NewProvider(testRPCURL)
	if err != nil {
		t.Fatalf("创建 Provider 失败: %v", err)
	}

	kit, err := NewKitWithComponents(privateKey, provider)
	if err != nil {
		t.Fatalf("使用组件创建 Kit 失败: %v", err)
	}
	defer kit.CloseWallet()

	address := kit.GetAddress()
	if address == (common.Address{}) {
		t.Error("地址不应该为空")
	}
}

// TestKitConversion 测试转换方法
func TestKitConversion(t *testing.T) {
	// 测试单位转换
	oneEth := ToWei(1.0, 18)
	expectedWei := new(big.Int)
	expectedWei.SetString("1000000000000000000", 10)

	if oneEth.Cmp(expectedWei) != 0 {
		t.Errorf("1 ETH 应该等于 %s Wei, 但得到 %s", expectedWei.String(), oneEth.String())
	}

	// 反向转换
	ethValue := ToDecimal(oneEth, 18)
	ethFloat, _ := ethValue.Float64()
	if ethFloat != 1.0 {
		t.Errorf("应该转换回 1.0 ETH, 但得到 %f", ethFloat)
	}
}

// 以下是需要实际 RPC 连接的测试，标记为跳过

func TestKitChainMethods(t *testing.T) {
	t.Skip("需要实际的 RPC 连接")

	pk, err := GeneratePrivateKey()
	if err != nil {
		t.Fatalf("生成私钥失败: %v", err)
	}
	testPrivateKey := GetHexPrivateKey(pk)
	testRPCURL := "https://eth-mainnet.g.alchemy.com/v2/your-api-key"

	kit, err := NewKit(testPrivateKey, testRPCURL)
	if err != nil {
		t.Fatalf("创建 Kit 失败: %v", err)
	}
	defer kit.CloseWallet()

	ctx := context.Background()

	// 测试直接调用 Provider 的方法
	chainID, err := kit.GetChainID(ctx)
	if err != nil {
		t.Errorf("获取 ChainID 失败: %v", err)
	}
	t.Logf("Chain ID: %s", chainID.String())

	blockNum, err := kit.GetBlockNumber(ctx)
	if err != nil {
		t.Errorf("获取区块高度失败: %v", err)
	}
	t.Logf("当前区块高度: %d", blockNum)

	// 测试便捷方法
	chainID2, networkID, blockNum2, err := kit.GetChainInfo(ctx)
	if err != nil {
		t.Errorf("获取链信息失败: %v", err)
	}
	t.Logf("ChainID: %s, NetworkID: %s, 区块: %s", chainID2, networkID, blockNum2)
}

// BenchmarkKitCreation 基准测试：创建 Kit
func BenchmarkKitCreation(b *testing.B) {
	pk, err := GeneratePrivateKey()
	if err != nil {
		b.Fatalf("生成私钥失败: %v", err)
	}
	testPrivateKey := GetHexPrivateKey(pk)
	testRPCURL := "https://eth-mainnet.g.alchemy.com/v2/demo"

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		kit, err := NewKit(testPrivateKey, testRPCURL)
		if err != nil {
			b.Fatalf("创建 Kit 失败: %v", err)
		}
		kit.CloseWallet()
	}
}
