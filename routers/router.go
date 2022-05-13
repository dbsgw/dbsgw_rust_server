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
	"dbsgw_rust_server/models"
	"encoding/json"
	"fmt"
	beego "github.com/beego/beego/v2/server/web"
	"github.com/beego/beego/v2/server/web/context"
)

func init() {
	ns := beego.NewNamespace("/v1",
		//beego.NSBefore(func(ctx *context.Context) {
		//	token := ctx.Request.Header["Authorization"]
		//	fmt.Println(len(token), token)
		//	// token 返回是 字符串切片  只能用数组判断
		//	if len(token) != 0 {
		//		// token 有值
		//		userinfo := ctx.Input.Session(token[0])
		//		user, ok := userinfo.(models.UserBase)
		//		if !ok {
		//			fmt.Println(user, ok)
		//			data := map[string]interface{}{
		//				"code": 401,
		//				"data": false,
		//				"msg":  "未登录/登录过期",
		//			}
		//			str, _ := json.Marshal(data)
		//			ctx.Output.SetStatus(401)
		//			ctx.Output.Body(str)
		//			return
		//		} else {
		//			fmt.Println(user.Uid, "--********--")
		//		}
		//	}
		//}),
		//beego.NSCond(func(ctx *context.Context) bool {
		//	token := ctx.Request.Header["Authorization"]
		//	fmt.Println(len(token), token)
		//	// token 返回是 字符串切片  只能用数组判断
		//	if len(token) != 0 {
		//		// token 有值
		//		userinfo := ctx.Input.Session(token[0])
		//		_, ok := userinfo.(models.UserBase)
		//		if !ok {
		//			data := map[string]interface{}{
		//				"code": 401,
		//				"data": false,
		//				"msg":  "未登录/登录过期",
		//			}
		//			str, _ := json.Marshal(data)
		//			ctx.Output.SetStatus(401)
		//			ctx.Output.Body(str)
		//			return true
		//		}
		//	}
		//	return false
		//}),
		// 登录/注册
		beego.NSNamespace("/user",

			beego.NSRouter("/info/:id", &v1.UserController{}, "get:Info"),    // 获取用户详情
			beego.NSRouter("/info/:id", &v1.UserController{}, "put:InfoPut"), // 更新用户详情

			beego.NSRouter("/login", &v1.UserController{}, "get:Login"),                    // 登录接口
			beego.NSRouter("/login/code", &v1.UserController{}, "get:Code"),                // 登录邮箱验证码
			beego.NSRouter("/login/oauth/gitee", &v1.UserController{}, "get:OauthGitee"),   // gitee回调接口
			beego.NSRouter("/login/oauth/githup", &v1.UserController{}, "get:OauthGitHup"), // githup回调接口

			beego.NSRouter("/article", &v1.UserController{}, "get:ArticleAll"),    // 获取所以文章 带分页
			beego.NSRouter("/article/:id", &v1.UserController{}, "get:ArticleId"), // 通过id获取文章

			beego.NSRouter("/logout", &v1.UserController{}, "get:Logout"), // 退出登录
		),
		//beego.NSBefore(middleware.Auth),

		// 权限页面
		beego.NSNamespace("/admin",
			beego.NSBefore(func(ctx *context.Context) {
				token := ctx.Request.Header["Authorization"]
				fmt.Println(len(token), token, "-------------------------------")
				// token 返回是 字符串切片  只能用数组判断
				if len(token) != 0 {
					// token 有值
					userinfo := ctx.Input.Session(token[0])
					user, ok := userinfo.(models.UserBase)
					if !ok {
						fmt.Println(user, ok)
						data := map[string]interface{}{
							"code": 401,
							"data": false,
							"msg":  "未登录/登录过期",
						}
						str, _ := json.Marshal(data)
						ctx.Output.SetStatus(401)
						ctx.Output.Body(str)
						return
					} else {
						fmt.Println(user.Uid, "--********--")
					}
				}
			}),
			beego.NSRouter("/article/:id", &v1.ArticleController{}, "get:ArticleId"),        // 通过id获取文章
			beego.NSRouter("/article/", &v1.ArticleController{}, "get:ArticleAll"),          // 获取所以文章 带分页
			beego.NSRouter("/article/", &v1.ArticleController{}, "post:ArticleCreate"),      // 创建文章
			beego.NSRouter("/article/:id", &v1.ArticleController{}, "delete:ArticleDelete"), // 删除文章
			beego.NSRouter("/article/:id", &v1.ArticleController{}, "put:ArticleEdit"),      // 修改文章
		),

		// 留言/回复
		beego.NSRouter("/comment", &v1.CommentController{}, "get:CommentAll"),     // 查看全部回复
		beego.NSRouter("/comment", &v1.CommentController{}, "post:CommentCreate"), // 添加评论
	)
	beego.AddNamespace(ns)
}
