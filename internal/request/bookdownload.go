package request

type BookDownloadDetails struct {
	session
	bookId uint64
}

func NewBookDownloadDetailsRequest(sessionId string, bookId uint64) BookDownloadDetails {
	return BookDownloadDetails{
		session: session{id: sessionId},
		bookId:  bookId,
	}
}
