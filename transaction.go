package etherkit

import (
	"encoding/hex"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/rlp"
)

//############ Transaction ############

// NewTx 新建一个tx
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

// NewTxWithHexData 基于hexData构建一个Tx
func NewTxWithHexData(to common.Address, nonce, gasLimit uint64, gasPrice, value *big.Int, hexData string) (*types.Transaction, error) {
	data, err := hex.DecodeString(hexData)
	if err != nil {
		return nil, err
	}
	return NewTx(to, nonce, gasLimit, gasPrice, value, data)
}

// DecodeRawTxHex 解析rawTx
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

// GetMaxUint256 获得合约中MaxUint256
func GetMaxUint256() *big.Int {
	return new(big.Int).Sub(new(big.Int).Exp(big.NewInt(2), big.NewInt(256), nil), big.NewInt(1))
}
