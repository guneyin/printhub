package sender

import (
	"errors"
	"github.com/go-gomail/gomail"
	"github.com/guneyin/printhub/market"
	"github.com/matcornic/hermes/v2"
	"net/mail"
	"strconv"
)

type smtpAuthentication struct {
	Server         string
	Port           string
	SenderEmail    string
	SenderIdentity string
	SMTPUser       string
	SMTPPassword   string
}

func SendMail(to, subject string, email hermes.Email) error {
	cfg := market.Get().Config
	if cfg.MailEnabled == false {
		return errors.New("mail sender disabled by default")
	}

	h := hermes.Hermes{
		Product: hermes.Product{
			Name:        cfg.AppName,
			Link:        cfg.AppURL,
			Logo:        "https://github.com/matcornic/hermes/blob/master/examples/gopher.png?raw=true",
			Copyright:   "portfoyum.com © 2020 - Tüm Hakları Saklıdır",
			TroubleText: "{ACTION} düğmesiyle ilgili sorun yaşıyorsanız, aşağıdaki URL'yi kopyalayıp web tarayıcınıza yapıştırın.",
		},
	}

	h.Theme = new(hermes.Flat)

	//email.Body.Name = u.Name + " " + u.Surname
	email.Body.Greeting = "Merhaba"
	email.Body.Signature = "Teşekkürler"

	htmlBytes, err := h.GenerateHTML(email)
	if err != nil {
		return err
	}
	txtBytes, err := h.GeneratePlainText(email)
	if err != nil {
		return err
	}

	return send(to, subject, htmlBytes, txtBytes)
}

// send sends the email
func send(to, subject, htmlBody, txtBody string) error {
	cfg := market.Get().Config

	smtpConfig := smtpAuthentication{
		Server:         cfg.EmailSmtpServer,
		Port:           cfg.EmailSmtpPort,
		SenderEmail:    cfg.EmailSender,
		SenderIdentity: cfg.EmailIdentity,
		SMTPPassword:   cfg.EmailPassword,
		SMTPUser:       cfg.EmailUser,
	}

	if smtpConfig.Server == "" {
		return errors.New("smtp server required")
	}

	port, err := strconv.Atoi(smtpConfig.Port)
	if err != nil {
		return errors.New("invalid smtp port")
	}

	if smtpConfig.SMTPUser == "" {
		return errors.New("smtp user is required")
	}

	if smtpConfig.SenderIdentity == "" {
		return errors.New("smtp sender identity is required")
	}

	if smtpConfig.SenderEmail == "" {
		return errors.New("sender email is required")
	}

	if to == "" {
		return errors.New("recipient email is required")
	}

	from := mail.Address{
		Name:    smtpConfig.SenderIdentity,
		Address: smtpConfig.SenderEmail,
	}

	m := gomail.NewMessage()
	m.SetHeader("From", from.String())
	m.SetHeader("To", to)
	m.SetHeader("Subject", subject)

	m.SetBody("text/plain", txtBody)
	m.AddAlternative("text/html", htmlBody)

	d := gomail.NewDialer(smtpConfig.Server, port, smtpConfig.SMTPUser, smtpConfig.SMTPPassword)

	return d.DialAndSend(m)
}
