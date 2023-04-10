# common_conf

example initialize credentials

```
package main

import (
	_ "embed"

	"github.com/pdcgo/common_conf/pdc_common"
)

//go:embed logger_credentials.json
var credentialsByte []byte

//go:embed version
var version []byte

func beforeRunning() {
	pdc_common.SetConfig("config.yml", string(version), "tiktok_chat", credentialsByte)

}

```