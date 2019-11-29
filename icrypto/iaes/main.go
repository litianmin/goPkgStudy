package main

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"io"
)

const (
	// PrivKey 私有键
	PrivKey = "kkkkkkkkkkkkkkkkkkkkkkkkkkkkkkkk"
)

// Padding 对明文进行填充
func Padding(plainText []byte, blockSize int) []byte {
	// 计算要填充的长度，如果长度刚好是16的倍数，那么就 填充 16 个 16
	// 去除填充的数据的时候直接切割 16 个byte就行了
	n := blockSize - len(plainText)%blockSize
	//对原来的明文填充n个n
	temp := bytes.Repeat([]byte{byte(n)}, n)
	plainText = append(plainText, temp...)
	return plainText
}

// UnPadding 对密文删除填充
func UnPadding(cipherText []byte) []byte {
	//取出密文最后一个字节end
	end := cipherText[len(cipherText)-1]
	//删除填充
	cipherText = cipherText[:len(cipherText)-int(end)]
	return cipherText
}

// AesCbcEncrypt AEC加密（CBC模式，密码分组链接模式）
func AesCbcEncrypt(plainText []byte, key []byte) []byte {
	//指定加密算法，返回一个AES算法的Block接口对象
	block, err := aes.NewCipher(key)
	if err != nil {
		panic(err)
	}

	//进行填充
	plainText = Padding(plainText, block.BlockSize())
	//指定初始向量vi,长度和block的块尺寸一致
	iv := []byte("12345678abcdefgh")
	//指定分组模式，返回一个BlockMode接口对象
	blockMode := cipher.NewCBCEncrypter(block, iv)
	//加密连续数据库
	cipherText := make([]byte, len(plainText))
	blockMode.CryptBlocks(cipherText, plainText)
	//返回密文
	return cipherText
}

// AesCbcDecrypt AEC解密（CBC模式，密码分组链接模式）
func AesCbcDecrypt(cipherText []byte, key []byte) []byte {
	//指定解密算法，返回一个AES算法的Block接口对象
	block, err := aes.NewCipher(key)
	if err != nil {
		panic(err)
	}
	//指定初始化向量IV,和加密的一致
	iv := []byte("12345678abcdefgh")
	//指定分组模式，返回一个BlockMode接口对象
	blockMode := cipher.NewCBCDecrypter(block, iv)
	//解密
	plainText := make([]byte, len(cipherText))
	blockMode.CryptBlocks(plainText, cipherText)

	//删除填充
	plainText = UnPadding(plainText)
	return plainText
}

// AesGCMEncrypt AES 的 GCM 加密
func AesGCMEncrypt() {
	key, _ := hex.DecodeString("6368616e676520746869732070617373776f726420746f206120736563726574")
	plaintext := []byte("exampleplaintext")

	block, err := aes.NewCipher(key)
	if err != nil {
		panic(err.Error())
	}

	// Never use more than 2^32 random nonces with a given key because of the risk of a repeat.
	nonce := make([]byte, 12)
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		panic(err.Error())
	}

	str := hex.EncodeToString(nonce)

	fmt.Println(string(str))

	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		panic(err.Error())
	}

	ciphertext := aesgcm.Seal(nil, nonce, plaintext, nil)
	fmt.Printf("%x\n", ciphertext)
}

// AesGCMDecrypt 解密
func AesGCMDecrypt() {
	key, _ := hex.DecodeString("6368616e676520746869732070617373776f726420746f206120736563726574")
	ciphertext, _ := hex.DecodeString("c3aaa29f002ca75870806e44086700f62ce4d43e902b3888e23ceff797a7a471")
	nonce, _ := hex.DecodeString("64a9433eae7ccceee2fc0eda")

	block, err := aes.NewCipher(key)
	if err != nil {
		panic(err.Error())
	}

	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		panic(err.Error())
	}

	plaintext, err := aesgcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		panic(err.Error())
	}

	fmt.Printf("%s\n", plaintext)
}

// 他们一直都是
func main() {
	// message := []byte("eeeareyou如果有这个呢ee")
	// //指定密钥h
	// key := []byte(PrivKey)
	// //加密
	// cipherText := AesCbcEncrypt(message, key)
	// fmt.Println("加密后为：", string(cipherText))
	// //解密
	// plainText := AesCbcDecrypt(cipherText, key)
	// fmt.Println("解密后为：", string(plainText))

	AesGCMEncrypt()

	// AesGCMDecrypt()
}
