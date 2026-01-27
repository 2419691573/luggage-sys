package utils

import (
	"math/rand"
	"time"
)

const digitCharset = "0123456789"

// GenerateRetrievalCode 生成6位随机数字取件码
func GenerateRetrievalCode() string {
	// 使用独立的随机数生成器，避免并发问题
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	b := make([]byte, 6)
	for i := range b {
		b[i] = digitCharset[r.Intn(len(digitCharset))]
	}
	return string(b)
}
