package common

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/x509"
	"encoding/hex"
	"encoding/pem"
	"errors"
	"hash/crc32"
	"hash/crc64"
	"io"
	"os"
	"strings"
)

// 公钥加密
func RsaEncryptPublic(origData []byte, publicKey []byte) ([]byte, error) {
	block, _ := pem.Decode(publicKey)
	if block == nil {
		return nil, errors.New("public key error")
	}
	pubInterface, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, err
	}
	pub := pubInterface.(*rsa.PublicKey)
	//return rsa.EncryptOAEP()
	return rsa.EncryptPKCS1v15(rand.Reader, pub, origData)
}

// 私钥解密
func RsaDecryptPrivate(cipherText []byte, privateKey []byte) ([]byte, error) {
	block, _ := pem.Decode(privateKey)
	if block == nil {
		return nil, errors.New("private key error!")
	}
	priv, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}
	//rsa.DecryptOAEP
	return rsa.DecryptPKCS1v15(rand.Reader, priv, cipherText)
}

// 文件校验
func HashFile(filePath string, mode string) (string, error) {
	var returnCRC32String string
	file, err := os.Open(filePath)
	if err != nil {
		return returnCRC32String, err
	}
	defer func() { _ = file.Close() }()
	// 关于CRC的表 https://mrwaggel.be/post/generate-crc32-hash-of-a-file-in-golang-turorial/
	switch mode {
	case "crc32":
		tablePolynomial := crc32.MakeTable(crc32.IEEE)
		hash := crc32.New(tablePolynomial)
		if _, err := io.Copy(hash, file); err != nil {
			return returnCRC32String, err
		}
		hashInBytes := hash.Sum(nil)[:]
		returnCRC32String = hex.EncodeToString(hashInBytes)
	case "crc64":
		tablePolynomial := crc64.MakeTable(crc64.ECMA)
		hash := crc64.New(tablePolynomial)
		if _, err := io.Copy(hash, file); err != nil {
			return returnCRC32String, err
		}
		hashInBytes := hash.Sum(nil)[:]
		returnCRC32String = hex.EncodeToString(hashInBytes)
	case "sha1":
		hash := sha1.New()
		if _, err := io.Copy(hash, file); err != nil {
			return returnCRC32String, err
		}
		hashInBytes := hash.Sum(nil)[:]
		returnCRC32String = hex.EncodeToString(hashInBytes)
	case "sha256":
		hash := sha256.New()
		if _, err := io.Copy(hash, file); err != nil {
			return returnCRC32String, err
		}
		hashInBytes := hash.Sum(nil)[:]
		returnCRC32String = hex.EncodeToString(hashInBytes)
	default:
		return returnCRC32String, errors.New("hash mode error")
	}

	return strings.ToUpper(returnCRC32String), nil

}
