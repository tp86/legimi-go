package service

import (
	"github.com/tp86/legimi-go/internal/api"
	"github.com/tp86/legimi-go/internal/api/protocol"
	"github.com/tp86/legimi-go/internal/config"
)

type SessionService interface {
	GetSession() (string, error)
}

func WithDefaultAccountService() config.ConfigFn[defaultSessionService] {
	return func(ss defaultSessionService) defaultSessionService {
		ss.accountService = NewDefaultAccountService()
		return ss
	}
}

func WithApiClient() config.ConfigFn[defaultSessionService] {
	return func(ss defaultSessionService) defaultSessionService {
		ss.client = api.GetClient()
		return ss
	}
}

func NewDefaultSessionService() SessionService {
	return config.New(
		WithDefaultAccountService(),
		WithApiClient(),
	)
}

type defaultSessionService struct {
	accountService AccountService
	client         api.Client
}

func (ss defaultSessionService) GetSession() (string, error) {
	login, password := ss.accountService.GetCredentials()
	kindleId := ss.accountService.GetKindleId()
	var session protocol.Session
	err := ss.client.Exchange(protocol.NewGetSessionRequest(login, password, kindleId), &session)
	if err != nil {
		return "", err
	}
	return session.Id, nil
}
