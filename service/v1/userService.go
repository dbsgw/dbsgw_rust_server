package v1

import (
	"dbsgw_rust_server/models"
	"fmt"
	"github.com/beego/beego/v2/client/orm"
	"github.com/beego/beego/v2/core/logs"
)

// GetUserInfo 通过 uid 获取用户信息的
func GetUserInfo(uid string) models.UserBase {
	fmt.Println(uid)
	o := orm.NewOrm()
	user := []models.UserBase{}
	_, err := o.Raw("select * from user_base where uid = ? limit 1", uid).QueryRows(&user)
	if err != nil {
		logs.Info("获取用户信息失败", err)
		return user[0]
	}
	return user[0]
}
