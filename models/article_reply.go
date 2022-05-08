package models

// ArticleReply 回复表
type ArticleReply struct {
	CommentId    string `json:"comment_id"`    // '被回复评论的id',
	ReplyId      string `json:"reply_id"`      // '回复id【unique唯一值】,这个表的id',
	TopicId      string `json:"topic_id"`      // '被回复文章id',
	TopicType    string `json:"topic_type"`    // '被回复文章类型',
	Content      string `json:"content"`       // '回复内容',
	FromUid      string `json:"from_uid"`      // '回复者uid',
	ToUid        string `json:"to_uid"`        // '被回复者uid',
	FromNickname string `json:"from_nickname"` // '回复者昵称',
	ToNickname   string `json:"to_nickname"`   // '被回复者昵称',
	FromAvatar   string `json:"from_avatar"`   // '回复者头像',
	ToAvatar     string `json:"to_avatar"`     // '被回复者头像',
	time         int    `json:"time"`          // '回复的时间',
	flag         int    `json:"flag"`          // '可选值，用于回复时是否要显示对谁的回复'
}

// TableName 自定义表名
func (ArticleReply) TableName() string {
	return "article_reply"
}
