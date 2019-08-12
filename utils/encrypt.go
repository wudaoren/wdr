package utils

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"errors"
)

//aes对称加密
func AesEncrypt(origData, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	blockSize := block.BlockSize()
	padding := blockSize - len(origData)%blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	origData = append(origData, padtext...)

	blockMode := cipher.NewCBCEncrypter(block, key[:blockSize])
	crypted := make([]byte, len(origData))
	blockMode.CryptBlocks(crypted, origData)
	return crypted, nil
}

//aes对称解密
func AesDecrypt(crypted, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	blockSize := block.BlockSize()
	blockMode := cipher.NewCBCDecrypter(block, key[:blockSize])
	if len(crypted)%blockSize != 0 {
		return nil, errors.New("parse error.")
	}
	origData := make([]byte, len(crypted))
	blockMode.CryptBlocks(origData, crypted)
	length := len(origData)
	unpadding := int(origData[length-1])
	origData = origData[:(length - unpadding)]
	return origData, nil
}
