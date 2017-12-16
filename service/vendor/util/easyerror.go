package util

import (
	log "util/logger"
)

func PanicIf(status interface{}, msg ...interface{}) {
	switch v := status.(type) {
	case error:
		log.Panicln(v, msg)
	case bool:
		if v {
			log.Panicf("... unknown panic: %v\n\t ... since %T:%v\n", msg, v, v)
		}
	default:
		if v != nil {
			log.Panicf("... unknown panic: %v\n\t ... since %T:%v\n", msg, v, v)
		}
	}
}
func WarnIf(status interface{}, msg ...interface{}) {
	switch v := status.(type) {
	case error:
		log.Warningln(v, msg)
	case bool:
		if v {
			log.Panicf("... unknown warning: %v\n\t ... since %T:%v\n", msg, v, v)
		}
	default:
		if v != nil {
			log.Panicf("... unknown warning: %v\n\t ... since %T:%v\n", msg, v, v)
		}
	}
}
