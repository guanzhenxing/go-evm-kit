package etherkit

import (
	"context"
	"encoding/hex"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rpc"
)

// EtherProvider 以太坊提供者接口
// 提供链上查询功能，但不需要账户信息（如私钥）
// 用于执行只读操作，如查询区块、交易、余额等
type EtherProvider interface {
	// GetEthClient 获取以太坊客户端实例
	// 返回底层的 ethclient.Client，可用于执行底层操作
	GetEthClient() *ethclient.Client
	// GetRpcClient 获取 RPC 客户端实例
	// 返回底层的 rpc.Client，可用于执行底层 RPC 调用
	GetRpcClient() *rpc.Client
	// Close 关闭客户端连接
	// 释放所有底层资源，包括 ethclient 和 rpc client
	Close()
	// GetNetworkID 获取网络 ID
	// 参数说明：
	//   - ctx: 上下文对象
	// 返回：
	//   - *big.Int: 网络 ID
	//   - error: 如果查询失败则返回错误
	GetNetworkID(ctx context.Context) (*big.Int, error)
	// GetChainID 获取链 ID
	// 参数说明：
	//   - ctx: 上下文对象
	// 返回：
	//   - *big.Int: 链 ID（如主网为 1，Goerli 为 5）
	//   - error: 如果查询失败则返回错误
	GetChainID(ctx context.Context) (*big.Int, error)
	// GetBlockByHash 根据区块哈希获取区块信息
	// 参数说明：
	//   - ctx: 上下文对象
	//   - hash: 区块哈希
	// 返回：
	//   - *types.Block: 区块对象，包含区块头、交易列表等信息
	//   - error: 如果查询失败则返回错误
	GetBlockByHash(ctx context.Context, hash common.Hash) (*types.Block, error)
	// GetBlockByNumber 根据区块号获取区块信息
	// 参数说明：
	//   - ctx: 上下文对象
	//   - number: 区块号（nil 表示最新区块）
	// 返回：
	//   - *types.Block: 区块对象，包含区块头、交易列表等信息
	//   - error: 如果查询失败则返回错误
	GetBlockByNumber(ctx context.Context, number *big.Int) (*types.Block, error)
	// GetBlockNumber 获取最新区块号
	// 参数说明：
	//   - ctx: 上下文对象
	// 返回：
	//   - uint64: 最新区块号
	//   - error: 如果查询失败则返回错误
	GetBlockNumber(ctx context.Context) (uint64, error)
	// GetSuggestGasPrice 获取建议的 Gas 价格
	// 返回网络建议的 Gas 价格（单位为 Wei）
	// 参数说明：
	//   - ctx: 上下文对象
	// 返回：
	//   - *big.Int: 建议的 Gas 价格（单位为 Wei）
	//   - error: 如果查询失败则返回错误
	GetSuggestGasPrice(ctx context.Context) (*big.Int, error)
	// GetTransactionByHash 根据交易哈希获取交易信息
	// 参数说明：
	//   - ctx: 上下文对象
	//   - hash: 交易哈希
	// 返回：
	//   - tx: 交易对象
	//   - isPending: 交易是否还在待处理状态
	//   - error: 如果查询失败则返回错误
	GetTransactionByHash(ctx context.Context, hash common.Hash) (tx *types.Transaction, isPending bool, err error)
	// GetTransactionReceipt 根据交易哈希获取交易收据
	// 交易收据包含交易执行结果、gas 使用情况、日志等信息
	// 参数说明：
	//   - ctx: 上下文对象
	//   - txHash: 交易哈希
	// 返回：
	//   - *types.Receipt: 交易收据，包含交易状态、gas 使用等信息
	//   - error: 如果查询失败则返回错误
	GetTransactionReceipt(ctx context.Context, txHash common.Hash) (*types.Receipt, error)
	// GetContractBytecode 获取合约字节码
	// 返回合约在链上的字节码（十六进制字符串，带 0x 前缀）
	// 参数说明：
	//   - ctx: 上下文对象
	//   - address: 合约地址
	// 返回：
	//   - string: 合约字节码（十六进制字符串）
	//   - error: 如果查询失败则返回错误
	GetContractBytecode(ctx context.Context, address common.Address) (string, error)
	// IsContractAddress 检查地址是否为合约地址
	// 通过检查地址的代码长度来判断是否为合约（合约代码长度 > 0）
	// 参数说明：
	//   - ctx: 上下文对象
	//   - address: 要检查的地址
	// 返回：
	//   - bool: true 表示是合约地址，false 表示是普通地址
	//   - error: 如果查询失败则返回错误
	IsContractAddress(ctx context.Context, address common.Address) (bool, error)
	// EstimateGas 估算交易所需的 Gas 数量
	// 通过模拟交易执行来估算 gas 消耗
	// 参数说明：
	//   - ctx: 上下文对象
	//   - from: 发送地址
	//   - to: 接收地址（合约地址或普通地址）
	//   - nonce: 交易 nonce
	//   - gasPrice: Gas 价格（nil 表示不设置）
	//   - value: 转账金额（nil 表示不转账）
	//   - data: 交易数据（合约调用数据或 nil）
	// 返回：
	//   - uint64: 估算的 Gas 数量
	//   - error: 如果估算失败则返回错误
	EstimateGas(ctx context.Context, from, to common.Address, nonce uint64, gasPrice, value *big.Int, data []byte) (uint64, error)
	// GetFromAddress 从交易中提取发送地址
	// 通过解析交易签名来获取发送者地址
	// 参数说明：
	//   - tx: 交易对象
	// 返回：
	//   - common.Address: 发送地址
	//   - error: 如果提取失败则返回错误
	GetFromAddress(tx *types.Transaction) (common.Address, error)
	// FilterLogs 查询事件日志
	// 用于查询指定区块范围内的事件日志，支持按合约地址、事件签名和 indexed 参数进行过滤
	// 参数说明：
	//   - ctx: 上下文对象
	//   - contractAddress: 合约地址（nil 表示查询所有合约）
	//   - eventTopic: 事件签名 topic（如 GetEventTopic("Transfer(address,address,uint256)")）
	//   - fromBlock: 起始区块号（nil 表示从最新区块开始）
	//   - toBlock: 结束区块号（nil 表示到最新区块）
	//   - indexedTopics: 可选的 indexed 参数过滤（nil 表示不过滤，每个元素对应一个 indexed 参数）
	// 返回：
	//   - []types.Log: 事件日志列表，用户需要自行解析 Data 和 Topics
	//   - error: 如果查询失败则返回错误
	FilterLogs(ctx context.Context, contractAddress *common.Address, eventTopic common.Hash, fromBlock, toBlock *big.Int, indexedTopics []common.Hash) ([]types.Log, error)
}

// Provider 以太坊提供者实现
// 封装了与以太坊节点通信的底层客户端
type Provider struct {
	rc      *rpc.Client       // RPC 客户端
	ec      *ethclient.Client // 以太坊客户端
	chainId *big.Int          // 链 ID（缓存，避免重复查询）
}

// NewProvider 创建新的以太坊提供者实例
// 连接到指定的以太坊节点 RPC URL
// 参数说明：
//   - rawUrl: 以太坊节点 RPC URL（如 "https://eth-mainnet.g.alchemy.com/v2/your-api-key" 或 "http://localhost:8545"）
//
// 返回：
//   - *Provider: 创建的 Provider 实例
//   - error: 如果连接失败则返回错误
func NewProvider(rawUrl string) (*Provider, error) {

	rpcClient, err := rpc.Dial(rawUrl)
	if err != nil {
		return nil, fmt.Errorf("failed to rpc.Dial(): %w", err)
	}

	return &Provider{
		rc: rpcClient,
		ec: ethclient.NewClient(rpcClient),
	}, nil
}

// NewProviderWithChainId 创建新的以太坊提供者实例（指定链 ID）
// 预先设置链 ID，避免首次调用时查询链 ID 的网络延迟
// 适用于已知链 ID 的场景，可以提高性能
// 参数说明：
//   - rawUrl: 以太坊节点 RPC URL
//   - chainId: 链 ID（如主网为 1，Goerli 为 5）
//
// 返回：
//   - *Provider: 创建的 Provider 实例
//   - error: 如果连接失败则返回错误
func NewProviderWithChainId(rawUrl string, chainId int64) (*Provider, error) {

	p, err := NewProvider(rawUrl)
	if err != nil {
		return nil, err
	}
	p.chainId = big.NewInt(chainId)

	return p, nil
}

// GetEthClient 获取以太坊客户端实例
// 返回底层的 ethclient.Client，可用于执行底层操作
// 返回：
//   - *ethclient.Client: 以太坊客户端实例
func (p *Provider) GetEthClient() *ethclient.Client {
	return p.ec
}

// GetRpcClient 获取 RPC 客户端实例
// 返回底层的 rpc.Client，可用于执行底层 RPC 调用
// 返回：
//   - *rpc.Client: RPC 客户端实例
func (p *Provider) GetRpcClient() *rpc.Client {
	return p.rc
}

// Close 关闭客户端连接
// 释放所有底层资源，包括 ethclient 和 rpc client
// 建议在程序退出或不再使用时调用此方法
func (p *Provider) Close() {
	p.ec.Close()
	p.rc.Close()
}

// GetNetworkID 获取网络 ID
// 网络 ID 用于标识不同的网络（主网、测试网等）
// 参数说明：
//   - ctx: 上下文对象
//
// 返回：
//   - *big.Int: 网络 ID
//   - error: 如果查询失败则返回错误
func (p *Provider) GetNetworkID(ctx context.Context) (*big.Int, error) {
	return p.ec.NetworkID(ctx)
}

// GetChainID 获取链 ID
// 链 ID 用于 EIP-155 签名，防止重放攻击
// 如果已缓存链 ID，则直接返回；否则查询链上并缓存
// 参数说明：
//   - ctx: 上下文对象
//
// 返回：
//   - *big.Int: 链 ID（如主网为 1，Goerli 为 5）
//   - error: 如果查询失败则返回错误
func (p *Provider) GetChainID(ctx context.Context) (*big.Int, error) {

	if p.chainId == nil {
		chainId, err := p.ec.ChainID(ctx)
		if err != nil {
			return nil, err
		}
		p.chainId = chainId
	}

	return p.chainId, nil
}

// GetBlockByHash 根据区块哈希获取区块信息
// 通过区块哈希查询完整的区块信息，包括区块头和所有交易
// 参数说明：
//   - ctx: 上下文对象
//   - blkHash: 区块哈希
//
// 返回：
//   - *types.Block: 区块对象，包含区块头、交易列表等信息
//   - error: 如果查询失败则返回错误
func (p *Provider) GetBlockByHash(ctx context.Context, blkHash common.Hash) (*types.Block, error) {
	return p.ec.BlockByHash(ctx, blkHash)
}

// GetBlockByNumber 根据区块号获取区块信息
// 通过区块号查询完整的区块信息，包括区块头和所有交易
// 参数说明：
//   - ctx: 上下文对象
//   - number: 区块号（nil 表示最新区块）
//
// 返回：
//   - *types.Block: 区块对象，包含区块头、交易列表等信息
//   - error: 如果查询失败则返回错误
func (p *Provider) GetBlockByNumber(ctx context.Context, number *big.Int) (*types.Block, error) {
	return p.ec.BlockByNumber(ctx, number)
}

// GetBlockNumber 获取最新区块号
// 返回当前链上的最新（最新打包的）区块号
// 参数说明：
//   - ctx: 上下文对象
//
// 返回：
//   - uint64: 最新区块号
//   - error: 如果查询失败则返回错误
func (p *Provider) GetBlockNumber(ctx context.Context) (uint64, error) {
	return p.ec.BlockNumber(ctx)
}

// GetSuggestGasPrice 获取建议的 Gas 价格
// 返回网络建议的 Gas 价格（单位为 Wei）
// 参数说明：
//   - ctx: 上下文对象
//
// 返回：
//   - *big.Int: 建议的 Gas 价格（单位为 Wei）
//   - error: 如果查询失败则返回错误
func (p *Provider) GetSuggestGasPrice(ctx context.Context) (*big.Int, error) {
	return p.ec.SuggestGasPrice(ctx)
}

// GetTransactionByHash 根据交易哈希获取交易信息
// 查询交易的详细信息，包括交易状态（是否已打包）
// 参数说明：
//   - ctx: 上下文对象
//   - txHash: 交易哈希
//
// 返回：
//   - tx: 交易对象，包含交易的所有字段
//   - isPending: 交易是否还在待处理状态（true 表示还在 mempool 中）
//   - error: 如果查询失败则返回错误
func (p *Provider) GetTransactionByHash(ctx context.Context, txHash common.Hash) (tx *types.Transaction, isPending bool, err error) {
	return p.ec.TransactionByHash(ctx, txHash)
}

// GetTransactionReceipt 根据交易哈希获取交易收据
// 交易收据包含交易执行结果、gas 使用情况、日志等信息
// 注意：只有已打包的交易才有收据，待处理的交易无法获取收据
// 参数说明：
//   - ctx: 上下文对象
//   - txHash: 交易哈希
//
// 返回：
//   - *types.Receipt: 交易收据，包含交易状态、gas 使用等信息
//   - error: 如果查询失败则返回错误（交易未打包时会返回错误）
func (p *Provider) GetTransactionReceipt(ctx context.Context, txHash common.Hash) (*types.Receipt, error) {
	return p.ec.TransactionReceipt(ctx, txHash)
}

// GetContractBytecode 根据合约地址获取字节码
// 返回合约在链上的字节码（十六进制字符串，带 0x 前缀）
// 参数说明：
//   - ctx: 上下文对象
//   - address: 合约地址
//
// 返回：
//   - string: 合约字节码（十六进制字符串，如 "0x608060405234801561001057600080fd5b50..."）
//   - error: 如果查询失败则返回错误
//
// 注意：如果地址不是合约（普通地址），返回的字节码为空字符串
func (p *Provider) GetContractBytecode(ctx context.Context, address common.Address) (string, error) {
	bytecode, err := p.ec.CodeAt(ctx, address, nil) // nil is the latest block
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(bytecode), nil
}

// IsContractAddress 检查地址是否为合约地址
// 通过检查地址的代码长度来判断是否为合约（合约代码长度 > 0）
// 参数说明：
//   - ctx: 上下文对象
//   - address: 要检查的地址
//
// 返回：
//   - bool: true 表示是合约地址，false 表示是普通地址（EOA）
//   - error: 如果查询失败则返回错误
func (p *Provider) IsContractAddress(ctx context.Context, address common.Address) (bool, error) {
	//获取一个代币智能合约的字节码并检查其长度以验证它是一个智能合约
	if bytecode, err := p.GetContractBytecode(ctx, address); err == nil {
		return len(bytecode) > 0, nil
	} else {
		return false, err
	}
}

// EstimateGas 估算交易所需的 Gas 数量
// 通过模拟交易执行来估算 gas 消耗，这对于确定交易的 gasLimit 很有用
// 参数说明：
//   - ctx: 上下文对象
//   - from: 发送地址
//   - to: 接收地址（合约地址或普通地址）
//   - nonce: 交易 nonce
//   - gasPrice: Gas 价格（nil 表示不设置）
//   - value: 转账金额（nil 表示不转账）
//   - data: 交易数据（合约调用数据或 nil）
//
// 返回：
//   - uint64: 估算的 Gas 数量
//   - error: 如果估算失败则返回错误（如合约执行失败、余额不足等）
func (p *Provider) EstimateGas(ctx context.Context, from, to common.Address, nonce uint64, gasPrice, value *big.Int, data []byte) (uint64, error) {
	return p.ec.EstimateGas(ctx, ethereum.CallMsg{
		From:       from,
		To:         &to,
		GasPrice:   gasPrice,
		Value:      value,
		Data:       data,
		Gas:        0,
		GasFeeCap:  nil,
		GasTipCap:  nil,
		AccessList: nil,
	})
}

// GetFromAddress 从交易中提取发送地址
// 通过解析交易签名来获取发送者地址（交易签名者的地址）
// 参数说明：
//   - tx: 交易对象
//
// 返回：
//   - common.Address: 发送地址
//   - error: 如果提取失败则返回错误（如签名无效）
func (p *Provider) GetFromAddress(tx *types.Transaction) (common.Address, error) {
	return types.Sender(types.NewLondonSigner(tx.ChainId()), tx)
}

// FilterLogs 查询事件日志
// 用于查询指定区块范围内的事件日志，支持按合约地址、事件签名和 indexed 参数进行过滤
// 参数说明：
//   - ctx: 上下文对象
//   - contractAddress: 合约地址（nil 表示查询所有合约）
//   - eventTopic: 事件签名 topic（如 GetEventTopic("Transfer(address,address,uint256)")）
//   - fromBlock: 起始区块号（nil 表示从最新区块开始）
//   - toBlock: 结束区块号（nil 表示到最新区块）
//   - indexedTopics: 可选的 indexed 参数过滤（nil 表示不过滤，每个元素对应一个 indexed 参数）
//
// 返回：
//   - []types.Log: 事件日志列表，用户需要自行解析 Data 和 Topics
//   - error: 如果查询失败则返回错误
//
// 使用示例：
//   - 查询单个合约的事件：FilterLogs(ctx, &contractAddr, topicHash, fromBlock, toBlock, nil)
//   - 查询所有合约的事件：FilterLogs(ctx, nil, topicHash, fromBlock, toBlock, nil)
//   - 带 indexed 参数过滤：FilterLogs(ctx, &contractAddr, topicHash, fromBlock, toBlock, []common.Hash{fromAddr.Hash(), toAddr.Hash()})
func (p *Provider) FilterLogs(ctx context.Context, contractAddress *common.Address, eventTopic common.Hash, fromBlock, toBlock *big.Int, indexedTopics []common.Hash) ([]types.Log, error) {
	query := ethereum.FilterQuery{
		FromBlock: fromBlock,
		ToBlock:   toBlock,
		Topics: [][]common.Hash{
			{eventTopic}, // 第一个 topic 是事件签名
		},
	}

	// 如果指定了合约地址，则添加到查询条件
	if contractAddress != nil {
		query.Addresses = []common.Address{*contractAddress}
	}

	// 如果有 indexed 参数过滤，添加到 Topics
	// Topics 的结构：Topics[0] 是事件签名，Topics[1] 是第一个 indexed 参数，Topics[2] 是第二个 indexed 参数，以此类推
	if len(indexedTopics) > 0 {
		for _, topic := range indexedTopics {
			query.Topics = append(query.Topics, []common.Hash{topic})
		}
	}

	return p.ec.FilterLogs(ctx, query)
}
