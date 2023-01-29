package global

import (
	"log"
)

var Debug = true

func Log(vals ...interface{}) {
	if Debug {
		log.Println(vals...)
	}
}
