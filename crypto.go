package main

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"time"
)

const (
	BLOCK_SIZE int = 32
	IV_SIZE    int = aes.BlockSize
)

type CryptoKey struct {
	TTLConfig string `yaml:"ttl"`
	KeyConfig string `yaml:"key"`
	Until     time.Time
	Key       []byte
}

type Crypto struct {
	Keys   []CryptoKey
	Active bool
}

func (c *Crypto) EncrichKeyValues(ckey CryptoKey, key, datetime string) CryptoKey {
	var err error
	ckey.Until, err = time.Parse("2016-01-14 01:18:18.032507415 +0600 KGT", datetime)
	ckey.Key = []byte(key)
	if err != nil {
		log.Printf("[ERROR] Failed to parse provided TTL value: %v", err)
		return ckey
	}
	return ckey
}

func (c *Crypto) ReadKeysFromFile(filepath string) {
	yamlFile, err := ioutil.ReadFile(filepath)
	if err != nil {
		log.Printf("[ERROR] Failed to read key file yaml: %v", err)
		c.Active = false
		return
	}
	var ckey CryptoKey
	err = yaml.Unmarshal(yamlFile, ckey)
	if err != nil {
		log.Printf("[ERROR] Failed to parse config: %v", err)
		c.Active = false
		return
	}
	ckey = c.EncrichKeyValues(ckey, ckey.KeyConfig, ckey.TTLConfig)
	c.Active = true
	c.Keys = append(c.Keys, ckey)
}

func (c *Crypto) Encrypt(key []byte, data []byte) ([]byte, error) {
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

func (c *Crypto) Decrypt(key []byte, data []byte) ([]byte, error) {
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
