package mail

import (
	"errors"
	"regexp"

	"gopkg.in/gomail.v2"
)

var (
	// ErrMailAddrCannotEmpty mail server addresses cannot empty
	ErrMailAddrCannotEmpty = errors.New("addresses cannot be empty")
	// ErrMailUserCannotEmpty username cannot empty
	ErrMailUserCannotEmpty = errors.New("username cannot empty")
	// ErrMailPasswordCannotEmpty password cannot empaty
	ErrMailPasswordCannotEmpty = errors.New("password cannot empaty")
	// ErrMailPortCannotEmpty port cannot empty
	ErrMailPortCannotEmpty = errors.New("port cannot empty")
)

// MailServer 邮件服务器配置信息
type MailServer struct {
	Addr     string
	Port     int
	Username string
	Password string
}

func New(username, password, addr string, port int) *MailServer {
	return &MailServer{
		Addr:     addr,
		Port:     port,
		Username: username,
		Password: password,
	}
}

// Send 发送普通文本、携带附件、HTML邮件
func (m *MailServer) Send(subject, body string, attach, to []string) error {
	err := checkMailServer(m)
	if err != nil {
		return err
	}

	return m.sendMail(subject, body, &attach, &to)
}

// 检查mail配置信息
func checkMailServer(m *MailServer) error {
	if m.Addr == "" {
		return ErrMailAddrCannotEmpty
	} else if m.Username == "" {
		return ErrMailUserCannotEmpty
	} else if m.Password == "" {
		return ErrMailPasswordCannotEmpty
	} else if m.Port == 0 {
		return ErrMailPortCannotEmpty
	}
	return nil
}

// 构造邮件内容并发送
func (m *MailServer) sendMail(subject, body string, attach, to *[]string) error {
	msg := gomail.NewMessage()
	msg.SetHeader("From", m.Username)
	msg.SetHeader("To", *to...)
	msg.SetHeader("Subject", subject)

	re, err := regexp.Compile(`.*<\S+>.*</\S+>.*`)
	if err != nil {
		return err
	}

	if re.MatchString(body) {
		msg.SetBody("text/html", body)
	} else {
		msg.SetBody("text/plain", body)
	}

	if attach != nil {
		for _, v := range *attach {
			msg.Attach(v)
		}
	}

	d := gomail.NewDialer(m.Addr, m.Port, m.Username, m.Password)

	if err := d.DialAndSend(msg); err != nil {
		return err
	}
	return nil
}
