package session

import (
	"github.com/tp86/legimi-go/internal/api"
	"github.com/tp86/legimi-go/internal/service"
)

func DefaultService(as service.Account, c api.Client) service.Session {
	return defaultSessionService{
		accountService: as,
		client:         c,
	}
}
