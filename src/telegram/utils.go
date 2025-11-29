package telegram

import (
	"crypto/sha256"

	"github.com/btcsuite/btcutil/base58"
)

// isValidTronAddress 校验 Tron Base58Check 地址是否合法
func isValidTronAddress(addr string) bool {
	// 基本过滤
	if len(addr) < 26 || len(addr) > 35 || addr[0] != 'T' {
		return false
	}

	decoded := base58.Decode(addr)
	if len(decoded) != 25 {
		return false
	}

	// TRON 主网地址必须以 0x41 开头
	if decoded[0] != 0x41 {
		return false
	}

	// Base58Check 校验
	payload := decoded[:21]  // 前 21 字节
	checksum := decoded[21:] // 后 4 字节

	hash := sha256.Sum256(payload)
	hash2 := sha256.Sum256(hash[:])

	return string(checksum) == string(hash2[:4])
}
