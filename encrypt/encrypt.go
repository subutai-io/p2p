package main

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"fmt"
)

/////////////////////////////////////////////////////

const (
	BLOCK_SIZE int = 32
	IV_SIZE    int = aes.BlockSize
)

/////////////////////////////////////////////////////

func main() {
	plaintext := "123456789012345678901234567890123456789"
	fmt.Printf("%s\n", plaintext)

	key := make([]byte, BLOCK_SIZE)
	for i := 0; i < BLOCK_SIZE; i++ {
		key[i] = byte(i)
	}
	enc_data, err := Encrypt(key, []byte(plaintext))
	if err != nil {
		panic(err)
	}
	fmt.Printf("encrypted data : %s\n", enc_data)

	dec_data, err := Decrypt(key, enc_data)
	if err != nil {
		panic(err)
	}
	fmt.Printf("decrypted data : %s\n", dec_data)
}

/////////////////////////////////////////////////////

func Encrypt(key []byte, data []byte) ([]byte, error) {
	cb, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	data_len := len(data)
	result_data_len := (data_len + BLOCK_SIZE - 1) & (^(BLOCK_SIZE - 1))
	encrypted_data := make([]byte, IV_SIZE+result_data_len)
	// The IV needs to be unique, but not secured.
	iv := make([]byte, aes.BlockSize)
	_, err = rand.Read(iv)
	if err != nil {
		return nil, err
	}

	copy(encrypted_data[:IV_SIZE], iv)
	count := result_data_len / BLOCK_SIZE
	for i := 0; i < count-1; i++ {
		mode := cipher.NewCBCEncrypter(cb, iv)
		mode.CryptBlocks(encrypted_data[i*BLOCK_SIZE+IV_SIZE:], data[i*BLOCK_SIZE:(i+1)*BLOCK_SIZE])
	}

	tmp_arr := make([]byte, BLOCK_SIZE)
	copy(tmp_arr, data[(count-1)*BLOCK_SIZE:])
	mode := cipher.NewCBCEncrypter(cb, iv)
	mode.CryptBlocks(encrypted_data[(count-1)*BLOCK_SIZE+IV_SIZE:], tmp_arr)
	return encrypted_data, nil
}

/////////////////////////////////////////////////////

func Decrypt(key []byte, data []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	iv := data[:IV_SIZE]
	data_len := len(data) - IV_SIZE
	decrypted_data := make([]byte, data_len)
	count := data_len / BLOCK_SIZE
	for i := 0; i < count; i++ {
		mode := cipher.NewCBCDecrypter(block, iv)
		mode.CryptBlocks(decrypted_data[i*BLOCK_SIZE:], data[i*BLOCK_SIZE+IV_SIZE:])
	}
	return decrypted_data, nil
}
