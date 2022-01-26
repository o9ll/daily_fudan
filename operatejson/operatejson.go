/**
 * @Author Oliver
 * @Date 1/26/22
 **/

package operatejson

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
)

func ReadFromJson(src string) map[string]interface{} {
	data, err := ioutil.ReadFile(src)
	CheckError(err)
	res := make(map[string]interface{})
	err = json.Unmarshal(data, &res)
	CheckError(err)
	return res
}

func WriteToJson(src string, res map[string]interface{}) {
	data, err := json.MarshalIndent(res, "", "	") // 第二个表示每行的前缀，这里不用，第三个是缩进符号，这里用tab
	CheckError(err)
	err = ioutil.WriteFile(src, data, 0777)
	CheckError(err)
}

func ReadJson(res []byte) string {
	var str bytes.Buffer
	json.Indent(&str, res, "", "    ")
	return str.String()
}

func CheckError(err error) {
	if err != nil {
		panic(err)
	}
}
