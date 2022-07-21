package helpers

import (
	"cecan_inventory/src/domain/models"

	"github.com/kataras/iris/v12"
)

func PrepareAndSendMessageResponse(c iris.Context, response models.Responser) {
	var mapResponse iris.Map
	var ok bool = true
	if response.StatusCode > 400 {
		ok = false
		if response.StatusCode == 500 {
			mapResponse = iris.Map{"ok": ok, "message": "Hable con el administrador."}
			c.JSON(mapResponse)
			return
		}
	}
	mapResponse = iris.Map{"ok": ok, "message": response.Message}
	c.StatusCode(response.StatusCode)
	c.JSON(mapResponse)
}

func PrepareAndSendDataResponse(c iris.Context, response models.Responser) {
	var mapResponse iris.Map
	var ok bool = true
	mapResponse = iris.Map{"ok": ok, "data": response.Data}
	c.StatusCode(response.StatusCode)
	c.JSON(mapResponse)
}
