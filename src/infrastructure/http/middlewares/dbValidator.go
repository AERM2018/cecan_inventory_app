package middlewares

import (
	"cecan_inventory/adapters/helpers"
	"cecan_inventory/domain/models"
	bodyreader "cecan_inventory/infrastructure/external/bodyReader"
	datasources "cecan_inventory/infrastructure/external/dataSources"
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
		httpRes     models.Responser
		medicine    models.Medicine
		medicineRes models.Medicine
	)
	medicineKey := ctx.Params().GetString("key")
	bodyreader.ReadBodyAsJson(ctx, &medicine, false)
	isMedicine := dbVal.MedicineDataSrc.DbPsql.Unscoped().Where("key != ? AND name = ?", medicineKey, medicine.Name).Find(&medicineRes).RowsAffected
	if isMedicine != 0 {
		var (
			resMessage string
		)
		if ctx.Request().Method == "PUT" {
			resMessage = fmt.Sprintf("No se actualizó el medicamento debido a que ya existe un medicamento con el nombre: %v.", medicine.Name)
		} else {
			resMessage = fmt.Sprintf("Ya existe un medicamento con el nombre: %v.", medicine.Name)
		}
		if medicineRes.DeletedAt.Valid {
			resMessage += " (DESABILITADO)"
		}

		httpRes = models.Responser{
			StatusCode: iris.StatusBadRequest,
			Message:    resMessage,
		}
		helpers.PrepareAndSendMessageResponse(ctx, httpRes)
		return
	}
	ctx.Next()
}

func (dbVal DbValidator) IsMedicineWithKey(ctx iris.Context) {
	var (
		httpRes  models.Responser
		medicine models.Medicine
	)
	bodyreader.ReadBodyAsJson(ctx, &medicine, false)
	medicineKey := ctx.Params().GetString("key")
	if medicine.Key != medicineKey {
		isMedicineWithKey := dbVal.MedicineDataSrc.DbPsql.Where("key = ?", medicine.Key).Find(&models.Medicine{}).RowsAffected
		if isMedicineWithKey == 1 {
			httpRes = models.Responser{
				StatusCode: iris.StatusBadRequest,
				Message:    fmt.Sprintf("No se actualizó el medicamento debido a que ya existe un medicamento con la clave: %v.", medicine.Key),
			}
			helpers.PrepareAndSendMessageResponse(ctx, httpRes)
		}
	}
	ctx.Next()
}

func (dbVal DbValidator) IsMedicineDeleted(ctx iris.Context) {
	var (
		httpRes models.Responser
	)
	medicineKey := ctx.Params().GetString("key")
	IsMedicineDeleted := dbVal.MedicineDataSrc.DbPsql.Where("key = ?", medicineKey).Find(&models.Medicine{}).RowsAffected
	if IsMedicineDeleted != 0 {
		httpRes = models.Responser{
			StatusCode: iris.StatusBadRequest,
			Message:    fmt.Sprintf("El medicamento con clave: %v no se reactivó debido a que no ha sido eliminado antes.", medicineKey),
		}
		helpers.PrepareAndSendMessageResponse(ctx, httpRes)
		return
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
