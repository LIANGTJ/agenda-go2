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
			log.Panicf("... unknown ... %T:%v\n", v, v)
		}
	}
}
