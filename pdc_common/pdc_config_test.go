package pdc_common_test

import (
	"errors"
	"os"
	"path/filepath"
	"testing"

	"github.com/pdcgo/common_conf/pdc_common"
	"github.com/stretchr/testify/assert"
)

func TestSendLog(t *testing.T) {
	cred, _ := filepath.Abs("../logger_credentials.json")
	credbyte, _ := os.ReadFile(cred)
	pdc_common.SetConfig("../config.yml", "beta", "golang_test_log", credbyte)
	// logg := pdc_common.GetLogger()

	// logg.Println("asdasdasdasdasdasdas tes")

	t.Run("zap logger", func(t *testing.T) {
		defer func() {
			recover()
		}()
		pdc_common.NewZapLogger()

		defer pdc_common.CapturePanicError()

		panic(errors.New("asdasd tes"))

	})
}

func TestSetConfig(t *testing.T) {
	cred, _ := filepath.Abs("../logger_credentials.json")
	credbyte, _ := os.ReadFile(cred)
	pdc_common.SetConfig("../config.yml", "beta", "golang_test_log", credbyte)

	config := pdc_common.GetConfig()

	assert.Equal(t, config.Version, "beta")
	assert.NotEmpty(t, config.Credential)
}
