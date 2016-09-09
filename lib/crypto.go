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
		Log(ERROR, "Failed to parse TTL. Falling back to default value of 1 hour")
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
	cb, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	if len(data) != IV_SIZE {
		padding := IV_SIZE - len(data)%IV_SIZE
		data = append(data, bytes.Repeat([]byte{byte(padding)}, padding)...)
	}

	encrypted_data := make([]byte, IV_SIZE+len(data))
	iv := encrypted_data[:IV_SIZE]
	if _, err = rand.Read(iv); err != nil {
		return nil, err
	}

	mode := cipher.NewCBCEncrypter(cb, iv)
	mode.CryptBlocks(encrypted_data[IV_SIZE:], data)

	return encrypted_data, nil
}

func (c Crypto) Decrypt(key []byte, data []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	encrypted_data := data[IV_SIZE:]

	mode := cipher.NewCBCDecrypter(block, data[:IV_SIZE])
	mode.CryptBlocks(encrypted_data, encrypted_data)

	return encrypted_data, nil
}
