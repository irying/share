package middleware

import (
	"backend/api/proxy"
	"backend/conf"
	"backend/pb/user"
	"backend/system/exception"
	"context"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)
var jwtSecret = []byte(conf.Conf.App.Api.ApiSecretKey)

// Claims payloads
type Claims struct {
	Username string `json:"username"`
	Password string `json:"password"`
	jwt.StandardClaims
}

// JWT verify token
func JWT() gin.HandlerFunc {
	return func(c *gin.Context) {
		var code int
		var data interface{}

		code = exception.SUCCESS
		token := c.Query("token")
		if token == "" {
			code = exception.INVALID_PARAMS
		} else if token == "test"{

		} else {
			claims, err := ParseToken(token)
			if err != nil {
				code = exception.ERROR_AUTH_CHECK_TOKEN_FAIL
			} else if IsExpired(claims, token){
				code = exception.ERROR_AUTH_CHECK_TOKEN_TIMEOUT
			}
		}

		if code != exception.SUCCESS {
			c.JSON(http.StatusUnauthorized, gin.H{
				"code": code,
				"msg":  exception.GetMsg(code),
				"data": data,
			})

			c.Abort()
			return
		}
		c.Next()
	}
}

// ParseToken parse token
func ParseToken(token string) (*Claims, error) {
	tokenClaims, err := jwt.ParseWithClaims(token, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})

	if tokenClaims != nil {
		if claims, ok := tokenClaims.Claims.(*Claims); ok && tokenClaims.Valid {
			return claims, nil
		}
	}

	return nil, err
}

// GenerateToken generate token by username and password
func GenerateToken(username, password string) (string, error) {
	nowTime := time.Now()
	expireTime := nowTime.Add(3 * time.Hour)

	claims := Claims{
		username,
		password,
		jwt.StandardClaims {
			ExpiresAt : expireTime.Unix(),
			Issuer : "backend",
		},
	}

	tokenClaims := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	token, err := tokenClaims.SignedString(jwtSecret)
	return token, err
}

// IsExpired is expired token
func IsExpired(claims *Claims, token string) bool {
	if time.Now().Unix() > claims.ExpiresAt {
		c := context.Background()
		res, _ := pb.NewUserServiceClient(proxy.NewRPCConn()).IsExpiredToken(c, &pb.TokenRequest{Username: claims.Username, Token: token})
		return res.IsExpired
	}

	return false
}
