package service

type Session interface {
	GetSession() (string, error)
}

type Account interface {
	GetCredentials() (string, string)
	GetKindleId() uint64
}
