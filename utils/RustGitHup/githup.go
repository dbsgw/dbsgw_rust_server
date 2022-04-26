package RustGitHup

import (
	"dbsgw_rust_server/utils/RustHttp"
	"encoding/json"
	"fmt"
	beego "github.com/beego/beego/v2/server/web"
	"net/url"
	"time"
)

//https://gitee.com/api/v5/oauth_doc#/

// token 信息

type Token struct {
	AccessToken           string `json:"access_token"`
	TokenType             string `json:"token_type"`
	Scope                 string `json:"scope"`
	ExpiresIn             string `json:"expires_in"`
	RefreshToken          string `json:"refresh_token"`
	RefreshTokenExpiresIn string `json:"refresh_token_expires_in"`
}

// 用户信息
type UserInfo struct {
	Login                   string    `json:"login"`
	Id                      int       `json:"id"`
	NodeId                  string    `json:"node_id"`
	AvatarUrl               string    `json:"avatar_url"`
	GravatarId              string    `json:"gravatar_id"`
	Url                     string    `json:"url"`
	HtmlUrl                 string    `json:"html_url"`
	FollowersUrl            string    `json:"followers_url"`
	FollowingUrl            string    `json:"following_url"`
	GistsUrl                string    `json:"gists_url"`
	StarredUrl              string    `json:"starred_url"`
	SubscriptionsUrl        string    `json:"subscriptions_url"`
	OrganizationsUrl        string    `json:"organizations_url"`
	ReposUrl                string    `json:"repos_url"`
	EventsUrl               string    `json:"events_url"`
	ReceivedEventsUrl       string    `json:"received_events_url"`
	Type                    string    `json:"type"`
	SiteAdmin               bool      `json:"site_admin"`
	Name                    string    `json:"name"`
	Company                 string    `json:"company"`
	Blog                    string    `json:"blog"`
	Location                string    `json:"location"`
	Email                   string    `json:"email"`
	Hireable                bool      `json:"hireable"`
	Bio                     string    `json:"bio"`
	TwitterUsername         string    `json:"twitter_username"`
	PublicRepos             int       `json:"public_repos"`
	PublicGists             int       `json:"public_gists"`
	Followers               int       `json:"followers"`
	Following               int       `json:"following"`
	CreatedAt               time.Time `json:"created_at"`
	UpdatedAt               time.Time `json:"updated_at"`
	PrivateGists            int       `json:"private_gists"`
	TotalPrivateRepos       int       `json:"total_private_repos"`
	OwnedPrivateRepos       int       `json:"owned_private_repos"`
	DiskUsage               int       `json:"disk_usage"`
	Collaborators           int       `json:"collaborators"`
	TwoFactorAuthentication bool      `json:"two_factor_authentication"`
	Plan                    struct {
		Name          string `json:"name"`
		Space         int    `json:"space"`
		PrivateRepos  int    `json:"private_repos"`
		Collaborators int    `json:"collaborators"`
	} `json:"plan"`
}

// 邮箱

type Email struct {
	Email string `json:"email"`
}

//  获取code  get	https://gitee.com/oauth/authorize?client_id=eb30f085980d8fea35284e6923a4c9213393c27172ce1a00df7ca0a88d7de9dd&redirect_uri=http://127.0.0.1:3000/login/oauth&response_type=code
func RedirectUrl() string {

	ClientID, _ := beego.AppConfig.String("GitHupClientID")
	RedirectUrl, _ := beego.AppConfig.String("GitHupRedirectUrl")
	return "https://github.com/login/oauth/authorize?client_id=" + ClientID + "&redirect_uri=" + RedirectUrl
}

// 获取 access_token post  https://gitee.com/oauth/token

func GetAccessToken(code string) Token {
	ClientID, _ := beego.AppConfig.String("GitHupClientID")
	RedirectUrl, _ := beego.AppConfig.String("GitHupRedirectUrl")
	ClientSecret, _ := beego.AppConfig.String("GitHupClientSecret")
	// 要 POST的 参数
	form := url.Values{
		"code":          {code},
		"client_id":     {ClientID},
		"redirect_uri":  {RedirectUrl},
		"client_secret": {ClientSecret},
	}
	userToken := Token{}
	str := RustHttp.PostJson("https://github.com/login/oauth/access_token", form)
	err := json.Unmarshal([]byte(str), &userToken)
	if err != nil {
		fmt.Println(err, "err------------")
	}
	return userToken
}

// 获取用户资料 get https://gitee.com/api/v5/user?access_token=1256465456444545

func GetUserInfos(userToken Token) UserInfo {
	userInfos := UserInfo{}
	data := RustHttp.GetJson("https://api.github.com/user", userToken.AccessToken)
	json.Unmarshal([]byte(data), &userInfos)
	return userInfos
}

// 获取用户邮箱 get  https://gitee.com/api/v5/emails?access_token=bdda2be031f36dccb82b8927b956c0a8

func GetEmails(userToken Token) string {
	data := RustHttp.GetJson("https://api.github.com/user/email/visibility", userToken.AccessToken)
	fmt.Println(data, "email----")
	return data
}
