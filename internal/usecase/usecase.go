package usecase

type BookLister interface {
	ListBooks() error
}

type BookDownloader interface {
	DownloadBooks([]uint64) error
}
