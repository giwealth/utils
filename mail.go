package utils

import (
	"fmt"

	"gopkg.in/gomail.v2"
)

// MailServer 邮件服务器配置信息
type MailServer struct {
	Addr     string
	Port     int
	User     string
	Password string
}

// Send 发送普通文本邮件
func (m *MailServer) Send(subject, body string, to []string) error {
	err := m.checkMailConfig()
	if err != nil {
		return err
	}

	return m.sendMail(subject, body, []string{}, to)
}

// SendAttach 发送带附件的普通文本邮件
func (m *MailServer) SendAttach(subject, body string, attach, to []string) error {
	err := m.checkMailConfig()
	if err != nil {
		return err
	}

	return m.sendMail(subject, body, attach, to)
}

// 检查mail配置信息
func (m *MailServer) checkMailConfig() error {
	if m.Addr == "" || m.Port == 0 || m.User == "" || m.Password == "" {
		return fmt.Errorf("check mail server [Addr|Port|User|Password] configration")
	}
	return nil
}

// 构造邮件内容并发送
func (m *MailServer) sendMail(subject, body string, attach, to []string) error {
	msg := gomail.NewMessage()
	msg.SetHeader("From", m.User)
	msg.SetHeader("To", to...)
	msg.SetHeader("Subject", subject)
	msg.SetBody("text/plain", body)

	for _, v := range attach {
		msg.Attach(v)
	}

	d := gomail.NewDialer(m.Addr, m.Port, m.User, m.Password)

	if err := d.DialAndSend(msg); err != nil {
		return err
	}
	return nil
}
