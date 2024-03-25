package repository

type AccountRepository interface {
	GetLogin() string
	GetPassword() string
	GetKindleId() uint64
	SaveLogin(login string)
	SavePassword(password string)
	SaveKindleId(kindleId uint64)
}
