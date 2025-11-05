package etherkit

import (
	"math/big"

	"github.com/shopspring/decimal"
)

//############ Cast ############

// ToDecimal 将最小单位（如 Wei）转换为带小数位的数值
// 将链上的最小单位（如 Wei、Satoshi 等）转换为可读的小数形式
// 参数说明：
//   - iValue: 要转换的值，可以是 string（十进制字符串）或 *big.Int
//   - decimals: 小数位数（如以太币为 18，USDT 为 6）
//
// 返回：
//   - decimal.Decimal: 转换后的十进制数值（保留完整精度）
//
// 示例：
//   - ToDecimal("1000000000000000000", 18)  // 1 ETH = 1.0
//   - ToDecimal("1000000", 6)               // 1 USDT = 1.0
//   - balance := big.NewInt(500000000000000000) // 0.5 ETH
//   - ToDecimal(balance, 18)               // 0.5
func ToDecimal(iValue interface{}, decimals int) decimal.Decimal {
	value := new(big.Int)
	switch v := iValue.(type) {
	case string:
		value.SetString(v, 10)
	case *big.Int:
		value = v
	}

	mul := decimal.NewFromFloat(float64(10)).Pow(decimal.NewFromFloat(float64(decimals)))
	num, _ := decimal.NewFromString(value.String())
	result := num.Div(mul)

	return result
}

// ToWei 将带小数位的数值转换为最小单位（如 Wei）
// 将可读的小数形式转换为链上的最小单位（如 Wei、Satoshi 等）
// 参数说明：
//   - iAmount: 要转换的数值，可以是 string、float64、int64、int、decimal.Decimal 或 *decimal.Decimal
//   - decimals: 小数位数（如以太币为 18，USDT 为 6）
//
// 返回：
//   - *big.Int: 转换后的最小单位值（如 Wei）
//
// 示例：
//   - ToWei(1.5, 18)        // 1.5 ETH = 1500000000000000000 Wei
//   - ToWei("0.1", 18)      // 0.1 ETH = 100000000000000000 Wei
//   - ToWei(100, 6)         // 100 USDT = 100000000 (最小单位)
func ToWei(iAmount interface{}, decimals int) *big.Int {
	amount := decimal.NewFromFloat(0)
	switch v := iAmount.(type) {
	case string:
		amount, _ = decimal.NewFromString(v)
	case float64:
		amount = decimal.NewFromFloat(v)
	case int64:
		amount = decimal.NewFromFloat(float64(v))
	case int:
		amount = decimal.NewFromFloat(float64(v))
	case decimal.Decimal:
		amount = v
	case *decimal.Decimal:
		amount = *v
	}

	mul := decimal.NewFromFloat(float64(10)).Pow(decimal.NewFromFloat(float64(decimals)))
	result := amount.Mul(mul)

	wei := new(big.Int)
	wei.SetString(result.String(), 10)

	return wei
}
