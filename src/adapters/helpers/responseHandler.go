package helpers

import "github.com/kataras/iris/v12"

type response struct {
	StatusCode int
	Error      error
	Message    string
	Data       iris.Map
}

func PrepareAndSendMessageResponse(c iris.Context,statusCode int, err error, message string) {
	var mapResponse iris.Map
	var ok bool = true
	if statusCode > 400{
		ok = false
		if statusCode == 500{
			mapResponse = iris.Map{"ok":ok,"message":"Hable con el administrador."}
			c.JSON(mapResponse)
			return
		}
	}
	mapResponse = iris.Map{"ok":ok,"message":message}	
	c.StatusCode(statusCode)
	c.JSON(mapResponse)
}

func PrepareAndSendDataResponse(c iris.Context,statusCode int, data iris.Map) {
	var mapResponse iris.Map
	var ok bool
	mapResponse = iris.Map{"ok":ok,"data":data}	
	c.StatusCode(statusCode)
	c.JSON(mapResponse)
}
