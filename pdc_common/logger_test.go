package pdc_common

import (
	"errors"
	"log"
	"testing"

	"github.com/rs/zerolog"
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

func TestLoggerEvent(t *testing.T) {

	ReportErrorCustom(errors.New("test error"), func(event *zerolog.Event) *zerolog.Event {
		return event.Str("payload", "asdasdasdasd")
	})

}
