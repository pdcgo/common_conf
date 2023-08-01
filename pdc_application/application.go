package pdc_application

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"log"
	"os"
	"path"
	"path/filepath"
	"time"

	"cloud.google.com/go/logging"
	zlg "github.com/mark-ignacio/zerolog-gcp"
	"github.com/pdcgo/common_conf/auth"
	"github.com/rs/zerolog"
	zlog "github.com/rs/zerolog/log"
	"google.golang.org/api/option"
	"gopkg.in/natefinch/lumberjack.v2"
	"gopkg.in/yaml.v3"
)

type Lisensi struct {
	Email string `json:"email" yaml:"email"`
	Pwd   string `json:"pwd" yaml:"pwd"`
}

type AppFileConfig struct {
	Lisensi   Lisensi `json:"lisensi" yaml:"lisensi"`
	ProjectID string  `json:"project_id" yaml:"project_id"`
	LogToFile string  `json:"log_to_file" yaml:"log_to_file"`
}

type BaseApplication interface {
	Path(path ...string) string
}

type PdcApplication struct {
	Base          BaseApplication
	Credential    []byte
	Version       string
	AppID         int
	LogHelper     *LogHelper
	ReplaceLogger bool

	Auth *auth.AuthClient
}

func (app *PdcApplication) RunWithLicenseFile(cfgname string, logname string, handle func(app *PdcApplication)) error {
	cfg, err := app.getAppFileConfig(cfgname)

	if err != nil {
		return err
	}

	if app.Auth == nil {
		app.Auth = auth.NewAuthClient("https://pdcoke.com/v2/login")
	}

	err = app.Auth.Login(cfg.Lisensi.Email, cfg.Lisensi.Pwd, app.AppID, app.Version)
	if err != nil {
		log.Println("[", cfg.Lisensi.Email, "]", err)
		time.Sleep(time.Hour)
		panic(err)
	}

	return app.Run(cfg, logname, handle)
}

func (app *PdcApplication) AuthenticateEmail() error {
	panic("not implemented")
}

func (app *PdcApplication) Run(cfg *AppFileConfig, logName string, handle func(app *PdcApplication)) error {
	logger, err := app.CreatingLogger(cfg, logName)
	if err != nil {
		return err
	}

	if app.ReplaceLogger {
		zlog.Logger = logger
	}
	app.LogHelper = &LogHelper{
		logger: &logger,
	}

	handle(app)

	defer app.LogHelper.CapturePanicError()
	return nil
}

func (app *PdcApplication) CreatingLogger(cfg *AppFileConfig, logName string) (zerolog.Logger, error) {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix

	hostname, _ := os.Hostname()

	opt := logging.CommonLabels(map[string]string{
		"version":      app.Version,
		"username":     cfg.Lisensi.Email,
		"hostname":     hostname,
		"environ_type": "bot",
	})

	logwriters := []io.Writer{
		zerolog.ConsoleWriter{
			Out: os.Stdout, TimeFormat: time.RFC3339,
		},
	}

	if len(app.Credential) != 0 {
		gcpWriter, err := zlg.NewCloudLoggingWriter(context.Background(), cfg.ProjectID, logName, zlg.CloudLoggingOptions{
			ClientOptions: []option.ClientOption{option.WithCredentialsJSON(app.Credential)},
			LoggerOptions: []logging.LoggerOption{opt},
		})
		if err != nil {
			panic("could not create a CloudLoggingWriter")
		}

		logwriters = append(logwriters, gcpWriter)
	} else {
		log.Println("[ warning ] data error not sending to cloud")
	}

	if cfg.LogToFile != "" {
		logwriters = append(logwriters,
			app.createLogRollingFile(cfg.LogToFile, logName),
		)
	}

	multi := zerolog.MultiLevelWriter(logwriters...)

	logger := zerolog.New(multi).With().Timestamp().Logger()

	return logger, nil

}

func (app *PdcApplication) createLogRollingFile(loglocation string, logname string) io.Writer {
	os.MkdirAll(loglocation, 0744)

	return &lumberjack.Logger{
		Filename:   path.Join(loglocation, logname),
		MaxBackups: 3,  // files
		MaxSize:    10, // megabytes
		MaxAge:     30, // days
	}
}

func (app *PdcApplication) getAppFileConfig(fname string) (*AppFileConfig, error) {
	cfg := AppFileConfig{
		ProjectID: "shopeepdc",
	}

	locfname := app.Base.Path(fname)
	if _, err := os.Stat(fname); errors.Is(err, os.ErrNotExist) {
		log.Println("config", locfname, "tidak ada....")
		return &cfg, errors.New("config " + locfname + " tidak ada....")
	}

	ext := filepath.Ext(locfname)
	f, err := os.Open(locfname)
	if err != nil {
		return &cfg, err
	}

	defer f.Close()

	switch ext {
	case "json":
		err = json.NewDecoder(f).Decode(&cfg)
	case "yaml":
		err = yaml.NewDecoder(f).Decode(&cfg)
	case "yml":
		err = yaml.NewDecoder(f).Decode(&cfg)
	default:
		err = errors.New(locfname + " configuration format not supported")
	}

	return &cfg, err
}
