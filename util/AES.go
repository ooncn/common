package util

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"errors"
	"time"
)

type AES struct {
}

func (a AES) PKCS7Padding(cipherText []byte, blockSize int) []byte {
	padding := blockSize - len(cipherText)%blockSize
	padText := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(cipherText, padText...)
}

func (a AES) PKCS7UnPadding(origData []byte) []byte {
	length := len(origData)
	unPadding := int(origData[length-1])
	return origData[:(length - unPadding)]
}

//加密
func (a AES) AesEncrypt(origData, key, IvParameterSpec []byte) (string, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}
	blockSize := block.BlockSize()
	origData = a.PKCS7Padding(origData, blockSize)
	blockMode := cipher.NewCBCEncrypter(block, IvParameterSpec[:blockSize])
	cryptEd := make([]byte, len(origData))
	blockMode.CryptBlocks(cryptEd, origData)
	return base64.StdEncoding.EncodeToString(cryptEd), nil
}

//解密
func (a AES) AesDecrypt(cryptEd string, key, IvParameterSpec []byte) (b []byte, e error) {
	if cryptEd == "" {
		return nil, errors.New("加密不能为空")
	}
	crypt, err := base64.StdEncoding.DecodeString(cryptEd)
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	blockSize := block.BlockSize()
	blockMode := cipher.NewCBCDecrypter(block, IvParameterSpec[:blockSize])
	b = make([]byte, len(crypt))
	blockMode.CryptBlocks(b, crypt)
	b = a.PKCS7UnPadding(b)
	return b, nil
}
func (a AES) AesDecryptByte(cryptEd string, key, IvParameterSpec []byte) (b []byte, e error) {
	if cryptEd == "" {
		return nil, errors.New("加密不能为空")
	}
	crypt, err := base64.StdEncoding.DecodeString(cryptEd)
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	blockSize := block.BlockSize()
	blockMode := cipher.NewCBCDecrypter(block, IvParameterSpec[:blockSize])
	b = make([]byte, len(crypt))
	blockMode.CryptBlocks(b, crypt)
	b = a.PKCS7UnPadding(b)
	return b, nil
}

func WsEncrypt(s string) string {
	return AesCBCEncrypt(s, "qweryu"+time.Now().Format("2006-01-02"))
}
func WsDecrypt(s string) string {
	return AesCBCDecrypt(s, "qweryu"+time.Now().Format("2006-01-02"))
}
