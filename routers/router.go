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
	"dbsgw_rust_server/middleware"
	"fmt"
	beego "github.com/beego/beego/v2/server/web"
	"github.com/beego/beego/v2/server/web/context"
)

func init() {
	ns := beego.NewNamespace("/v1",
		beego.NSCond(func(ctx *context.Context) bool {
			fmt.Println(ctx.Input.Domain())
			if ctx.Input.Domain() == "127.0.0.1" {
				return true
			}
			return false
		}),
		beego.NSBefore(middleware.Auth),
		// 登录/注册
		beego.NSNamespace("/user",
			beego.NSRouter("/login", &v1.UserController{}, "get:Login"),                    // 登录接口
			beego.NSRouter("/login/code", &v1.UserController{}, "get:Code"),                // 登录邮箱验证码
			beego.NSRouter("/login/oauth/gitee", &v1.UserController{}, "get:OauthGitee"),   // gitee回调接口
			beego.NSRouter("/login/oauth/githup", &v1.UserController{}, "get:OauthGitHup"), // githup回调接口
		),
	)
	beego.AddNamespace(ns)
}
