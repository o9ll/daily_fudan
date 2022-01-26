/**
 * @Author Oliver
 * @Date 1/26/22
 **/

package baiduAPI

import (
	. "daily_fudan/util"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

var (
	file        = "api.json"
	url         = "https://aip.baidubce.com/oauth/2.0/token"
	apiKey      = "API_key"
	secretKey   = "secret_key"
	accessToken = "access_token"
)

type API struct {
	API_key    string `json:"API_key"`
	Secret_key string `json:"secret_key"`
}

func createAPIJson(src string) {
	res := &API{}
	fmt.Println(`请输入API_key`)
	fmt.Scanln(&res.API_key)
	fmt.Println(`请输入secret_key`)
	fmt.Scanln(&res.Secret_key)
	data, err := json.MarshalIndent(res, "", "	") // 第二个表示每行的前缀，这里不用，第三个是缩进符号，这里用tab
	CheckError(err)
	err = ioutil.WriteFile(src, data, 0777)
	CheckError(err)
}

func getAPI() map[string]interface{} {
	api := ReadFromJson(file)
	if api == nil {
		createAPIJson(file)
		api = ReadFromJson(file)
	}
	return api
}

func getAccessToken() string {
	data, _ := ioutil.ReadFile(file)
	if data == nil {
		createAPIJson(file)
	}
	api := getAPI()
	resp, _ := http.Get(url + "?grant_type=client_credentials&client_id=" + api[apiKey].(string) + "&client_secret=" + api[secretKey].(string))
	res, _ := ioutil.ReadAll(resp.Body)
	mp := Json2Map(res)
	return mp[accessToken].(string)
}
