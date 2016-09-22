package ptp

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"strconv"
	"time"
)

type CryptoKey struct {
	TTLConfig string `yaml:"ttl"`
	KeyConfig string `yaml:"key"`
	Until     time.Time
	Key       []byte
}

type Crypto struct {
	Keys      []CryptoKey
	ActiveKey CryptoKey
	Active    bool
}

func (c Crypto) EnrichKeyValues(ckey CryptoKey, key, datetime string) CryptoKey {
	var err error
	i, err := strconv.ParseInt(datetime, 10, 64)
	ckey.Until = time.Now()
	// Default value is +1 hour
	ckey.Until = ckey.Until.Add(60 * time.Minute)
	if err != nil {
		Log(WARNING, "Failed to parse TTL. Falling back to default value of 1 hour")
	} else {
		ckey.Until = time.Unix(i, 0)
	}
	ckey.Key = []byte(key)
	if err != nil {
		Log(ERROR, "Failed to parse provided TTL value: %v", err)
		return ckey
	}
	return ckey
}

func (c Crypto) ReadKeysFromFile(filepath string) {
	yamlFile, err := ioutil.ReadFile(filepath)
	if err != nil {
		Log(ERROR, "Failed to read key file yaml: %v", err)
		c.Active = false
		return
	}
	var ckey CryptoKey
	err = yaml.Unmarshal(yamlFile, ckey)
	if err != nil {
		Log(ERROR, "Failed to parse config: %v", err)
		c.Active = false
		return
	}
	ckey = c.EnrichKeyValues(ckey, ckey.KeyConfig, ckey.TTLConfig)
	c.Active = true
	c.Keys = append(c.Keys, ckey)
}

func (c Crypto) Encrypt(key []byte, data []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	if len(data) != aes.BlockSize {
		padding := aes.BlockSize - len(data)%aes.BlockSize
		data = append(data, bytes.Repeat([]byte{byte(padding)}, padding)...)
	}

	encrypted_data := make([]byte, aes.BlockSize+len(data))
	iv := encrypted_data[:aes.BlockSize]
	if _, err = rand.Read(iv); err != nil {
		return nil, err
	}

	mode := cipher.NewCBCEncrypter(block, iv)
	mode.CryptBlocks(encrypted_data[aes.BlockSize:], data)

	return encrypted_data, nil
}

func (c Crypto) Decrypt(key []byte, data []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	encrypted_data := data[aes.BlockSize:]

	mode := cipher.NewCBCDecrypter(block, data[:aes.BlockSize])
	mode.CryptBlocks(encrypted_data, encrypted_data)

	return encrypted_data, nil
}
