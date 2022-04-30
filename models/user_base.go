package models

import "github.com/beego/beego/v2/client/orm"

// UserBase 用户信息表
type UserBase struct {
	Id             int    `orm:"-"`
	Uid            string `json:"uid"`              // '用户ID'
	UserRole       int8   `json:"user_role"`        // '2正常用户 3禁言用户 4虚拟用户 5运营'
	RegisterSource int8   `json:"register_source"`  // '注册来源：1邮箱 2gitee 3githup '
	UserName       string `json:"user_name"`        // '用户账号，必须唯一'
	NickName       string `json:"nick_name"`        // '用户昵称'
	Gender         int8   `json:"gender"`           // '用户性别 0-female 1-male'
	Birthday       int    `json:"birthday"`         // '用户生日'
	Signature      string `json:"signature"`        // '用户个人签名'
	Mobile         string `json:"mobile"`           // '手机号码(唯一)'
	MobileBindTime int    `json:"mobile_bind_time"` // '手机号码绑定时间'
	Email          string `json:"email"`            // '邮箱(唯一)'
	EmailBindTime  int    `json:"email_bind_time"`  // '邮箱绑定时间'
	Face           string `json:"face"`             // '头像'
	Face200        string `json:"face200"`          // '头像 200x200x80'
	Srcface        string `json:"srcface"`          // '原图头像'
	CreateTime     int    `json:"create_time"`      // '创建时间'
	UpdateTime     int    `json:"update_time"`      // '修改时间'
	//PushToken      string `json:"push_token"`       // '用户设备push_token'
}

func init() {
	orm.RegisterModel(new(UserBase))
}
