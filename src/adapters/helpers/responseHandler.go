package helpers

import (
	"cecan_inventory/domain/models"
	"fmt"

	"github.com/kataras/iris/v12"
)

func PrepareAndSendMessageResponse(c iris.Context, response models.Responser) {
	var mapResponse iris.Map
	var ok bool = true
	var tag = "message"
	if response.StatusCode >= 400 {
		tag = "error"
		ok = false
		if response.StatusCode == 500 {
			mapResponse = iris.Map{"ok": ok, "message": "Hable con el administrador!"}
			fmt.Printf("Error: %v", response.Err)
			c.StopWithJSON(response.StatusCode, mapResponse)
			return
		}
	}
	mapResponse = iris.Map{"ok": ok, tag: response.Message}
	c.StopWithJSON(response.StatusCode, mapResponse)
}

func PrepareAndSendDataResponse(c iris.Context, response models.Responser) {
	var mapResponse iris.Map
	var ok bool = true
	var responseTag = "data"
	if response.StatusCode >= 400 {
		ok = false
		responseTag = "error"
	}
	mapResponse = iris.Map{"ok": ok, responseTag: response.Data}
	c.StatusCode(response.StatusCode)
	c.StopWithJSON(response.StatusCode, mapResponse)
}
