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

// EtherProvider 需要通过链上查询的，但是不需要账户的
type EtherProvider interface {
	GetEthClient() *ethclient.Client
	GetRpcClient() *rpc.Client
	Close()
	GetNetworkID(ctx context.Context) (*big.Int, error)
	GetChainID(ctx context.Context) (*big.Int, error)
	GetBlockByHash(ctx context.Context, hash common.Hash) (*types.Block, error)
	GetBlockByNumber(ctx context.Context, number *big.Int) (*types.Block, error)
	GetBlockNumber(ctx context.Context) (uint64, error)
	GetSuggestGasPrice(ctx context.Context) (*big.Int, error)
	GetTransactionByHash(ctx context.Context, hash common.Hash) (tx *types.Transaction, isPending bool, err error)
	GetTransactionReceipt(ctx context.Context, txHash common.Hash) (*types.Receipt, error)
	GetContractBytecode(ctx context.Context, address common.Address) (string, error)
	IsContractAddress(ctx context.Context, address common.Address) (bool, error)
	EstimateGas(ctx context.Context, from, to common.Address, nonce uint64, gasPrice, value *big.Int, data []byte) (uint64, error)
	GetFromAddress(tx *types.Transaction) (common.Address, error)
	FilterLogs(ctx context.Context, contractAddress *common.Address, eventTopic common.Hash, fromBlock, toBlock *big.Int, indexedTopics []common.Hash) ([]types.Log, error)
}

type Provider struct {
	rc      *rpc.Client
	ec      *ethclient.Client
	chainId *big.Int
}

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

func NewProviderWithChainId(rawUrl string, chainId int64) (*Provider, error) {

	p, err := NewProvider(rawUrl)
	if err != nil {
		return nil, err
	}
	p.chainId = big.NewInt(chainId)

	return p, nil
}

// GetEthClient 获得ethClient客户端
func (p *Provider) GetEthClient() *ethclient.Client {
	return p.ec
}

// GetRpcClient 获得rpcClient客户端
func (p *Provider) GetRpcClient() *rpc.Client {
	return p.rc
}

// Close 关闭ethClient客户端和rpcClient客户端
func (p *Provider) Close() {
	p.ec.Close()
	p.rc.Close()
}

// GetNetworkID 获得NetworkId
func (p *Provider) GetNetworkID(ctx context.Context) (*big.Int, error) {
	return p.ec.NetworkID(ctx)
}

// GetChainID 获得ChainId
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

// GetBlockByHash 根据区块Hash获得区块信息
func (p *Provider) GetBlockByHash(ctx context.Context, blkHash common.Hash) (*types.Block, error) {
	return p.ec.BlockByHash(ctx, blkHash)
}

// GetBlockByNumber 根据区块号获得区块信息
func (p *Provider) GetBlockByNumber(ctx context.Context, number *big.Int) (*types.Block, error) {
	return p.ec.BlockByNumber(ctx, number)
}

// GetBlockNumber 获得最新区块
func (p *Provider) GetBlockNumber(ctx context.Context) (uint64, error) {
	return p.ec.BlockNumber(ctx)
}

// GetSuggestGasPrice 获得建议的Gas
func (p *Provider) GetSuggestGasPrice(ctx context.Context) (*big.Int, error) {
	return p.ec.SuggestGasPrice(ctx)
}

// GetTransactionByHash 根据txHash获得交易信息
func (p *Provider) GetTransactionByHash(ctx context.Context, txHash common.Hash) (tx *types.Transaction, isPending bool, err error) {
	return p.ec.TransactionByHash(ctx, txHash)
}

// GetTransactionReceipt 根据txHash获得交易Receipt
func (p *Provider) GetTransactionReceipt(ctx context.Context, txHash common.Hash) (*types.Receipt, error) {
	return p.ec.TransactionReceipt(ctx, txHash)
}

// GetContractBytecode 根据合约地址获得bytecode
func (p *Provider) GetContractBytecode(ctx context.Context, address common.Address) (string, error) {
	bytecode, err := p.ec.CodeAt(ctx, address, nil) // nil is the latest block
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(bytecode), nil
}

// IsContractAddress 是否是合约地址。
func (p *Provider) IsContractAddress(ctx context.Context, address common.Address) (bool, error) {
	//获取一个代币智能合约的字节码并检查其长度以验证它是一个智能合约
	if bytecode, err := p.GetContractBytecode(ctx, address); err == nil {
		return len(bytecode) > 0, nil
	} else {
		return false, err
	}
}

// EstimateGas 预估手续费
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

// GetFromAddress 获得交易的fromAddress
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
