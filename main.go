package main

import (
	"fmt"

	"github.com/tp86/legimi-go/internal/service"
)

var sessionService service.Session

func main() {
	sessionId, err := sessionService.GetSession()
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Printf("session id: %s\n", sessionId)
}
