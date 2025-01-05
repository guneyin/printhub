package mail

import (
	"github.com/guneyin/printhub/mail/auth"
	"github.com/matcornic/hermes/v2"
)

type EMail interface {
	Send(to, subject string) error
	Generate() (hermes.Email, error)
}

func NewRecoverPasswordEmail(token string) EMail {
	return auth.RecoverPassword{Token: token}
}

func NewVerifyUserEmail(token string) EMail {
	return auth.VerifyUser{Token: token}
}
