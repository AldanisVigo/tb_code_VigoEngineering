package lib

import (
	"bitbucket.org/taubyte/go-sdk/event"
)

//export ping
func readjson(e event.Event) uint32 {
	h, err := e.HTTP()
	if err != nil {
		return 1
	}

	h.Write([]byte("PONG"))

	return 0
}
