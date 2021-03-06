package main

import (
	"bufio"
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"errors"
	"flag"
	"fmt"
	"math/rand"
	"os"
	"time"
)

var PwdKey = []byte("ABCDABCDABCDABCD")

func pkcs7Padding(data []byte, blockSize int) []byte {
	padding := blockSize - len(data)%blockSize
	padText := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(data, padText...)
}

func pkcs7UnPadding(data []byte) ([]byte, error) {
	length := len(data)
	if length == 0 {
		return nil, errors.New("加密字符串错误！")
	}
	unPadding := int(data[length-1])
	return data[:(length - unPadding)], nil
}

func AesEncrypt(data []byte, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	blockSize := block.BlockSize()
	encryptBytes := pkcs7Padding(data, blockSize)
	crypted := make([]byte, len(encryptBytes))
	blockMode := cipher.NewCBCEncrypter(block, key[:blockSize])
	blockMode.CryptBlocks(crypted, encryptBytes)
	return crypted, nil
}

func AesDecrypt(data []byte, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	blockSize := block.BlockSize()
	blockMode := cipher.NewCBCDecrypter(block, key[:blockSize])
	crypted := make([]byte, len(data))
	blockMode.CryptBlocks(crypted, data)
	crypted, err = pkcs7UnPadding(crypted)
	if err != nil {
		return nil, err
	}
	return crypted, nil
}

func EncryptByAes(data []byte) (string, error) {
	res, err := AesEncrypt(data, PwdKey)
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(res), nil
}

func DecryptByAes(data string) ([]byte, error) {
	dataByte, err := base64.StdEncoding.DecodeString(data)
	if err != nil {
		return nil, err
	}
	return AesDecrypt(dataByte, PwdKey)
}

func EncryptFile(filePath string, fName string) (err error) {
	f, err := os.Open(filePath)
	if err != nil {
		fmt.Println("未找到文件")
		return
	}
	defer f.Close()

	fInfo, _ := f.Stat()
	fLen := fInfo.Size()
	fmt.Println("待处理文件大小:", fLen)
	maxLen := 1024 * 1024 * 100
	var forNum int64 = 0
	getLen := fLen

	if fLen > int64(maxLen) {
		getLen = int64(maxLen)
		forNum = fLen / int64(maxLen)
		fmt.Println("需要加密次数：", forNum+1)
	}
	ff, err := os.OpenFile(fName, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		fmt.Println("文件写入错误")
		return err
	}
	defer ff.Close()
	for i := 0; i < int(forNum+1); i++ {
		a := make([]byte, getLen)
		n, err := f.Read(a)
		if err != nil {
			fmt.Println("文件读取错误")
			return err
		}
		getByte, err := EncryptByAes(a[:n])
		if err != nil {
			fmt.Println("加密错误")
			return err
		}

		getBytes := append([]byte(getByte), []byte("\n")...)
		buf := bufio.NewWriter(ff)
		buf.WriteString(string(getBytes[:]))
		buf.Flush()
	}
	ffInfo, _ := ff.Stat()
	fmt.Printf("文件加密成功，生成文件名为：%s，文件大小为：%v Byte \n", ffInfo.Name(), ffInfo.Size())
	return nil
}

func DecryptFile(filePath, fName string) (err error) {
	f, err := os.Open(filePath)
	if err != nil {
		fmt.Println("未找到文件")
		return
	}
	defer f.Close()
	fInfo, _ := f.Stat()
	fmt.Println("待处理文件大小:", fInfo.Size())

	br := bufio.NewReader(f)
	ff, err := os.OpenFile(fName, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		fmt.Println("文件写入错误")
		return err
	}
	defer ff.Close()
	num := 0
	for {
		num = num + 1
		a, err := br.ReadString('\n')
		if err != nil {
			break
		}
		getByte, err := DecryptByAes(a)
		if err != nil {
			fmt.Println("解密错误")
			return err
		}

		buf := bufio.NewWriter(ff)
		buf.Write(getByte)
		buf.Flush()
	}
	fmt.Println("解密次数：", num)
	ffInfo, _ := ff.Stat()
	fmt.Printf("文件解密成功，生成文件名为：%s，文件大小为：%v Byte \n", ffInfo.Name(), ffInfo.Size())
	return
}
func CreateCaptcha() string {
	return fmt.Sprintf("%06v", rand.New(rand.NewSource(time.Now().UnixNano())).Int31n(1000000))
}
func main() {
	var encrypt_file = flag.String("encrypt_file", "", "string类型参数")
	var encrypt_out = flag.String("encrypt_out", "", "string类型参数")
	var decrypt_file = flag.String("decrypt_file", "", "string类型参数")
	var decrypt_out = flag.String("decrypt_out", "", "string类型参数")
	flag.Parse()
	fmt.Println("---- auth:kill02lc ----")
	fmt.Println("-encrypt_file 要加密的文件路径", *encrypt_file)
	fmt.Println("-encrypt_out 加密后的文件路径(无需后缀)", *encrypt_out)
	fmt.Println("-decrypt_file 要解密的文件路径", *decrypt_file)
	fmt.Println("-decrypt_out 解密后的文件路径(需要原文件的后缀)", *decrypt_out)
	fmt.Println("\n")
	fmt.Println("使用样例:")
	fmt.Println(".\\aes_en_de.exe -encrypt_file test.txt -encrypt_out encrypt_123456")
	fmt.Println(".\\aes_en_de.exe -decrypt_file encrypt_123456 -decrypt_out temp.txt")
	startTime := time.Now()
	if len(*encrypt_file) != 0 && len(*encrypt_out) != 0 {
		fmt.Println("\n")
		EncryptFile(*encrypt_file, *encrypt_out)
	} else if len(*decrypt_file) != 0 && len(*decrypt_out) != 0 {
		fmt.Println("\n")
		DecryptFile(*decrypt_file, *decrypt_out)
	} else {
		fmt.Println("\n")
		fmt.Println("plase input your parameter!")
		os.Exit(-1)
	}
	fmt.Printf("耗时：%v", time.Since(startTime))
}
