package main

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"fmt"
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

// AesCbcEncrypt AEC加密（CBC模式，密码分组链接模式， Cipher Block Chaining）
// step1：根据key（长度 16、24、32）值的长度 new 一个 cipher， 相当于开启一个加密，然后返回一个 block
// step2：把要加密的字符串填充成 key 值的整数倍, 填充的数值为填充的长度值，便于去填充
// step3：定义 初始化向量 ，创建 blockMode
// step4: 创建一个跟填充后的明文长度一样的 bytes 接收进行加密后的块信息
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

// 他们一直都是
func main() {

	// CBC 模式
	msg := []byte("这里是我的详细内容")
	cipherText := AesCbcEncrypt(msg, []byte("1234567891234567"))
	fmt.Println(string(cipherText))

	plainText := AesCbcDecrypt(cipherText, []byte("1234567891234567"))
	fmt.Println(string(plainText))

	// CFB （Cipher Feedback）

}
