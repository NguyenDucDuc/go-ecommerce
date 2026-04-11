package util

import (
	"crypto/rand"
	"fmt"
	"math/big"
)

// GenerateOTP tạo ra chuỗi số ngẫu nhiên gồm 6 chữ số
func GenerateOTP() (string, error) {
	otp := ""
	for i := 0; i < 6; i++ {
		// Tạo số ngẫu nhiên từ 0 đến 9
		num, err := rand.Int(rand.Reader, big.NewInt(10))
		if err != nil {
			return "", err
		}
		otp += fmt.Sprintf("%d", num)
	}
	return otp, nil
}