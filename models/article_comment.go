package models

// ArticleComment 评论表
type ArticleComment struct {
	CommentId    string `json:"comment_id"`    // '评论人的id,这个表的id',
	TopicId      string `json:"topic_id"`      // '被评论文章id',
	TopicType    string `json:"topic_type"`    // '被评论文章类型',
	Content      string `json:"content"`       // '评论内容',
	FromUid      string `json:"from_uid"`      // '评论者id',
	FromNickname string `json:"from_nickname"` // '评论者昵称',
	FromAvatar   string `json:"from_avatar"`   // '评论者头像',
	Time         string `json:"time"`          // '评论的时间'
}

// TableName 自定义表名
func (ArticleComment) TableName() string {
	return "article_comment"
}
