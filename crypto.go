package etherkit

import (
	"crypto/ecdsa"
	"fmt"
	"strings"

	"github.com/ethereum/go-ethereum/accounts"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
	hdwallet "github.com/miguelmota/go-ethereum-hdwallet"
	"github.com/pkg/errors"
)

//############ Account ############

// GeneratePrivateKey 生成新的随机私钥
// 使用加密安全的随机数生成器创建 ECDSA 私钥
// 适用于创建新钱包的场景
//
// 返回：
//   - *ecdsa.PrivateKey: 新生成的 ECDSA 私钥
//   - error: 如果生成失败则返回错误
//
// 注意：
//   - 请妥善保存生成的私钥，丢失私钥将无法恢复钱包
//   - 生成的私钥是随机的，每次调用都会创建新的私钥
func GeneratePrivateKey() (*ecdsa.PrivateKey, error) {
	return crypto.GenerateKey()
}

// GetHexPrivateKey 获取私钥的十六进制字符串表示
// 将私钥转换为十六进制字符串，不包含 0x 前缀
// 参数说明：
//   - privateKey: ECDSA 私钥对象
//
// 返回：
//   - string: 私钥的十六进制字符串（不包含 0x 前缀，如 "abc123..."）
//
// 示例：
//   - hexPk := GetHexPrivateKey(privateKey) // 返回 "abc123..."
func GetHexPrivateKey(privateKey *ecdsa.PrivateKey) string {
	return hexutil.Encode(crypto.FromECDSA(privateKey))[2:]
}

// PrivateKeyToAddress 从私钥派生以太坊地址
// 以太坊地址是从私钥对应的公钥派生出来的
// 参数说明：
//   - privateKey: ECDSA 私钥对象
//
// 返回：
//   - common.Address: 派生出的以太坊地址
//
// 注意：
//   - 同一个私钥总是对应同一个地址
//   - 地址是从公钥的 Keccak256 哈希的后 20 字节派生出来的
func PrivateKeyToAddress(privateKey *ecdsa.PrivateKey) common.Address {
	publicKey := privateKey.Public()
	publicKeyECDSA := publicKey.(*ecdsa.PublicKey)
	return crypto.PubkeyToAddress(*publicKeyECDSA)
}

// GetHexPublicKey 获取公钥的十六进制字符串表示
// 将私钥对应的公钥转换为十六进制字符串，不包含 0x 前缀和 0x04 前缀
// 参数说明：
//   - privateKey: ECDSA 私钥对象
//
// 返回：
//   - string: 公钥的十六进制字符串（不包含 0x 和 0x04 前缀）
//
// 注意：
//   - 公钥的完整格式通常包含 0x04 前缀（表示未压缩格式）
//   - 此方法返回的是去掉 0x04 前缀后的公钥
func GetHexPublicKey(privateKey *ecdsa.PrivateKey) string {
	publicKey := privateKey.Public()
	publicKeyECDSA := publicKey.(*ecdsa.PublicKey)

	publicKeyBytes := crypto.FromECDSAPub(publicKeyECDSA)
	return hexutil.Encode(publicKeyBytes)[4:]
}

// BuildPrivateKeyFromHex 从十六进制字符串构建私钥对象
// 将十六进制私钥字符串转换为 ECDSA 私钥对象
// 参数说明：
//   - privateKeyHex: 十六进制私钥字符串（带或不带 0x 前缀）
//
// 返回：
//   - *ecdsa.PrivateKey: ECDSA 私钥对象
//   - error: 如果格式无效则返回错误
//
// 示例：
//   - pk, err := BuildPrivateKeyFromHex("abc123...")      // 不带 0x 前缀
//   - pk, err := BuildPrivateKeyFromHex("0xabc123...")    // 带 0x 前缀
func BuildPrivateKeyFromHex(privateKeyHex string) (*ecdsa.PrivateKey, error) {
	privateKeyHex = strings.TrimPrefix(privateKeyHex, "0x")

	privateKey, err := crypto.HexToECDSA(privateKeyHex)
	if err != nil {
		return nil, err
	}

	return privateKey, nil
}

// BuildPrivateKeyFromMnemonic 从助记词构建私钥对象（默认账户）
// 使用 BIP-44 标准从助记词派生私钥，使用默认账户索引（0）
// 参数说明：
//   - mnemonic: BIP-39 助记词字符串（12 或 24 个单词）
//
// 返回：
//   - *ecdsa.PrivateKey: ECDSA 私钥对象
//   - error: 如果助记词无效或派生失败则返回错误
//
// 注意：
//   - 使用 BIP-44 标准路径：m/44'/60'/0'/0/0（以太坊主网，第一个账户）
//   - 如果需要其他账户，使用 BuildPrivateKeyFromMnemonicAndAccountId
func BuildPrivateKeyFromMnemonic(mnemonic string) (*ecdsa.PrivateKey, error) {
	return BuildPrivateKeyFromMnemonicAndAccountId(mnemonic, 0)
}

// BuildPrivateKeyFromMnemonicAndAccountId 从助记词和账户索引构建私钥对象
// 使用 BIP-44 标准从助记词派生指定账户的私钥
// 参数说明：
//   - mnemonic: BIP-39 助记词字符串（12 或 24 个单词）
//   - accountId: 账户索引（0 表示第一个账户，1 表示第二个账户，以此类推）
//
// 返回：
//   - *ecdsa.PrivateKey: ECDSA 私钥对象
//   - error: 如果助记词无效或派生失败则返回错误
//
// 注意：
//   - 使用 BIP-44 标准路径：m/44'/60'/0'/0/{accountId}（以太坊主网）
//   - 同一个助记词配合不同的 accountId 可以生成不同的私钥和地址
//   - 这是 HD 钱包的标准做法，允许从一个助记词管理多个账户
func BuildPrivateKeyFromMnemonicAndAccountId(mnemonic string, accountId uint32) (*ecdsa.PrivateKey, error) {
	wallet, err := hdwallet.NewFromMnemonic(mnemonic)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create HD wallet from mnemonic")
	}
	path, err := accounts.ParseDerivationPath(fmt.Sprintf("m/44'/60'/0'/0/%d", accountId))
	if err != nil {
		return nil, errors.Wrap(err, "failed to parse derivation path")
	}
	account, err := wallet.Derive(path, true)
	if err != nil {
		return nil, errors.Wrap(err, "failed to derive account from HD wallet")
	}
	pk, err := wallet.PrivateKey(account)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get account's private key from HD wallet")
	}
	return pk, nil
}

// VerifySignature 验证签名是否由指定地址创建
// 验证给定的数据和签名是否由指定地址对应的私钥签名
// 参数说明：
//   - address: 用于签名的地址（十六进制字符串，带或不带 0x 前缀）
//   - data: 原始数据（字节）
//   - signature: 签名数据（65 字节，包含 r、s、v）
//
// 返回：
//   - bool: true 表示签名有效（由指定地址创建），false 表示签名无效
//
// 注意：
//   - 会对数据进行 Keccak256 哈希，然后验证签名
//   - 签名格式必须是 65 字节（r、s 各 32 字节，v 1 字节）
func VerifySignature(address string, data, signature []byte) bool {

	digestHash := crypto.Keccak256Hash(data)
	//returns the public key that created the given signature.
	sigPublicKeyECDSA, err := crypto.SigToPub(digestHash.Bytes(), signature)
	if err != nil {
		return false
	}

	sigAddress := crypto.PubkeyToAddress(*sigPublicKeyECDSA)
	return sigAddress.String() == address
}
