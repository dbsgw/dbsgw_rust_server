package v1

import (
	"dbsgw_rust_server/controllers"
	"dbsgw_rust_server/initialize"
	"dbsgw_rust_server/models"
	"dbsgw_rust_server/utils"
	"encoding/json"
	"fmt"
	"github.com/beego/beego/v2/core/logs"
	gonanoid "github.com/matoous/go-nanoid/v2"
)

// CommentController  Operations about Users
type CommentController struct {
	controllers.BaseController
}

// CommentAll  查全部评论
func (u *CommentController) CommentAll() {
	articleID := u.GetString("articleID")
	comment := []models.ArticleComment{}
	err := initialize.DB.Where("topic_id = ?", articleID).Find(&comment).Error
	if err != nil {
		logs.Info(err)
		u.Fail("查询评论失败", 500)
		return
	}
	var count int64
	initialize.DB.Table("article_comment").Where("topic_id = ?", articleID).Count(&count)
	u.Ok(map[string]interface{}{
		"result": comment,
		"total":  count,
	})
}

type RustComment struct {
	Comment   string `json:"comment"`
	ArticleID string `json:"articleID"`
}

// CommentCreate   添加评论
func (u *CommentController) CommentCreate() {
	result := RustComment{}
	err := json.Unmarshal(u.Ctx.Input.RequestBody, &result)
	fmt.Println(result, string(u.Ctx.Input.RequestBody))
	if err != nil {
		logs.Info("创建文章转换json失败", err)
	}

	token := u.Ctx.Request.Header["Authorization"]
	if len(token) != 1 {
		u.Fail("获取用户token失败", 500)
		logs.Info("获取用户token失败")
		return
	}
	sessionResult := u.GetSession(token[0])
	userInfo := sessionResult.(models.UserBase)
	CommentId, _ := gonanoid.New()

	comment := models.ArticleComment{
		CommentId:    CommentId,
		TopicId:      result.ArticleID,
		TopicType:    "",
		Content:      result.Comment,
		FromUid:      userInfo.Uid,
		FromNickname: userInfo.NickName,
		FromAvatar:   userInfo.Face,
		Time:         int(utils.GetUnix()),
	}

	err = initialize.DB.Create(&comment).Error
	if err != nil {
		logs.Info("插入评论失败")
		u.Fail("插入评论失败", 500)
		return
	}
	u.Ok("评论成功")
}
