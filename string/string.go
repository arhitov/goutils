package string

import (
	cryptorand "crypto/rand"
	"math/big"
	"math/rand"
)

const (
	GenerateCodeCharsLower        = "abcdefghijklmnopqrstuvwxyz"
	GenerateCodeCharsUpper        = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	GenerateCodeCharsNumber       = "0123456789"
	GenerateCodeCharsLetterNumber = GenerateCodeCharsLower + GenerateCodeCharsUpper + GenerateCodeCharsNumber
	GenerateCodeCharsSpecial      = "~!@#$%&*();:"
	GenerateCodeCharsAll          = GenerateCodeCharsLower + GenerateCodeCharsUpper +
		GenerateCodeCharsNumber + GenerateCodeCharsSpecial
)

// GenerateCode - Строка состояния в виде случайной строки, состоящей из символов: "a-z, A-Z, 0-9, _, -", длиной не
// менее 32 символа. Передается на старте авторизации и должна возвращаться клиентскому приложению без изменения.
// Иначе ответ можно считать подмененным
func GenerateCode(length int, chars string) string {
	result := make([]byte, length)

	var idx int
	for i := 0; i < length; i++ {
		if num, err := cryptorand.Int(cryptorand.Reader, big.NewInt(int64(len(chars)))); err == nil {
			idx = int(num.Int64())
		} else {
			// Если произошла ошибка, генерируем без криптографической безопасности
			idx = rand.Intn(len(chars))
		}
		result[i] = chars[idx]
	}

	return string(result)
}

func StringPtrToString(str *string) string {
	if str == nil {
		return ""
	}
	return *str
}
