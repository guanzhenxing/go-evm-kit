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
