/**
 * @Author Oliver
 * @Date 1/24/22
 **/

package mail

import (
	"crypto/tls"
	. "daily_fudan/util"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/smtp"
)

type Header struct {
	From        string `json:"From"`         //发送方名字
	To          string `json:"To"`           //接收方邮箱
	Subject     string `json:"Subject"`      //标题
	ContentType string `json:"Content-Type"` //内容格式
}

type Email struct {
	Host     string `json:"Host"`     //smtp服务器
	Port     string `json:"Port"`     //smtp服务器端口
	Email    string `json:"Email"`    // 这里是你的邮箱地址
	Password string `json:"Password"` // 这里填你的授权码
	ToEmail  string `json:"ToEmail"`  // 目标地址
	Header   Header `json:"Header"`
	Body     string `json:"Body"` //邮件内容
}

//return a smtp client
func Dial(addr string) (*smtp.Client, error) {
	conn, err := tls.Dial("tcp", addr, nil)
	if err != nil {
		log.Panicln("Dialing Error:", err)
		return nil, err
	}
	//分解主机端口字符串
	host, _, _ := net.SplitHostPort(addr)
	return smtp.NewClient(conn, host)
}

//参考net/smtp的func SendMail()
//使用net.Dial连接tls(ssl)端口时,smtp.NewClient()会卡住且不提示err
//len(to)>1时,to[1]开始提示是密送
func SendMailUsingTLS(addr string, auth smtp.Auth, from string, to string, msg []byte) (err error) {
	//create smtp client
	c, err := Dial(addr)
	if err != nil {
		log.Println("Create smpt client error:", err)
		return err
	}
	defer c.Close()

	if auth != nil {
		if ok, _ := c.Extension("AUTH"); ok {
			if err = c.Auth(auth); err != nil {
				log.Println("Error during AUTH", err)
				return err
			}
		}
	}

	if err = c.Mail(from); err != nil {
		return err
	}

	if err = c.Rcpt(to); err != nil {
		return err
	}

	w, err := c.Data()
	if err != nil {
		return err
	}

	_, err = w.Write(msg)
	if err != nil {
		return err
	}

	err = w.Close()
	if err != nil {
		return err
	}
	return c.Quit()
}

func createEmailJson(src string) {
	res := &Email{
		Host:     "smtp.qq.com",  //smtp服务器
		Port:     "465",          //端口
		Email:    "xxxx@xxx.xxx", //发送方的邮箱
		Password: "xxx",          //发送方的密钥
		ToEmail:  "xxxx@xxx.xxx", //接收方邮箱
		Header: Header{
			From: "xxxxx", //发送方昵称
			To:   "xxxxx", //接收方邮箱

			Subject:     "xxxxx",                    //邮件标题
			ContentType: "text/html;chartset=UTF-8", //邮件格式
		},
		Body: "xxxxx", //邮件
	}
	fmt.Println(`请输入发送方smtp服务器（直接回车为默认"smtp.qq.com"）`)
	fmt.Scanln(&res.Host)
	fmt.Println(`请输入发送方端口（直接回车为默认"465"）`)
	fmt.Scanln(&res.Port)
	fmt.Println("请输入发送方邮箱")
	fmt.Scanln(&res.Email)
	fmt.Println("请输入发送方密钥")
	fmt.Scanln(&res.Password)
	fmt.Println(`请输入接收方邮箱`)
	fmt.Scanln(&res.ToEmail)
	fmt.Println("请输入发送方昵称")
	fmt.Scanln(&res.Header.From)
	res.Header.From = res.Header.From + "<" + res.Email + ">"
	res.Header.To = res.ToEmail
	fmt.Println("请输入邮件标题")
	fmt.Scanln(&res.Header.Subject)
	fmt.Println(`请输入邮件格式（直接回车为默认"text/html;chartset=UTF-8"）`)
	fmt.Scanln(&res.Header.ContentType)
	data, err := json.MarshalIndent(res, "", "	") // 第二个表示每行的前缀，这里不用，第三个是缩进符号，这里用tab
	CheckError(err)
	err = ioutil.WriteFile(src, data, 0777)
	CheckError(err)
}

func Mail() {
	data, _ := ioutil.ReadFile("mail.json")
	if data == nil {
		createEmailJson("mail.json")
	}
	email := ReadFromJsonFile("mail.json")
	header := email["Header"].(map[string]interface{})
	message := ""
	for k, v := range header {
		message += fmt.Sprintf("%s:%s\r\n", k, v)
	}
	message += "\r\n" + email["Body"].(string)

	auth := smtp.PlainAuth(
		"",
		email["Email"].(string),
		email["Password"].(string),
		email["Host"].(string),
	)

	err := SendMailUsingTLS(
		email["Host"].(string)+":"+email["Port"].(string),
		auth,
		email["Email"].(string),
		email["ToEmail"].(string),
		[]byte(message),
	)

	if err != nil {
		panic(err)
	}
}

func MailTo(to, msg string) {
	data, _ := ioutil.ReadFile("mail.json")
	if data == nil {
		createEmailJson("mail.json")
	}
	email := ReadFromJsonFile("mail.json")
	header := email["Header"].(map[string]interface{})
	message := ""
	for k, v := range header {
		message += fmt.Sprintf("%s:%s\r\n", k, v)
	}
	if to != "" {
		email["ToEmail"] = to
	}
	if msg != "" {
		email["Body"] = msg
	}
	message += "\r\n" + email["Body"].(string)

	auth := smtp.PlainAuth(
		"",
		email["Email"].(string),
		email["Password"].(string),
		email["Host"].(string),
	)

	err := SendMailUsingTLS(
		email["Host"].(string)+":"+email["Port"].(string),
		auth,
		email["Email"].(string),
		email["ToEmail"].(string),
		[]byte(message),
	)

	if err != nil {
		panic(err)
	}
}
