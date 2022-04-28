package v1

import (
	"dbsgw_rust_server/controllers"
	"dbsgw_rust_server/initialize"
	"dbsgw_rust_server/models"
	"dbsgw_rust_server/utils"
	"dbsgw_rust_server/utils/RustEmail"
	"dbsgw_rust_server/utils/RustGitHup"
	"dbsgw_rust_server/utils/RustGitee"
	"dbsgw_rust_server/utils/RustJwt"
	"github.com/beego/beego/v2/client/orm"
	"github.com/beego/beego/v2/core/logs"
	"github.com/beego/beego/v2/core/validation"
	beego "github.com/beego/beego/v2/server/web"
	gonanoid "github.com/matoous/go-nanoid/v2"
	"log"
	"strconv"
	"time"
)

// UserController Operations about Users
type UserController struct {
	controllers.BaseController
}

// Code  邮箱验证码
func (u *UserController) Code() {
	// 获取邮箱地址发送验证码
	email := u.GetString("email")
	valid := validation.Validation{}
	valid.Required(email, "email")

	if valid.HasErrors() {
		// 如果有错误信息，证明验证没通过
		// 打印错误信息
		for _, err := range valid.Errors {
			log.Println(err.Key, err.Message)
			logs.Info("验证码---表单验证错误：", err)

		}
		u.Fail("表单验证错误", 500)
		return
	} else {

		// 生成验证码  放到redis里面  可以 已邮箱做key

		randstr := utils.RandString(6)

		// 1 分钟  设置到  redis里面去
		err := initialize.Rdb.Set(email, randstr, time.Minute*1).Err()
		if err != nil {
			logs.Info("验证码---插入redis失败:", email)
			u.Fail("验证设置失败", 500)
			return
		}

		// 发送邮箱验证码
		rustEmail := RustEmail.NewDefaultSendEmail()
		rustEmail.Send([]string{email}, "Rust中文网", "<h1>来自Rust中文网验证码："+randstr+"</h1>")
		u.Ok("发送成功")
	}

}

// Login  邮箱登录
func (u *UserController) Login() {
	email := u.GetString("email")
	code := u.GetString("code")

	valid := validation.Validation{}
	valid.Required(code, "code")
	valid.Required(email, "email")

	if valid.HasErrors() {
		// 如果有错误信息，证明验证没通过
		// 打印错误信息
		for _, err := range valid.Errors {
			logs.Info("邮箱登录---表单验证错误：", err)
		}
		u.Fail("表单验证错误", 500)
		return
	} else {
		val, err := initialize.Rdb.Get(email).Result()
		if err != nil {
			logs.Error("邮箱登录---获取redis失败：", err)
			u.Fail("获取验证码失败", 500)
			return
		}
		if code == val {

			// 登录注册  有就登录  没有就注册
			// 插入数据库
			o := orm.NewOrm()
			GetAuth := []models.UserAuth{}
			num, err := o.Raw("select * from user_auth where identity_type = ? and identifier = ? limit 100", 1, email).QueryRows(&GetAuth)
			if err != nil {
				logs.Error("用户查询auth报错---系统错误", err)
				u.Fail("用户查询auth报错---系统错误", 500)
				return
			}

			switch num {
			case 0:
				uid, _ := gonanoid.New()
				UserAuth := models.UserAuth{
					Uid:          uid,
					IdentityType: 1,
					Identifier:   email,
					Certificate:  "",
					CreateTime:   int(utils.GetUnix()),
					UpdateTime:   int(utils.GetUnix()),
				}
				UserBase := models.UserBase{
					Uid:            uid,
					UserRole:       0,
					RegisterSource: 1,
					UserName:       "",
					NickName:       "",
					Gender:         0,
					Birthday:       0,
					Signature:      "",
					Mobile:         "",
					MobileBindTime: 0,
					Email:          "",
					EmailBindTime:  0,
					Face:           "",
					Face200:        "",
					Srcface:        "",
					CreateTime:     int(utils.GetUnix()),
					UpdateTime:     int(utils.GetUnix()),
				}
				_, err = o.Insert(&UserBase)
				if err != nil {
					logs.Error("插入UserBase失败", err)
					u.Fail("插入UserBase失败", 500)
					return
				}
				_, err = o.Insert(&UserAuth)
				if err != nil {
					logs.Error("插入UserAuth失败", err)
					u.Fail("插入UserAuth失败", 500)
					return
				}
				Secrect, err := beego.AppConfig.String("Secrect")
				if err != nil {
					logs.Info("获取jwt密钥错误", err)
				}
				token, err := RustJwt.CreateToken(uid, Secrect)
				if err != nil {
					logs.Info("token生成失败", err)
					u.Fail("token生成失败", 500)
					return
				}

				u.SetSession(token, UserBase)
				u.Ok(map[string]interface{}{
					"msg":   "注册成功",
					"token": token,
				})
				break
			case 1:
				GetBase := []models.UserBase{}
				num, err = o.Raw("select * from user_base where uid =?", GetAuth[0].Uid).QueryRows(&GetBase)
				if err != nil {
					logs.Error("用户查询base报错---系统错误", err)
					u.Fail("用户查询base报错---系统错误", 500)
					return
				}
				Secrect, err := beego.AppConfig.String("Secrect")
				if err != nil {
					logs.Info("获取jwt密钥错误", err)
				}
				token, err := RustJwt.CreateToken(GetBase[0].Uid, Secrect)
				if err != nil {
					logs.Info("token生成失败", err)
					u.Fail("token生成失败", 500)
					return
				}

				u.SetSession(token, GetBase[0])
				u.Ok(map[string]interface{}{
					"msg":   "登录成功",
					"token": token,
				})

			default:
				u.Fail("用户存在多个账号", 500)
			}

		} else {
			logs.Info("验证码不正确")
			u.Fail("验证码不正确", 500)
			return
		}
	}
}

// OauthGitee  授权gitee
func (u *UserController) OauthGitee() {
	// 获取code
	code := u.GetString("code")
	bl := utils.IsEmpty(code)
	if !bl {
		logs.Info("code为空---系统错误")
		u.Fail("code为空---系统错误", 500)
		return
	}
	// 通过code 获取token
	userToken := RustGitee.GetAccessToken(code)
	bl = utils.IsEmpty(userToken.AccessToken)
	if !bl {
		logs.Info("userToken为空---系统错误")
		u.Fail("userToken为空---系统错误", 500)
		return
	}
	// 通过token 获取用户信息
	user := RustGitee.GetUserInfos(userToken)
	bl = utils.IsEmpty(user.Id)
	if !bl {
		u.Fail("user.Id为空---系统错误", 500)
		return
	}
	//// 通过token 获取邮箱
	//emailsReust := RustGitee.GetEmails(userToken)
	//// 序列化  取邮箱
	//emailsAll := []map[string]interface{}{}
	//json.Unmarshal([]byte(emailsReust), &emailsAll)
	//
	//email := ""
	//if len(emailsAll) != 0 {
	//	email = emailsAll[0]["email"].(string)
	//}

	// 插入数据库
	o := orm.NewOrm()
	GetAuth := []models.UserAuth{}
	num, err := o.Raw("select * from user_auth where identity_type = ? and identifier = ? limit 100", 2, user.Id).QueryRows(&GetAuth)
	if err != nil {
		logs.Error("用户查询auth报错---系统错误", err)
		u.Fail("用户查询auth报错---系统错误", 500)
		return
	}

	switch num {
	case 0:
		uid, _ := gonanoid.New()
		UserAuth := models.UserAuth{
			Uid:          uid,
			IdentityType: 2,
			Identifier:   strconv.Itoa(user.Id),
			Certificate:  "",
			CreateTime:   int(utils.GetUnix()),
			UpdateTime:   int(utils.GetUnix()),
		}
		UserBase := models.UserBase{
			Uid:            uid,
			UserRole:       0,
			RegisterSource: 2,
			UserName:       user.Login,
			NickName:       user.Name,
			Gender:         0,
			Birthday:       0,
			Signature:      "",
			Mobile:         "",
			MobileBindTime: 0,
			Email:          "",
			EmailBindTime:  0,
			Face:           user.AvatarUrl,
			Face200:        "",
			Srcface:        "",
			CreateTime:     int(utils.GetUnix()),
			UpdateTime:     int(utils.GetUnix()),
		}
		_, err = o.Insert(&UserBase)
		if err != nil {
			logs.Error("插入UserBase失败", err)
			u.Fail("插入UserBase失败", 500)
			return
		}
		_, err = o.Insert(&UserAuth)
		if err != nil {
			logs.Error("插入UserAuth失败", err)
			u.Fail("插入UserAuth失败", 500)
			return
		}
		Secrect, err := beego.AppConfig.String("Secrect")
		if err != nil {
			logs.Info("获取jwt密钥错误", err)
		}
		token, err := RustJwt.CreateToken(uid, Secrect)
		if err != nil {
			logs.Info("token生成失败", err)
			u.Fail("token生成失败", 500)
			return
		}

		u.SetSession(token, UserBase)
		u.Ok(map[string]interface{}{
			"msg":   "注册成功",
			"token": token,
		})
		break
	case 1:
		GetBase := []models.UserBase{}
		num, err = o.Raw("select * from user_base where uid =?", GetAuth[0].Uid).QueryRows(&GetBase)
		if err != nil {
			logs.Error("用户查询base报错---系统错误", err)
			u.Fail("用户查询base报错---系统错误", 500)
			return
		}
		Secrect, err := beego.AppConfig.String("Secrect")
		if err != nil {
			logs.Info("获取jwt密钥错误", err)
		}
		token, err := RustJwt.CreateToken(GetBase[0].Uid, Secrect)
		if err != nil {
			logs.Info("token生成失败", err)
			u.Fail("token生成失败", 500)
			return
		}

		u.SetSession(token, GetBase[0])
		u.Ok(map[string]interface{}{
			"msg":   "登录成功",
			"token": token,
		})

	default:
		u.Fail("用户存在多个账号", 500)
	}
}

// OauthGitHup 授权githup
func (u *UserController) OauthGitHup() {
	// 获取code
	code := u.GetString("code")
	bl := utils.IsEmpty(code)
	if !bl {
		logs.Info("code为空---系统错误")
		u.Fail("code为空---系统错误", 500)
		return
	}
	// 通过code 获取token
	userToken := RustGitHup.GetAccessToken(code)
	bl = utils.IsEmpty(userToken.AccessToken)
	if !bl {
		logs.Info("userToken为空---系统错误")
		u.Fail("userToken为空---系统错误", 500)
		return
	}
	// 通过token 获取用户信息
	user := RustGitHup.GetUserInfos(userToken)
	bl = utils.IsEmpty(user.Id)
	if !bl {
		logs.Info("user.Id为空---系统错误")
		u.Fail("user.Id为空---系统错误", 500)
		return
	}
	//// 通过token 获取邮箱
	//emailsReust := RustGitHup.GetEmails(userToken)
	//// 序列化  取邮箱
	//emailsAll := []map[string]interface{}{}
	//json.Unmarshal([]byte(emailsReust), &emailsAll)
	//
	//email := ""
	//if len(emailsAll) != 0 {
	//	email = emailsAll[0]["email"].(string)
	//}

	// 插入数据库
	o := orm.NewOrm()
	GetAuth := []models.UserAuth{}
	num, err := o.Raw("select * from user_auth where identity_type = ? and identifier = ? limit 100", 3, user.Id).QueryRows(&GetAuth)
	if err != nil {
		logs.Error("用户查询auth报错---系统错误", err)
		u.Fail("用户查询auth报错---系统错误", 500)
		return
	}

	switch num {
	case 0:
		uid, _ := gonanoid.New()
		UserAuth := models.UserAuth{
			Uid:          uid,
			IdentityType: 3,
			Identifier:   strconv.Itoa(user.Id),
			Certificate:  "",
			CreateTime:   int(utils.GetUnix()),
			UpdateTime:   int(utils.GetUnix()),
		}
		UserBase := models.UserBase{
			Uid:            uid,
			UserRole:       0,
			RegisterSource: 3,
			UserName:       user.Login,
			NickName:       user.Name,
			Gender:         0,
			Birthday:       0,
			Signature:      "",
			Mobile:         "",
			MobileBindTime: 0,
			Email:          "",
			EmailBindTime:  0,
			Face:           user.AvatarUrl,
			Face200:        "",
			Srcface:        "",
			CreateTime:     int(utils.GetUnix()),
			UpdateTime:     int(utils.GetUnix()),
		}
		_, err = o.Insert(&UserBase)
		if err != nil {
			logs.Error("插入UserBase失败---系统错误", err)
			u.Fail("插入UserBase失败", 500)
			return
		}
		_, err = o.Insert(&UserAuth)
		if err != nil {
			logs.Error("插入UserAuth失败---系统错误", err)
			u.Fail("插入UserAuth失败", 500)
			return
		}

		Secrect, err := beego.AppConfig.String("Secrect")
		if err != nil {
			logs.Info("获取jwt密钥错误", err)
		}
		token, err := RustJwt.CreateToken(uid, Secrect)
		if err != nil {
			logs.Info("token生成失败", err)
			u.Fail("token生成失败", 500)
			return
		}

		u.SetSession(token, UserBase)
		u.Ok(map[string]interface{}{
			"msg":   "注册成功",
			"token": token,
		})
		break
	case 1:
		GetBase := []models.UserBase{}
		num, err = o.Raw("select * from user_base where uid =?", GetAuth[0].Uid).QueryRows(&GetBase)
		if err != nil {
			logs.Error("用户查询base报错---系统错误", err)
			u.Fail("用户查询base报错---系统错误", 500)
			return
		}
		Secrect, err := beego.AppConfig.String("Secrect")
		if err != nil {
			logs.Info("获取jwt密钥错误", err)
		}
		token, err := RustJwt.CreateToken(GetBase[0].Uid, Secrect)
		if err != nil {
			logs.Info("token生成失败", err)
			u.Fail("token生成失败", 500)
			return
		}

		u.SetSession(token, GetBase[0])
		u.Ok(map[string]interface{}{
			"msg":   "登录成功",
			"token": token,
		})

	default:
		u.Fail("用户存在多个账号", 500)
	}

}
