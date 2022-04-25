// @APIVersion 1.0.0
// @Title beego Test API
// @Description beego has a very cool tools to autogenerate documents for your API
// @Contact astaxie@gmail.com
// @TermsOfServiceUrl http://beego.me/
// @License Apache 2.0
// @LicenseUrl http://www.apache.org/licenses/LICENSE-2.0.html
package routers

import (
	"dbsgw_rust_server/controllers/v1"

	beego "github.com/beego/beego/v2/server/web"
)

func init() {
	ns := beego.NewNamespace("/v1",
		// 登录/注册
		beego.NSNamespace("/user",
			beego.NSRouter("/", &v1.UserController{}, "get:GetAll"),
			beego.NSRouter("/login", &v1.UserController{}, "get:Login"),
			beego.NSRouter("/login/oauth/gitee", &v1.UserController{}, "get:OauthGitee"),
		),
	)
	beego.AddNamespace(ns)
}
