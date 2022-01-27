/**
 * @Author Oliver
 * @Date 1/26/22
 **/

package util

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
)

func ReadFromJsonFile(src string) map[string]interface{} {
	data, err := ioutil.ReadFile(src)
	CheckError(err)
	res := make(map[string]interface{})
	err = json.Unmarshal(data, &res)
	CheckError(err)
	return res
}

func WriteToJsonFile(src string, res map[string]interface{}) {
	data, err := json.MarshalIndent(res, "", "	") // 第二个表示每行的前缀，这里不用，第三个是缩进符号，这里用tab
	CheckError(err)
	err = ioutil.WriteFile(src, data, 0777)
	CheckError(err)
}

func Json2Map(src []byte) map[string]interface{} {
	res := make(map[string]interface{})
	json.Unmarshal([]byte(src), &res)
	return res
}

func ReadJson(res []byte) string {
	var str bytes.Buffer
	json.Indent(&str, res, "", "    ")
	return str.String()
}
