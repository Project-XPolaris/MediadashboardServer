package httpapi

import (
	"mediadashboard/module"
	"mediadashboard/service"
	"net/http"

	"github.com/allentom/haruka"
)

var helloHandler haruka.RequestHandler = func(context *haruka.Context) {
	context.JSON(haruka.JSON{
		"message": service.ExampleService(),
	})
}
var makeErrorHandler haruka.RequestHandler = func(context *haruka.Context) {
	appErrorHandler.RaiseHttpError(context, &ResourceNotFoundError{})
}

var serviceInfoHandler haruka.RequestHandler = func(context *haruka.Context) {
	authMaps, err := module.Auth.GetAuthConfig()
	if err != nil {
		AbortError(context, err, http.StatusInternalServerError)
		return
	}
	context.JSON(haruka.JSON{
		"success": true,
		"name":    "YouPhoto service",
		"oauth":   true,
		"auth":    authMaps,
	})
}
