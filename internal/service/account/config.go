package account

import (
	"github.com/tp86/legimi-go/internal/repository"
	"github.com/tp86/legimi-go/internal/service"
)

func DefaultService(r repository.Account) service.Account {
	return defaultAccountService{accountRepository: r}
}
