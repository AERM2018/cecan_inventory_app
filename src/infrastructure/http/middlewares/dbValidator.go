package middlewares

import (
	"cecan_inventory/src/adapters/helpers"
	"cecan_inventory/src/domain/models"
	bodyreader "cecan_inventory/src/infrastructure/external/bodyReader"
	datasources "cecan_inventory/src/infrastructure/external/dataSources"
	"fmt"

	"github.com/kataras/iris/v12"
)

type DbValidator struct {
	MedicineDataSrc datasources.MedicinesDataSource
	PharmacyDataSrc datasources.PharmacyStocksDataSource
}

func (dbVal DbValidator) IsMedicineInCatalogByKey(ctx iris.Context) {
	var (
		httpRes  models.Responser
		medicine models.Medicine
	)
	medicineKey := ctx.Params().GetString("key")
	if medicineKey == "" { // Get medicine key from body when it's not found in the url
		bodyreader.ReadBodyAsJson(ctx, &medicine, false)
		medicineKey = medicine.Key
	}
	isMedicine := dbVal.MedicineDataSrc.DbPsql.Unscoped().First(&medicine, medicineKey).RowsAffected
	if ctx.Request().Method == "POST" {
		if isMedicine == 1 {
			httpRes = models.Responser{
				StatusCode: iris.StatusNotFound,
				Message:    fmt.Sprintf("El medicamento con clave: %v ya existe.", medicineKey),
			}
			helpers.PrepareAndSendMessageResponse(ctx, httpRes)
			return
		}
	} else {
		if isMedicine == 0 {
			httpRes = models.Responser{
				StatusCode: iris.StatusNotFound,
				Message:    fmt.Sprintf("El medicamento con clave: %v no existe.", medicineKey),
			}
			helpers.PrepareAndSendMessageResponse(ctx, httpRes)
		}

	}
	ctx.Next()
}

func (dbVal DbValidator) IsMedicineWithName(ctx iris.Context) {
	var (
		httpRes  models.Responser
		medicine models.Medicine
	)
	medicineKey := ctx.Params().GetString("key")
	bodyreader.ReadBodyAsJson(ctx, &medicine, false)
	isMedicine := dbVal.MedicineDataSrc.DbPsql.Where("key != ? AND name = ?", medicineKey, medicine.Name).Find(&models.Medicine{}).RowsAffected
	if isMedicine != 0 {
		httpRes = models.Responser{
			StatusCode: iris.StatusNotFound,
			Message:    fmt.Sprintf("Ya existe un medicamento con el nombre: %v.", medicine.Name),
		}
		helpers.PrepareAndSendMessageResponse(ctx, httpRes)
	}
	ctx.Next()
}

// func (dbVal DbValidator) AreStocksOfMedicine(ctx iris.Context) {
// 	var (
// 		httpRes models.Responser
// 		numStocks int64
// 	)
// 	medicine_key := ctx.Params().GetString("key")
// 	numStocks, _ = dbVal.PharmacyDataSrc.GetPharmacyStocksByMedicineKey(medicine_key)
// 	if numStocks > 0 {
// 		httpRes = models.Responser{
// 				StatusCode: iris.StatusNotFound,
// 				Message:    fmt.Sprintf("El medicamento con clave: %v no pudó ser eliminado ya que existen registros en el inventario con fecha de vencimiento", medicineKey),
// 			}
// 			helpers.PrepareAndSendMessageResponse(ctx, httpRes)
// 	}
// 	ctx.Next()
// }
