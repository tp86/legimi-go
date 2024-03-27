package presenter

import "github.com/tp86/legimi-go/internal/service"

func DefaultService() service.BookListPresenter {
	return defaultBookListPresenter{}
}
