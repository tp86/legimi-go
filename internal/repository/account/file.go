package account

import (
	"fmt"
	"os"
	"path"

	"github.com/tp86/legimi-go/internal/repository"
	"gopkg.in/ini.v1"
)

type fileAccountRepository struct {
	filePath string
	file     *ini.File
	config   *ini.Section
}

func newFileAccountRepository(configFile string) repository.Account {
	file, err := ini.Load(configFile)
	if err != nil {
		file = ini.Empty()
		os.MkdirAll(path.Dir(configFile), 0755)
		file.SaveTo(configFile)
	}
	return &fileAccountRepository{
		filePath: configFile,
		file:     file,
		config:   file.Section(""),
	}
}

func (far fileAccountRepository) GetLogin() string {
	key, err := far.config.GetKey("login")
	if err != nil {
		return ""
	}
	return key.MustString("")
}

func (far fileAccountRepository) GetPassword() string {
	key, err := far.config.GetKey("password")
	if err != nil {
		return ""
	}
	return key.MustString("")
}

func (far fileAccountRepository) GetKindleId() uint64 {
	key, err := far.config.GetKey("kindleId")
	if err != nil {
		return 0
	}
	return key.MustUint64(0)
}

func (far fileAccountRepository) SaveLogin(login string) {
	key := far.config.Key("login")
	key.SetValue(login)
	far.file.SaveTo(far.filePath)
}

func (far fileAccountRepository) SavePassword(password string) {
	key := far.config.Key("password")
	key.SetValue(password)
	far.file.SaveTo(far.filePath)
}

func (far fileAccountRepository) SaveKindleId(kindleId uint64) {
	key := far.config.Key("kindleId")
	key.SetValue(fmt.Sprint(kindleId))
	far.file.SaveTo(far.filePath)
}
