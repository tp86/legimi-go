package presenter

import "github.com/tp86/legimi-go/internal/service"

func DefaultBookListPresenter() service.BookListPresenter {
	return defaultBookListPresenter{}
}

func DefaultBookDownloadPresenter() service.DownloadPresenter {
	return defaultBookDownloadPresenter{}
}
