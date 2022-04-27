package v1

import (
	"dbsgw_rust_server/controllers"
	"dbsgw_rust_server/models"
	"dbsgw_rust_server/utils"
	"dbsgw_rust_server/utils/RustConstant"
	"dbsgw_rust_server/utils/RustGitHup"
	"dbsgw_rust_server/utils/RustGitee"
	"github.com/beego/beego/v2/client/orm"
	"github.com/beego/beego/v2/core/validation"
	gonanoid "github.com/matoous/go-nanoid/v2"
	"log"
	"strconv"
)

// UserController Operations about Users
type UserController struct {
	controllers.BaseController
}

// Login  登录
func (u *UserController) Login() {
	// 通过登录 source 方式 来判读 是邮箱还是第三方
	// 邮箱： 通过  Email 和 密码登录
	// 第三方：统一通过第三方授权登录
	//username := u.GetString("email")
	//password := u.GetString("password")
	source, _ := u.GetInt("source")

	valid := validation.Validation{}
	valid.Required(source, "source")

	if valid.HasErrors() {
		// 如果有错误信息，证明验证没通过
		// 打印错误信息
		for _, err := range valid.Errors {
			log.Println(err.Key, err.Message)
		}
		u.Fail("表单验证错误", 500)
	} else {
		switch source {
		case RustConstant.EMAIL:
			break
		case RustConstant.GITEE:
			u.Redirect(RustGitee.RedirectUrl(), 302)
			break
		case RustConstant.GITHUP:
			u.Redirect(RustGitHup.RedirectUrl(), 302)
			break
		default:
			u.Fail("登录方式不存在", 500)
		}
		u.Ok("成功")
	}
}

// OauthGitee  授权gitee
func (u *UserController) OauthGitee() {
	// 获取code
	code := u.GetString("code")
	bl := utils.IsEmpty(code)
	if !bl {
		u.Fail("code为空---系统错误", 500)
		return
	}
	// 通过code 获取token
	userToken := RustGitee.GetAccessToken(code)
	bl = utils.IsEmpty(userToken.AccessToken)
	if !bl {
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
			u.Fail("插入UserBase失败", 500)
			return
		}
		_, err = o.Insert(&UserAuth)
		if err != nil {
			u.Fail("插入UserAuth失败", 500)
			return
		}
		u.SetSession(uid, UserBase)
		u.Ok(map[string]interface{}{
			"msg":   "注册成功",
			"token": uid,
		})
		break
	case 1:
		GetBase := []models.UserBase{}
		num, err = o.Raw("select * from user_base where uid =?", GetAuth[0].Uid).QueryRows(&GetBase)
		if err != nil {
			u.Fail("用户查询base报错---系统错误", 500)
			return
		}
		u.SetSession(GetBase[0].Uid, GetBase[0])
		u.Ok(map[string]interface{}{
			"msg":   "登录成功",
			"token": GetBase[0].Uid,
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
		u.Fail("code为空---系统错误", 500)
		return
	}
	// 通过code 获取token
	userToken := RustGitHup.GetAccessToken(code)
	bl = utils.IsEmpty(userToken.AccessToken)
	if !bl {
		u.Fail("userToken为空---系统错误", 500)
		return
	}
	// 通过token 获取用户信息
	user := RustGitHup.GetUserInfos(userToken)
	bl = utils.IsEmpty(user.Id)
	if !bl {
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
			u.Fail("插入UserBase失败", 500)
			return
		}
		_, err = o.Insert(&UserAuth)
		if err != nil {
			u.Fail("插入UserAuth失败", 500)
			return
		}
		u.SetSession(uid, UserBase)
		u.Ok(map[string]interface{}{
			"msg":   "注册成功",
			"token": uid,
		})
		break
	case 1:
		GetBase := []models.UserBase{}
		num, err = o.Raw("select * from user_base where uid =?", GetAuth[0].Uid).QueryRows(&GetBase)
		if err != nil {
			u.Fail("用户查询base报错---系统错误", 500)
			return
		}
		u.SetSession(GetBase[0].Uid, GetBase[0])
		u.Ok(map[string]interface{}{
			"msg":   "登录成功",
			"token": GetBase[0].Uid,
		})

	default:
		u.Fail("用户存在多个账号", 500)
	}

}
