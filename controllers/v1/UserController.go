package v1

import (
	"dbsgw_rust_server/controllers"
	"dbsgw_rust_server/initialize"
	"dbsgw_rust_server/models"
	v1 "dbsgw_rust_server/service/v1"
	"dbsgw_rust_server/utils"
	"dbsgw_rust_server/utils/RustEmail"
	"dbsgw_rust_server/utils/RustGitHup"
	"dbsgw_rust_server/utils/RustGitee"
	"encoding/json"
	"github.com/beego/beego/v2/core/logs"
	"github.com/beego/beego/v2/core/validation"
	gonanoid "github.com/matoous/go-nanoid/v2"
	"log"
	"strconv"
	"time"
)

// UserController Operations about Users
type UserController struct {
	controllers.BaseController
}

// ArticleAll 用户所有 文章的
func (u *UserController) ArticleAll() {

	size, _ := u.GetInt("size")
	page, _ := u.GetInt("page")

	token := u.Ctx.Request.Header["Authorization"]
	if len(token) != 1 {
		u.Fail("获取用户token失败", 500)
		logs.Info("获取用户token失败")
		return
	}
	sessionResult := u.GetSession(token[0])

	if size == 0 {
		size = 20
	}
	if page == 0 {
		page = 1
	}
	var count int64
	GetArticle := []models.Article{}
	err := initialize.DB.Where("user_id = ?", sessionResult.(models.UserBase).Uid).Offset((page - 1) * size).Limit(size).Order("article_time desc").Find(&GetArticle).Error
	if err != nil {
		logs.Error("用户查询article报错---系统错误", err)
		u.Fail("用户查询article报错---系统错误", 500)
		return
	}
	initialize.DB.Table("article").Where("user_id = ?", sessionResult.(models.UserBase).Uid).Count(&count)
	u.Ok(map[string]interface{}{
		"result": GetArticle,
		"total":  count,
	})
}

// Logout 退出登录
func (u *UserController) Logout() {
	token := u.Ctx.Request.Header["Authorization"]
	if len(token) != 1 {
		u.Fail("获取用户token失败", 500)
		logs.Info("获取用户token失败")
		return
	}
	err := u.DelSession(token[0])
	if err != nil {
		logs.Error("删除token失败")
		u.Fail("删除toke失败", 500)
		return
	}
	u.Ok("退出成功")
}

// Info 用户详情
func (u *UserController) Info() {
	id := u.Ctx.Input.Param(":id")
	userinfo, err := v1.GetUserInfo(id)
	if err != nil {
		u.Fail("获取用户信息失败", 500)
		return
	}
	u.Ok(userinfo)
}

// InfoPut 更新详情
func (u *UserController) InfoPut() {
	id := u.Ctx.Input.Param(":id")
	var userInfo map[string]interface{}
	json.Unmarshal(u.Ctx.Input.RequestBody, &userInfo)
	//valid := validation.Validation{}

	err := v1.GetUserUpdateInfo(userInfo["Mobile"].(string), userInfo["NickName"].(string), id)
	if err != nil {
		u.Fail("修改失败", 500)
	}
	u.Ok("修改成功")
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
		err = rustEmail.Send([]string{email}, "Rust中文网", "<h1>来自Rust中文网验证码："+randstr+"</h1>")
		if err != nil {
			logs.Error(err, "邮箱发送错误")
			u.Fail("邮箱发送错误", 500)
			return
		}
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
				UserName:       "新用户",
				NickName:       "新用户",
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

			u.RustCreate(&UserAuth, &UserBase, email, 1)

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

	id := strconv.Itoa(user.Id)
	u.RustCreate(&UserAuth, &UserBase, id, 2)

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

	id := strconv.Itoa(user.Id)
	u.RustCreate(&UserAuth, &UserBase, id, 3)

}

func (u UserController) RustCreate(UserAuth *models.UserAuth, UserBase *models.UserBase, uid string, types int) {
	// 插入数据库
	GetAuth := []models.UserAuth{}
	err := initialize.DB.Raw("select * from user_auth where identity_type = ? and identifier = ? limit 100", types, uid).Scan(&GetAuth).Error
	if err != nil {
		logs.Error("用户查询auth报错---系统错误", err)
		u.Fail("用户查询auth报错---系统错误", 500)
		return
	}

	switch len(GetAuth) {
	case 0:
		err = initialize.DB.Create(&UserBase).Error
		if err != nil {
			logs.Error("插入UserBase失败---系统错误", err)
			u.Fail("插入UserBase失败", 500)
			return
		}
		err = initialize.DB.Create(&UserAuth).Error
		if err != nil {
			logs.Error("插入UserAuth失败---系统错误", err)
			u.Fail("插入UserAuth失败", 500)
			return
		}

		token, err := v1.RustCreateToken(UserBase.Uid)
		if err != nil {
			logs.Error("token生成失败")
			u.Fail("token生成失败", 500)
			return
		}
		userinfo, err := v1.GetUserInfo(UserBase.Uid)
		if err != nil {
			logs.Error("获取用户信息失败")
			u.Fail("获取用户信息失败", 500)
			return
		}
		err = u.SetSession(token, userinfo)
		if err != nil {
			logs.Error("SetSession失败")
			u.Fail("SetSession失败", 500)
			return
		}
		u.Ok(map[string]interface{}{
			"msg":   "注册成功",
			"token": token,
			"data":  userinfo,
		})
		break
	case 1:
		token, err := v1.RustCreateToken(GetAuth[0].Uid)
		if err != nil {
			logs.Error("token生成失败")
			u.Fail("token生成失败", 500)
			return
		}

		userinfo, err := v1.GetUserInfo(GetAuth[0].Uid)
		if err != nil {
			logs.Error("获取用户信息失败")
			u.Fail("获取用户信息失败", 500)
			return
		}
		err = u.SetSession(token, userinfo)
		if err != nil {
			logs.Error("SetSession失败")
			u.Fail("SetSession失败", 500)
			return
		}
		u.Ok(map[string]interface{}{
			"msg":   "登录成功",
			"token": token,
			"data":  userinfo,
		})

	default:
		u.Fail("用户存在多个账号", 500)
	}
}
