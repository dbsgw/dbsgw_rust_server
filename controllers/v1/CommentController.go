package v1

import (
	"dbsgw_rust_server/controllers"
)

// CommentController  Operations about Users
type CommentController struct {
	controllers.BaseController
}

// CommentAll  查全部评论
func (u *CommentController) CommentAll() {
	u.Ok("ok")
}

// CommentCreate   添加评论
func (u *CommentController) CommentCreate() {
	u.Ok("ok")
}
