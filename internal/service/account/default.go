package account

import (
	"fmt"
	"os"

	"golang.org/x/term"

	"github.com/tp86/legimi-go/internal/api"
	"github.com/tp86/legimi-go/internal/model"
	"github.com/tp86/legimi-go/internal/options"
	"github.com/tp86/legimi-go/internal/repository"
)

// TODO handle errors

type defaultAccountService struct {
	accountRepository repository.Account
	client            api.Client
	userCredentials   options.Credentials
}

func (as defaultAccountService) getLogin() string {
	login := as.userCredentials.GetLogin()
	if login == "" {
		login = as.accountRepository.GetLogin()
	}
	if login == "" {
		fmt.Print("Enter legimi login: ")
		if _, err := fmt.Scanln(&login); err == nil {
			as.accountRepository.SaveLogin(login)
		}
	}
	return login
}

func (as defaultAccountService) getPassword() string {
	password := as.userCredentials.GetPassword()
	if password == "" {
		password = as.accountRepository.GetPassword()
	}
	if password == "" {
		fmt.Print("Enter legimi password: ")
		if passwordBytes, err := term.ReadPassword(int(os.Stdin.Fd())); err == nil {
			password = string(passwordBytes)
			as.accountRepository.SavePassword(password)
		}
		fmt.Println()
	}
	return password
}

func (as defaultAccountService) GetCredentials() (string, string) {
	login := as.getLogin()
	password := as.getPassword()
	return login, password
}

func (as defaultAccountService) GetKindleId() uint64 {
	kindleId := as.accountRepository.GetKindleId()
	// if kindle id is not in repository
	if kindleId == 0 {
		// ask user for Kindle Serial No
		fmt.Print("Enter Kindle Serial Number: ")
		var kindleSerialNumber string
		fmt.Scanln(&kindleSerialNumber)
		// then query api for kindle id
		var registered model.Register
		login, password := as.GetCredentials()
		as.client.Exchange(model.NewRegisterRequest(login, password, kindleSerialNumber), &registered)
		kindleId = registered.KindleId
		// and store result in repository
		as.accountRepository.SaveKindleId(kindleId)
	}
	return kindleId
}
