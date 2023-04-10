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


cara login
```
client := auth.NewAuthClient("https://hostname/v2/login")
	err := client.Login(license.Email, license.Pwd, botID, string(version))

	if err != nil {
		log.Error().Msg(err.Error())
		time.Sleep(time.Minute)
		panic(err)
	}
```