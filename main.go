package main

import (
	log "github.com/sirupsen/logrus"
	"github.com/vapor-ware/synse-juniper-jti-plugin/pkg"
	"github.com/vapor-ware/synse-sdk/sdk"
)

const (
	pluginName       = "juniper jti"
	pluginMaintainer = "vaporio"
	pluginDesc       = "Streamed networking data with Juniper JTI over UDP"
	pluginVcs        = "github.com/vapor-ware/synse-juniper-jti-plugin"
)

func main() {
	sdk.SetPluginInfo(
		pluginName,
		pluginMaintainer,
		pluginDesc,
		pluginVcs,
	)

	plugin, err := pkg.MakePlugin()
	if err != nil {
		log.Fatal(err)
	}

	if err := plugin.Run(); err != nil {
		log.Fatal(err)
	}
}
