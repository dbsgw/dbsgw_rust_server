package v1

import (
	"dbsgw_rust_server/initialize"
	"dbsgw_rust_server/models"
	"dbsgw_rust_server/utils/RustJwt"
	"errors"
	"github.com/beego/beego/v2/adapter/context"
	"github.com/beego/beego/v2/core/logs"
	beego "github.com/beego/beego/v2/server/web"
)

// GetUserInfo 通过 uid 获取用户信息的
func GetUserInfo(uid string) (models.UserBase, error) {
	user := []models.UserBase{}
	err := initialize.DB.Raw("select * from user_base where uid = ? limit 1", uid).Scan(&user).Error
	if len(user) == 0 {
		if err != nil {
			logs.Info("获取用户信息失败", err)
			return models.UserBase{}, err
		}
		return models.UserBase{}, nil
	} else {
		if err != nil {
			logs.Info("获取用户信息失败", err)
			return models.UserBase{}, err
		}
		return user[0], nil
	}

}

// GetUserUpdateInfo 通过 uid 更新用户信息的
func GetUserUpdateInfo(mobild, nick_name, uid string) error {
	result := map[string]interface{}{}
	err := initialize.DB.Raw("update user_base set mobile = ?,nick_name=? where uid = ?", mobild, nick_name, uid).Scan(&result).Error
	if err != nil {
		logs.Info("更新用户信息失败", err)
		return errors.New("更新用户信息失败")
	}

	return nil
}

// RustCreateToken 传入uid 生成生成 token
func RustCreateToken(uid string) (string, error) {
	Secrect, err := beego.AppConfig.String("Secrect")
	if err != nil {
		logs.Info("获取jwt密钥错误", err)
	}
	token, err := RustJwt.CreateToken(uid, Secrect)
	if err != nil {
		logs.Info("token生成失败", err)
		return "", err
	}
	return token, nil
}

// RustTokenGetUserInfo 通过token获取用户信息
func RustTokenGetUserInfo() models.UserBase {
	u := context.NewContext()
	token := u.Request.Header["Authorization"]

	userinfo := models.UserBase{}
	if len(token) != 1 {
		logs.Info("获取用户token失败")
		return userinfo
	}

	sessionResult := u.Input.Session(token[0])
	return sessionResult.(models.UserBase)
}
