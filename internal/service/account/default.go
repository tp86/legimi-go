package account

import (
	"github.com/tp86/legimi-go/internal/repository"
)

type defaultAccountService struct {
	accountRepository repository.Account
}

func (as defaultAccountService) GetCredentials() (string, string) {
	login := as.accountRepository.GetLogin()
	password := as.accountRepository.GetPassword()
	return login, password
}

func (as defaultAccountService) GetKindleId() uint64 {
	return as.accountRepository.GetKindleId()
}
