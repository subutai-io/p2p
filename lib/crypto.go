package ptp

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"fmt"
	"io/ioutil"
	"strconv"
	"time"

	"gopkg.in/yaml.v2"
)

// CryptoKey represents a key and it's expiration date
type CryptoKey struct {
	TTLConfig string `yaml:"ttl"`
	KeyConfig string `yaml:"key"`
	Until     time.Time
	Key       []byte
}

// Crypto is a object used by crypto subsystem
type Crypto struct {
	Keys      []CryptoKey
	ActiveKey CryptoKey
	Active    bool
}

// EnrichKeyValues update information about current and feature keys
func (c Crypto) EnrichKeyValues(ckey CryptoKey, key, datetime string) CryptoKey {
	var err error
	i, err := strconv.ParseInt(datetime, 10, 64)
	ckey.Until = time.Now()
	// Default value is +1 hour
	ckey.Until = ckey.Until.Add(60 * time.Minute)
	if err != nil {
		Warning("Failed to parse TTL. Falling back to default value of 1 hour")
	} else {
		ckey.Until = time.Unix(i, 0)
	}
	ckey.Key = []byte(key)
	if err != nil {
		Error("Failed to parse provided TTL value: %v", err)
		return ckey
	}
	return ckey
}

// ReadKeysFromFile read a file stored in a file system and extracts keys to be used
func (c Crypto) ReadKeysFromFile(filepath string) {
	yamlFile, err := ioutil.ReadFile(filepath)
	if err != nil {
		Error("Failed to read key file yaml: %v", err)
		c.Active = false
		return
	}
	var ckey CryptoKey
	err = yaml.Unmarshal(yamlFile, ckey)
	if err != nil {
		Error("Failed to parse config: %v", err)
		c.Active = false
		return
	}
	ckey = c.EnrichKeyValues(ckey, ckey.KeyConfig, ckey.TTLConfig)
	c.Active = true
	c.Keys = append(c.Keys, ckey)
}

// Encrypt encrypts data
func (c Crypto) encrypt(key []byte, data []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	if len(data) != aes.BlockSize {
		padding := aes.BlockSize - len(data)%aes.BlockSize
		data = append(data, bytes.Repeat([]byte{byte(padding)}, padding)...)
	}

	encData := make([]byte, aes.BlockSize+len(data))
	iv := encData[:aes.BlockSize]
	if _, err = rand.Read(iv); err != nil {
		return nil, err
	}

	mode := cipher.NewCBCEncrypter(block, iv)
	mode.CryptBlocks(encData[aes.BlockSize:], data)

	return encData, nil
}

// Decrypt decrypts data
func (c Crypto) decrypt(key []byte, data []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	encData := data[aes.BlockSize:]
	if len(data)%aes.BlockSize != 0 {
		return nil, fmt.Errorf("Input not full blocks: %s", string(data))
	}
	mode := cipher.NewCBCDecrypter(block, data[:aes.BlockSize])
	mode.CryptBlocks(encData, encData)

	return encData, nil
}
