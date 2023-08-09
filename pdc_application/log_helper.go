package pdc_application

import (
	"runtime/debug"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

type LogHelper struct {
	Logger *zerolog.Logger
}

func (h *LogHelper) CapturePanicError() {
	if r := recover(); r != nil {
		err := r.(error)
		h.Logger.Panic().Err(err).Str("stacktrace", string(debug.Stack())).Msg(err.Error())
	}
}

func (h *LogHelper) CapturePanicErrorCustom(handles ...func(error)) {
	if r := recover(); r != nil {
		err := r.(error)
		if h.Logger == nil {
			log.Err(err).Str("stacktrace", string(debug.Stack())).Str("send", "false").Msg(err.Error())
		} else {
			h.Logger.Err(err).Str("stacktrace", string(debug.Stack())).Msg(err.Error())
		}

		for _, handle := range handles {
			handle(err)
		}
	}
}

func (h *LogHelper) ReportError(err error) error {

	h.Logger.Error().Err(err).Str("stacktrace", string(debug.Stack())).Msg(err.Error())

	return err
}

func (h *LogHelper) ReportErrorCustom(err error, handler func(event *zerolog.Event) *zerolog.Event) error {
	event := h.Logger.Error().Err(err).Str("stacktrace", string(debug.Stack()))
	event = handler(event)
	event.Msg(err.Error())
	return err
}
