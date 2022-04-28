package RustJwt

import (
	"errors"
	"github.com/dgrijalva/jwt-go"
	"time"
)

// CreateToken 设置token
func CreateToken(uid, secret string) (string, error) {
	at := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"uid": uid,
		"exp": time.Now().Add(time.Minute * 43200).Unix(), // 15分钟  30天
	})
	token, err := at.SignedString([]byte(secret))
	if err != nil {
		return "", err
	}
	return token, nil
}

// DeleteToken 删除token 就是设置token的过期时间是0
func DeleteToken(uid, secret string) (string, error) {
	at := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"uid": uid,
		"exp": time.Now().Add(time.Second * 1).Unix(), // 1秒
	})
	token, err := at.SignedString([]byte(secret))
	if err != nil {
		return "", err
	}
	return token, nil
}

// ParseToken 解析token
func ParseToken(token string, secret string) (string, error) {
	claim, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})
	if err != nil {
		return "", err
	}
	return claim.Claims.(jwt.MapClaims)["uid"].(string), errors.New("正确")
}

func init() {
	// 创建 用法
	//Secrect, _ := beego.AppConfig.String("Secrect")
	//token, _ := utils.CreateToken(fmt.Sprintf("%d", user.Id), Secrect)

	// 解析token
	//Secrect, _ := beego.AppConfig.String("Secrect")
	//_, err := utils.ParseToken(cookie, Secrect)
	// if err != nil { 解析token失败 }
}
