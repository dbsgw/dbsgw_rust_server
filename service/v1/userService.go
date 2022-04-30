package v1

import (
	"dbsgw_rust_server/models"
	"dbsgw_rust_server/utils/RustJwt"
	"fmt"
	"github.com/beego/beego/v2/client/orm"
	"github.com/beego/beego/v2/core/logs"
	beego "github.com/beego/beego/v2/server/web"
)

// GetUserInfo 通过 uid 获取用户信息的
func GetUserInfo(uid string) (models.UserBase, error) {
	fmt.Println(uid)
	o := orm.NewOrm()
	user := []models.UserBase{}
	_, err := o.Raw("select * from user_base where uid = ? limit 1", uid).QueryRows(&user)
	if err != nil {
		logs.Info("获取用户信息失败", err)
		return user[0], err
	}
	return user[0], nil
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
