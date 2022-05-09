package RustEmail

// 此处配置 没有放到app.conf 是为了 保护个人邮箱安全

import (
	"gopkg.in/gomail.v2"
)

type RustSendEmail struct {
	Host     string
	Port     int
	Username string
	Password string
}

// 阿里云默认屏蔽的25端口  请用465
const (
	RustEmailHost     = "smtp.qiye.aliyun.com"
	RustEmailPort     = 465
	RustEmailUsername = "admin@dbsgw.cn"
	RustEmailPassword = "ly8334721.="
)

func (e RustSendEmail) Send(mailTo []string, subject string, body string) error {
	m := gomail.NewMessage()
	// 发件人信息
	m.SetHeader("From", m.FormatAddress("admin@dbsgw.cn", "Rust中文网"))
	// 收件人
	m.SetHeader("To", mailTo...)
	// 主题
	m.SetHeader("Subject", subject)
	// 内容
	m.SetBody("text/html", body)

	d := gomail.NewDialer(e.Host, e.Port, e.Username, e.Password)
	err := d.DialAndSend(m)
	if err != nil {
		// 处理错误
		return err
	}
	return nil
}

func NewRustSendEmail(host, username, password string, port int) *RustSendEmail {
	return &RustSendEmail{
		Host:     host,
		Port:     port,
		Username: username,
		Password: password,
	}
}

// NewDefaultSendEmail 默认配置发送邮件
func NewDefaultSendEmail() *RustSendEmail {
	return &RustSendEmail{
		Host:     RustEmailHost,
		Port:     RustEmailPort,
		Username: RustEmailUsername,
		Password: RustEmailPassword,
	}
}
func init() {

	//rustEmail := NewRustSendEmail(RustEmailHost, RustEmailUsername, RustEmailPassword, RustEmailPort)
	//rustEmail := NewDefaultSendEmail()
	//rustEmail.Send([]string{"1578347363@qq.com"}, "Rust中文网", "<h1>来自Rust中文网验证码：654230222222</h1>")
}
