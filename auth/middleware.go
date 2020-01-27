package auth

import (
	"gadmin/config"
	"github.com/gin-gonic/gin"
	"math/rand"
	"net/http"
	"time"
)

var tokenBuf string
var tokenExpire time.Time
var tokenType int
var tmpTokenDuration time.Duration
var tmpTokenSeconds int
var longTokenDuration time.Duration
var longTokenSeconds int

func init() {
	tokenBuf = ""
	tokenExpire = time.Now()
	tokenType = 0
	tmpTokenDuration = 10 * time.Minute
	tmpTokenSeconds = int(tmpTokenDuration.Seconds())
	longTokenDuration = 30 * 24 * time.Hour
	longTokenSeconds = int(longTokenDuration.Seconds())
}

func CookieCheck() gin.HandlerFunc {
	return func(c *gin.Context) {
		token, err := c.Cookie("token")

		if err != nil || token == "" || token != tokenBuf || time.Now().After(tokenExpire) {
			c.Redirect(http.StatusSeeOther, config.Prefix+"/signin")
			c.Abort()
			return
		}

		c.Next()
	}
}

// refresh token if token type is 1 time.Now().Unix()
func CookieUpdate() gin.HandlerFunc {
	return func(c *gin.Context) {
		if tokenType == 1 {
			generateToken(tmpTokenDuration)
			tmpCookieToken(c)
		}
		c.Next()
	}
}

func SetToken(c *gin.Context, tType int) {
	tokenType = tType
	if tokenType == 1 {
		generateToken(tmpTokenDuration)
		tmpCookieToken(c)
	} else {
		generateToken(longTokenDuration)
		longCookieToken(c)
	}
}

func DelToken(c *gin.Context) {
	c.SetCookie("token", "", 0, "/", "localhost", false, true)
	tokenBuf = ""
	tokenExpire = time.Now()
	tokenType = 0
}

func generateToken(duration time.Duration) {
	tokenExpire = time.Now().Add(duration)
	tokenBuf = randomString(32, tokenExpire.Unix())
}

// todo domain
func tmpCookieToken(c *gin.Context) {
	c.SetCookie("token", tokenBuf, tmpTokenSeconds, "/", "localhost", false, true)
}
func longCookieToken(c *gin.Context) {
	c.SetCookie("token", tokenBuf, longTokenSeconds, "/", "localhost", false, true)
}

const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

// generate random string
func randomString(length int, seed int64) string {
	r := rand.New(rand.NewSource(seed))
	b := make([]byte, length)
	for i := range b {
		b[i] = letters[r.Intn(len(letters))]
	}
	return string(b)
}
