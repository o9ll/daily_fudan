/**
 * @Author Oliver
 * @Date 1/25/22
 **/

package main

import (
	"bytes"
	. "daily_fudan/operatejson"
	"github.com/antchfx/htmlquery"
	"io/ioutil"
	"net/http"
	"net/http/cookiejar"
	"strings"
)

var (
	client        *http.Client
	fudanDailyUrl = "https://zlapp.fudan.edu.cn/site/ncov/fudanDaily"
	loginUrl      = "https://uis.fudan.edu.cn/authserver/login?service=https%3A%2F%2Fzlapp.fudan.edu.cn%2Fa_fudanzlapp%2Fapi%2Fsso%2Findex%3Fredirect%3Dhttps%253A%252F%252Fzlapp.fudan.edu.cn%252Fsite%252Fncov%252FfudanDaily%26from%3Dwap"
	getInfoUrl    = "https://zlapp.fudan.edu.cn/ncov/wap/fudan/get-info"
	saveLogUrl    = "https://zlapp.fudan.edu.cn/wap/log/save-log"
	saveUrl       = "https://zlapp.fudan.edu.cn/ncov/wap/fudan/save"
	userAgent     = "Mozilla/5.0 (iPhone; CPU iPhone OS 14_2 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Mobile/15E148 MicroMessenger/7.0.18(0x17001229) NetType/WIFI Language/zh_CN miniProgram"
	origin        = "https://zlapp.fudan.edu.cn"
	Referer       = fudanDailyUrl
	ContentType   = "application/x-www-form-urlencoded"
	gCurCookies   []*http.Cookie
	gCurCookieJar *cookiejar.Jar
)

type userInfo struct {
	Username string
	Password string
	Email    string
}

/*设置请求头*/
func setHeader(r *http.Request) {
	r.Header.Add("User-Agent", userAgent)
	r.Header.Add("Origin", origin)
	r.Header.Add("Referer", Referer)
	r.Header.Add("Content-Type", ContentType)
}

/*初始化client*/
func init() {
	gCurCookieJar, _ = cookiejar.New(nil)
	client = &http.Client{
		CheckRedirect: nil,
		Jar:           gCurCookieJar,
	}
}

/*登陆*/
func login(info userInfo) {
	req, _ := http.NewRequest("GET", loginUrl, nil)
	setHeader(req)
	resp, _ := client.Do(req)
	body, _ := ioutil.ReadAll(resp.Body)
	//找到表单中的所有参数按照默认填入
	uv := ""
	h, _ := htmlquery.Parse(strings.NewReader(string(body)))
	a := htmlquery.Find(h, "//input")
	for i := range a {
		name := htmlquery.SelectAttr(a[i], "name")
		value := htmlquery.SelectAttr(a[i], "value")
		if name != "" && name != "captchaResponse" {
			if name == "username" {
				uv += "&" + name + "=" + info.Username
			} else if name == "password" {
				uv += "&" + name + "=" + info.Password
			} else {
				uv += "&" + name + "=" + value
			}
		}
	}
	uv = uv[1:]
	req, _ = http.NewRequest("POST", loginUrl, bytes.NewReader([]byte(uv)))
	setHeader(req)
	resp, _ = client.Do(req)
	gCurCookies = gCurCookieJar.Cookies(req.URL)
}

func getHistoryInfo() string {
	req, _ := http.NewRequest("GET", getInfoUrl, nil)
	setHeader(req)
	resp, _ := client.Do(req)
	res, _ := ioutil.ReadAll(resp.Body)
	return ReadJson(res)
}

func main() {
	user := userInfo{
		Username: "20210240194",
		Password: "Liu159632",
	}
	login(user)
	history := getHistoryInfo()
	ioutil.WriteFile(user.Username+".json", []byte(history), 0777)
}
