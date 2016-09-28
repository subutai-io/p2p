// +build ignore

package ptp

import (
	"math/rand"
	"testing"
	"time"
)

func TestEncrypt(t *testing.T) {
	crypto := new(Crypto)
	var key CryptoKey
	crypto.EnrichKeyValues(key, "keylessthan32", "1")
}

func RandomString(size int) string {
	var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")
	rand.Seed(time.Now().UnixNano())
	b := make([]rune, size)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

func BenchmarkEncrypt(b *testing.B) {
	var data []string
	for i := 1; i < 10; i++ {
		data = append(data, RandomString(i*10))
	}
	crypto := new(Crypto)
	var key CryptoKey
	crypto.EnrichKeyValues(key, "keylessthan32", "1")
	for i := 0; i < b.N; i++ {
		for _, str := range data {
			crypto.Encrypt(crypto.ActiveKey.Key, []byte(str))
		}
	}
}
