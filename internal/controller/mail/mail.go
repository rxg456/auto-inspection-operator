package mail

import (
	"crypto/tls"
	"fmt"
	"net/smtp"
	"strings"

	devopsv1 "github.com/rxg456/auto-inspection-operator/api/v1"
	"sigs.k8s.io/controller-runtime/pkg/log"
)

// Sender 邮件发送器
type Sender struct {
	Config devopsv1.SMTP
}

// NewSender 创建一个新的邮件发送器
func NewSender(config devopsv1.SMTP) *Sender {
	return &Sender{
		Config: config,
	}
}

// SendMail 发送邮件
func (s *Sender) SendMail(to []string, subject, body string) error {
	logger := log.Log.WithName("mail-sender")

	// 设置邮件头
	header := make(map[string]string)
	header["From"] = s.Config.From
	header["To"] = strings.Join(to, ";")
	header["Subject"] = subject
	header["Content-Type"] = "text/html; charset=UTF-8"

	// 组装邮件内容
	message := ""
	for k, v := range header {
		message += fmt.Sprintf("%s: %s\r\n", k, v)
	}
	message += "\r\n" + body

	// 邮件服务器地址
	addr := fmt.Sprintf("%s:%d", s.Config.Server, s.Config.Port)
	logger.Info("准备发送邮件", "server", s.Config.Server, "port", s.Config.Port)

	// 根据端口选择合适的发送方式
	switch s.Config.Port {
	case 465:
		// SSL 连接方式
		return s.sendMailWithSSL(addr, to, []byte(message))
	case 587:
		// TLS 连接方式
		return s.sendMailWithTLS(addr, to, []byte(message))
	default:
		// 标准连接方式
		// 认证信息
		auth := smtp.PlainAuth("", s.Config.Username, s.Config.Password, s.Config.Server)
		logger.Info("使用标准SMTP连接发送邮件")

		// 发送邮件
		err := smtp.SendMail(addr, auth, s.Config.From, to, []byte(message))
		if err != nil {
			logger.Error(err, "标准SMTP连接发送邮件失败")
			return fmt.Errorf("标准SMTP连接发送邮件失败: %w", err)
		}
		return nil
	}
}

// sendMailWithSSL 使用SSL发送邮件
func (s *Sender) sendMailWithSSL(addr string, to []string, msg []byte) error {
	logger := log.Log.WithName("mail-sender")
	logger.Info("使用SSL连接发送邮件")

	// 跳过TLS证书验证
	tlsConfig := &tls.Config{
		InsecureSkipVerify: true,
		ServerName:         s.Config.Server,
	}

	// 建立SSL连接
	conn, err := tls.Dial("tcp", addr, tlsConfig)
	if err != nil {
		logger.Error(err, "SSL连接失败")
		return fmt.Errorf("SSL连接失败: %w", err)
	}
	defer conn.Close()

	// 创建SMTP客户端
	client, err := smtp.NewClient(conn, s.Config.Server)
	if err != nil {
		logger.Error(err, "创建SMTP客户端失败")
		return fmt.Errorf("创建SMTP客户端失败: %w", err)
	}
	defer client.Close()

	// 认证
	auth := smtp.PlainAuth("", s.Config.Username, s.Config.Password, s.Config.Server)
	if err = client.Auth(auth); err != nil {
		logger.Error(err, "SMTP认证失败")
		return fmt.Errorf("SMTP认证失败: %w", err)
	}

	// 设置发件人
	if err = client.Mail(s.Config.From); err != nil {
		logger.Error(err, "设置发件人失败")
		return fmt.Errorf("设置发件人失败: %w", err)
	}

	// 设置收件人
	for _, recipient := range to {
		if err = client.Rcpt(recipient); err != nil {
			logger.Error(err, "设置收件人失败", "recipient", recipient)
			return fmt.Errorf("设置收件人失败: %w", err)
		}
	}

	// 设置邮件内容
	w, err := client.Data()
	if err != nil {
		logger.Error(err, "准备写入邮件内容失败")
		return fmt.Errorf("准备写入邮件内容失败: %w", err)
	}
	defer w.Close()

	_, err = w.Write(msg)
	if err != nil {
		logger.Error(err, "写入邮件内容失败")
		return fmt.Errorf("写入邮件内容失败: %w", err)
	}

	return nil
}

// sendMailWithTLS 使用TLS发送邮件
func (s *Sender) sendMailWithTLS(addr string, to []string, msg []byte) error {
	logger := log.Log.WithName("mail-sender")
	logger.Info("使用TLS连接发送邮件")

	// 建立普通连接
	client, err := smtp.Dial(addr)
	if err != nil {
		logger.Error(err, "连接SMTP服务器失败")
		return fmt.Errorf("连接SMTP服务器失败: %w", err)
	}
	defer client.Close()

	// 启动TLS
	tlsConfig := &tls.Config{
		InsecureSkipVerify: true,
		ServerName:         s.Config.Server,
	}
	if err = client.StartTLS(tlsConfig); err != nil {
		logger.Error(err, "启动TLS失败")
		return fmt.Errorf("启动TLS失败: %w", err)
	}

	// 认证
	auth := smtp.PlainAuth("", s.Config.Username, s.Config.Password, s.Config.Server)
	if err = client.Auth(auth); err != nil {
		logger.Error(err, "SMTP认证失败")
		return fmt.Errorf("SMTP认证失败: %w", err)
	}

	// 设置发件人
	if err = client.Mail(s.Config.From); err != nil {
		logger.Error(err, "设置发件人失败")
		return fmt.Errorf("设置发件人失败: %w", err)
	}

	// 设置收件人
	for _, recipient := range to {
		if err = client.Rcpt(recipient); err != nil {
			logger.Error(err, "设置收件人失败", "recipient", recipient)
			return fmt.Errorf("设置收件人失败: %w", err)
		}
	}

	// 设置邮件内容
	w, err := client.Data()
	if err != nil {
		logger.Error(err, "准备写入邮件内容失败")
		return fmt.Errorf("准备写入邮件内容失败: %w", err)
	}
	defer w.Close()

	_, err = w.Write(msg)
	if err != nil {
		logger.Error(err, "写入邮件内容失败")
		return fmt.Errorf("写入邮件内容失败: %w", err)
	}

	return nil
}
