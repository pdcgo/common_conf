package scenario

import (
	"encoding/json"
	"os"
	"testing"

	"github.com/pdcgo/common_conf/pdc_application"
	"github.com/stretchr/testify/assert"
)

type DirScenario struct {
	t    *testing.T
	Base string
}

func NewDirScenario(t *testing.T) *DirScenario {
	scen := DirScenario{
		t:    t,
		Base: GetBaseTestAsset(""),
	}

	return &scen
}

func (scen *DirScenario) WithConfig(cfg *pdc_application.AppFileConfig, handler func(fname string, cfg *pdc_application.AppFileConfig, scen *DirScenario)) {
	fname := "config.json"

	f, err := os.OpenFile(scen.Path(fname), os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0755)
	assert.Nil(scen.t, err)

	err = json.NewEncoder(f).Encode(cfg)
	assert.Nil(scen.t, err)

	f.Close()

	handler(fname, cfg, scen)

}
