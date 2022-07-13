package v1

import (
	"dbsgw_rust_server/controllers"
	v1 "dbsgw_rust_server/service/v1"
	"encoding/json"
)

// InfoController Operations about Users
type InfoController struct {
	controllers.BaseController
}

// Info 用户详情
func (u *InfoController) Info() {
	id := u.Ctx.Input.Param(":id")
	userinfo, err := v1.GetUserInfo(id)
	if err != nil {
		u.Fail("获取用户信息失败", 500)
		return
	}
	u.Ok(userinfo)
}

// InfoPut 更新详情
func (u *InfoController) InfoPut() {
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
