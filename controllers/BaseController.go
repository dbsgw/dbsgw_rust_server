package controllers

import beego "github.com/beego/beego/v2/server/web"

type BaseController struct {
	beego.Controller
}

type any = interface{}

type Result struct {
	Data   interface{} `json:"data"`
	Msg    string      `json:"msg"`
	Status int         `json:"status"`
}

// Ok 成功
func (c *BaseController) Ok(data any) {
	c.Data["json"] = SuccessData(data)
	c.ServeJSON()
}

// Fail 失败
func (c *BaseController) Fail(msg string, status int) {
	c.Data["json"] = ErrMsg(msg, status)
	c.ServeJSON()
}

func ErrMsg(msg string, status ...int) Result {
	var r Result
	if len(status) > 0 {
		r.Status = status[0]
	} else {
		r.Status = 500
	}
	r.Msg = msg
	r.Data = nil

	return r
}

func SuccessData(data any) Result {
	var r Result

	r.Status = 200
	r.Msg = "ok"
	r.Data = data

	return r
}
