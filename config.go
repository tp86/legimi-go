package main

import (
	"github.com/tp86/legimi-go/internal/api"
	ar "github.com/tp86/legimi-go/internal/repository/account"
	as "github.com/tp86/legimi-go/internal/service/account"
	"github.com/tp86/legimi-go/internal/service/session"
)

func init() {
	accountRepository := ar.GetMemoryRepository()
	accountService := as.DefaultService(accountRepository)
	apiClient := api.GetClient()
	sessionService = session.DefaultService(accountService, apiClient)
}
