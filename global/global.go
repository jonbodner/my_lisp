package global

import (
	"log"
)

var DEBUG bool = true

func Log(vals ...interface{}) {
	if DEBUG {
		log.Println(vals...)
	}
}