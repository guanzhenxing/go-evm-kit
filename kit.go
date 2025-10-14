package etherkit

import (
	"context"
	"crypto/ecdsa"
	"math/big"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

// Kit 以太坊开发工具包，提供最便捷的使用方式
// 通过嵌入接口，所有 Provider 和 Wallet 的方法都可以直接调用
type Kit struct {
	*Wallet       // 嵌入 Wallet，获得所有钱包方法（包括 GetAddress、GetPrivateKey）
	EtherProvider // 嵌入 Provider 接口，直接调用所有 Provider 方法！
}

// NewKit 创建以太坊开发工具包
func NewKit(hexPk string, rawUrl string) (*Kit, error) {
	wallet, err := NewWallet(hexPk, rawUrl)
	if err != nil {
		return nil, err
	}
	return &Kit{
		Wallet:        wallet,
		EtherProvider: wallet.GetEthProvider(),
	}, nil
}

// NewKitWithComponents 使用已有组件创建 Kit
func NewKitWithComponents(privateKey *ecdsa.PrivateKey, ep EtherProvider) (*Kit, error) {
	wallet, err := NewWalletWithComponents(privateKey, ep)
	if err != nil {
		return nil, err
	}
	return &Kit{
		Wallet:        wallet,
		EtherProvider: ep,
	}, nil
}

// ============ 以下是增强功能 ============

// WaitForReceipt 等待交易被打包，带超时控制
func (k *Kit) WaitForReceipt(ctx context.Context, txHash common.Hash, timeout time.Duration) (*types.Receipt, error) {
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		case <-ticker.C:
			receipt, err := k.GetTransactionReceipt(ctx, txHash)
			if err == nil && receipt != nil {
				return receipt, nil
			}
		}
	}
}

// SendTxAndWait 发送交易并等待确认
func (k *Kit) SendTxAndWait(ctx context.Context, to common.Address, nonce, gasLimit uint64, gasPrice, value *big.Int, data []byte, timeout time.Duration) (*types.Receipt, error) {
	txHash, err := k.SendTx(ctx, to, nonce, gasLimit, gasPrice, value, data)
	if err != nil {
		return nil, err
	}
	return k.WaitForReceipt(ctx, txHash, timeout)
}

// GetBalanceInEther 获取以太币余额（以 ETH 为单位）
func (k *Kit) GetBalanceInEther(ctx context.Context) (float64, error) {
	balance, err := k.GetBalance(ctx)
	if err != nil {
		return 0, err
	}
	// 使用 ToDecimal 转换，以太币的 decimals 是 18
	ethBalance := ToDecimal(balance, 18)
	result, _ := ethBalance.Float64()
	return result, nil
}

// TransferEther 转账以太币（便捷方法）
func (k *Kit) TransferEther(ctx context.Context, to common.Address, valueInEther float64) (common.Hash, error) {
	// 使用 ToWei 转换，以太币的 decimals 是 18
	value := ToWei(valueInEther, 18)
	return k.SendTx(ctx, to, 0, 0, nil, value, nil)
}

// GetChainInfo 获取链信息（便捷方法）
func (k *Kit) GetChainInfo(ctx context.Context) (chainID, networkID, blockNumber *big.Int, err error) {
	chainID, err = k.GetChainID(ctx)
	if err != nil {
		return nil, nil, nil, err
	}

	networkID, err = k.GetNetworkID(ctx)
	if err != nil {
		return nil, nil, nil, err
	}

	blockNum, err := k.GetBlockNumber(ctx)
	if err != nil {
		return nil, nil, nil, err
	}

	return chainID, networkID, big.NewInt(int64(blockNum)), nil
}

// ============ 便捷的合约交互方法 ============

// CallContractSimple 简化的合约调用（传入 ABI JSON 字符串）
func (k *Kit) CallContractSimple(ctx context.Context, contractAddress common.Address, abiJSON string, functionName string, params ...interface{}) ([]interface{}, error) {
	contractAbi, err := abi.JSON(strings.NewReader(abiJSON))
	if err != nil {
		return nil, err
	}
	return k.CallContract(ctx, contractAddress, contractAbi, functionName, params...)
}
