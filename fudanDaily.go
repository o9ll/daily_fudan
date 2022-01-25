/**
 * @Author Oliver
 * @Date 1/25/22
 **/

package main

import (
	"fmt"
	"github.com/antchfx/htmlquery"
	"github.com/kirinlabs/HttpRequest"
	"net/http"
	"strings"
)

var (
	client        = &http.Client{}
	fudanDailyUrl = "https://zlapp.fudan.edu.cn/site/ncov/fudanDaily"
	loginUrl      = "https://uis.fudan.edu.cn/authserver/login?service=https%3A%2F%2Fzlapp.fudan.edu.cn%2Fa_fudanzlapp%2Fapi%2Fsso%2Findex%3Fredirect%3Dhttps%253A%252F%252Fzlapp.fudan.edu.cn%252Fsite%252Fncov%252FfudanDaily%26from%3Dwap"
	getInfoUrl    = "https://zlapp.fudan.edu.cn/ncov/wap/fudan/get-info"
	saveLogUrl    = "https://zlapp.fudan.edu.cn/wap/log/save-log"
	saveUrl       = "https://zlapp.fudan.edu.cn/ncov/wap/fudan/save"
	userAgent     = "Mozilla/5.0 (iPhone; CPU iPhone OS 14_2 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Mobile/15E148 MicroMessenger/7.0.18(0x17001229) NetType/WIFI Language/zh_CN miniProgram"
	origin        = "https://zlapp.fudan.edu.cn"
	Referer       = fudanDailyUrl
)

type userInfo struct {
	Username string
	Password string
	Email    string
}

func setHeader(r *http.Request) {
	r.Header.Add("User-Agent", userAgent)
	r.Header.Add("Origin", origin)
	r.Header.Add("Referer", Referer)
}

func login(info userInfo) *HttpRequest.Request {
	req := HttpRequest.NewRequest()
	req.SetHeaders(map[string]string{
		"User-Agent": userAgent,
		"Origin":     origin,
		"Referer":    Referer,
	})
	resp, _ := req.Get(loginUrl)
	uv := ""
	body, _ := resp.Body()
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
	resp, _ = req.Post(loginUrl, uv)
	body, _ = resp.Body()
	fmt.Println(string(body))
	return req
}

func main() {
	login(userInfo{
		Username: "20210240194",
		Password: "Liu159632",
	})
}
