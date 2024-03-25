package service

import (
	"github.com/tp86/legimi-go/internal/config"
	"github.com/tp86/legimi-go/internal/repository"
)

type AccountService interface {
	GetCredentials() (string, string)
	GetKindleId() uint64
}

func WithMemoryAccountRepository() config.ConfigFn[defaultAccountService] {
	mar := repository.GetMemoryRepository()
	return func(as defaultAccountService) defaultAccountService {
		as.accountRepository = mar
		return as
	}
}

func NewDefaultAccountService() AccountService {
	return config.New(
		WithMemoryAccountRepository(),
	)
}

type defaultAccountService struct {
	accountRepository repository.AccountRepository
}

func (as defaultAccountService) GetCredentials() (string, string) {
	login := as.accountRepository.GetLogin()
	password := as.accountRepository.GetPassword()
	return login, password
}

func (as defaultAccountService) GetKindleId() uint64 {
	return as.accountRepository.GetKindleId()
}
