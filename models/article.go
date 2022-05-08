package models

// Article 文章表
type Article struct {
	Id               int    `json:"id" orm:"-"`         // id 主键
	ArticleId        string `json:"article_id"`         // '文章自增ID',
	AcId             int    `json:"ac_id"`              // '分类id',
	ArticleUrl       string `json:"article___url"`      // '文章跳转链接',
	ArticleShow      int    `json:"article_show"`       // '文章是否显示，0为否，1为是，默认为1',
	ArticleSort      int    `json:"article_sort"`       // '文章排序',
	ArticleTitle     string `json:"article_title"`      // '文章标题',
	ArticleContent   string `json:"article_content"`    // '内容',
	ArticleContentMk string `json:"article_content_mk"` // 'mk内容',
	ArticleTime      int    `json:"article_time"`       // '文章发布时间',
	ArticlePic       string `json:"article_pic"`        // '文章主图'
	UserId           string `json:"user_id"`            // '用户id'
	UserName         string `json:"user_name"`          // '用户昵称'
}

// TableName 自定义表名
func (Article) TableName() string {
	return "article"
}
