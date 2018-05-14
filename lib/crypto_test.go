package ptp

import (
	"math/rand"
	"testing"
	"time"
	"reflect"
)

func TestEncrypt(t *testing.T) {
	crypto := new(Crypto)
	_, err := crypto.encrypt([]byte{}, []byte{})
	if err == nil {
		t.Errorf("Encrypt didn't return error on empty key")
	}
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
			crypto.encrypt(crypto.ActiveKey.Key, []byte(str))
		}
	}
}
/*
Generated TestCrypto_EnrichKeyValues
Generated TestCrypto_ReadKeysFromFile
Generated TestCrypto_encrypt
Generated TestCrypto_decrypt

package ptp

import (
	"crypto/rand"
	"math/rand"
	"reflect"
	"testing"
	"time"
)
*/

/*
func TestEncrypt(t *testing.T) {
	crypto := new(Crypto)
	_, err := crypto.encrypt([]byte{}, []byte{})
	if err == nil {
		t.Errorf("Encrypt didn't return error on empty key")
	}
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
			crypto.encrypt(crypto.ActiveKey.Key, []byte(str))
		}
	}
}
*/

func TestCrypto_EnrichKeyValues(t *testing.T) {
	type fields struct {
		Keys      []CryptoKey
		ActiveKey CryptoKey
		Active    bool
	}
	type args struct {
		ckey     CryptoKey
		key      string
		datetime string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   CryptoKey
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := Crypto{
				Keys:      tt.fields.Keys,
				ActiveKey: tt.fields.ActiveKey,
				Active:    tt.fields.Active,
			}
			if got := c.EnrichKeyValues(tt.args.ckey, tt.args.key, tt.args.datetime); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Crypto.EnrichKeyValues() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCrypto_ReadKeysFromFile(t *testing.T) {
	type fields struct {
		Keys      []CryptoKey
		ActiveKey CryptoKey
		Active    bool
	}
	type args struct {
		filepath string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := Crypto{
				Keys:      tt.fields.Keys,
				ActiveKey: tt.fields.ActiveKey,
				Active:    tt.fields.Active,
			}
			c.ReadKeysFromFile(tt.args.filepath)
		})
	}
}

func TestCrypto_encrypt(t *testing.T) {
	type fields struct {
		Keys      []CryptoKey
		ActiveKey CryptoKey
		Active    bool
	}
	type args struct {
		key  []byte
		data []byte
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []byte
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := Crypto{
				Keys:      tt.fields.Keys,
				ActiveKey: tt.fields.ActiveKey,
				Active:    tt.fields.Active,
			}
			got, err := c.encrypt(tt.args.key, tt.args.data)
			if (err != nil) != tt.wantErr {
				t.Errorf("Crypto.encrypt() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Crypto.encrypt() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCrypto_decrypt(t *testing.T) {
	type fields struct {
		Keys      []CryptoKey
		ActiveKey CryptoKey
		Active    bool
	}
	type args struct {
		key  []byte
		data []byte
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []byte
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := Crypto{
				Keys:      tt.fields.Keys,
				ActiveKey: tt.fields.ActiveKey,
				Active:    tt.fields.Active,
			}
			got, err := c.decrypt(tt.args.key, tt.args.data)
			if (err != nil) != tt.wantErr {
				t.Errorf("Crypto.decrypt() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Crypto.decrypt() = %v, want %v", got, tt.want)
			}
		})
	}
}
