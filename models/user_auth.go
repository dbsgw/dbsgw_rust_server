package models

import "github.com/beego/beego/v2/client/orm"

// UserAuth 用户授权表
type UserAuth struct {
	Id           int    `json:"id"`            // id
	Uid          string `json:"uid"`           // '用户id',
	IdentityType int    `json:"identity_type"` // '1邮箱 2gitee 3githup ',
	Identifier   string `json:"identifier"`    // '手机号 邮箱 用户名或第三方应用的唯一标识',
	Certificate  string `json:"certificate"`   // '密码凭证(站内的保存密码，站外的不保存或保存token)',
	CreateTime   int    `json:"create_time"`   // '绑定时间',
	UpdateTime   int    `json:"update_time"`   // '更新绑定时间',
}

func init() {
	orm.RegisterModel(new(UserAuth))
}
