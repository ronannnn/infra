package email_test

import (
	"testing"

	"github.com/ronannnn/infra/email"
	"github.com/stretchr/testify/require"
)

func TestSend(t *testing.T) {
	cfg := email.Cfg{
		SmtpAddr:     "smtp.qq.com",
		SmtpPort:     465,
		EmailAccount: "",
		EmailPasswd:  "",
	}
	client := email.NewClient(cfg)
	payload := email.EmailPayload{
		Subject:  "test1",
		To:       []string{"853879506@qq.com", "crnoogle@gmail.com"},
		Body:     "<div>123</div>",
		HtmlType: true,
	}
	err := client.SendEmail(payload)
	require.NoError(t, err)
}
