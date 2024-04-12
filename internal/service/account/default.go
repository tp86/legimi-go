package account

import (
	"fmt"
	"os"

	"golang.org/x/term"

	"github.com/tp86/legimi-go/internal/api"
	"github.com/tp86/legimi-go/internal/model"
	"github.com/tp86/legimi-go/internal/repository"
)

// TODO handle errors

type defaultAccountService struct {
	accountRepository repository.Account
	client            api.Client
	login             string
	password          string
}

func (as defaultAccountService) GetCredentials() (string, string) {
	var (
		login, password string
	)
	// if credentials are given as command line options, use them
	if as.login != "" {
		login = as.login
	}
	if as.password != "" {
		password = as.password
	}
	// if credentials are not given as command line options, get them from repository
	if login == "" {
		login = as.accountRepository.GetLogin()
	}
	if password == "" {
		password = as.accountRepository.GetPassword()
	}
	// if credentials do not exist in repository, ask user for them interactively
	// and store them in repository
	if login == "" {
		fmt.Print("Enter legimi login: ")
		if _, err := fmt.Scanln(&login); err == nil {
			as.accountRepository.SaveLogin(login)
		}
	}
	if password == "" {
		fmt.Print("Enter legimi password: ")
		if passwordBytes, err := term.ReadPassword(int(os.Stdin.Fd())); err == nil {
			password = string(passwordBytes)
			as.accountRepository.SavePassword(password)
		}
		fmt.Println()
	}
	return login, password
}

func (as defaultAccountService) GetKindleId() uint64 {
	kindleId := as.accountRepository.GetKindleId()
	// if kindle id is not in repository, ask user for Kindle Serial No
	// then query api for kindle id
	// and store result in repository
	if kindleId == 0 {
		fmt.Print("Enter Kindle Serial Number: ")
		var kindleSerialNumber string
		fmt.Scanln(&kindleSerialNumber)
		var registered model.Register
		login, password := as.GetCredentials()
		as.client.Exchange(model.NewRegisterRequest(login, password, kindleSerialNumber), &registered)
		kindleId = registered.KindleId
		as.accountRepository.SaveKindleId(kindleId)
	}
	return kindleId
}
