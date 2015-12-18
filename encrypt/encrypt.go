package main

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"fmt"
)

func main() {
	plaintext := "Hello, world"
	fmt.Printf("%s\n", plaintext)

	key := make([]byte, aes.BlockSize)
	for i := 0; i < aes.BlockSize; i++ {
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
	result_data_len := (data_len + aes.BlockSize - 1) & (^(aes.BlockSize - 1))
	encrypted_data := make([]byte, aes.BlockSize+result_data_len)
	// The IV needs to be unique, but not secured.
	iv := make([]byte, aes.BlockSize)
	_, err = rand.Read(iv)
	if err != nil {
		return nil, err
	}

	copy(encrypted_data[:aes.BlockSize], iv)
	count := result_data_len / aes.BlockSize
	for i := 0; i < count-1; i++ {
		mode := cipher.NewCBCEncrypter(cb, iv)
		mode.CryptBlocks(encrypted_data[(i+1)*aes.BlockSize:], data[i*aes.BlockSize:(i+1)*aes.BlockSize])
	}

	tmp_arr := make([]byte, aes.BlockSize)
	copy(tmp_arr, data[(count-1)*aes.BlockSize:])
	mode := cipher.NewCBCEncrypter(cb, iv)
	mode.CryptBlocks(encrypted_data[(count)*aes.BlockSize:], tmp_arr)
	return encrypted_data, nil
}

/////////////////////////////////////////////////////

func Decrypt(key []byte, data []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	iv := data[:aes.BlockSize]
	data_len := len(data) - aes.BlockSize
	decrypted_data := make([]byte, data_len)
	count := data_len / aes.BlockSize
	for i := 0; i < count; i++ {
		mode := cipher.NewCBCDecrypter(block, iv)
		mode.CryptBlocks(decrypted_data[i*aes.BlockSize:], data[(i+1)*aes.BlockSize:])
	}
	return decrypted_data, nil
}
