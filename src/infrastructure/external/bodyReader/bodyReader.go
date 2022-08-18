package bodyreader

import (
	"bytes"
	"encoding/json"
	"io"

	"github.com/kataras/iris/v12"
)

func ReadBodyAsJson(ctx iris.Context, dest interface{}, consumeBody bool) {
	body, _ := io.ReadAll(ctx.Request().Body)
	json.Unmarshal(body, dest)
	if !consumeBody {
		ctx.Request().Body = io.NopCloser(bytes.NewBuffer(body))
	}
}
