package account

import (
	"github.com/tp86/legimi-go/internal/api"
	"github.com/tp86/legimi-go/internal/repository"
	"github.com/tp86/legimi-go/internal/service"
)

func DefaultService(r repository.Account, client api.Client, login string, password string) service.Account {
	return defaultAccountService{accountRepository: r, client: client, login: login, password: password}
}
