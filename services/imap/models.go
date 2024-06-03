package imap

import (
	"fmt"
	"io"
	"time"

	"github.com/emersion/go-message/mail"
	"github.com/k3a/html2text"
)

type EmailAttachment struct {
	Filename string
	Content  []byte
}

type EmailEntity struct {
	Id          string
	From        mail.Address
	Sender      mail.Address
	ReplyTo     mail.Address
	To          mail.Address
	Cc          mail.Address
	Bcc         mail.Address
	Subject     string
	Date        time.Time
	Content     string
	Attachments []EmailAttachment
}

func (ee *EmailEntity) String() string {
	str := fmt.Sprintf(`
Id: %v,
From: %v, To: %v, Cc: %v, Bcc: %v
Sender: %v, ReplyTo: %v
------------------
Subject: %v
Date: %v
------------------
Content: %v
------------------
`,
		ee.Id,
		ee.From.Address,
		ee.To.Address,
		ee.Cc.Address,
		ee.Bcc.Address,
		ee.Sender.Address,
		ee.ReplyTo.Address,
		ee.Subject,
		ee.Date,
		ee.Content,
	)
	for i, attachment := range ee.Attachments {
		str += fmt.Sprintf("Attachment %v: %v\n", i+1, attachment.Filename)
	}
	return str
}

func ParseEmailLiteral(r io.Reader) (entity *EmailEntity, err error) {
	var mr *mail.Reader
	if mr, err = mail.CreateReader(r); err != nil {
		return
	}
	addressFn := func(key string) (address mail.Address) {
		var err error
		var addressList []*mail.Address
		if addressList, err = mr.Header.AddressList(key); err == nil && len(addressList) > 0 && addressList[0] != nil {
			address = *addressList[0]
		}
		return
	}
	entity = &EmailEntity{
		From:    addressFn("From"),
		Sender:  addressFn("Sender"),
		ReplyTo: addressFn("Reply-To"),
		To:      addressFn("To"),
		Cc:      addressFn("Cc"),
		Bcc:     addressFn("Bcc"),
	}
	// id
	var messageId string
	if messageId, err = mr.Header.MessageID(); err == nil {
		entity.Id = messageId
	}
	// subject
	var subject string
	if subject, err = mr.Header.Subject(); err == nil {
		entity.Subject = subject
	}
	// date
	var date time.Time
	if date, err = mr.Header.Date(); err == nil {
		entity.Date = date
	}
	// content or attachment
	for {
		var p *mail.Part
		if p, err = mr.NextPart(); err == io.EOF {
			err = nil
			break
		} else if err != nil {
			return
		}

		switch h := p.Header.(type) {
		case *mail.InlineHeader:
			// This is the message's text (can be plain-text or HTML)
			var b []byte
			if b, err = io.ReadAll(p.Body); err == nil {
				entity.Content = html2text.HTML2Text(string(b))
			}
		case *mail.AttachmentHeader:
			// This is an attachment
			var filename string
			if filename, err = h.Filename(); err == nil {
				var b []byte
				if b, err = io.ReadAll(p.Body); err == nil {
					entity.Attachments = append(entity.Attachments, EmailAttachment{
						Filename: filename,
						Content:  b,
					})
				}
			}
		}
	}
	return
}
