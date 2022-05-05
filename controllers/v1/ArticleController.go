package v1

import (
	"dbsgw_rust_server/controllers"
	"dbsgw_rust_server/initialize"
	"dbsgw_rust_server/models"
	"dbsgw_rust_server/utils"
	"github.com/beego/beego/v2/core/logs"
	gonanoid "github.com/matoous/go-nanoid/v2"
)

// ArticleController Operations about Users
type ArticleController struct {
	controllers.BaseController
}

// ArticleAll 查全部文章
func (u *ArticleController) ArticleAll() {
	GetArticle := []models.Article{}
	err := initialize.DB.Find(&GetArticle).Limit(100).Error
	if err != nil {
		logs.Error("用户查询article报错---系统错误", err)
		u.Fail("用户查询article报错---系统错误", 500)
		return
	}
	u.Ok(GetArticle)
}

// ArticleId 通过id查询文章
func (u *ArticleController) ArticleId() {
	id := u.Ctx.Input.Param(":id")
	GetArticle := models.Article{}
	err := initialize.DB.Find(&GetArticle, "article_id = ?", id).Limit(100).Error
	if err != nil {
		logs.Error("用户查询article报错---系统错误", err)
		u.Fail("用户查询article报错---系统错误", 500)
		return
	}
	u.Ok(GetArticle)
}

// ArticleCreate 创建文章
func (u *ArticleController) ArticleCreate() {
	uid, _ := gonanoid.New()
	article := models.Article{
		ArticleId:      uid,
		AcId:           0,
		ArticleUrl:     "",
		ArticleShow:    0,
		ArticleSort:    0,
		ArticleTitle:   "",
		ArticleContent: "",
		ArticleTime:    int(utils.GetUnix()),
		ArticlePic:     "",
	}

	err := initialize.DB.Create(&article).Error
	if err != nil {
		logs.Error("用户插入article报错---系统错误", err)
		u.Fail("用户插入article报错---系统错误", 500)
		return
	}
	u.Ok("ok")
}

// ArticleDelete 删除文章
func (u *ArticleController) ArticleDelete() {
	id := u.Ctx.Input.Param(":id")
	article := models.Article{}
	err := initialize.DB.Where("article_id = ?", id).Delete(&article).Error
	if err != nil {
		logs.Error("用户删除article报错---系统错误", err)
		u.Fail("用户删除article报错---系统错误", 500)
		return
	}
	u.Ok("get")
}

// ArticleEdit 修改文章
func (u *ArticleController) ArticleEdit() {
	id := u.Ctx.Input.Param(":id")
	article := models.Article{}
	err := initialize.DB.Model(&article).Where("article_id = ?", id).Updates(models.Article{
		ArticleTime: int(utils.GetUnix()),
	}).Error
	if err != nil {
		logs.Error("用户更新article报错---系统错误", err)
		u.Fail("用户更新article报错---系统错误", 500)
		return
	}
	u.Ok("get")
}
