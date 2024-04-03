package book

import (
	"github.com/tp86/legimi-go/internal/api"
	"github.com/tp86/legimi-go/internal/service"
)

func DefaultService(
	sessionService service.Session,
	apiClient api.Client,
	bookDownloadPresenter service.DownloadPresenter,
) service.Book {
	return defaultBookService{
		sessionService:    sessionService,
		client:            apiClient,
		downloadPresenter: bookDownloadPresenter,
	}
}
