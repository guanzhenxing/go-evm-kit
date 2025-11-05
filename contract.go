package etherkit

import (
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
)

//############ Contract ############

// GetABI 从 ABI JSON 字符串中解析 ABI 对象
// ABI（Application Binary Interface）定义了合约的接口，包括函数签名、事件等
// 参数说明：
//   - abiStr: ABI JSON 字符串（完整的合约 ABI 或单个函数的 ABI）
//
// 返回：
//   - abi.ABI: 解析后的 ABI 对象，可用于合约调用
//   - error: 如果 JSON 格式无效则返回错误
//
// 示例：
//   - abiStr := `[{"constant":false,"inputs":[{"name":"to","type":"address"},{"name":"amount","type":"uint256"}],"name":"transfer","outputs":[{"name":"","type":"bool"}],"type":"function"}]`
//   - abiObj, err := GetABI(abiStr)
func GetABI(abiStr string) (abi.ABI, error) {
	abiContract, err := abi.JSON(strings.NewReader(abiStr))
	return abiContract, err
}

// GetContractMethodId 获取合约方法的函数选择器（method ID）
// 函数选择器是函数签名的 Keccak256 哈希的前 4 字节，用于标识要调用的函数
// 参数说明：
//   - method: 函数签名字符串，格式为 "函数名(参数类型1,参数类型2,...)"
//     例如："transfer(address,uint256)"、"balanceOf(address)"
//
// 返回：
//   - string: 函数选择器（十六进制字符串，带 0x 前缀，如 "0xa9059cbb"）
//
// 示例：
//   - GetContractMethodId("transfer(address,uint256)") // 返回 "0xa9059cbb"
//   - GetContractMethodId("balanceOf(address)")         // 返回 "0x70a08231"
func GetContractMethodId(method string) string {
	methodId := hexutil.Encode(crypto.Keccak256([]byte(method))[:4])
	return methodId
}

// GetEventTopic 获取事件的 topic（事件签名的哈希）
// 事件的 topic 是事件签名的 Keccak256 哈希，用于在日志中识别事件
// 参数说明：
//   - event: 事件签名字符串，格式为 "事件名(参数类型1,参数类型2,...)"
//     例如："Transfer(address,address,uint256)"、"Approval(address,address,uint256)"
//
// 返回：
//   - string: 事件 topic（十六进制字符串，带 0x 前缀）
//
// 示例：
//   - GetEventTopic("Transfer(address,address,uint256)") // 返回 "0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef"
func GetEventTopic(event string) string {
	return crypto.Keccak256Hash([]byte(event)).String()
}

// BuildContractInputData 构建合约调用的输入数据
// 将函数名和参数打包成合约调用所需的字节数据
// 参数说明：
//   - contract: 合约 ABI 对象
//   - name: 函数名（如 "transfer", "balanceOf"）
//   - args: 函数参数（按函数定义顺序传入）
//
// 返回：
//   - []byte: 合约调用数据（包含函数选择器和编码后的参数）
//   - error: 如果打包失败则返回错误（如参数类型不匹配、参数数量不对等）
//
// 示例：
//   - data, err := BuildContractInputData(abi, "transfer", toAddress, amount)
//   - data, err := BuildContractInputData(abi, "balanceOf", userAddress)
func BuildContractInputData(contract abi.ABI, name string, args ...interface{}) ([]byte, error) {
	return contract.Pack(name, args...)
}
