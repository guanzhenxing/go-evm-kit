package etherkit

import (
	"context"
	"crypto/ecdsa"
	"errors"
	"math/big"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

// Kit 相关常量
const (
	// GweiDecimals Gwei 的小数位数
	GweiDecimals = 9
	// DefaultWaitInterval 默认等待交易确认的轮询间隔
	DefaultWaitInterval = time.Second
)

// Kit 以太坊开发工具包，提供最便捷的使用方式
// 通过嵌入接口，所有 Provider 和 Wallet 的方法都可以直接调用
type Kit struct {
	*Wallet       // 嵌入 Wallet，获得所有钱包方法（包括 GetAddress、GetPrivateKey）
	EtherProvider // 嵌入 Provider 接口，直接调用所有 Provider 方法！
}

// NewKit 创建以太坊开发工具包
// 参数说明：
//   - hexPk: 十六进制私钥字符串（带或不带 0x 前缀）
//   - rawUrl: 以太坊节点 RPC URL（如 "https://eth-mainnet.g.alchemy.com/v2/your-api-key"）
//
// 返回：
//   - *Kit: 创建的 Kit 实例
//   - error: 如果创建失败则返回错误
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

// NewKitWithGeneratedKey 创建以太坊开发工具包（自动生成随机私钥）
// 适用于不需要导入已有私钥的场景，如测试、临时钱包等
// 会自动生成一个随机私钥并创建对应的 Kit 实例
// 参数说明：
//   - rawUrl: 以太坊节点 RPC URL（如 "https://eth-mainnet.g.alchemy.com/v2/your-api-key"）
//
// 返回：
//   - *Kit: 创建的 Kit 实例（包含新生成的私钥和地址）
//   - error: 如果创建失败则返回错误
//
// 注意：
//   - 生成的私钥是随机的，每次调用都会创建新的钱包
//   - 请妥善保存生成的私钥，可通过 kit.GetPrivateKey() 获取私钥对象，或使用 GetHexPrivateKey(kit.GetPrivateKey()) 获取十六进制字符串
//   - 适用于临时场景，生产环境建议使用 NewKit 导入已有私钥
func NewKitWithGeneratedKey(rawUrl string) (*Kit, error) {
	pk, err := GeneratePrivateKey()
	if err != nil {
		return nil, err
	}
	ep, err := NewProvider(rawUrl)
	if err != nil {
		return nil, err
	}
	return NewKitWithComponents(pk, ep)
}

// NewKitWithComponents 使用已有组件创建 Kit
// 适用于已经创建好私钥和 Provider 的情况，避免重复创建
// 参数说明：
//   - privateKey: 已存在的 ECDSA 私钥
//   - ep: 已存在的 EtherProvider 实例
//
// 返回：
//   - *Kit: 创建的 Kit 实例
//   - error: 如果创建失败则返回错误
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
// 按指定间隔轮询交易收据，直到交易被打包或超时
// 参数说明：
//   - ctx: 上下文对象
//   - txHash: 交易哈希
//   - timeout: 超时时间（如 30*time.Second）
//
// 返回：
//   - *types.Receipt: 交易收据，包含交易状态、gas 使用等信息
//   - error: 如果超时或查询失败则返回错误
func (k *Kit) WaitForReceipt(ctx context.Context, txHash common.Hash, timeout time.Duration) (*types.Receipt, error) {
	return k.WaitForReceiptWithInterval(ctx, txHash, timeout, DefaultWaitInterval)
}

// WaitForReceiptWithInterval 等待交易被打包，带超时控制和自定义轮询间隔
// 按指定间隔轮询交易收据，直到交易被打包或超时
// 参数说明：
//   - ctx: 上下文对象
//   - txHash: 交易哈希
//   - timeout: 超时时间（如 30*time.Second）
//   - interval: 轮询间隔（如 2*time.Second，建议不小于 1 秒以避免频繁请求）
//
// 返回：
//   - *types.Receipt: 交易收据，包含交易状态、gas 使用等信息
//   - error: 如果超时或查询失败则返回错误
func (k *Kit) WaitForReceiptWithInterval(ctx context.Context, txHash common.Hash, timeout time.Duration, interval time.Duration) (*types.Receipt, error) {
	if interval < time.Second {
		interval = DefaultWaitInterval // 最小间隔为 1 秒
	}

	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	ticker := time.NewTicker(interval)
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
// 这是 SendTx 和 WaitForReceipt 的组合方法，发送交易后自动等待打包
// 参数说明：
//   - ctx: 上下文对象
//   - to: 接收地址（合约地址或普通地址）
//   - nonce: 交易 nonce（0 表示自动计算）
//   - gasLimit: Gas 限制（0 表示自动估算）
//   - gasPrice: Gas 价格（nil 表示自动获取）
//   - value: 转账金额（nil 表示不转账）
//   - data: 交易数据（合约调用数据或 nil）
//   - timeout: 等待超时时间（如 30*time.Second）
//
// 返回：
//   - *types.Receipt: 交易收据
//   - error: 如果发送失败或超时则返回错误
func (k *Kit) SendTxAndWait(ctx context.Context, to common.Address, nonce, gasLimit uint64, gasPrice, value *big.Int, data []byte, timeout time.Duration) (*types.Receipt, error) {
	txHash, err := k.SendTx(ctx, to, nonce, gasLimit, gasPrice, value, data)
	if err != nil {
		return nil, err
	}
	return k.WaitForReceipt(ctx, txHash, timeout)
}

// GetBalanceInEther 获取以太币余额（以 ETH 为单位）
// 将 Wei 余额转换为 ETH，保留小数精度
// 参数说明：
//   - ctx: 上下文对象
//
// 返回：
//   - float64: 余额（以 ETH 为单位，如 1.5 表示 1.5 ETH）
//   - error: 如果查询失败或转换失败则返回错误
func (k *Kit) GetBalanceInEther(ctx context.Context) (float64, error) {
	balance, err := k.GetBalance(ctx)
	if err != nil {
		return 0, err
	}
	// 使用 ToDecimal 转换，以太币的 decimals 是 18
	ethBalance := ToDecimal(balance, EthDecimals)
	result, ok := ethBalance.Float64()
	if !ok {
		return 0, errors.New("failed to convert balance to float64")
	}
	return result, nil
}

// TransferEther 转账以太币（便捷方法）
// 将 ETH 金额转换为 Wei 并发送交易，自动计算 nonce、gasLimit 和 gasPrice
// 参数说明：
//   - ctx: 上下文对象
//   - to: 接收地址
//   - valueInEther: 转账金额（以 ETH 为单位，如 0.1 表示 0.1 ETH）
//
// 返回：
//   - common.Hash: 交易哈希
//   - error: 如果转账失败则返回错误
func (k *Kit) TransferEther(ctx context.Context, to common.Address, valueInEther float64) (common.Hash, error) {
	// 验证接收地址
	if !IsValidAddress(to) {
		return common.Hash{}, errors.New("invalid receiver address")
	}

	// 验证金额
	if valueInEther < 0 {
		return common.Hash{}, errors.New("transfer amount cannot be negative")
	}

	// 使用 ToWei 转换，以太币的 decimals 是 18
	value := ToWei(valueInEther, EthDecimals)
	return k.SendTx(ctx, to, 0, 0, nil, value, nil)
}

// GetChainInfo 获取链信息（便捷方法）
// 一次性获取链 ID、网络 ID 和当前区块号
// 参数说明：
//   - ctx: 上下文对象
//
// 返回：
//   - chainID: 链 ID
//   - networkID: 网络 ID
//   - blockNumber: 当前区块号
//   - error: 如果查询失败则返回错误
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

// StaticCall 静态调用合约方法（不花费 gas，不发送交易）
// 可以调用 view/pure 函数，也可以模拟调用非 view/pure 函数来查看执行结果
// 适用于读取合约状态、查询数据等场景
// 参数说明：
//   - ctx: 上下文对象
//   - contractAddress: 合约地址
//   - contractAbi: 合约 ABI 对象
//   - functionName: 函数名（如 "balanceOf", "totalSupply"）
//   - blockNumber: 区块号（nil 表示最新区块，可用于查询历史状态）
//   - from: 调用者地址（nil 表示使用 Kit 的地址）
//   - value: 模拟转账金额（nil 表示不转账，用于模拟 payable 函数）
//   - params: 函数参数（按函数定义顺序传入）
//
// 返回：
//   - []interface{}: 函数返回值数组（按函数定义顺序）
//   - error: 如果调用失败则返回错误
func (k *Kit) StaticCall(ctx context.Context, contractAddress common.Address, contractAbi abi.ABI, functionName string, blockNumber *big.Int, from *common.Address, value *big.Int, params ...interface{}) ([]interface{}, error) {
	// 输入验证
	if !IsValidAddress(contractAddress) {
		return nil, errors.New("invalid contract address")
	}
	if functionName == "" {
		return nil, errors.New("function name cannot be empty")
	}

	// 如果没有指定 from，使用 Kit 的地址
	callFrom := k.GetAddress()
	if from != nil {
		if !IsValidAddress(*from) {
			return nil, errors.New("invalid from address")
		}
		callFrom = *from
	}

	// 使用增强后的 CallContract 方法
	return k.CallContract(ctx, blockNumber, &callFrom, value, contractAddress, contractAbi, functionName, params...)
}

// StaticCallWithABIString 使用 ABI JSON 字符串进行静态调用（不花费 gas，不发送交易）
// 这是 StaticCall 的便捷版本，接受 ABI JSON 字符串而不是 ABI 对象
// 适用于从配置文件或 API 获取 ABI 的场景
// 使用示例：
//   - 简单调用（使用默认值）：StaticCallWithABIString(ctx, addr, abiJSON, "balanceOf", nil, nil, nil, userAddress)
//   - 指定区块号：StaticCallWithABIString(ctx, addr, abiJSON, "balanceOf", blockNum, nil, nil, userAddress)
//
// 参数说明：
//   - ctx: 上下文对象
//   - contractAddress: 合约地址
//   - abiJSON: ABI JSON 字符串（完整的合约 ABI 或单个函数的 ABI）
//   - functionName: 函数名（如 "balanceOf", "totalSupply"）
//   - blockNumber: 区块号（nil 表示最新区块）
//   - from: 调用者地址（nil 表示使用 Kit 的地址）
//   - value: 模拟转账金额（nil 表示不转账）
//   - params: 函数参数（按函数定义顺序传入）
//
// 返回：
//   - []interface{}: 函数返回值数组（按函数定义顺序）
//   - error: 如果调用失败则返回错误
func (k *Kit) StaticCallWithABIString(ctx context.Context, contractAddress common.Address, abiJSON string, functionName string, blockNumber *big.Int, from *common.Address, value *big.Int, params ...interface{}) ([]interface{}, error) {
	// 输入验证
	if abiJSON == "" {
		return nil, errors.New("ABI JSON string cannot be empty")
	}
	if functionName == "" {
		return nil, errors.New("function name cannot be empty")
	}

	contractAbi, err := abi.JSON(strings.NewReader(abiJSON))
	if err != nil {
		return nil, err
	}
	return k.StaticCall(ctx, contractAddress, contractAbi, functionName, blockNumber, from, value, params...)
}

// InvokeContract 调用合约方法并发送交易（花费 gas，会修改链上状态）
// 用于调用非 view/pure 函数，需要发送交易来执行
// 适用于转账代币、调用状态修改函数等场景
// 参数说明：
//   - ctx: 上下文对象
//   - contractAddress: 合约地址
//   - contractAbi: 合约 ABI 对象
//   - functionName: 函数名（如 "transfer", "approve"）
//   - nonce: 交易 nonce（0 表示自动计算）
//   - gasLimit: Gas 限制（0 表示自动估算）
//   - gasPrice: Gas 价格（nil 表示自动获取）
//   - value: 转账金额（nil 表示不转账，用于 payable 函数）
//   - params: 函数参数（按函数定义顺序传入）
//
// 返回：
//   - common.Hash: 交易哈希，可用于查询交易状态
//   - error: 如果发送失败则返回错误
func (k *Kit) InvokeContract(ctx context.Context, contractAddress common.Address, contractAbi abi.ABI, functionName string, nonce, gasLimit uint64, gasPrice, value *big.Int, params ...interface{}) (common.Hash, error) {
	// 输入验证
	if !IsValidAddress(contractAddress) {
		return common.Hash{}, errors.New("invalid contract address")
	}
	if functionName == "" {
		return common.Hash{}, errors.New("function name cannot be empty")
	}

	// 构建合约调用数据
	inputData, err := BuildContractInputData(contractAbi, functionName, params...)
	if err != nil {
		return common.Hash{}, err
	}

	// 发送交易
	return k.SendTx(ctx, contractAddress, nonce, gasLimit, gasPrice, value, inputData)
}

// InvokeContractWithABIString 使用 ABI JSON 字符串调用合约方法并发送交易（花费 gas）
// 这是 InvokeContract 的便捷版本，接受 ABI JSON 字符串而不是 ABI 对象
// 适用于从配置文件或 API 获取 ABI 的场景
// 使用示例：
//   - 调用转账：InvokeContractWithABIString(ctx, addr, abiJSON, "transfer", 0, 0, nil, nil, toAddress, amount)
//
// 参数说明：
//   - ctx: 上下文对象
//   - contractAddress: 合约地址
//   - abiJSON: ABI JSON 字符串（完整的合约 ABI 或单个函数的 ABI）
//   - functionName: 函数名（如 "transfer", "approve"）
//   - nonce: 交易 nonce（0 表示自动计算）
//   - gasLimit: Gas 限制（0 表示自动估算）
//   - gasPrice: Gas 价格（nil 表示自动获取）
//   - value: 转账金额（nil 表示不转账，用于 payable 函数）
//   - params: 函数参数（按函数定义顺序传入）
//
// 返回：
//   - common.Hash: 交易哈希，可用于查询交易状态
//   - error: 如果发送失败则返回错误
func (k *Kit) InvokeContractWithABIString(ctx context.Context, contractAddress common.Address, abiJSON string, functionName string, nonce, gasLimit uint64, gasPrice, value *big.Int, params ...interface{}) (common.Hash, error) {
	// 输入验证
	if abiJSON == "" {
		return common.Hash{}, errors.New("ABI JSON string cannot be empty")
	}
	if functionName == "" {
		return common.Hash{}, errors.New("function name cannot be empty")
	}

	contractAbi, err := abi.JSON(strings.NewReader(abiJSON))
	if err != nil {
		return common.Hash{}, err
	}
	return k.InvokeContract(ctx, contractAddress, contractAbi, functionName, nonce, gasLimit, gasPrice, value, params...)
}

// IsContract 检查地址是否为合约地址
// 通过检查地址的代码长度来判断是否为合约（合约代码长度 > 0）
// 参数说明：
//   - ctx: 上下文对象
//   - address: 要检查的地址
//
// 返回：
//   - bool: true 表示是合约地址，false 表示是普通地址
//   - error: 如果查询失败则返回错误
func (k *Kit) IsContract(ctx context.Context, address common.Address) (bool, error) {
	return k.EtherProvider.IsContractAddress(ctx, address)
}

// GetContractBytecode 获取合约字节码
// 返回合约在链上的字节码（十六进制字符串，带 0x 前缀）
// 参数说明：
//   - ctx: 上下文对象
//   - address: 合约地址
//
// 返回：
//   - string: 合约字节码（十六进制字符串，如 "0x608060405234801561001057600080fd5b50..."）
//   - error: 如果查询失败则返回错误
func (k *Kit) GetContractBytecode(ctx context.Context, address common.Address) (string, error) {
	return k.EtherProvider.GetContractBytecode(ctx, address)
}

// ============ 交易增强方法 ============

// SendTx 发送交易（不等待确认）
// 构建、签名并发送交易，返回交易哈希后立即返回，不等待交易打包
// 参数说明：
//   - ctx: 上下文对象
//   - to: 接收地址（合约地址或普通地址）
//   - nonce: 交易 nonce（0 表示自动计算）
//   - gasLimit: Gas 限制（0 表示自动估算）
//   - gasPrice: Gas 价格（nil 表示自动获取）
//   - value: 转账金额（nil 表示不转账）
//   - data: 交易数据（合约调用数据或 nil）
//
// 返回：
//   - common.Hash: 交易哈希，可用于后续查询交易状态
//   - error: 如果发送失败则返回错误
//
// 注意：此方法通过嵌入的 Wallet 提供，如需等待交易确认，请使用 SendTxAndWait
func (k *Kit) SendTx(ctx context.Context, to common.Address, nonce, gasLimit uint64, gasPrice, value *big.Int, data []byte) (common.Hash, error) {
	return k.Wallet.SendTx(ctx, to, nonce, gasLimit, gasPrice, value, data)
}

// SendTxWithHexInput 发送十六进制输入的交易（不等待确认）
// 构建、签名并发送交易，输入数据为十六进制字符串，返回交易哈希后立即返回
// 参数说明：
//   - ctx: 上下文对象
//   - to: 接收地址（合约地址或普通地址）
//   - nonce: 交易 nonce（0 表示自动计算）
//   - gasLimit: Gas 限制（0 表示自动估算）
//   - gasPrice: Gas 价格（nil 表示自动获取）
//   - value: 转账金额（nil 表示不转账）
//   - input: 十六进制输入数据（带或不带 0x 前缀，如 "0x1234..." 或 "1234..."）
//
// 返回：
//   - common.Hash: 交易哈希，可用于后续查询交易状态
//   - error: 如果发送失败则返回错误
//
// 注意：此方法通过嵌入的 Wallet 提供，如需等待交易确认，请使用 SendTxWithHexInputAndWait
func (k *Kit) SendTxWithHexInput(ctx context.Context, to common.Address, nonce, gasLimit uint64, gasPrice, value *big.Int, input string) (common.Hash, error) {
	return k.Wallet.SendTxWithHexInput(ctx, to, nonce, gasLimit, gasPrice, value, input)
}

// SendTxWithHexInputAndWait 发送十六进制输入的交易并等待确认
// 这是 SendTxWithHexInput 和 WaitForReceipt 的组合方法
// 参数说明：
//   - ctx: 上下文对象
//   - to: 接收地址（合约地址或普通地址）
//   - nonce: 交易 nonce（0 表示自动计算）
//   - gasLimit: Gas 限制（0 表示自动估算）
//   - gasPrice: Gas 价格（nil 表示自动获取）
//   - value: 转账金额（nil 表示不转账）
//   - input: 十六进制输入数据（带或不带 0x 前缀）
//   - timeout: 等待超时时间（如 30*time.Second）
//
// 返回：
//   - *types.Receipt: 交易收据
//   - error: 如果发送失败或超时则返回错误
func (k *Kit) SendTxWithHexInputAndWait(ctx context.Context, to common.Address, nonce, gasLimit uint64, gasPrice, value *big.Int, input string, timeout time.Duration) (*types.Receipt, error) {
	txHash, err := k.SendTxWithHexInput(ctx, to, nonce, gasLimit, gasPrice, value, input)
	if err != nil {
		return nil, err
	}
	return k.WaitForReceipt(ctx, txHash, timeout)
}

// TransferEtherAndWait 转账以太币并等待确认
// 这是 TransferEther 和 WaitForReceipt 的组合方法
// 参数说明：
//   - ctx: 上下文对象
//   - to: 接收地址
//   - valueInEther: 转账金额（以 ETH 为单位，如 0.1 表示 0.1 ETH）
//   - timeout: 等待超时时间（如 30*time.Second）
//
// 返回：
//   - *types.Receipt: 交易收据
//   - error: 如果转账失败或超时则返回错误
func (k *Kit) TransferEtherAndWait(ctx context.Context, to common.Address, valueInEther float64, timeout time.Duration) (*types.Receipt, error) {
	txHash, err := k.TransferEther(ctx, to, valueInEther)
	if err != nil {
		return nil, err
	}
	return k.WaitForReceipt(ctx, txHash, timeout)
}

// ============ 区块和交易查询增强方法 ============

// GetLatestBlock 获取最新区块
// 获取当前链上的最新区块信息，包括区块号、时间戳、交易列表等
// 参数说明：
//   - ctx: 上下文对象
//
// 返回：
//   - *types.Block: 区块对象，包含区块头、交易列表等信息
//   - error: 如果查询失败则返回错误
func (k *Kit) GetLatestBlock(ctx context.Context) (*types.Block, error) {
	blockNumber, err := k.GetBlockNumber(ctx)
	if err != nil {
		return nil, err
	}
	return k.GetBlockByNumber(ctx, big.NewInt(int64(blockNumber)))
}

// ============ Gas 和费用相关增强方法 ============

// GetSuggestedGasPriceInGwei 获取建议的 Gas 价格（以 Gwei 为单位）
// 将网络建议的 Gas 价格从 Wei 转换为 Gwei，便于阅读和使用
// 参数说明：
//   - ctx: 上下文对象
//
// 返回：
//   - float64: Gas 价格（以 Gwei 为单位，如 20.5 表示 20.5 Gwei）
//   - error: 如果查询失败或转换失败则返回错误
func (k *Kit) GetSuggestedGasPriceInGwei(ctx context.Context) (float64, error) {
	gasPrice, err := k.GetSuggestGasPrice(ctx)
	if err != nil {
		return 0, err
	}
	// 1 Gwei = 10^9 Wei
	gweiPrice := ToDecimal(gasPrice, GweiDecimals)
	result, ok := gweiPrice.Float64()
	if !ok {
		return 0, errors.New("failed to convert gas price to float64")
	}
	return result, nil
}

// EstimateGasForTransfer 估算以太币转账的 Gas
// 估算一笔以太币转账交易所需的 Gas 数量
// 参数说明：
//   - ctx: 上下文对象
//   - to: 接收地址
//   - valueInEther: 转账金额（以 ETH 为单位，如 0.1 表示 0.1 ETH）
//
// 返回：
//   - uint64: 估算的 Gas 数量
//   - error: 如果估算失败则返回错误
func (k *Kit) EstimateGasForTransfer(ctx context.Context, to common.Address, valueInEther float64) (uint64, error) {
	// 输入验证
	if !IsValidAddress(to) {
		return 0, errors.New("invalid receiver address")
	}
	if valueInEther < 0 {
		return 0, errors.New("transfer amount cannot be negative")
	}

	value := ToWei(valueInEther, EthDecimals)
	return k.EstimateGas(ctx, k.GetAddress(), to, 0, nil, value, nil)
}

// ============ 签名和验证增强方法 ============

// SignMessage 对消息进行签名
// 使用 Kit 的私钥对消息进行 ECDSA 签名
// 参数说明：
//   - ctx: 上下文对象（当前未使用，保留用于未来扩展）
//   - message: 要签名的消息（原始字节）
//
// 返回：
//   - []byte: 签名结果（65 字节，包含 r、s、v）
//   - error: 如果签名失败则返回错误
func (k *Kit) SignMessage(ctx context.Context, message []byte) ([]byte, error) {
	return k.Signature(message)
}

// VerifyMessage 验证消息签名
// 验证消息签名是否由指定的地址（Kit 的地址）签名
// 参数说明：
//   - ctx: 上下文对象（当前未使用，保留用于未来扩展）
//   - message: 原始消息（字节）
//   - signature: 签名结果（65 字节）
//
// 返回：
//   - bool: true 表示签名有效，false 表示签名无效
func (k *Kit) VerifyMessage(ctx context.Context, message, signature []byte) bool {
	return VerifySignature(k.GetAddress().Hex(), message, signature)
}

// ============ 钱包管理增强方法 ============

// GetBalanceInGwei 获取余额（以 Gwei 为单位）
// 将 Wei 余额转换为 Gwei，保留小数精度
// 参数说明：
//   - ctx: 上下文对象
//
// 返回：
//   - float64: 余额（以 Gwei 为单位，如 1000000.5 表示 1000000.5 Gwei）
//   - error: 如果查询失败或转换失败则返回错误
func (k *Kit) GetBalanceInGwei(ctx context.Context) (float64, error) {
	balance, err := k.GetBalance(ctx)
	if err != nil {
		return 0, err
	}
	gweiBalance := ToDecimal(balance, GweiDecimals)
	result, ok := gweiBalance.Float64()
	if !ok {
		return 0, errors.New("failed to convert balance to float64")
	}
	return result, nil
}

// GetFormattedBalance 获取格式化的余额字符串
// 将余额格式化为易读的字符串，显示全部有效数字，并添加 " ETH" 后缀
// 参数说明：
//   - ctx: 上下文对象
//
// 返回：
//   - string: 格式化的余额字符串（如 "1.234567890123456789 ETH"）
//   - error: 如果查询失败则返回错误
func (k *Kit) GetFormattedBalance(ctx context.Context) (string, error) {
	balance, err := k.GetBalance(ctx)
	if err != nil {
		return "", err
	}
	ethBalance := ToDecimal(balance, EthDecimals)
	return ethBalance.String() + " ETH", nil
}

// ============ 网络状态增强方法 ============

// GetNetworkStatus 获取网络状态信息
// 一次性获取链 ID、网络 ID、当前区块号和 Gas 价格等网络状态信息
// 参数说明：
//   - ctx: 上下文对象
//
// 返回：
//   - map[string]interface{}: 网络状态信息，包含以下键：
//   - "chain_id": 链 ID (*big.Int)
//   - "network_id": 网络 ID (*big.Int)
//   - "block_number": 当前区块号 (uint64)
//   - "gas_price": Gas 价格 (*big.Int，单位为 Wei)
//   - error: 如果查询失败则返回错误
func (k *Kit) GetNetworkStatus(ctx context.Context) (map[string]interface{}, error) {
	chainID, err := k.GetChainID(ctx)
	if err != nil {
		return nil, err
	}

	networkID, err := k.GetNetworkID(ctx)
	if err != nil {
		return nil, err
	}

	blockNumber, err := k.GetBlockNumber(ctx)
	if err != nil {
		return nil, err
	}

	gasPrice, err := k.GetSuggestGasPrice(ctx)
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"chain_id":     chainID,
		"network_id":   networkID,
		"block_number": blockNumber,
		"gas_price":    gasPrice,
	}, nil
}
