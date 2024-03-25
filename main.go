package main

import (
	"fmt"

	"github.com/tp86/legimi-go/internal/service"
)

func main() {
	sessionId, err := service.NewSessionService().GetSession()
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Printf("session id: %s\n", sessionId)
}
