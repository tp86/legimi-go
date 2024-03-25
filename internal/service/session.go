package service

import (
	"github.com/tp86/legimi-go/internal/protocol"
)

type SessionService interface {
	GetSession() (string, error)
}

func NewSessionService() SessionService {
	return sessionService{}
}

type sessionService struct{}

func (ss sessionService) GetSession() (string, error) {
	l, p, i := getSessionRequestData()
	var session protocol.Session
	err := protocol.Exchange(protocol.NewGetSessionRequest(l, p, i), &session)
	if err != nil {
		return "", err
	}
	return session.Id, nil
}
