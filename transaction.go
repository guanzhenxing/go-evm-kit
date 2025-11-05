package etherkit

import (
	"encoding/hex"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/rlp"
)

//############ Transaction ############

// NewTx 创建新的交易对象
// 构建一个以太坊交易，使用传统交易类型（Legacy Transaction）
// 参数说明：
//   - to: 接收地址（合约地址或普通地址，nil 表示合约部署）
//   - nonce: 交易 nonce（用于防止重放攻击）
//   - gasLimit: Gas 限制（交易最多消耗的 gas）
//   - gasPrice: Gas 价格（单位为 Wei）
//   - value: 转账金额（单位为 Wei，nil 表示不转账）
//   - data: 交易数据（合约调用数据或 nil）
//
// 返回：
//   - *types.Transaction: 交易对象（未签名）
//   - error: 如果创建失败则返回错误
//
// 注意：
//   - 此方法创建的是传统交易类型，不支持 EIP-1559 的动态 gas 费用
//   - 交易需要签名后才能发送
func NewTx(to common.Address, nonce, gasLimit uint64, gasPrice, value *big.Int, data []byte) (*types.Transaction, error) {
	return types.NewTx(&types.LegacyTx{
		Nonce:    nonce,
		To:       &to,
		Value:    value,
		Gas:      gasLimit,
		GasPrice: gasPrice,
		Data:     data,
	}), nil
}

// NewTxWithHexData 基于十六进制数据创建交易对象
// 与 NewTx 类似，但接受十六进制字符串作为交易数据
// 参数说明：
//   - to: 接收地址
//   - nonce: 交易 nonce
//   - gasLimit: Gas 限制
//   - gasPrice: Gas 价格（单位为 Wei）
//   - value: 转账金额（单位为 Wei，nil 表示不转账）
//   - hexData: 十六进制交易数据（带或不带 0x 前缀）
//
// 返回：
//   - *types.Transaction: 交易对象（未签名）
//   - error: 如果创建失败则返回错误（如十六进制格式无效）
func NewTxWithHexData(to common.Address, nonce, gasLimit uint64, gasPrice, value *big.Int, hexData string) (*types.Transaction, error) {
	data, err := hex.DecodeString(hexData)
	if err != nil {
		return nil, err
	}
	return NewTx(to, nonce, gasLimit, gasPrice, value, data)
}

// DecodeRawTxHex 解析原始交易十六进制字符串
// 将 RLP 编码的原始交易数据解码为交易对象
// 参数说明：
//   - rawTx: 原始交易的十六进制字符串（RLP 编码，带或不带 0x 前缀）
//
// 返回：
//   - *types.Transaction: 解析后的交易对象
//   - error: 如果解析失败则返回错误（如格式无效、RLP 编码错误等）
//
// 使用场景：
//   - 从其他系统接收原始交易数据并解析
//   - 从链上获取原始交易并重新构建交易对象
func DecodeRawTxHex(rawTx string) (*types.Transaction, error) {

	tx := new(types.Transaction)
	rawTxBytes, err := hex.DecodeString(rawTx)
	if err != nil {
		return nil, err
	}
	err = rlp.DecodeBytes(rawTxBytes, &tx)
	if err != nil {
		return nil, err
	}

	return tx, nil
}

// GetMaxUint256 获取 uint256 的最大值
// 返回 2^256 - 1，这是 Solidity 中 uint256 类型能表示的最大值
// 常用于 ERC20 代币的 approve 操作，表示无限授权
//
// 返回：
//   - *big.Int: uint256 的最大值（2^256 - 1）
//
// 示例：
//   - maxUint := GetMaxUint256()
//   - txHash, err := kit.InvokeContract(ctx, tokenAddr, abi, "approve", 0, 0, nil, nil, spenderAddr, maxUint)
func GetMaxUint256() *big.Int {
	return new(big.Int).Sub(new(big.Int).Exp(big.NewInt(2), big.NewInt(256), nil), big.NewInt(1))
}
