package device

import (
	"crypto/md5"
	"crypto/rand"
	"fmt"
)

// Generate 生成设备指纹，格式: WIN + 32位MD5十六进制
// 如: WINe5efa67df895bd6cba85709f8df09dcf
func Generate() string {
	b := make([]byte, 16)
	_, _ = rand.Read(b)
	hash := md5.Sum(b)
	return fmt.Sprintf("WIN%x", hash)
}

// Validate 检查设备指纹格式是否合法
func Validate(id string) bool {
	if len(id) != 35 {
		return false
	}
	if id[:3] != "WIN" {
		return false
	}
	for _, c := range id[3:] {
		if !((c >= '0' && c <= '9') || (c >= 'a' && c <= 'f')) {
			return false
		}
	}
	return true
}
