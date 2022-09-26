package middlewares

import (
	"cecan_inventory/adapters/helpers"
	"cecan_inventory/domain/common"
	"cecan_inventory/domain/models"
	bodyreader "cecan_inventory/infrastructure/external/bodyReader"
	datasources "cecan_inventory/infrastructure/external/dataSources"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"github.com/kataras/iris/v12"
)

type DbValidator struct {
	MedicineDataSrc        datasources.MedicinesDataSource
	PharmacyDataSrc        datasources.PharmacyStocksDataSource
	RolesDataSource        datasources.RolesDataSource
	UserDataSource         datasources.UserDataSource
	PrescriptionDataSource datasources.PrescriptionsDataSource
}

func (dbVal DbValidator) IsRoleId(ctx iris.Context) {
	var (
		httpRes models.Responser
		user    models.User
	)
	bodyreader.ReadBodyAsJson(ctx, &user, false)
	_, err := dbVal.RolesDataSource.GetRoleById(user.RoleId)
	if err != nil {
		httpRes = models.Responser{
			StatusCode: iris.StatusNotFound,
			Message:    fmt.Sprintf("El rol con id: %v no existe.", user.RoleId),
		}
		helpers.PrepareAndSendMessageResponse(ctx, httpRes)
	}
	ctx.Next()
}

func (dbVal DbValidator) CanUserDoAction(roleNamesAllowed ...string) func(ctx iris.Context) {
	return func(ctx iris.Context) {
		var httpRes models.Responser
		// This is just for testing, include Admin role in all request
		roleNamesAllowed = append(roleNamesAllowed, "Admin")
		roleName := fmt.Sprintf("%v", ctx.Values().Get("roleName"))
		if !common.FindElementInSlice(roleName, roleNamesAllowed) {
			httpRes = models.Responser{
				StatusCode: iris.StatusForbidden,
				Message:    "Acción denegada, no cuenta con los permisos necesarios.",
			}
			helpers.PrepareAndSendMessageResponse(ctx, httpRes)
		}
		ctx.Next()
	}
}
func (dbVal DbValidator) IsEmail(ctx iris.Context) {
	var (
		httpRes models.Responser
		user    models.User
	)
	bodyreader.ReadBodyAsJson(ctx, &user, false)
	_, err := dbVal.UserDataSource.GetUserByEmail(user.Email)
	if err == nil {
		httpRes = models.Responser{
			StatusCode: iris.StatusBadRequest,
			Message:    fmt.Sprintf("El email %v ya está siendo usado por otro usuario.", user.Email),
		}
		helpers.PrepareAndSendMessageResponse(ctx, httpRes)
	}
	ctx.Next()
}

func (dbVal DbValidator) IsMedicineInCatalogByKey(ctx iris.Context) {
	var (
		httpRes               models.Responser
		medicine              models.Medicine
		pharmacyStockToUpdate models.PharmacyStockToUpdate
		medicineKey           string
	)
	routePath := ctx.Path()
	if strings.Contains(routePath, "pharmacy_inventory") {
		if ctx.Request().Method == "PUT" {
			bodyreader.ReadBodyAsJson(ctx, &pharmacyStockToUpdate, false)
			medicineKey = pharmacyStockToUpdate.MedicineKey
		} else if ctx.Request().Method == "POST" {
			medicineKey = ctx.Params().GetString("key")
		}
	} else {
		medicineKey = ctx.Params().GetString("key")
		if medicineKey == "" { // Get medicine key from body when it's not found in the url
			bodyreader.ReadBodyAsJson(ctx, &medicine, false)
			medicineKey = medicine.Key
		}
	}
	isMedicine := dbVal.MedicineDataSrc.DbPsql.Unscoped().Where("key = ?", medicineKey).First(&medicine).RowsAffected
	if isMedicine == 0 {
		httpRes = models.Responser{
			StatusCode: iris.StatusNotFound,
			Message:    fmt.Sprintf("El medicamento con clave: %v no existe.", medicineKey),
		}
		helpers.PrepareAndSendMessageResponse(ctx, httpRes)
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
		var resMessage string
		isMedicineWithKey := dbVal.MedicineDataSrc.DbPsql.Where("key = ?", medicine.Key).Find(&models.Medicine{}).RowsAffected
		if isMedicineWithKey == 1 {
			if ctx.Request().Method == "PUT" {
				resMessage = fmt.Sprintf("No se actualizó el medicamento debido a que ya existe un medicamento con la clave: %v.", medicine.Key)
			} else {
				resMessage = fmt.Sprintf("El medicamento con clave: %v ya existe.", medicine.Key)
			}
			httpRes = models.Responser{
				StatusCode: iris.StatusBadRequest,
				Message:    resMessage,
			}
			helpers.PrepareAndSendMessageResponse(ctx, httpRes)
		}
	}
	ctx.Next()
}

func (dbVal DbValidator) IsMedicineDeleted(ctx iris.Context) {
	var (
		httpRes models.Responser
		isError bool
		me      models.Medicine
	)
	medicineKey := ctx.Params().GetString("key")
	reqPath := ctx.Path()
	IsMedicineDeleted := dbVal.MedicineDataSrc.DbPsql.Where("key = ?", medicineKey).First(&me).RowsAffected
	var message string
	if strings.Contains(reqPath, "pharmacy_inventory") && ctx.Request().Method == "POST" && IsMedicineDeleted == 0 {
		isError = true
		message = fmt.Sprintf("No se pudó ingresar el stock a farmacia del medicamento con clave: %v debido a que esta inactivo, activelo y vuelvalo a intentar.", medicineKey)
	} else if strings.Contains(reqPath, "reactivate") && ctx.Request().Method == "PUT" && IsMedicineDeleted != 0 {
		isError = true
		message = fmt.Sprintf("El medicamento con clave: %v no se reactivó debido a que no ha sido eliminado antes.", medicineKey)
	}
	if isError {
		httpRes = models.Responser{
			StatusCode: iris.StatusBadRequest,
			Message:    message,
		}
		helpers.PrepareAndSendMessageResponse(ctx, httpRes)
		return
	}
	ctx.Next()
}

func (dbVal DbValidator) IsPharmacyStockUsed(ctx iris.Context) {
	var (
		httpRes models.Responser
	)
	pharmacyStockId := ctx.Params().GetString("id")
	pharmacyStockUuid, _ := uuid.Parse(pharmacyStockId)
	isStockUsed, _ := dbVal.PharmacyDataSrc.IsStockUsed(pharmacyStockUuid)
	if isStockUsed {
		httpRes = models.Responser{
			StatusCode: iris.StatusBadRequest,
			Message:    "No se puede eliminar un stock de farmacia cuando ha sido utilizado.",
		}
		helpers.PrepareAndSendMessageResponse(ctx, httpRes)
		return
	}
	ctx.Next()
}

func (dbVal DbValidator) IsPharmacyStockById(ctx iris.Context) {
	var httpRes models.Responser
	pharmacyStockId := ctx.Params().GetString("id")
	pharmacyStockUuid, _ := uuid.Parse(pharmacyStockId)
	_, err := dbVal.PharmacyDataSrc.GetPharmacyStockById(pharmacyStockUuid)
	if err != nil {
		httpRes = models.Responser{
			StatusCode: iris.StatusBadRequest,
			Message:    fmt.Sprintf("El stock de farmacia con id: %v no existe.", pharmacyStockId),
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
