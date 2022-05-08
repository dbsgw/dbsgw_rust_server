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

// ArticleController Operations about Users
type ArticleController struct {
	controllers.BaseController
}
type resultArticleUse struct {
	Title     string `json:"title"`
	Connect   string `json:"connect"`
	ConnectMk string `json:"connect_mk"`
}

// ArticleAll 查全部文章
func (u *ArticleController) ArticleAll() {

	size, _ := u.GetInt("size")
	page, _ := u.GetInt("page")
	if size == 0 {
		size = 20
	}
	if page == 0 {
		page = 1
	}
	var count int64
	GetArticle := []models.Article{}
	err := initialize.DB.Offset((page - 1) * size).Limit(size).Order("article_time desc").Find(&GetArticle).Error
	if err != nil {
		logs.Error("用户查询article报错---系统错误", err)
		u.Fail("用户查询article报错---系统错误", 500)
		return
	}
	initialize.DB.Table("article").Count(&count)
	fmt.Println(count, "-----")
	u.Ok(map[string]interface{}{
		"result": GetArticle,
		"total":  count,
	})
}

// ArticleId 通过id查询文章
func (u *ArticleController) ArticleId() {
	id := u.Ctx.Input.Param(":id")
	GetArticle := models.Article{}
	err := initialize.DB.Limit(100).Find(&GetArticle, "article_id = ?", id).Error
	if err != nil {
		logs.Error("用户查询article报错---系统错误", err)
		u.Fail("用户查询article报错---系统错误", 500)
		return
	}
	u.Ok(GetArticle)
}

// ArticleCreate 创建文章
func (u *ArticleController) ArticleCreate() {

	result := resultArticleUse{}
	err := json.Unmarshal(u.Ctx.Input.RequestBody, &result)
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
	fmt.Println(sessionResult, "------", token)
	uid, _ := gonanoid.New()
	article := models.Article{
		ArticleId:        uid,
		AcId:             0,
		ArticleUrl:       "",
		ArticleShow:      0,
		ArticleSort:      0,
		ArticleTitle:     result.Title,
		ArticleContent:   result.Connect,
		ArticleContentMk: result.ConnectMk,
		ArticleTime:      int(utils.GetUnix()),
		ArticlePic:       "",
		UserId:           sessionResult.(models.UserBase).Uid,
		UserName:         sessionResult.(models.UserBase).NickName,
	}

	err = initialize.DB.Create(&article).Error
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
	result := resultArticleUse{}
	json.Unmarshal(u.Ctx.Input.RequestBody, &result)
	article := models.Article{}
	err := initialize.DB.Model(&article).Where("article_id = ?", id).Updates(models.Article{
		ArticleTitle:   result.Title,
		ArticleContent: result.Connect,
	}).Error
	if err != nil {
		logs.Error("用户更新article报错---系统错误", err)
		u.Fail("用户更新article报错---系统错误", 500)
		return
	}
	u.Ok("get")
}
