package pdc_common

import (
	"errors"
	"log"
	"testing"
)

func RunDenganPanic() {
	defer CapturePanicError()

	log.Println("running test panic")
	err := errors.New("test error")

	panic(err)
}

func TestLogger(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {

			log.Println(r)
		}
	}()

	RunDenganPanic()
}
