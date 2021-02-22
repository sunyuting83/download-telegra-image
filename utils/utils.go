package utils

import (
	"crypto/md5"
	"encoding/base64"

	"github.com/gin-gonic/gin"
)

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

// DeBase decode base
func DeBase(s string) (k string) {
	var str []byte = []byte(s)
	decodeBytes := make([]byte, base64.StdEncoding.DecodedLen(len(str)))
	base64.StdEncoding.Decode(decodeBytes, str)
	k = string(decodeBytes)
	return
}

// MakeMD5 make md5
func MakeMD5(s string) (m string) {
	data := []byte(s)
	hash := md5.New()
	m = string(hash.Sum(data))
	return
}
