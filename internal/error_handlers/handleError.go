package errorhandlers

import "log"

func FailOnError(err error) {
	if err != nil {
		log.Panicln(err)
	}
}

func LogOnError(err error) {
	if err != nil {
		log.Println(err)
	}
}
