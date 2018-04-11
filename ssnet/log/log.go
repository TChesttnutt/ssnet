//wrapping the standard log functions with error
//or debug strings. preserving interfaces.
package log

import (
	"log"
)

func LogErrln(v ...interface{}) {
	var args []interface{}
	args = append(args, "ERROR\t")
	args = append(args, v...)
	log.Println(args)
}

func LogErr(fmt string, vars ...interface{}) {
	log.Printf("ERROR\t"+fmt, vars...)
}

func LogDebln(v ...interface{}) {
	var args []interface{}
	args = append(args, "DEBUG\t")
	args = append(args, v...)
	log.Println(args)
}

func LogDebl(fmt string, vars ...interface{}) {
	log.Printf("DEBUG\t"+fmt, vars...)
}
