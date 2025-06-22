package main

import (
	"mediadashboard/application/httpapi"
	"mediadashboard/config"
	"mediadashboard/database"
	"mediadashboard/module"
	"mediadashboard/plugins"
	"mediadashboard/youplus"
	"net/http"

	"github.com/allentom/haruka"
	"github.com/allentom/harukap"
	"github.com/allentom/harukap/cli"
	"github.com/sirupsen/logrus"
)

func main() {
	err := config.InitConfigProvider()
	if err != nil {
		logrus.Fatal(err)
	}
	err = plugins.DefaultYouLogPlugin.OnInit(config.DefaultConfigProvider)
	if err != nil {
		logrus.Fatal(err)
	}
	appEngine := harukap.NewHarukaAppEngine()
	appEngine.ConfigProvider = config.DefaultConfigProvider
	appEngine.LoggerPlugin = plugins.DefaultYouLogPlugin
	appEngine.UsePlugin(youplus.DefaultYouPlusPlugin)
	appEngine.UsePlugin(database.DefaultPlugin)
	if config.Instance.YouAuthConfig != nil {
		plugins.CreateYouAuthPlugin()
		plugins.DefaultYouAuthOauthPlugin.ConfigPrefix = config.Instance.YouAuthConfigPrefix
		appEngine.UsePlugin(plugins.DefaultYouAuthOauthPlugin)
	}
	// init module
	module.CreateAuthModule()
	module.Auth.AuthMiddleware.OnError = func(ctx *haruka.Context, err error) {
		httpapi.AbortError(ctx, err, http.StatusForbidden)
		ctx.Abort()
	}
	appEngine.HttpService = httpapi.GetEngine()
	if err != nil {
		logrus.Fatal(err)
	}
	appWrap, err := cli.NewWrapper(appEngine)
	if err != nil {
		logrus.Fatal(err)
	}
	appWrap.RunApp()
}
