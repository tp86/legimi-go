package account

type MemoryAccountRepository struct {
	Login, Password string
	KindleId        uint64
}

func (mar MemoryAccountRepository) GetLogin() string {
	return mar.Login
}

func (mar MemoryAccountRepository) GetPassword() string {
	return mar.Password
}

func (mar MemoryAccountRepository) GetKindleId() uint64 {
	return mar.KindleId
}

func (mar *MemoryAccountRepository) SaveLogin(login string) {
	mar.Login = login
}

func (mar *MemoryAccountRepository) SavePassword(password string) {
	mar.Password = password
}

func (mar *MemoryAccountRepository) SaveKindleId(id uint64) {
	mar.KindleId = id
}
