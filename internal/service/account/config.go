package account

import (
	"github.com/tp86/legimi-go/internal/api"
	"github.com/tp86/legimi-go/internal/options"
	"github.com/tp86/legimi-go/internal/repository"
	"github.com/tp86/legimi-go/internal/service"
)

func DefaultService(r repository.Account, client api.Client, opts options.Credentials) service.Account {
	return defaultAccountService{accountRepository: r, client: client, userCredentials: opts}
}
