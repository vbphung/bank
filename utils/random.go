package utils

import (
	"fmt"
	"math/rand"
	"time"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

func RandomEmail() string {
	return fmt.Sprintf("%s@%s.com", pickStr(16), pickStr(7))
}

func RandomAmount(min, max int64) int64 {
	return rand.Int63n(max-min+1) + min
}

func pickStr(strLen int64) string {
	chars := []rune("abcdefghijklmnopqrstuvwxyz")

	str := make([]rune, strLen)
	for i := range str {
		str[i] = chars[rand.Intn(len(chars))]
	}

	return string(str)
}
