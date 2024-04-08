package account

import (
	"github.com/tp86/legimi-go/internal/repository"
)

type defaultAccountService struct {
	accountRepository repository.Account
	login             string
	password          string
}

func (as defaultAccountService) GetCredentials() (string, string) {
	var (
		login, password string
	)
	if as.login != "" {
		login = as.login
	}
	if as.password != "" {
		password = as.password
	}
	if login == "" {
		login = as.accountRepository.GetLogin()
	}
	if password == "" {
		password = as.accountRepository.GetPassword()
	}
	return login, password
}

func (as defaultAccountService) GetKindleId() uint64 {
	return as.accountRepository.GetKindleId()
}
