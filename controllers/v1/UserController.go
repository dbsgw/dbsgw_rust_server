package v1

import (
	"dbsgw_rust_server/controllers"
	"dbsgw_rust_server/models"
	"dbsgw_rust_server/utils"
	"dbsgw_rust_server/utils/RustConstant"
	"dbsgw_rust_server/utils/RustGitHup"
	"dbsgw_rust_server/utils/RustGitee"
	"encoding/json"
	"fmt"
	"github.com/beego/beego/v2/client/orm"
	"github.com/beego/beego/v2/core/validation"
	gonanoid "github.com/matoous/go-nanoid/v2"
	"log"
)

// Operations about Users
type UserController struct {
	controllers.BaseController
}

// @Title 登录
// @Description Logs user into the system
// @Param	username		query 	string	true		"The username for login"
// @Param	password		query 	string	true		"The password for login"
// @Success 200 {string} login success
// @Failure 403 user not exist
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
			fmt.Println("email登录")
			break
		case RustConstant.GITEE:
			fmt.Println("gitee登录", RustGitee.RedirectUrl())
			u.Redirect(RustGitee.RedirectUrl(), 302)
			break
		case RustConstant.GITHUP:
			fmt.Println("githup登录", RustGitHup.RedirectUrl())
			u.Redirect(RustGitHup.RedirectUrl(), 302)
			break
		default:
			u.Fail("登录方式不存在", 500)
		}
		u.Ok("成功")
	}
}

// @Title 授权gitee
// @Description Logs user into the system
// @Param	username		query 	string	true		"The username for login"
// @Param	password		query 	string	true		"The password for login"
// @Success 200 {string} login success
// @Failure 403 user not exist
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
	fmt.Println(user, "user---user")
	bl = utils.IsEmpty(user.Id)
	if !bl {
		u.Fail("user.Id为空---系统错误", 500)
		return
	}
	// 通过token 获取邮箱
	emailsReust := RustGitee.GetEmails(userToken)
	// 序列化  取邮箱
	emailsAll := []map[string]interface{}{}
	json.Unmarshal([]byte(emailsReust), &emailsAll)

	email := ""
	if len(emailsAll) != 0 {
		email = emailsAll[0]["email"].(string)
	}
	fmt.Println(email)

	uid, _ := gonanoid.New()
	UserAuth := models.UserAuth{
		Uid:          uid,
		IdentityType: 0,
		Identifier:   "",
		Certificate:  "",
		CreateTime:   int(utils.GetUnix()),
		UpdateTime:   int(utils.GetUnix()),
	}

	UserBase := models.UserBase{
		Uid:            uid,
		UserRole:       0,
		RegisterSource: 0,
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

	// TODO  需要加限定 是注册  还是 登录，登录是不需要插入的
	// 插入数据库
	o := orm.NewOrm()
	_, err := o.Insert(&UserBase)
	if err == nil {
		fmt.Println(err)
	}
	_, err = o.Insert(&UserAuth)
	if err == nil {
		fmt.Println(err)
	}

	u.Ok(user)

	//// 登录的  uid 也就是  gitee的id  不能为空  为空代表 取消授权，取消登录
	//if string(user.Id) == "" || user.Id == 0 {
	//	u.Ctx.Redirect(301, "/")
	//	return
	//}
	//
	//// 查询 数据库  是否有  用户存在
	//Users := []models.User{}
	//err := utils.DB.Raw("select * from user where uid = ?", user.Id).Find(&Users).Error
	//if err != nil {
	//	panic(err)
	//}
	//
	//// 不存在 插入
	//if len(Users) == 0 {
	//	res := map[string]interface{}{}
	//	err = utils.DB.Raw("insert into user (uid, username, login, email, avatar_url, html_url, created_at, status,addTime) VALUES (?,?,?,?,?,?,?,?)", user.Id, user.Name, user.Login, email, user.AvatarUrl, user.HtmlUrl, user.CreatedAt, utils.GetDate()).Find(&res).Error
	//	if err != nil {
	//		panic(err)
	//	}
	//} else {
	//	// 存在  更新数据
	//	res := map[string]interface{}{}
	//	err = utils.DB.Raw("update user set username = ?,login=?,email=?,avatar_url=?,html_url=?,created_at=?,updateTime=? where uid = ?", user.Name, user.Login, email, user.AvatarUrl, user.HtmlUrl, user.CreatedAt, utils.GetDate(), user.Id).Find(&res).Error
	//	if err != nil {
	//		panic(err)
	//	}
	//}
	//// 创建token 发放到cookie
	//Secrect, _ := beego.AppConfig.String("Secrect")
	//token, _ := utils.CreateToken(fmt.Sprintf("%d", user.Id), Secrect)
	//u.Ctx.SetCookie("token", token)
	//u.Redirect("/", 301)
}

// @Title 授权gitee
// @Description Logs user into the system
// @Param	username		query 	string	true		"The username for login"
// @Param	password		query 	string	true		"The password for login"
// @Success 200 {string} login success
// @Failure 403 user not exist
func (u *UserController) OauthGitHup() {
	// 获取code
	code := u.GetString("code")
	// 通过code 获取token
	userToken := RustGitHup.GetAccessToken(code)
	// 通过token 获取用户信息
	user := RustGitHup.GetUserInfos(userToken)
	// 通过token 获取邮箱
	emailsReust := RustGitHup.GetEmails(userToken)
	// 序列化  取邮箱
	emailsAll := []map[string]interface{}{}
	json.Unmarshal([]byte(emailsReust), &emailsAll)

	email := ""
	if len(emailsAll) != 0 {
		email = emailsAll[0]["email"].(string)
	}
	fmt.Println(email)

	//o := orm.NewOrm()
	//var maps []orm.Params
	//r, _ := o.Raw("show tables").Values(&maps)
	//fmt.Println(r, maps[0]["Tables_in_rust"])

	u.Ok(user)

	//// 登录的  uid 也就是  gitee的id  不能为空  为空代表 取消授权，取消登录
	//if string(user.Id) == "" || user.Id == 0 {
	//	u.Ctx.Redirect(301, "/")
	//	return
	//}
	//
	//// 查询 数据库  是否有  用户存在
	//Users := []models.User{}
	//err := utils.DB.Raw("select * from user where uid = ?", user.Id).Find(&Users).Error
	//if err != nil {
	//	panic(err)
	//}
	//
	//// 不存在 插入
	//if len(Users) == 0 {
	//	res := map[string]interface{}{}
	//	err = utils.DB.Raw("insert into user (uid, username, login, email, avatar_url, html_url, created_at, status,addTime) VALUES (?,?,?,?,?,?,?,?)", user.Id, user.Name, user.Login, email, user.AvatarUrl, user.HtmlUrl, user.CreatedAt, utils.GetDate()).Find(&res).Error
	//	if err != nil {
	//		panic(err)
	//	}
	//} else {
	//	// 存在  更新数据
	//	res := map[string]interface{}{}
	//	err = utils.DB.Raw("update user set username = ?,login=?,email=?,avatar_url=?,html_url=?,created_at=?,updateTime=? where uid = ?", user.Name, user.Login, email, user.AvatarUrl, user.HtmlUrl, user.CreatedAt, utils.GetDate(), user.Id).Find(&res).Error
	//	if err != nil {
	//		panic(err)
	//	}
	//}
	//// 创建token 发放到cookie
	//Secrect, _ := beego.AppConfig.String("Secrect")
	//token, _ := utils.CreateToken(fmt.Sprintf("%d", user.Id), Secrect)
	//u.Ctx.SetCookie("token", token)
	//u.Redirect("/", 301)
}
