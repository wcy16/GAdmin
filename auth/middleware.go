package auth

import (
	"github.com/gin-gonic/gin"
	"math/rand"
	"net/http"
	"time"
)

var tokenBuf string
var tokenExpire time.Time
var tokenType int
var tmpTokenDuration time.Duration
var tmpTokenMinutes int

func init() {
	tokenBuf = ""
	tokenExpire = time.Now()
	tokenType = 0
	tmpTokenDuration = 10 * time.Minute
	tmpTokenMinutes = int(tmpTokenDuration.Minutes())
}

func CookieCheck() gin.HandlerFunc {
	return func(c *gin.Context) {
		token, err := c.Cookie("token")

		if err != nil || token == "" || token != tokenBuf || time.Now().Equal(tokenExpire) {
			c.Redirect(http.StatusSeeOther, "/signin")
			c.Abort()
			return
		}

		c.Next()
	}
}

// refresh token if token type is 1 time.Now().Unix()
func CookieUpdate() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()
		if tokenType == 1 {
			generateToken()
			cookieToken(c)
		}
	}
}

func SetToken(c *gin.Context, tType int) {
	generateToken()
	tokenType = tType
	cookieToken(c)
}

func generateToken() {
	tokenExpire = time.Now().Add(tmpTokenDuration)
	tokenBuf = randomString(32, tokenExpire.Unix())
}

// todo domain
func cookieToken(c *gin.Context) {
	c.SetCookie("token", tokenBuf, tmpTokenMinutes, "/", "localhost", false, true)
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
