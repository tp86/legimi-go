package account

import "github.com/tp86/legimi-go/internal/repository"

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

func GetFileRepository(configFile string) repository.Account {
	if ar == nil {
		ar = newFileAccountRepository(configFile)
	}
	return ar
}
