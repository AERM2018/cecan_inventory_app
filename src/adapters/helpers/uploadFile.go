package helpers

import (
	"fmt"
	"os"
	"path"

	"github.com/kataras/iris/v12"
)

func UploadFile(ctx iris.Context, folderName string, fileName string, formFieldName string) (string, error) {
	fmt.Printf("form field name %v",formFieldName)
	ctx.SetMaxRequestBodySize(86 * iris.MB)
	cwd, _ := os.Getwd()
	_, fileHeader, err := ctx.FormFile(formFieldName)
	fmt.Printf("file %v",fileHeader)
	if err != nil {
		return "", err
	}
	dest := path.Join(cwd, "domain", folderName, fileName)
	ctx.SaveFormFile(fileHeader, dest)
	return dest, nil
}
