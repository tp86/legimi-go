package protocol

import "fmt"

const BookDownloadDetailsPreparingError = 295

// TODO list known errors and their descriptions
var errorResponses = map[uint16]string{
	133: "invalid credentials",
	163: "invalid kindle id",
}

type ErrorResponse struct {
	Type uint16
}

func (er ErrorResponse) Error() string {
	if message, ok := errorResponses[er.Type]; ok {
		return fmt.Sprintf(message)
	}
	return fmt.Sprintf("error response received: %d", er.Type)
}

func isErrorResponse(responseType uint16) bool {
	return responseType%2 == 1
}
