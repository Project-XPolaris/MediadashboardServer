package httpapi

import (
	"mediadashboard/service"
	"net/http"

	"github.com/allentom/haruka"
	"github.com/project-xpolaris/youplustoolkit/youlink"
)

var generateAccessCodeWithYouAuthHandler haruka.RequestHandler = func(context *haruka.Context) {
	code := context.GetQueryString("code")
	accessToken, username, err := service.GenerateYouAuthToken(code)
	if err != nil {
		youlink.AbortErrorWithStatus(err, context, http.StatusInternalServerError)
		return
	}
	context.JSON(haruka.JSON{
		"success": true,
		"data": haruka.JSON{
			"accessToken": accessToken,
			"username":    username,
		},
	})
}

type LoginUserRequestBody struct {
	Username    string `json:"username"`
	Password    string `json:"password"`
	WithYouPlus bool   `json:"withYouPlus"`
}

type UserAuthRequestBody struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

var generateAccessCodeWithYouAuthPasswordHandler haruka.RequestHandler = func(context *haruka.Context) {
	var requestBody UserAuthRequestBody
	err := context.ParseJson(&requestBody)
	if err != nil {
		AbortError(context, err, http.StatusBadRequest)
		return
	}
	accessToken, username, err := service.GenerateYouAuthTokenByPassword(requestBody.Username, requestBody.Password)
	if err != nil {
		AbortError(context, err, http.StatusInternalServerError)
		return
	}
	context.JSON(haruka.JSON{
		"success": true,
		"data": haruka.JSON{
			"accessToken": accessToken,
			"username":    username,
		},
	})
}
var checkTokenHandler haruka.RequestHandler = func(context *haruka.Context) {
	context.JSON(haruka.JSON{
		"success": true,
		"data": haruka.JSON{
			"isValid": true,
		},
	})
}
