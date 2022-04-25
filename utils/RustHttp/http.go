package RustHttp

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"net/url"
)

// 直接拼rul 路径

func Get(requestUrl string) string {
	client := &http.Client{}
	resp, errGet := client.Get(requestUrl)
	if errGet != nil {
		panic(errGet)
	}
	defer resp.Body.Close()
	body, errAll := ioutil.ReadAll(resp.Body)
	if errAll != nil {
		panic(errAll)
	}
	return string(body)
}

/**
requestUrl 请求url  // https://gitee.com/oauth/token
form  为参数：
	form := url.Values{
		"grant_type":    {"authorization_code"},
		"code":          {code},
		"client_id":     {config.InitConfig.ClientID},
		"redirect_uri":  {config.InitConfig.RedirectUrl},
		"client_secret": {config.InitConfig.ClientSecret},
	}
*/

func Post(requestUrl string, form url.Values) string {
	//对form进行编码
	data := bytes.NewBufferString(form.Encode())
	rsp, errPost := http.Post(requestUrl, "application/x-www-form-urlencoded", data)
	if errPost != nil {
		panic(errPost)
	}
	defer rsp.Body.Close()
	body, errAll := ioutil.ReadAll(rsp.Body)
	if errAll != nil {
		panic(errAll)
	}
	return string(body)
}
