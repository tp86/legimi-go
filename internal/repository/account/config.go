package account

import (
	"github.com/tp86/legimi-go/internal/options"
	"github.com/tp86/legimi-go/internal/repository"
)

var (
	mar *MemoryAccountRepository
	ar  repository.Account
)

func GetMemoryRepository() repository.Account {
	if mar == nil {
		mar = &MemoryAccountRepository{}
	}
	return mar
}

func GetFileRepository(opts options.Configuration) repository.Account {
	if ar == nil {
		ar = newFileAccountRepository(opts.GetFile())
	}
	return ar
}
