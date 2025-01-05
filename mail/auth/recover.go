package auth

import (
	"github.com/guneyin/printhub/mail/sender"
	"github.com/guneyin/printhub/market"
	"github.com/matcornic/hermes/v2"
)

type RecoverPassword struct {
	Token string `json:"token"`
}

func (r RecoverPassword) Send(to, subject string) error {
	body, err := r.Generate()
	if err != nil {
		return err
	}

	return sender.SendMail(to, subject, body)
}

func (r RecoverPassword) Generate() (hermes.Email, error) {
	cfg := market.Get().Config
	return hermes.Email{
		Body: hermes.Body{
			Intros: []string{
				"Bu e-postayı, PrintHub hesabı için bir şifre sıfırlama talebi alındığı için aldınız.",
			},
			Actions: []hermes.Action{
				{
					Instructions: "Şifrenizi sıfırlamak için aşağıdaki düğmeyi tıklayın:",
					Button: hermes.Button{
						Color: "#DC4D2F",
						Text:  "Şifrenizi sıfırlayın",
						Link:  cfg.AppURL + "/recover?token=" + r.Token,
					},
				},
			},
			Outros: []string{
				"Parola sıfırlama talebinde bulunmadıysanız, başka bir işlem yapmanız gerekmez.",
			},
		},
	}, nil
}
