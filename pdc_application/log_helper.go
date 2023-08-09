package pdc_application

import (
	"runtime/debug"

	"github.com/rs/zerolog"
)

type LogHelper struct {
	logger *zerolog.Logger
}

func (h *LogHelper) CapturePanicError() {
	if r := recover(); r != nil {
		err := r.(error)
		h.logger.Panic().Err(err).Str("stacktrace", string(debug.Stack())).Msg(err.Error())
	}
}

func (h *LogHelper) CapturePanicErrorCustom(handles ...func(error)) {
	if r := recover(); r != nil {
		err := r.(error)
		h.logger.Panic().Err(err).Str("stacktrace", string(debug.Stack())).Msg(err.Error())

		for _, handle := range handles {
			handle(err)
		}
	}
}

func (h *LogHelper) ReportError(err error) error {

	h.logger.Error().Err(err).Str("stacktrace", string(debug.Stack())).Msg(err.Error())

	return err
}

func (h *LogHelper) ReportErrorCustom(err error, handler func(event *zerolog.Event) *zerolog.Event) error {
	event := h.logger.Error().Err(err).Str("stacktrace", string(debug.Stack()))
	event = handler(event)
	event.Msg(err.Error())
	return err
}
