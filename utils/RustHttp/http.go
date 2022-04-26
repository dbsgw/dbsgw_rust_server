package RustHttp

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
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

// PostJson 真对githup 设置的 Accept 请求头 返回json数据的
func PostJson(requestUrl string, form url.Values) string {
	//对form进行编码
	data := bytes.NewBufferString(form.Encode())
	req, err := http.NewRequest("POST", requestUrl, data)
	// req.Header.Set("X-Custom-Header", "myvalue")
	req.Header.Set("Accept", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}

	defer resp.Body.Close()
	body, errAll := ioutil.ReadAll(resp.Body)
	if errAll != nil {
		panic(errAll)
	}
	return string(body)
}

func GetJson(requestUrl, token string) string {

	req, err := http.NewRequest("GET", requestUrl, nil)
	if err != nil {
		log.Print(err)
		os.Exit(1)
	}

	q := req.URL.Query()
	req.Header.Set("Accept", "application/json")
	//Authorization: token 361507da
	req.Header.Set("Authorization", "token "+token)
	req.URL.RawQuery = q.Encode()
	fmt.Println(req.URL.String())

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Print(err)
	}
	defer resp.Body.Close()
	body, errAll := ioutil.ReadAll(resp.Body)
	if errAll != nil {
		panic(errAll)
	}
	return string(body)
}
