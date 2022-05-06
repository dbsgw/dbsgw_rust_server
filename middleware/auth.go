package middleware

import (
	"encoding/json"
	"fmt"
	"github.com/beego/beego/v2/server/web/context"
	"net/url"
	"strings"
)

func Auth(ctx *context.Context) {
	pathname := ctx.Request.URL.String()
	urlPath, _ := url.Parse(pathname) //urlPath.Path  /role/edit

	// 包含/admin路径的时候  就经行 验证token
	StrContainers := strings.Contains(urlPath.Path, "/admin")
	if StrContainers {
		fmt.Println(1111)
		toekn := ctx.Request.Header["Authorization"]

		fmt.Println(len(toekn), toekn)
		if len(toekn) == 0 {
			data := map[string]interface{}{
				"code": 401,
				"data": false,
				"msg":  "未登录/登录过期",
			}
			str, _ := json.Marshal(data)
			ctx.Output.SetStatus(401)
			ctx.Output.Body(str)
			return
		}
	}

	//fmt.Println("2222222222222------11111", toekn)
	//if string(urlPath.Path) != "/admin/login" {
	//	cookie := ctx.Input.Cookie("token")
	//
	//	Secrect, _ := beego.AppConfig.String("Secrect")
	//	_, err := utils.ParseToken(cookie, Secrect)
	//	if err != nil {
	//		data := map[string]interface{}{
	//			"code": 50014,
	//			"data": false,
	//			"msg":  "cookie登录失败",
	//		}
	//		//str, _ := json.Marshal(data)
	//		//ctx.Output.Body(str)
	//		fmt.Println(data)
	//		ctx.SetCookie("token", "", -1)
	//		ctx.Redirect(302, "/")
	//	}
	//}
}
