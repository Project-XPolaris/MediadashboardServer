package httpapi

import "github.com/allentom/harukap/module/errorhandler"

const (
	CustomErrorCode           = "1001"
	ResourceNotFoundErrorCode = "1002"
)

var appErrorHandler = errorhandler.NewErrorModule()

type ResourceNotFoundError struct {
}

func (r *ResourceNotFoundError) Error() string {
	return "resource not found"
}

func InitErrorHandler() {
	appErrorHandler.RegisterHandler(errorhandler.ErrorHandler{
		Match:  &ResourceNotFoundError{},
		Code:   ResourceNotFoundErrorCode,
		Status: 404,
	})
}
