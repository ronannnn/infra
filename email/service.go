package email

import (
	"crypto/tls"
	"fmt"
	"net/mail"
	"net/smtp"
	"strings"
)

func NewClient(
	cfg *Cfg,
) *Client {
	return &Client{
		SmtpHost:     cfg.SmtpAddr,
		SmtpPort:     cfg.SmtpPort,
		EmailAccount: cfg.EmailAccount,
		EmailPasswd:  cfg.EmailPasswd,
	}
}

type Client struct {
	SmtpHost     string
	SmtpPort     uint16
	EmailAccount string
	EmailPasswd  string
}

type EmailPayload struct {
	Subject  string
	To       []string
	Cc       []string
	Bcc      []string
	Body     string
	HtmlType bool
}

func (c *Client) Send(payload EmailPayload) (err error) {
	// 发件人和收件人信息
	from := mail.Address{Name: "", Address: c.EmailAccount}
	fromStr := from.String()
	filteredTo, _ := FilterOutDistinctValidAndInvalidEmails(payload.To)
	tos := make([]mail.Address, 0)
	for _, to := range filteredTo {
		tos = append(tos, mail.Address{Name: "", Address: to})
	}
	toStrList := make([]string, 0)
	for _, to := range tos {
		toStrList = append(toStrList, to.String())
	}

	// 邮件内容
	subject := payload.Subject
	body := payload.Body

	// 邮件头
	header := make(map[string]string)
	header["From"] = fromStr
	header["To"] = strings.Join(toStrList, ",")
	header["Subject"] = subject
	if payload.HtmlType {
		header["MIME-Version"] = "1.0"
		header["Content-Type"] = "text/html; charset=\"UTF-8\""
	}

	// 邮件体
	message := ""
	for k, v := range header {
		message += k + ": " + v + "\r\n"
	}
	message += "\r\n" + body

	// SMTP 连接配置
	auth := smtp.PlainAuth("", c.EmailAccount, c.EmailPasswd, c.SmtpHost)
	tlsConfig := &tls.Config{
		ServerName:         c.SmtpHost,
		InsecureSkipVerify: true,
	}

	// 连接到 SMTP 服务器
	var conn *tls.Conn
	if conn, err = tls.Dial("tcp", fmt.Sprintf("%s:%d", c.SmtpHost, c.SmtpPort), tlsConfig); err != nil {
		return
	}
	defer conn.Close()

	var client *smtp.Client
	if client, err = smtp.NewClient(conn, c.SmtpHost); err != nil {
		return
	}
	defer client.Close()

	// 验证并发送邮件
	if err = client.Auth(auth); err != nil {
		return
	}

	if err = client.Mail(c.EmailAccount); err != nil {
		return
	}

	for _, to := range payload.To {
		if err = client.Rcpt(to); err != nil {
			return
		}
	}

	w, err := client.Data()
	if err != nil {
		return
	}
	defer w.Close()

	_, err = w.Write([]byte(message))
	if err != nil {
		return
	}

	return
}
