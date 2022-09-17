package middlewares

import (
	"cecan_inventory/adapters/helpers"
	"cecan_inventory/domain/models"
	bodyreader "cecan_inventory/infrastructure/external/bodyReader"
	"strings"

	"github.com/kataras/iris/v12"
)

func ValidateRequest(validator func(mapObject interface{}, omitedFields ...string) error, omitedFields ...string) func(ctx iris.Context) {
	return func(ctx iris.Context) {
		var intialStruct interface{}
		bodyreader.ReadBodyAsJson(ctx, &intialStruct, false)
		validatorError := validator(intialStruct, strings.Join(omitedFields, ","))
		if validatorError != nil {
			helpers.PrepareAndSendDataResponse(ctx, parseErrorToStruct(validatorError))
			return
		}
		ctx.Next()
	}
}

func parseErrorToStruct(err error) models.Responser {
	errorsMessage := make(map[string]string, 0)
	errorsLines := make([]string, 0)
	if strings.Contains(err.Error(), ";") {
		errorsLines = strings.Split(err.Error(), ";")
	} else {
		errorsLines = append(errorsLines, err.Error())
	}
	for _, errorLine := range errorsLines {
		keyPair := strings.Split(errorLine, ":")
		errorsMessage[strings.Trim(keyPair[0], " ")] = strings.Trim(keyPair[1], " ")
	}
	return models.Responser{
		StatusCode: iris.StatusBadRequest,
		Data:       iris.Map{"errors": errorsMessage},
	}
}
