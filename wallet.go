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

// EtherWallet 钱包信息
type EtherWallet interface {
	GetEthProvider() EtherProvider
	GetClient() *ethclient.Client
	GetAddress() common.Address
	GetPrivateKey() *ecdsa.PrivateKey
	CloseWallet()
	GetNonce(ctx context.Context) (uint64, error)
	GetBalance(ctx context.Context) (*big.Int, error)
	NewTx(ctx context.Context, to common.Address, nonce, gasLimit uint64, gasPrice, value *big.Int, data []byte) (*types.Transaction, error)
	SendTx(ctx context.Context, to common.Address, nonce, gasLimit uint64, gasPrice, value *big.Int, data []byte) (common.Hash, error)
	NewTxWithHexInput(ctx context.Context, to common.Address, nonce, gasLimit uint64, gasPrice, value *big.Int, input string) (*types.Transaction, error)
	SendTxWithHexInput(ctx context.Context, to common.Address, nonce, gasLimit uint64, gasPrice, value *big.Int, input string) (common.Hash, error)
	BuildTxOpts(ctx context.Context, value, nonce, gasPrice *big.Int) (*bind.TransactOpts, error)
	SignTx(ctx context.Context, tx *types.Transaction) (*types.Transaction, error)
	SendSignedTx(ctx context.Context, signedTx *types.Transaction) (common.Hash, error)
	Signature(data []byte) ([]byte, error)
	CallContract(ctx context.Context, blockNumber *big.Int, from *common.Address, value *big.Int, contractAddress common.Address, contractAbi abi.ABI, functionName string, params ...interface{}) ([]interface{}, error)
}

type Wallet struct {
	privateKey *ecdsa.PrivateKey
	address    common.Address
	ep         EtherProvider
}

// NewWallet 新建一个Wallet
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

// NewWalletWithComponents creates a new Wallet with given private key and provider
func NewWalletWithComponents(privateKey *ecdsa.PrivateKey, ep EtherProvider) (*Wallet, error) {
	return &Wallet{
		privateKey: privateKey,
		address:    PrivateKeyToAddress(privateKey),
		ep:         ep,
	}, nil
}

// GetEthProvider 获得EthProvider
func (w *Wallet) GetEthProvider() EtherProvider {
	return w.ep
}

func (w *Wallet) GetClient() *ethclient.Client {
	return w.ep.GetEthClient()
}

// GetAddress 获得地址
func (w *Wallet) GetAddress() common.Address {
	return w.address
}

// GetPrivateKey 获得私钥
func (w *Wallet) GetPrivateKey() *ecdsa.PrivateKey {
	return w.privateKey
}

// CloseWallet 关闭Wallet
func (w *Wallet) CloseWallet() {
	w.ep.Close()
}

// GetNonce 获得nonce
func (w *Wallet) GetNonce(ctx context.Context) (uint64, error) {
	return w.GetClient().PendingNonceAt(ctx, w.GetAddress())
}

// GetBalance 获得本位币的约
func (w *Wallet) GetBalance(ctx context.Context) (*big.Int, error) {
	return w.GetClient().BalanceAt(ctx, w.GetAddress(), nil)
}

// NewTx 构建一笔交易。nonce传0表示字段计算；gasLimit传0表示字段计算；gasPrice穿nil或者big.NewInt(0)表示gasPrice自动计算。
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

// SendTx 发送交易。nonce传0表示字段计算；gasLimit传0表示字段计算；gasPrice穿nil或者big.NewInt(0)表示gasPrice自动计算。
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

// NewTxWithHexInput 构建一笔交易，使用0x开头的input。nonce传0表示字段计算；gasLimit传0表示字段计算；gasPrice穿nil或者big.NewInt(0)表示gasPrice自动计算。
func (w *Wallet) NewTxWithHexInput(ctx context.Context, to common.Address, nonce, gasLimit uint64, gasPrice, value *big.Int, input string) (*types.Transaction, error) {
	data, err := hexutil.Decode(input)
	if err != nil {
		return nil, err
	}
	return w.NewTx(ctx, to, nonce, gasLimit, gasPrice, value, data)
}

// SendTxWithHexInput 发送一笔交易，使用0x开头的input。nonce传0表示字段计算；gasLimit传0表示字段计算；gasPrice穿nil或者big.NewInt(0)表示gasPrice自动计算。
func (w *Wallet) SendTxWithHexInput(ctx context.Context, to common.Address, nonce, gasLimit uint64, gasPrice, value *big.Int, input string) (common.Hash, error) {
	data, err := hexutil.Decode(input)
	if err != nil {
		return [32]byte{}, err
	}
	return w.SendTx(ctx, to, nonce, gasLimit, gasPrice, value, data)
}

// BuildTxOpts 构建交易的选项
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

// SendSignedTx 发送签名后的Tx
func (w *Wallet) SendSignedTx(ctx context.Context, signedTx *types.Transaction) (common.Hash, error) {
	err := w.GetClient().SendTransaction(ctx, signedTx)
	if err != nil {
		return [32]byte{}, err
	}
	return signedTx.Hash(), nil
}

// Signature 生成一个签名
func (w *Wallet) Signature(data []byte) ([]byte, error) {
	hash := crypto.Keccak256Hash(data)
	return crypto.Sign(hash.Bytes(), w.privateKey)
}

// CallContract 调用合约的方法，无需创建交易
// 参数说明：
//   - blockNumber: 区块号（nil 表示最新区块）
//   - from: 调用者地址（nil 表示不设置）
//   - value: 模拟转账金额（nil 表示不转账）
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
