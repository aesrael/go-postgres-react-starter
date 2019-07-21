package errors

import "log"

// var ValidationErrors = []string{}

//HandleErr //generic error handler, logs error and Os.Exit(1)
func HandleErr(err error) error {
	if err != nil {
		log.Fatal(err)
	}
	return err
}
