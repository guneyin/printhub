package auth

import (
	"github.com/guneyin/printhub/mail/sender"
	"github.com/guneyin/printhub/market"
	"github.com/matcornic/hermes/v2"
)

type VerifyUser struct {
	Token string `json:"token"`
}

func (v VerifyUser) Send(to, subject string) error {
	body, err := v.Generate()
	if err != nil {
		return err
	}

	return sender.SendMail(to, subject, body)
}

func (v VerifyUser) Generate() (hermes.Email, error) {
	cfg := market.Get().Config
	return hermes.Email{
		Body: hermes.Body{
			Intros: []string{
				"Bu e-postayı, yeni bir PrintHub hesabı oluşturduğunuz için aldınız.",
			},
			Actions: []hermes.Action{
				{
					Instructions: "Hesabınızı doğrulamak için aşağıdaki düğmeyi tıklayın:",
					Button: hermes.Button{
						Color: "#DC4D2F",
						Text:  "Hesabımı doğrula",
						Link:  cfg.AppURL + "/validate?token=" + v.Token,
					},
				},
			},
			Outros: []string{
				"Yeni bir hesap açmadıysanız başka bir işlem yapmanız gerekmez.",
			},
		},
	}, nil
}
