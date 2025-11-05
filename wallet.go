package etherkit

import (
	"context"
	"crypto/ecdsa"
	"math/big"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
)

// EtherWallet 以太坊钱包接口
// 提供钱包管理、交易构建、签名和发送等功能
type EtherWallet interface {
	// GetEthProvider 获取以太坊提供者实例
	// 返回：
	//   - EtherProvider: 以太坊提供者接口
	GetEthProvider() EtherProvider
	// GetClient 获取以太坊客户端实例
	// 返回：
	//   - *ethclient.Client: 以太坊客户端
	GetClient() *ethclient.Client
	// GetAddress 获取钱包地址
	// 返回：
	//   - common.Address: 钱包地址
	GetAddress() common.Address
	// GetPrivateKey 获取私钥
	// 返回：
	//   - *ecdsa.PrivateKey: ECDSA 私钥对象
	GetPrivateKey() *ecdsa.PrivateKey
	// CloseWallet 关闭钱包连接
	// 释放所有底层资源
	CloseWallet()
	// GetNonce 获取账户的 nonce
	// nonce 用于防止交易重放，每个交易必须使用唯一的 nonce
	// 参数说明：
	//   - ctx: 上下文对象
	// 返回：
	//   - uint64: 下一个可用的 nonce
	//   - error: 如果查询失败则返回错误
	GetNonce(ctx context.Context) (uint64, error)
	// GetBalance 获取账户余额
	// 返回账户的本位币余额（单位为 Wei）
	// 参数说明：
	//   - ctx: 上下文对象
	// 返回：
	//   - *big.Int: 余额（单位为 Wei）
	//   - error: 如果查询失败则返回错误
	GetBalance(ctx context.Context) (*big.Int, error)
	// NewTx 构建一笔交易
	// 自动计算 nonce、gasLimit 和 gasPrice（如果未提供）
	// 参数说明：
	//   - ctx: 上下文对象
	//   - to: 接收地址（合约地址或普通地址）
	//   - nonce: 交易 nonce（0 表示自动计算）
	//   - gasLimit: Gas 限制（0 表示自动估算）
	//   - gasPrice: Gas 价格（nil 或 0 表示自动获取）
	//   - value: 转账金额（nil 表示不转账）
	//   - data: 交易数据（合约调用数据或 nil）
	// 返回：
	//   - *types.Transaction: 交易对象
	//   - error: 如果构建失败则返回错误
	NewTx(ctx context.Context, to common.Address, nonce, gasLimit uint64, gasPrice, value *big.Int, data []byte) (*types.Transaction, error)
	// SendTx 发送交易
	// 构建、签名并发送交易，返回交易哈希
	// 参数说明：
	//   - ctx: 上下文对象
	//   - to: 接收地址（合约地址或普通地址）
	//   - nonce: 交易 nonce（0 表示自动计算）
	//   - gasLimit: Gas 限制（0 表示自动估算）
	//   - gasPrice: Gas 价格（nil 或 0 表示自动获取）
	//   - value: 转账金额（nil 表示不转账）
	//   - data: 交易数据（合约调用数据或 nil）
	// 返回：
	//   - common.Hash: 交易哈希，可用于查询交易状态
	//   - error: 如果发送失败则返回错误
	SendTx(ctx context.Context, to common.Address, nonce, gasLimit uint64, gasPrice, value *big.Int, data []byte) (common.Hash, error)
	// NewTxWithHexInput 构建一笔交易，使用十六进制输入数据
	// 与 NewTx 类似，但接受十六进制字符串作为输入数据
	// 参数说明：
	//   - ctx: 上下文对象
	//   - to: 接收地址
	//   - nonce: 交易 nonce（0 表示自动计算）
	//   - gasLimit: Gas 限制（0 表示自动估算）
	//   - gasPrice: Gas 价格（nil 或 0 表示自动获取）
	//   - value: 转账金额（nil 表示不转账）
	//   - input: 十六进制输入数据（带或不带 0x 前缀）
	// 返回：
	//   - *types.Transaction: 交易对象
	//   - error: 如果构建失败则返回错误
	NewTxWithHexInput(ctx context.Context, to common.Address, nonce, gasLimit uint64, gasPrice, value *big.Int, input string) (*types.Transaction, error)
	// SendTxWithHexInput 发送一笔交易，使用十六进制输入数据
	// 与 SendTx 类似，但接受十六进制字符串作为输入数据
	// 参数说明：
	//   - ctx: 上下文对象
	//   - to: 接收地址
	//   - nonce: 交易 nonce（0 表示自动计算）
	//   - gasLimit: Gas 限制（0 表示自动估算）
	//   - gasPrice: Gas 价格（nil 或 0 表示自动获取）
	//   - value: 转账金额（nil 表示不转账）
	//   - input: 十六进制输入数据（带或不带 0x 前缀）
	// 返回：
	//   - common.Hash: 交易哈希
	//   - error: 如果发送失败则返回错误
	SendTxWithHexInput(ctx context.Context, to common.Address, nonce, gasLimit uint64, gasPrice, value *big.Int, input string) (common.Hash, error)
	// BuildTxOpts 构建交易的选项
	// 用于与 go-ethereum 的 bind 包配合使用，生成 TransactOpts
	// 参数说明：
	//   - ctx: 上下文对象
	//   - value: 转账金额（nil 表示不转账）
	//   - nonce: 交易 nonce（nil 或 <= 0 表示自动计算）
	//   - gasPrice: Gas 价格（nil 或 <= 0 表示自动获取）
	// 返回：
	//   - *bind.TransactOpts: 交易选项，可用于合约交互
	//   - error: 如果构建失败则返回错误
	BuildTxOpts(ctx context.Context, value, nonce, gasPrice *big.Int) (*bind.TransactOpts, error)
	// SignTx 对交易进行签名
	// 使用钱包的私钥对交易进行 EIP-155 签名
	// 参数说明：
	//   - ctx: 上下文对象
	//   - tx: 未签名的交易对象
	// 返回：
	//   - *types.Transaction: 已签名的交易对象
	//   - error: 如果签名失败则返回错误
	SignTx(ctx context.Context, tx *types.Transaction) (*types.Transaction, error)
	// SendSignedTx 发送已签名的交易
	// 将已签名的交易发送到网络
	// 参数说明：
	//   - ctx: 上下文对象
	//   - signedTx: 已签名的交易对象
	// 返回：
	//   - common.Hash: 交易哈希
	//   - error: 如果发送失败则返回错误
	SendSignedTx(ctx context.Context, signedTx *types.Transaction) (common.Hash, error)
	// Signature 对数据进行签名
	// 使用钱包的私钥对数据进行 ECDSA 签名
	// 参数说明：
	//   - data: 要签名的原始数据（字节）
	// 返回：
	//   - []byte: 签名结果（65 字节，包含 r、s、v）
	//   - error: 如果签名失败则返回错误
	Signature(data []byte) ([]byte, error)
	// CallContract 调用合约方法（静态调用，不发送交易）
	// 可以调用 view/pure 函数，也可以模拟调用非 view/pure 函数来查看执行结果
	// 参数说明：
	//   - ctx: 上下文对象
	//   - blockNumber: 区块号（nil 表示最新区块）
	//   - from: 调用者地址（nil 表示不设置）
	//   - value: 模拟转账金额（nil 表示不转账）
	//   - contractAddress: 合约地址
	//   - contractAbi: 合约 ABI 对象
	//   - functionName: 函数名
	//   - params: 函数参数（按函数定义顺序传入）
	// 返回：
	//   - []interface{}: 函数返回值数组（按函数定义顺序）
	//   - error: 如果调用失败则返回错误
	CallContract(ctx context.Context, blockNumber *big.Int, from *common.Address, value *big.Int, contractAddress common.Address, contractAbi abi.ABI, functionName string, params ...interface{}) ([]interface{}, error)
}

// Wallet 以太坊钱包实现
// 封装了私钥、地址和提供者，提供钱包管理、交易构建、签名和发送等功能
type Wallet struct {
	privateKey *ecdsa.PrivateKey // ECDSA 私钥
	address    common.Address    // 钱包地址（从私钥派生）
	ep         EtherProvider     // 以太坊提供者
}

// NewWallet 创建新的钱包实例
// 从十六进制私钥字符串创建钱包，并连接到指定的以太坊节点
// 参数说明：
//   - hexPk: 十六进制私钥字符串（带或不带 0x 前缀）
//   - rawUrl: 以太坊节点 RPC URL（如 "https://eth-mainnet.g.alchemy.com/v2/your-api-key"）
//
// 返回：
//   - *Wallet: 创建的钱包实例
//   - error: 如果创建失败则返回错误
func NewWallet(hexPk string, rawUrl string) (*Wallet, error) {
	privateKey, err := BuildPrivateKeyFromHex(hexPk)
	if err != nil {
		return nil, err
	}

	ep, err := NewProvider(rawUrl)
	if err != nil {
		return nil, err
	}

	return &Wallet{
		privateKey: privateKey,
		address:    PrivateKeyToAddress(privateKey),
		ep:         ep,
	}, nil
}

// NewWalletWithComponents 使用已有组件创建钱包实例
// 适用于已经创建好私钥和 Provider 的情况，避免重复创建
// 参数说明：
//   - privateKey: 已存在的 ECDSA 私钥
//   - ep: 已存在的 EtherProvider 实例
//
// 返回：
//   - *Wallet: 创建的钱包实例
//   - error: 如果创建失败则返回错误
func NewWalletWithComponents(privateKey *ecdsa.PrivateKey, ep EtherProvider) (*Wallet, error) {
	return &Wallet{
		privateKey: privateKey,
		address:    PrivateKeyToAddress(privateKey),
		ep:         ep,
	}, nil
}

// GetEthProvider 获取以太坊提供者实例
// 返回：
//   - EtherProvider: 以太坊提供者接口
func (w *Wallet) GetEthProvider() EtherProvider {
	return w.ep
}

// GetClient 获取以太坊客户端实例
// 返回：
//   - *ethclient.Client: 以太坊客户端
func (w *Wallet) GetClient() *ethclient.Client {
	return w.ep.GetEthClient()
}

// GetAddress 获取钱包地址
// 返回：
//   - common.Address: 钱包地址
func (w *Wallet) GetAddress() common.Address {
	return w.address
}

// GetPrivateKey 获取私钥
// 返回：
//   - *ecdsa.PrivateKey: ECDSA 私钥对象
//
// 注意：请妥善保管私钥，泄露私钥将导致资产丢失
func (w *Wallet) GetPrivateKey() *ecdsa.PrivateKey {
	return w.privateKey
}

// CloseWallet 关闭钱包连接
// 释放所有底层资源，包括 Provider 的连接
// 建议在程序退出或不再使用时调用此方法
func (w *Wallet) CloseWallet() {
	w.ep.Close()
}

// GetNonce 获取账户的 nonce
// nonce 用于防止交易重放，每个交易必须使用唯一的 nonce
// 返回待处理状态的 nonce（pending nonce），即下一个可用的 nonce
// 参数说明：
//   - ctx: 上下文对象
//
// 返回：
//   - uint64: 下一个可用的 nonce
//   - error: 如果查询失败则返回错误
func (w *Wallet) GetNonce(ctx context.Context) (uint64, error) {
	return w.GetClient().PendingNonceAt(ctx, w.GetAddress())
}

// GetBalance 获取账户余额
// 返回账户的本位币余额（单位为 Wei）
// 参数说明：
//   - ctx: 上下文对象
//
// 返回：
//   - *big.Int: 余额（单位为 Wei）
//   - error: 如果查询失败则返回错误
func (w *Wallet) GetBalance(ctx context.Context) (*big.Int, error) {
	return w.GetClient().BalanceAt(ctx, w.GetAddress(), nil)
}

// NewTx 构建一笔交易
// 自动计算 nonce、gasLimit 和 gasPrice（如果未提供）
// 参数说明：
//   - ctx: 上下文对象
//   - to: 接收地址（合约地址或普通地址）
//   - nonce: 交易 nonce（0 表示自动计算）
//   - gasLimit: Gas 限制（0 表示自动估算）
//   - gasPrice: Gas 价格（nil 或 big.NewInt(0) 表示自动获取）
//   - value: 转账金额（nil 表示不转账）
//   - data: 交易数据（合约调用数据或 nil）
//
// 返回：
//   - *types.Transaction: 交易对象（未签名）
//   - error: 如果构建失败则返回错误
func (w *Wallet) NewTx(ctx context.Context, to common.Address, nonce, gasLimit uint64, gasPrice, value *big.Int, data []byte) (*types.Transaction, error) {

	if nonce == 0 {
		var err error
		nonce, err = w.GetNonce(ctx)
		if err != nil {
			return nil, err
		}
	}

	if gasPrice == nil || gasPrice.Sign() == 0 {
		var err error
		gasPrice, err = w.GetEthProvider().GetSuggestGasPrice(ctx)
		if err != nil {
			return nil, err
		}
	}

	if gasLimit == 0 {
		var err error
		gasLimit, err = w.ep.EstimateGas(ctx, w.GetAddress(), to, nonce, gasPrice, value, data)
		if err != nil {
			return nil, err
		}
	}

	return NewTx(to, nonce, gasLimit, gasPrice, value, data)
}

// SendTx 发送交易
// 构建、签名并发送交易，返回交易哈希
// 参数说明：
//   - ctx: 上下文对象
//   - to: 接收地址（合约地址或普通地址）
//   - nonce: 交易 nonce（0 表示自动计算）
//   - gasLimit: Gas 限制（0 表示自动估算）
//   - gasPrice: Gas 价格（nil 或 big.NewInt(0) 表示自动获取）
//   - value: 转账金额（nil 表示不转账）
//   - data: 交易数据（合约调用数据或 nil）
//
// 返回：
//   - common.Hash: 交易哈希，可用于查询交易状态
//   - error: 如果发送失败则返回错误
func (w *Wallet) SendTx(ctx context.Context, to common.Address, nonce, gasLimit uint64, gasPrice, value *big.Int, data []byte) (common.Hash, error) {

	tx, err := w.NewTx(ctx, to, nonce, gasLimit, gasPrice, value, data)
	if err != nil {
		return [32]byte{}, err
	}

	signedTx, err := w.SignTx(ctx, tx)
	if err != nil {
		return [32]byte{}, err
	}

	return w.SendSignedTx(ctx, signedTx)
}

// NewTxWithHexInput 构建一笔交易，使用十六进制输入数据
// 与 NewTx 类似，但接受十六进制字符串作为输入数据
// 参数说明：
//   - ctx: 上下文对象
//   - to: 接收地址
//   - nonce: 交易 nonce（0 表示自动计算）
//   - gasLimit: Gas 限制（0 表示自动估算）
//   - gasPrice: Gas 价格（nil 或 big.NewInt(0) 表示自动获取）
//   - value: 转账金额（nil 表示不转账）
//   - input: 十六进制输入数据（带或不带 0x 前缀，如 "0x1234..." 或 "1234..."）
//
// 返回：
//   - *types.Transaction: 交易对象（未签名）
//   - error: 如果构建失败则返回错误
func (w *Wallet) NewTxWithHexInput(ctx context.Context, to common.Address, nonce, gasLimit uint64, gasPrice, value *big.Int, input string) (*types.Transaction, error) {
	data, err := hexutil.Decode(input)
	if err != nil {
		return nil, err
	}
	return w.NewTx(ctx, to, nonce, gasLimit, gasPrice, value, data)
}

// SendTxWithHexInput 发送一笔交易，使用十六进制输入数据
// 与 SendTx 类似，但接受十六进制字符串作为输入数据
// 参数说明：
//   - ctx: 上下文对象
//   - to: 接收地址
//   - nonce: 交易 nonce（0 表示自动计算）
//   - gasLimit: Gas 限制（0 表示自动估算）
//   - gasPrice: Gas 价格（nil 或 big.NewInt(0) 表示自动获取）
//   - value: 转账金额（nil 表示不转账）
//   - input: 十六进制输入数据（带或不带 0x 前缀）
//
// 返回：
//   - common.Hash: 交易哈希
//   - error: 如果发送失败则返回错误
func (w *Wallet) SendTxWithHexInput(ctx context.Context, to common.Address, nonce, gasLimit uint64, gasPrice, value *big.Int, input string) (common.Hash, error) {
	data, err := hexutil.Decode(input)
	if err != nil {
		return [32]byte{}, err
	}
	return w.SendTx(ctx, to, nonce, gasLimit, gasPrice, value, data)
}

// BuildTxOpts 构建交易的选项
// 用于与 go-ethereum 的 bind 包配合使用，生成 TransactOpts
// 参数说明：
//   - ctx: 上下文对象
//   - value: 转账金额（nil 表示不转账）
//   - nonce: 交易 nonce（nil 或 <= 0 表示自动计算）
//   - gasPrice: Gas 价格（nil 或 <= 0 表示自动获取）
//
// 返回：
//   - *bind.TransactOpts: 交易选项，可用于合约交互
//   - error: 如果构建失败则返回错误
func (w *Wallet) BuildTxOpts(ctx context.Context, value, nonce, gasPrice *big.Int) (*bind.TransactOpts, error) {

	chainId, err := w.ep.GetChainID(ctx)
	if err != nil {
		return nil, err
	}

	txOpts, _ := bind.NewKeyedTransactorWithChainID(w.privateKey, chainId)

	txOpts.Value = value

	if gasPrice != nil && gasPrice.Sign() == 1 {
		txOpts.GasPrice = gasPrice
	} else {
		_gasPrice, err := w.GetEthProvider().GetSuggestGasPrice(ctx)
		if err != nil {
			return nil, err
		}
		txOpts.GasPrice = _gasPrice
	}

	// 如果nonce不为nil，就用传入的值（这里默认 nonce >0 ）
	if nonce != nil && nonce.Sign() > 0 {
		txOpts.Nonce = nonce
	} else {
		_nonce, err := w.GetNonce(ctx)
		if err != nil {
			return nil, err
		}
		txOpts.Nonce = big.NewInt(int64(_nonce))
	}

	return txOpts, nil
}

// SignTx 对交易进行签名
// 使用钱包的私钥对交易进行 EIP-155 签名（伦敦签名）
// 参数说明：
//   - ctx: 上下文对象
//   - tx: 未签名的交易对象
//
// 返回：
//   - *types.Transaction: 已签名的交易对象
//   - error: 如果签名失败则返回错误
func (w *Wallet) SignTx(ctx context.Context, tx *types.Transaction) (*types.Transaction, error) {

	chainId, err := w.ep.GetChainID(ctx)
	if err != nil {
		return nil, err
	}

	// 使用伦敦签名
	signer := types.NewLondonSigner(chainId)
	signedTx, err := types.SignTx(tx, signer, w.privateKey)
	if err != nil {
		return &types.Transaction{}, err
	}

	return signedTx, nil
}

// SendSignedTx 发送已签名的交易
// 将已签名的交易发送到网络
// 参数说明：
//   - ctx: 上下文对象
//   - signedTx: 已签名的交易对象
//
// 返回：
//   - common.Hash: 交易哈希
//   - error: 如果发送失败则返回错误（如余额不足、nonce 错误等）
func (w *Wallet) SendSignedTx(ctx context.Context, signedTx *types.Transaction) (common.Hash, error) {
	err := w.GetClient().SendTransaction(ctx, signedTx)
	if err != nil {
		return [32]byte{}, err
	}
	return signedTx.Hash(), nil
}

// Signature 对数据进行签名
// 使用钱包的私钥对数据进行 ECDSA 签名
// 先对数据进行 Keccak256 哈希，然后使用私钥签名
// 参数说明：
//   - data: 要签名的原始数据（字节）
//
// 返回：
//   - []byte: 签名结果（65 字节，包含 r、s、v）
//   - error: 如果签名失败则返回错误
func (w *Wallet) Signature(data []byte) ([]byte, error) {
	hash := crypto.Keccak256Hash(data)
	return crypto.Sign(hash.Bytes(), w.privateKey)
}

// CallContract 调用合约方法（静态调用，不发送交易）
// 可以调用 view/pure 函数，也可以模拟调用非 view/pure 函数来查看执行结果
// 参数说明：
//   - ctx: 上下文对象
//   - blockNumber: 区块号（nil 表示最新区块，可用于查询历史状态）
//   - from: 调用者地址（nil 表示不设置）
//   - value: 模拟转账金额（nil 表示不转账，用于模拟 payable 函数）
//   - contractAddress: 合约地址
//   - contractAbi: 合约 ABI 对象
//   - functionName: 函数名（如 "balanceOf", "totalSupply"）
//   - params: 函数参数（按函数定义顺序传入）
//
// 返回：
//   - []interface{}: 函数返回值数组（按函数定义顺序）
//   - error: 如果调用失败则返回错误
func (w *Wallet) CallContract(ctx context.Context, blockNumber *big.Int, from *common.Address, value *big.Int, contractAddress common.Address, contractAbi abi.ABI, functionName string, params ...interface{}) ([]interface{}, error) {

	inputData, err := BuildContractInputData(contractAbi, functionName, params...)
	if err != nil {
		return nil, err
	}

	// 构建 CallMsg
	callMsg := ethereum.CallMsg{
		To:   &contractAddress,
		Data: inputData,
	}

	// 如果指定了 from，则设置
	if from != nil {
		callMsg.From = *from
	}

	// 如果指定了 value，则设置
	if value != nil {
		callMsg.Value = value
	}

	// 执行静态调用（blockNumber 为 nil 表示最新区块）
	res, err := w.GetClient().CallContract(ctx, callMsg, blockNumber)
	if err != nil {
		return nil, err
	}

	response, err := contractAbi.Unpack(functionName, res)
	if err != nil {
		return nil, err
	}
	return response, nil
}
