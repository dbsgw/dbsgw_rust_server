package RustGitee

import (
	"dbsgw_rust_server/utils/RustHttp"
	"encoding/json"
	beego "github.com/beego/beego/v2/server/web"
	"net/url"
	"time"
)

//https://gitee.com/api/v5/oauth_doc#/

// token 信息

type Token struct {
	AccessToken  string `json:"access_token"`
	TokenType    string `json:"token_type"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    string `json:"expires_in"`
	Scope        string `json:"scope"`
	CreatedAt    string `json:"created_at"`
}

// 用户信息

type UserInfo struct {
	Id                int       `json:"id"`
	Login             string    `json:"login"`
	Name              string    `json:"name"`
	AvatarUrl         string    `json:"avatar_url"`
	Url               string    `json:"url"`
	HtmlUrl           string    `json:"html_url"`
	Remark            string    `json:"remark"`
	FollowersUrl      string    `json:"followers_url"`
	FollowingUrl      string    `json:"following_url"`
	GistsUrl          string    `json:"gists_url"`
	StarredUrl        string    `json:"starred_url"`
	SubscriptionsUrl  string    `json:"subscriptions_url"`
	OrganizationsUrl  string    `json:"organizations_url"`
	ReposUrl          string    `json:"repos_url"`
	EventsUrl         string    `json:"events_url"`
	ReceivedEventsUrl string    `json:"received_events_url"`
	Types             string    `json:"type"`
	Blog              string    `json:"blog"`
	CreatedAt         time.Time `json:"created_at"`
	UpdatedAt         time.Time `json:"updated_at"`
	Email             string    `json:"email"`
	Bio               string    `json:"bio"`
	Company           string    `json:"company"`
	Location          string    `json:"location"`
}

// 邮箱

type Email struct {
	Email string `json:"email"`
}

//  获取code  get	https://gitee.com/oauth/authorize?client_id=eb30f085980d8fea35284e6923a4c9213393c27172ce1a00df7ca0a88d7de9dd&redirect_uri=http://127.0.0.1:3000/login/oauth&response_type=code
func RedirectUrl() string {

	ClientID, _ := beego.AppConfig.String("ClientID")
	RedirectUrl, _ := beego.AppConfig.String("RedirectUrl")
	return "https://gitee.com/oauth/authorize?client_id=" + ClientID + "&redirect_uri=" + RedirectUrl + "&response_type=code"
}

// 获取 access_token post  https://gitee.com/oauth/token

func GetAccessToken(code string) Token {
	ClientID, _ := beego.AppConfig.String("ClientID")
	RedirectUrl, _ := beego.AppConfig.String("RedirectUrl")
	ClientSecret, _ := beego.AppConfig.String("ClientSecret")
	// 要 POST的 参数
	form := url.Values{
		"grant_type":    {"authorization_code"},
		"code":          {code},
		"client_id":     {ClientID},
		"redirect_uri":  {RedirectUrl},
		"client_secret": {ClientSecret},
	}
	userToken := Token{}
	str := RustHttp.Post("https://gitee.com/oauth/token", form)
	json.Unmarshal([]byte(str), &userToken)
	return userToken
}

// 获取用户资料 get https://gitee.com/api/v5/user?access_token=1256465456444545

func GetUserInfos(userToken Token) UserInfo {
	userInfos := UserInfo{}
	data := RustHttp.Get("https://gitee.com/api/v5/user?access_token=" + userToken.AccessToken)
	json.Unmarshal([]byte(data), &userInfos)
	return userInfos
}

// 获取用户邮箱 get  https://gitee.com/api/v5/emails?access_token=bdda2be031f36dccb82b8927b956c0a8

func GetEmails(userToken Token) string {
	data := RustHttp.Get("https://gitee.com/api/v5/emails?access_token=" + userToken.AccessToken)
	return data
}
