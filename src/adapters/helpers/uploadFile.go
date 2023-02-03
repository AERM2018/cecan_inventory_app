package helpers

import (
	"errors"
	"os"
	"path"

	"github.com/kataras/iris/v12"
)

func UploadFile(ctx iris.Context, folderName string, fileName string, formFieldName string) (string, error) {
	ctx.SetMaxRequestBodySize(86 * iris.MB)
	cwd, _ := os.Getwd()
	_, fileHeader, err := ctx.FormFile(formFieldName)
	if err != nil {
		return "", errors.New("No se ha podido leer el archivo.")
	}
	dest := path.Join(cwd, "domain", folderName, fileName)
	ctx.SaveFormFile(fileHeader, dest)
	return dest, nil
}
