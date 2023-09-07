package jwtutil

import (
	"errors"
	"time"

	"github.com/dgrijalva/jwt-go"
)

const (
	EXPIRE_TIME = 60 * 60 * 24 * 30 //过期时间，一个月后过期
	START_TIME  = 60 * 5            //五分钟前生效
	KEY_BYTES   = "toktik"
	ISSUER      = "toktik" //签发者
)

var (
	TokenExpired   = errors.New("token已过期")           //token过期
	TokenMalformed = errors.New("token is malformed") //token格式不对
	TokenInvalid   = errors.New("token无效")
	TokenEmpty     = errors.New("token为空")
)

type JwtUtil struct {
	keyBytes []byte
}

type UserClaims struct {
	UserId int64 `json:"user_id"`
	jwt.StandardClaims
}

var jwtUtil *JwtUtil

func NewJwtUtil() *JwtUtil {
	return &JwtUtil{
		keyBytes: []byte(KEY_BYTES),
	}
}

func (j *JwtUtil) GenerateToken(userClaims *UserClaims) (string, error) {
	return jwt.NewWithClaims(jwt.SigningMethodHS256, userClaims).SignedString(j.keyBytes)
}

func (j *JwtUtil) ParseToken(tokenStr string) (*UserClaims, error) {
	if tokenStr == "" {
		return nil, TokenEmpty
	}

	token, err := jwt.ParseWithClaims(tokenStr, &UserClaims{}, func(token *jwt.Token) (interface{}, error) {
		return j.keyBytes, nil
	})

	if err != nil {
		//如果是由于token不可用引发的错误
		if ve, ok := err.(*jwt.ValidationError); ok {
			//判断是何种错误
			if ve.Errors == jwt.ValidationErrorMalformed {
				return nil, TokenMalformed
			} else if ve.Errors == jwt.ValidationErrorExpired {
				return nil, TokenExpired
			}
		} else {
			return nil, TokenInvalid
		}
		return nil, err
	}

	//类型转化一下
	if claims, ok := token.Claims.(*UserClaims); ok {
		return claims, nil
	} else {
		return nil, TokenInvalid
	}
}

func CreateClaims(id int64) *UserClaims {
	return &UserClaims{
		UserId: id,
		StandardClaims: jwt.StandardClaims{
			Issuer:    ISSUER,
			ExpiresAt: time.Now().Unix() + EXPIRE_TIME,
			NotBefore: time.Now().Unix() - START_TIME,
		},
	}
}

func GenerateTokenWithUserId(userId int64) (string, error) {
	if jwtUtil == nil {
		jwtUtil = NewJwtUtil()
	}
	return jwtUtil.GenerateToken(CreateClaims(userId))
}

func ParseToken(tokenStr string) (*UserClaims, error) {
	if jwtUtil == nil {
		jwtUtil = NewJwtUtil()
	}
	return jwtUtil.ParseToken(tokenStr)
}
