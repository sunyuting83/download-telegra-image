package utils

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/md5"
	"encoding/base64"
	"encoding/json"
	"io/ioutil"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
)

// Config config
type Config struct {
	RunPath  string `json:"runPath"`
	PathList []Path `json:"pathList"`
}

// Path path
type Path struct {
	Title string `json:"title"`
	Path  string `json:"path"`
}

// CORSMiddleware cors middleware
func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "OPTIONS, GET")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	}
}

// SetConfigMiddleWare set config
func SetConfigMiddleWare(d string) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set("runPath", d)
		c.Writer.Status()
	}
}

// GetConfig get config
func GetConfig(d string) (j *Config) {
	jsonFile := strings.Join([]string{d, "config.json"}, "/")
	config, _ := ioutil.ReadFile(jsonFile)
	var (
		index int = len(config)
	)
	index = bytes.IndexByte(config, 0)
	if index != -1 {
		config = config[:index]
	}
	if err := json.Unmarshal(config, &j); err != nil {
		return
	}
	return
}

// PathExists PathExists
func PathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

// DecodeBytes decode Bytes
func DecodeBytes(a string) (b string) {
	// bb, _ := base64.RawURLEncoding.DecodeString(a)
	bb, _ := base64.StdEncoding.DecodeString(a)
	return string(bb)
}

// DeCodeBytes De Code Bytes
func DeCodeBytes(a string) (b string) {
	var str []byte = []byte(a)
	decodeBytes := make([]byte, base64.StdEncoding.DecodedLen(len(str))) // 计算解码后的长度
	base64.StdEncoding.Decode(decodeBytes, str)
	return string(decodeBytes)
}

// PKCS5Padding PKCS5Padding
func PKCS5Padding(ciphertext []byte, blockSize int) []byte {
	padding := blockSize - len(ciphertext)%blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(ciphertext, padtext...)
}

// PKCS5UnPadding PKCS5UnPadding
func PKCS5UnPadding(origData []byte) []byte {
	length := len(origData)
	unpadding := int(origData[length-1])
	return origData[:(length - unpadding)]
}

// AesEncrypt AesEncrypt
func AesEncrypt(origData, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	blockSize := block.BlockSize()
	origData = PKCS5Padding(origData, blockSize)
	blockMode := cipher.NewCBCEncrypter(block, key[:blockSize])
	crypted := make([]byte, len(origData))
	blockMode.CryptBlocks(crypted, origData)
	return crypted, nil
}

// AesDecrypt AesDecrypt
func AesDecrypt(crypted, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	blockSize := block.BlockSize()
	blockMode := cipher.NewCBCDecrypter(block, key[:blockSize])
	origData := make([]byte, len(crypted))
	blockMode.CryptBlocks(origData, crypted)
	origData = PKCS5UnPadding(origData)
	return origData, nil
}

// MakeMD5 make md5
func MakeMD5(s string) (m string) {
	data := []byte(s)
	hash := md5.New()
	m = string(hash.Sum(data))
	return
}
