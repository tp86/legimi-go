package options

type Credentials interface {
	GetLogin() string
	GetPassword() string
}

type Configuration interface {
	GetFile() string
}
