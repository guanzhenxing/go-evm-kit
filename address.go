package etherkit

import (
	"encoding/hex"
	"regexp"

	"github.com/ethereum/go-ethereum/common"
	"golang.org/x/crypto/sha3"
)

//############ Address ############

// IsValidAddress 验证是否是有效的以太坊地址
// 验证地址格式是否正确（必须以 0x 开头，后跟 40 个十六进制字符）
// 参数说明：
//   - iAddress: 要验证的地址，可以是 string 或 common.Address 类型
//
// 返回：
//   - bool: true 表示地址格式有效，false 表示格式无效
//
// 示例：
//   - IsValidAddress("0x742d35Cc6634C0532925a3b844Bc9e7595f0bEb") // 返回 true
//   - IsValidAddress("0x742d35Cc6634C0532925a3b844Bc9e7595f0bEb") // 返回 false（长度不对）
func IsValidAddress(iAddress interface{}) bool {
	re := regexp.MustCompile("^0x[0-9a-fA-F]{40}$")
	switch v := iAddress.(type) {
	case string:
		return re.MatchString(v)
	case common.Address:
		return re.MatchString(v.Hex())
	default:
		return false
	}
}

// PublicKeyBytesToAddress 从公钥字节转换为以太坊地址
// 以太坊地址是从公钥派生出来的：对公钥进行 Keccak256 哈希，然后取后 20 字节
// 参数说明：
//   - publicKey: 公钥字节（通常包含 0x04 前缀，长度为 65 字节）
//
// 返回：
//   - common.Address: 派生出的以太坊地址
//
// 注意：
//   - 公钥字节的第一个字节（0x04）会被移除，然后对剩余 64 字节进行哈希
//   - 哈希结果的后 20 字节即为地址
func PublicKeyBytesToAddress(publicKey []byte) common.Address {
	var buf []byte

	hash := sha3.NewLegacyKeccak256()
	hash.Write(publicKey[1:]) // remove EC prefix 04
	buf = hash.Sum(nil)
	address := buf[12:]

	return common.HexToAddress(hex.EncodeToString(address))
}
