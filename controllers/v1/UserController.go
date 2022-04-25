package v1

import (
	"dbsgw_rust_server/controllers"
	"dbsgw_rust_server/models"
	"github.com/beego/beego/v2/core/validation"
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
	username := u.GetString("email")
	//password := u.GetString("password")
	valid := validation.Validation{}
	valid.Required(username, "email")
	if valid.HasErrors() {
		// 如果有错误信息，证明验证没通过
		// 打印错误信息
		for _, err := range valid.Errors {
			log.Println(err.Key, err.Message)
		}
		u.Fail("表单验证错误", 300)

	} else {
		u.Ok("成功")
	}
}

// @Title GetAll
// @Description get all Users
// @Success 200 {object} models.User
func (u *UserController) GetAll() {
	//users := models.GetAllUsers()
	//u.Data["json"] = users
	//u.ServeJSON()

	u.Fail("错误", 123)
}

// @Title Delete
// @Description delete the user
// @Param	uid		path 	string	true		"The uid you want to delete"
// @Success 200 {string} delete success!
// @Failure 403 uid is empty
func (u *UserController) Delete() {
	uid := u.GetString(":uid")
	models.DeleteUser(uid)
	u.Data["json"] = "delete success!"
	u.ServeJSON()
}
