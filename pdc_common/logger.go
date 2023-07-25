package pdc_common

import (
	"context"
	"errors"
	stdlog "log"
	"os"
	"runtime/debug"
	"sync"
	"time"

	"cloud.google.com/go/logging"
	zlg "github.com/mark-ignacio/zerolog-gcp"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"google.golang.org/api/option"
)

// var pdclog *log.Logger
var clogcreate sync.Once

func InitializeLogger() {
	clogcreate.Do(func() {
		NewZapLogger()
	})
}

func init() {

}

func NewZapLogger() {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix

	ctx := context.Background()

	config := GetConfig()

	consoleW := zerolog.ConsoleWriter{
		Out: os.Stdout, TimeFormat: time.RFC3339,
	}

	if string(config.Credential) == "" {

		log.Logger = zerolog.New(consoleW).With().Timestamp().Logger()

		log.Info().Msg("no use credentials")
		log.Printf("creating basic logger")
		return
	}

	opt := logging.CommonLabels(map[string]string{
		"version":      config.Version,
		"username":     config.Lisensi.Email,
		"hostname":     config.Hostname,
		"environ_type": "bot",
	})

	gcpWriter, err := zlg.NewCloudLoggingWriter(ctx, config.ProjectID, config.logname, zlg.CloudLoggingOptions{
		ClientOptions: []option.ClientOption{config.CredOption()},
		LoggerOptions: []logging.LoggerOption{opt},
	})
	if err != nil {
		log.Panic().Err(err).Msg("could not create a CloudLoggingWriter")
	}

	multi := zerolog.MultiLevelWriter(consoleW, gcpWriter)
	log.Logger = zerolog.New(multi).With().Timestamp().Logger()

	// log.Print("asdasd asdasdasd asd asdasteast")
}

// func NewLogger() *log.Logger {
// 	ctx := context.Background()

// 	config := GetConfig()

// 	if _, err := os.Stat(config.CredentialPath); errors.Is(err, os.ErrNotExist) {
// 		log.Println(config.CredentialPath, "not exists")
// 		log.Println("creating basic logger")
// 		return log.Default()
// 	}

// 	// Creates a client.
// 	client, err := logging.NewClient(ctx, config.ProjectID, option.WithCredentialsFile(config.CredentialPath))
// 	if err != nil {
// 		log.Fatalf("Failed to create client: %v", err)
// 	}
// 	// defer client.Close()

// 	// Sets the name of the log to write to.
// 	opt := logging.CommonLabels(map[string]string{
// 		"version":       config.Version,
// 		"username":      config.Lisensi.Email,
// 		"hostname":      config.Hostname,
// 		"environt_type": "bot",
// 	})
// 	logger := client.Logger(config.logname, opt).StandardLogger(logging.Info)
// 	return logger
// }

func CapturePanicError() {
	if r := recover(); r != nil {
		err := r.(error)
		log.Panic().Err(err).Str("stacktrace", string(debug.Stack())).Msg(err.Error())
	}
}

func ReportError(err error) error {
	if errors.Is(err, context.Canceled) {
		stdlog.Println(err)
		return err
	}
	log.Error().Err(err).Str("stacktrace", string(debug.Stack())).Msg(err.Error())

	return err
}

func ReportErrorCustom(err error, handler func(event *zerolog.Event) *zerolog.Event) error {
	event := log.Error().Err(err).Str("stacktrace", string(debug.Stack()))
	event = handler(event)
	event.Msg(err.Error())
	return err
}
