package pdc_application_test

import (
	"log"
	"os"
	"testing"

	"github.com/pdcgo/common_conf/pdc_application"
	"github.com/pdcgo/common_conf/scenario"
	"github.com/stretchr/testify/assert"
)

type MockAuth struct {
}

func (auth *MockAuth) Login(email string, password string, botID int, version string) error {
	return nil
}

func TestApplication(t *testing.T) {

	scen := scenario.NewDirScenario(t)

	data, err := os.ReadFile(scenario.GetBaseTestAsset("../logger_credentials.json"))
	assert.Nil(t, err)
	app := pdc_application.PdcApplication{
		Base:       scen,
		Credential: data,
		Version:    "development",
		AppID:      1,
		Auth:       &MockAuth{},
	}

	scen.WithConfig(&pdc_application.AppFileConfig{
		ProjectID: "shopeepdc",
	}, func(fname string, cfg *pdc_application.AppFileConfig, scen *scenario.DirScenario) {
		err = app.RunWithLicenseFile(fname, "common_conf", func(app *pdc_application.PdcApplication) error {

			log.Println("running common conf")

			return nil
		})

		assert.Nil(t, err)
	})

	// pdc_application.PdcApplication{
	// 	Base          BaseApplication
	// 	Credential    []byte
	// 	Version       string
	// 	AppID         int
	// 	LogHelper     *LogHelper
	// 	ReplaceLogger bool
	// 	OnPanic       []func(err error)
	// 	OnError       []func(err error)
	// 	OnStartup     []func(app *PdcApplication) error

	// 	Auth *auth.AuthClient
	// }
}
