package utils

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/md5"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"math"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
)

// Config config
type Config struct {
	RunPath     string `json:"runPath"`
	PathList    []Path `json:"pathList"`
	Cors        string `json:"cors"`
	RootURL     string `json:"rooturl"`
	RunningFile string `json:"runningFile"`
	DoneFile    string `json:"doneFile"`
}

// Path path
type Path struct {
	Title string `json:"title"`
	Path  string `json:"path"`
}

// SetConfigMiddleWare set config
func SetConfigMiddleWare(d, port string) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set("runPath", d)
		c.Set("port", port)
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
	w := md5.New()
	io.WriteString(w, s)
	m = fmt.Sprintf("%x", w.Sum(nil))
	return
}

// IsDir Is Dir
func IsDir(path string) bool {
	s, err := os.Stat(path)
	if err != nil {
		return false
	}
	return s.IsDir()
}

// Round round
func Round(x float64) int {
	return int(math.Floor(x + 0/5))
}
