package helpers

import (
	"errors"
	"os"
	"path/filepath"

	"github.com/kataras/iris/v12"
)

func UploadFile(ctx iris.Context, path string, fileName string) (string, error) {
	ctx.SetMaxRequestBodySize(86 * iris.MB)
	cwd, _ := os.Getwd()
	formDataFileField := "excel_file"
	_, fileHeader, err := ctx.FormFile(formDataFileField)
	if err != nil {
		return "", errors.New("No se ha podido leer el archivo.")
	}
	dest := filepath.Join(cwd, "domain", path, fileName)
	ctx.SaveFormFile(fileHeader, dest)
	return dest, nil
}
