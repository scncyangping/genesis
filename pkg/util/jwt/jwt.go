package jwt

import (
	"github.com/dgrijalva/jwt-go"
	"time"
)

type Info struct {
	Username string `json:"username"`
}

type Claims struct {
	Base  jwt.StandardClaims
	Extra Info
}

func (c Claims) Valid() error {
	return c.Base.Valid()
}

func GenerateToken(name, issuer, secret string, et int) (string, error) {
	nowTime := time.Now()
	expireTime := nowTime.Add(time.Duration(et) * time.Second)

	claims := Claims{
		Extra: Info{Username: name},
		Base: jwt.StandardClaims{
			ExpiresAt: expireTime.Unix(),
			Issuer:    issuer,
		},
	}

	tokenClaims := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	//  该方法内部生成签名字符串，再用于获取完整、已签名的token
	token, err := tokenClaims.SignedString([]byte(secret))
	return token, err
}

func ParseToken(token, secret string) (*Claims, error) {
	// 用于解析鉴权的声明，方法内部主要是具体的解码和校验的过程，最终返回*Token
	tokenClaims, err := jwt.ParseWithClaims(token, &Claims{}, func(token *jwt.Token) (any, error) {
		return []byte(secret), nil
	})

	if tokenClaims != nil {
		// 验证基于时间的声明exp, iat, nbf，注意如果没有任何声明在令牌中
		// 仍然会被认为是有效的。并且对于时区偏差没有计算方法
		if claims, ok := tokenClaims.Claims.(*Claims); ok && tokenClaims.Valid {
			return claims, nil
		}
	}

	return nil, err
}
