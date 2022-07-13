// Package routers @APIVersion 1.0.0
// @Title beego Test API
// @Description beego has a very cool tools to autogenerate documents for your API
// @Contact astaxie@gmail.com
// @TermsOfServiceUrl http://beego.me/
// @License Apache 2.0
// @LicenseUrl http://www.apache.org/licenses/LICENSE-2.0.html
package routers

import (
	UserArticle "dbsgw_rust_server/controllers/v1/user/article"
	UserComment "dbsgw_rust_server/controllers/v1/user/comment"
	UserInfo "dbsgw_rust_server/controllers/v1/user/info"
	WebArticle "dbsgw_rust_server/controllers/v1/web/article"
	WebComment "dbsgw_rust_server/controllers/v1/web/comment"
	WebLogin "dbsgw_rust_server/controllers/v1/web/login"
	beego "github.com/beego/beego/v2/server/web"
)

func init() {
	ns := beego.NewNamespace("/v1",
		beego.NSRouter("/login", &WebLogin.LoginController{}, "get:Login"),                    // 登录接口
		beego.NSRouter("/login/code", &WebLogin.LoginController{}, "get:Code"),                // 登录邮箱验证码
		beego.NSRouter("/login/oauth/gitee", &WebLogin.LoginController{}, "get:OauthGitee"),   // gitee回调接口
		beego.NSRouter("/login/oauth/githup", &WebLogin.LoginController{}, "get:OauthGitHup"), // githup回调接口
		beego.NSRouter("/login/logout", &WebLogin.LoginController{}, "get:Logout"),            // 退出登录

		beego.NSRouter("/article", &WebArticle.ArticleController{}, "get:ArticleAll"),    // 获取所以文章 带分页
		beego.NSRouter("/article/:id", &WebArticle.ArticleController{}, "get:ArticleId"), // 通过id获取文章

		// 留言/回复
		beego.NSRouter("/comment", &WebComment.CommentController{}, "get:CommentAll"), // 查看全部回复

		// 个人中心 需要登录访问的权限
		beego.NSNamespace("/user",

			beego.NSRouter("/info/:id", &UserInfo.InfoController{}, "get:Info"),    // 获取用户详情
			beego.NSRouter("/info/:id", &UserInfo.InfoController{}, "put:InfoPut"), // 更新用户详情

			beego.NSRouter("/article", &UserArticle.ArticleController{}, "get:ArticleAll"),    // 获取所以文章 带分页
			beego.NSRouter("/article/:id", &UserArticle.ArticleController{}, "get:ArticleId"), // 通过id获取文章

			beego.NSRouter("/comment", &UserComment.CommentController{}, "post:CommentCreate"), // 添加评论

		),
		//beego.NSBefore(middleware.Auth),

		// 后台管理系统 ，项目管理的
		beego.NSNamespace("/admin",
			// TODO
			beego.NSRouter("/article", &UserArticle.ArticleController{}, "get:ArticleAll"), // 获取所以文章 带分页
		),
	)
	beego.AddNamespace(ns)
}
