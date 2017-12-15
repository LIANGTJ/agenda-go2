package util

import (
	log "util/logger"
)

func PanicIf(status interface{}) {
	switch v := status.(type) {
	case error:
		panic(v)
	default:
		if v != nil {
			log.Panicf("... unknown panic ... %T:%v\n", v, v)
		}
	}
}
func WarnIf(status interface{}) {
	switch v := status.(type) {
	case error:
		log.Warningln(v)
	default:
		if v != nil {
			log.Panicf("... unknown warning ... %T:%v\n", v, v)
		}
	}
}
