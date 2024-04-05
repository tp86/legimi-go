package account

import "github.com/tp86/legimi-go/internal/repository"

var (
	mar *MemoryAccountRepository
)

func GetMemoryRepository() repository.Account {
	if mar == nil {
		mar = &MemoryAccountRepository{}
	}
	return mar
}
